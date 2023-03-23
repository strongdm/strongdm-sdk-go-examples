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
	"context"
	"encoding/json"
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
			fmt.Printf("Skipping encrypted query made by %v at %v\n", user.Email, q.Timestamp)
			fmt.Println("See encrypted_query_replay for an example of query decryption.")
		} else if q.Replayable {
			fmt.Printf("Replaying query made by %v at %v\n", user.Email, q.Timestamp)
			replayParts, err := client.Replays().List(ctx, "id:?", q.ID)
			if err != nil {
				log.Fatalf("failed to scan replay: %v", err)
			}
			for replayParts.Next() {
				next := replayParts.Value()
				for _, ev := range next.Events {
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
				fmt.Printf("failed to unmarshal query JSON %v: %v", q.QueryBody, err)
			} else {
				fmt.Printf("Command run by %v at %v: %v\n", user.Email, q.Timestamp, capture.Command)
			}
		}
	}
}
