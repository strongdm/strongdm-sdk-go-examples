// Copyright 2020 StrongDM Inc
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	sdm "github.com/strongdm/strongdm-sdk-go/v2"
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

	// Define the SSH server
	server := &sdm.SSH{
		Name:     "Example SSH Server",
		Hostname: "203.0.113.23",
		Username: "example",
		Port:     22,
		Tags: sdm.Tags{
			"example": "example",
		},
	}

	// Create the server
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	createResponse, err := client.Resources().Create(ctx, server)
	if err != nil {
		log.Fatalf("Could not create SSH server: %v", err)
	}

	fmt.Println("Successfully created SSH server.")
	fmt.Println("\tID:", createResponse.Resource.GetID())
	fmt.Println("\tName:", createResponse.Resource.GetName())

//Example: View resource history, replay queries
func exampleSSHReplay(client *sdm.Client) {
	ctx := context.Background()

	// You need an SSH resource that has had queries made against it
	resourceName := "Example"
	resourceResp, err := client.Resources().List(ctx, "name:?", resourceName)
	if err != nil {
		log.Fatalf("failed to list resources: %v", err)
	}
	if !resourceResp.Next() {
		log.Fatalf("couldn't find resource named %v (error: %v)", resourceName, resourceResp.Err())
	}
	resource := resourceResp.Value()

	fmt.Printf("Queries made against %v\n", resourceName)
	queries, err := client.Queries().List(ctx, "resource_id:?", resource.GetID())
	if err != nil {
		log.Fatalf("failed to list queries: %v", err)
	}
	for queries.Next() {
		q := queries.Value()
		if q.Replayable {
			fmt.Printf("Replaying query made at %v\n", q.Timestamp)
			replayParts, err := client.Replays().List(ctx, "id:?", q.ID)
			if err != nil {
				log.Fatalf("failed to scan replay: %v", err)
			}
			for replayParts.Next() {
				next := replayParts.Value()
				for _, ev := range next.Events {
					// This won't handle some characters as expected like deletes
					fmt.Println(string(ev.Data))
					time.Sleep(ev.Duration)
				}
			}
			if err := replayParts.Err(); err != nil {
				log.Fatalf("failed to iterate replay: %v", err)
			}
		}
	}

	//Example output for resource history example:
	
	//Queries made against SSH Example
	//Replaying query made at 2023-03-14 21:51:16.75671 +0000 UTC
	//cannot start SSH Example playback: dial tcp 127.0.0.1:23: connect: connection refused
	//Replaying query made at 2023-03-14 21:51:34.307327 +0000 UTC
	//cannot start SSH Example playback: dial tcp 127.0.0.1:23: connect: connection refused
	//Replaying query made at 2023-03-14 21:51:35.019356 +0000 UTC
	//cannot start SSH Example playback: dial tcp 127.0.0.1:23: connect: connection refused
}
