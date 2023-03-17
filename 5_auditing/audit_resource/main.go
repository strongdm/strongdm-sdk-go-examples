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

	start := time.Now()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	/* * * */
	// Set up some audit records to read
	fmt.Println("Creating, renaming, and deleting a test resource...")

	redisResp, err := client.Resources().Create(ctx, &sdm.Redis{
		Name:     "example-redis",
		Hostname: "example-redis",
		Username: "example-username",
	})
	if err != nil {
		log.Fatalf("failed to create redis: %v", err)
	}

	createdAt := time.Now()

	resourceID := redisResp.Resource.GetID()

	redisResp.Resource.SetName("example-redis-renamed")
	_, err = client.Resources().Update(ctx, redisResp.Resource)
	if err != nil {
		log.Fatalf("failed to rename redis: %v", err)
	}

	renamedAt := time.Now()

	_, err = client.Resources().Delete(ctx, resourceID)
	if err != nil {
		log.Fatalf("failed to delete redis: %v", err)
	}

	deletedAt := time.Now()
	/* * * */

	// Audit records may take a few seconds to be processed.
	time.Sleep(4 * time.Second)

	_, err = client.SnapshotAt(start).Resources().Get(ctx, resourceID)
	fmt.Printf("Attempting to retrieve resource before creation (%v): %v\n", start, err) // Does Not Exist

	getResp, err := client.SnapshotAt(createdAt).Resources().Get(ctx, resourceID)
	if err != nil {
		log.Fatalf("failed to retrieve created redis: %v", err)
	}

	fmt.Printf("Resource name after creation (%v): %v", createdAt, getResp.Resource.GetName()) // example-redis

	getResp, err = client.SnapshotAt(renamedAt).Resources().Get(ctx, resourceID)
	if err != nil {
		log.Fatalf("failed to retrieve renamed redis: %v", err)
	}

	fmt.Printf("Resource name after rename (%v): %v", renamedAt, getResp.Resource.GetName()) // example-redis-renamed

	_, err = client.SnapshotAt(deletedAt).Resources().Get(ctx, resourceID)
	fmt.Printf("Attempting to retrieve resource after deletion (%v): %v\n", deletedAt, err) // Does Not Exist

	fmt.Println("Full history of the resource:")
	history, err := client.ResourcesHistory().List(ctx, "id:?", resourceID)
	if err != nil {
		log.Fatalf("failed to list resource history: %v", err)
	}
	for history.Next() {
		v := history.Value()
		activity, err := client.Activities().Get(ctx, v.ActivityID)
		if err != nil {
			log.Fatalf("failed to lookup history: %v", err)
		}
		fmt.Println(activity.Activity.Description) // created, updated, deleted resource; in order
	}
	if err := history.Err(); err != nil {
		log.Fatalf("failed to finish listing resource history: %v", err)
	}
}
