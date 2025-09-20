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
	"fmt"
	"log"
	"os"
	"time"

	sdm "github.com/strongdm/strongdm-sdk-go/v15"
)

func main() {
	log.SetFlags(0)
	// Load the SDM API keys from the environment.
	// If these values are not set in your environment,
	// please follow the documentation here:
	// https://www.strongdm.com/docs/api/api-keys/
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
		log.Fatal("failed to create strongDM client:", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// === CREATE ===
	fmt.Println("=== Groups CRUD Example ===")
	fmt.Println("\n1. CREATE - Creating two test groups...")

	group1 := &sdm.Group{
		Name: fmt.Sprintf("Test Group 1 - %d", time.Now().Unix()),
	}

	group2 := &sdm.Group{
		Name: fmt.Sprintf("Test Group 2 - %d", time.Now().Unix()),
	}

	createGroup1Resp, err := client.Groups().Create(ctx, group1)
	if err != nil {
		log.Fatalf("failed to create group 1: %v", err)
	}
	createdGroup1 := createGroup1Resp.Group
	fmt.Printf("Created Group 1: %s (ID: %s)\n", createdGroup1.Name, createdGroup1.ID)

	createGroup2Resp, err := client.Groups().Create(ctx, group2)
	if err != nil {
		log.Fatalf("failed to create group 2: %v", err)
	}
	createdGroup2 := createGroup2Resp.Group
	fmt.Printf("Created Group 2: %s (ID: %s)\n", createdGroup2.Name, createdGroup2.ID)

	// === READ (List) ===
	fmt.Println("\n2. READ - Listing all groups...")
	groupIter, err := client.Groups().List(ctx, "")
	if err != nil {
		log.Fatalf("failed to list groups: %v", err)
	}

	fmt.Println("All Groups:")
	groupCount := 0
	for groupIter.Next() {
		group := groupIter.Value()
		groupCount++
		fmt.Printf("  %d. Name: %s, ID: %s\n", groupCount, group.Name, group.ID)
	}

	if groupIter.Err() != nil {
		log.Fatalf("error during group iteration: %v", groupIter.Err())
	}
	fmt.Printf("Total groups found: %d\n", groupCount)

	// === UPDATE ===
	fmt.Println("\n3. UPDATE - Updating group names...")

	createdGroup1.Name = fmt.Sprintf("Updated Test Group 1 - %d", time.Now().Unix())
	updateGroup1Resp, err := client.Groups().Update(ctx, createdGroup1)
	if err != nil {
		log.Fatalf("failed to update group 1: %v", err)
	}
	updatedGroup1 := updateGroup1Resp.Group
	fmt.Printf("Updated Group 1: %s (ID: %s)\n", updatedGroup1.Name, updatedGroup1.ID)

	createdGroup2.Name = fmt.Sprintf("Updated Test Group 2 - %d", time.Now().Unix())
	updateGroup2Resp, err := client.Groups().Update(ctx, createdGroup2)
	if err != nil {
		log.Fatalf("failed to update group 2: %v", err)
	}
	updatedGroup2 := updateGroup2Resp.Group
	fmt.Printf("Updated Group 2: %s (ID: %s)\n", updatedGroup2.Name, updatedGroup2.ID)

	// Verify updates by listing our specific groups
	fmt.Println("\nVerifying updates by filtering for our test groups...")
	testGroupIter, err := client.Groups().List(ctx, "name:\"Updated Test Group*\"")
	if err != nil {
		log.Fatalf("failed to list updated groups: %v", err)
	}

	fmt.Println("Updated Test Groups:")
	testGroupCount := 0
	for testGroupIter.Next() {
		group := testGroupIter.Value()
		testGroupCount++
		fmt.Printf("  %d. Name: %s, ID: %s\n", testGroupCount, group.Name, group.ID)
	}

	if testGroupIter.Err() != nil {
		log.Fatalf("error during updated group iteration: %v", testGroupIter.Err())
	}

	// === DELETE ===
	fmt.Println("\n4. DELETE - Cleaning up created groups...")

	_, err = client.Groups().Delete(ctx, updatedGroup1.ID)
	if err != nil {
		log.Fatalf("failed to delete group 1: %v", err)
	}
	fmt.Printf("Deleted Group 1: %s\n", updatedGroup1.Name)

	_, err = client.Groups().Delete(ctx, updatedGroup2.ID)
	if err != nil {
		log.Fatalf("failed to delete group 2: %v", err)
	}
	fmt.Printf("Deleted Group 2: %s\n", updatedGroup2.Name)

	// Verify deletion
	fmt.Println("\nVerifying deletion by attempting to list our deleted groups...")
	verifyIter, err := client.Groups().List(ctx, "name:\"Updated Test Group*\"")
	if err != nil {
		log.Fatalf("failed to verify group deletion: %v", err)
	}

	deletedGroupCount := 0
	for verifyIter.Next() {
		deletedGroupCount++
	}

	if verifyIter.Err() != nil {
		log.Fatalf("error during deletion verification: %v", verifyIter.Err())
	}

	fmt.Printf("Groups remaining after deletion: %d\n", deletedGroupCount)

	fmt.Println("\n=== Groups CRUD Example Completed Successfully! ===")
}
