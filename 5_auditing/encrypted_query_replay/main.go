// Copyright 2023 StrongDM Inc
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package main

import (
	"bytes"
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"log"
	"os"
	"time"

	sdm "github.com/strongdm/strongdm-sdk-go/v3"
)

func main() {
	//	Load the SDM API keys from the environment.
	//	If these values are not set in your environment,
	//	please follow the documentation here:
	//	https://www.strongdm.com/docs/api/api-keys/
	accessKey := os.Getenv("SDM_API_ACCESS_KEY")
	secretKey := os.Getenv("SDM_API_SECRET_KEY")
	if accessKey == "" || secretKey == "" {
		log.Fatal("SDM_API_ACCESS_KEY and SDM_API_SECRET_KEY must be provided")
	}

	// Load the private key for query and replay decryption.
	// This environment variable should contain the path to the private encryption
	// key configured for StrongDM remote log encryption.
	privateKeyFile := os.Getenv("SDM_LOG_PRIVATE_KEY_FILE")
	if privateKeyFile == "" {
		log.Fatal("SDM_LOG_PRIVATE_KEY_FILE must be provided for this example")
	}
	privateKey, err := loadPrivateKeyFromFile(privateKeyFile)
	if err != nil {
		log.Fatal("failed to load private key: %v", err)
	}

	// Create the client
	client, err := sdm.New(
		accessKey,
		secretKey,
	)
	if err != nil {
		log.Fatalf("could not create client: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// You'll need an SSH resource that has had queries made against it, provide its name:
	resourceName := "Example"
	resourceResp, err := client.Resources().List(ctx, "name:?", resourceName)
	if err != nil {
		log.Fatalf("failed to list resources: %v", err)
	}
	if !resourceResp.Next() {
		log.Fatalf("couldn't find resource named %v (error: %v)", resourceName, resourceResp.Err())
	}
	resource := resourceResp.Value()

	fmt.Printf("Queries made against %v:\n", resourceName)
	queries, err := client.Queries().List(ctx, "resource_id:?", resource.GetID())
	if err != nil {
		log.Fatalf("failed to list queries: %v", err)
	}
	for queries.Next() {
		q := queries.Value()
		accountResp, err := client.SnapshotAt(q.Timestamp).Accounts().Get(ctx, q.AccountID)
		if err != nil {
			log.Fatalf("failed to get account: %v", err)
		}
		user := accountResp.Account.(*sdm.User)

		if q.Encrypted {
			fmt.Println("Decrypting encrypted query")
			queryBody, err := base64.StdEncoding.DecodeString(q.QueryBody)
			if err != nil {
				log.Fatalf("failed to decode query body: %v", err)
			}
			q.QueryBody, err = decryptQueryData(privateKey, q.QueryKey, queryBody)
			if err != nil {
				log.Fatalf("failed to decrypt query body: %v", err)
			}
			var capture struct{ Type string }
			if err := json.Unmarshal([]byte(q.QueryBody), &capture); err != nil {
				log.Fatalf("failed to unmarshal query JSON %v: %v", q.QueryBody, err)
			}
			q.Replayable = capture.Type == "shell"
		}

		if q.Replayable {
			fmt.Printf("Replaying query made by %v at %v\n", user.Email, q.Timestamp)
			replayParts, err := client.Replays().List(ctx, "id:?", q.ID)
			if err != nil {
				log.Fatalf("failed to scan replay: %v", err)
			}
			for replayParts.Next() {
				part := replayParts.Value()
				if q.Encrypted {
					partData, err := decryptQueryData(privateKey, q.QueryKey, part.Data)
					if err != nil {
						log.Fatalf("failed to decrypt replay data: %v", err)
					}
					var events []struct {
						Data     []byte
						Duration int64
					}
					if err := json.Unmarshal([]byte(partData), &events); err != nil {
						log.Fatalf("failed to unmarshal events JSON %v: %v", partData, err)
					}
					for _, e := range events {
						event := &sdm.ReplayChunkEvent{
							Data:     e.Data,
							Duration: time.Millisecond * time.Duration(e.Duration),
						}
						part.Events = append(part.Events, event)
					}
				}

				for _, ev := range part.Events {
					// Some characters may not be printed cleanly by this method
					fmt.Print(string(ev.Data))
					time.Sleep(ev.Duration)
				}
			}
			if err := replayParts.Err(); err != nil {
				log.Fatalf("failed to iterate replay: %v", err)
			}
			fmt.Println("")
		} else {
			var capture struct{ Command string }
			if err := json.Unmarshal([]byte(q.QueryBody), &capture); err != nil {
				log.Fatalf("failed to unmarshal query JSON %v: %v", q.QueryBody, err)
			}
			fmt.Printf("Command run by %v at %v: %v\n", user.Email, q.Timestamp, capture.Command)
		}
	}
}

func loadPrivateKeyFromFile(privateKeyFile string) (*rsa.PrivateKey, error) {
	privateKeyBytes, err := os.ReadFile(privateKeyFile)
	if err != nil {
		return nil, err
	}
	pemBlock, _ := pem.Decode(privateKeyBytes)
	return x509.ParsePKCS1PrivateKey(pemBlock.Bytes)
}

// This method demonstrates how to decrypt encrypted query/replay data
func decryptQueryData(privateKey *rsa.PrivateKey, encryptedQueryKey string, encryptedData []byte) (string, error) {
	// Use the organization's private key to decrypt the symmetric key
	queryKeyBytes, err := base64.StdEncoding.DecodeString(encryptedQueryKey)
	if err != nil {
		return "", err
	}
	symmetricKey, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, privateKey, queryKeyBytes, nil)
	if err != nil {
		return "", err
	}

	// Use the symmetric key to decrypt the data
	block, err := aes.NewCipher(symmetricKey)
	if err != nil {
		return "", err
	}
	if len(encryptedData) < block.BlockSize() {
		return "", fmt.Errorf("ciphertext is smaller than AES block size %v", block.BlockSize())
	}
	iv := encryptedData[:block.BlockSize()]
	ciphertext := encryptedData[block.BlockSize():]

	plaintext := make([]byte, len(ciphertext))
	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(plaintext, ciphertext)

	return string(bytes.TrimRight(plaintext, "\x00")), nil
}
