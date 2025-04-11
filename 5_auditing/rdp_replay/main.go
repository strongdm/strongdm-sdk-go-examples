// Copyright 2025 StrongDM Inc
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
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
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

	// You'll need an RDP resource that has had queries made against it, provide its name:
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
			fmt.Printf("Skipping encrypted query made by %v at %v\n", user.Email, q.Timestamp)
			fmt.Println("See encrypted_query_replay for an example of query decryption.")
		} else if q.ResourceType == "rdp" && q.Duration > 0 {
			// Skipping Start query (duration = 0), as it won't have the metadata we need
			// for the RDP replay
			fmt.Printf("RDP query made by %v at %v\n", user.Email, q.Timestamp)
			replayParts, err := client.Replays().List(ctx, "id:?", q.ID)
			if err != nil {
				log.Fatalf("failed to scan replay: %v", err)
			}

			tempDir, err := os.MkdirTemp("", q.ID)
			if err != nil {
				log.Fatalf("failed to create tempory directory: %v", err)
			}
			defer os.RemoveAll(tempDir)

			// Massage the query into the expected format (https://www.strongdm.com/docs/admin/logs/references/post-start/)
			data := fmt.Sprintf(`{"type":"postStart","uuid":"%v","query":%q}`, q.ID, q.QueryBody)
			if err != nil {
				log.Fatalf("failed to marshal query data: %v", err)
			}

			if err := os.WriteFile(filepath.Join(tempDir, "relay.0000000000.log"), []byte(data), 0644); err != nil {
				log.Fatalf("failed to write query data: %v", err)
			}
			chunkId := 1
			for replayParts.Next() {
				chunkEvents, err := json.Marshal(replayParts.Value().Events)
				if err != nil {
					log.Fatalf("failed to marshal chunk data: %v", err)
				}

				// Massage the chunk into the expected format (https://www.strongdm.com/docs/admin/logs/references/replays/)
				chunkData := fmt.Sprintf(`{"type":"chunk","uuid":"%v","chunkId":"%v","events":%s}`, q.ID, chunkId, chunkEvents)
				if err := os.WriteFile(filepath.Join(tempDir, fmt.Sprintf("relay.%010d.log", chunkId)), []byte(chunkData), 0644); err != nil {
					log.Fatalf("failed to write chunk data: %v", err)
				}
				chunkId++
			}
			if err := replayParts.Err(); err != nil {
				log.Fatalf("failed to iterate replay: %v", err)
			}
			logs, err := filepath.Glob(filepath.Join(tempDir, "*"))
			if err != nil {
				log.Fatalf("failed to glob log directory: %v", err)
			}
			// Run the sdm cli, this must be in the path
			cmd := exec.Command("sdm", append([]string{"replay", "rdp", q.ID}, logs...)...)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr

			if err := cmd.Run(); err != nil {
				log.Fatalf("failed to execute sdm replay: %v", err)
			}

			fmt.Println("")
		}
	}
}
