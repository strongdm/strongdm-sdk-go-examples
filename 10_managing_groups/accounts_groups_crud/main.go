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

	fmt.Println("=== AccountsGroups CRUD Example ===")

	// Setup: Create prerequisite resources
	fmt.Println("\nSetup: Creating prerequisite account and groups...")

	// Create test account
	account := &sdm.User{
		FirstName: "Test",
		LastName:  "User",
		Email:     fmt.Sprintf("test-user-%d@example.com", time.Now().Unix()),
	}

	createAccountResp, err := client.Accounts().Create(ctx, account)
	if err != nil {
		log.Fatalf("failed to create account: %v", err)
	}
	createdAccount := createAccountResp.Account
	accountID := createdAccount.GetID()
	fmt.Printf("Created Account: %s (ID: %s)\n", account.Email, accountID)

	// Create test groups
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

	// === CREATE ===
	fmt.Println("\n1. CREATE - Linking account to groups...")

	accountGroup1 := &sdm.AccountGroup{
		AccountID: accountID,
		GroupID:   createdGroup1.ID,
	}

	accountGroup2 := &sdm.AccountGroup{
		AccountID: accountID,
		GroupID:   createdGroup2.ID,
	}

	createAccountGroup1Resp, err := client.AccountsGroups().Create(ctx, accountGroup1)
	if err != nil {
		log.Fatalf("failed to create account group link 1: %v", err)
	}
	createdAccountGroup1 := createAccountGroup1Resp.AccountGroup
	fmt.Printf("Created AccountGroup 1: Account %s linked to Group %s (ID: %s)\n", accountID, createdGroup1.Name, createdAccountGroup1.ID)

	createAccountGroup2Resp, err := client.AccountsGroups().Create(ctx, accountGroup2)
	if err != nil {
		log.Fatalf("failed to create account group link 2: %v", err)
	}
	createdAccountGroup2 := createAccountGroup2Resp.AccountGroup
	fmt.Printf("Created AccountGroup 2: Account %s linked to Group %s (ID: %s)\n", accountID, createdGroup2.Name, createdAccountGroup2.ID)

	// === READ (List) ===
	fmt.Println("\n2. READ - Listing account group relationships...")

	// List all account groups
	fmt.Println("\nAll AccountGroups:")
	allAccountGroupIter, err := client.AccountsGroups().List(ctx, "")
	if err != nil {
		log.Fatalf("failed to list all account groups: %v", err)
	}

	allAccountGroupCount := 0
	for allAccountGroupIter.Next() {
		accountGroup := allAccountGroupIter.Value()
		allAccountGroupCount++
		fmt.Printf("  %d. ID: %s, Account ID: %s, Group ID: %s\n", allAccountGroupCount, accountGroup.ID, accountGroup.AccountID, accountGroup.GroupID)
	}

	if allAccountGroupIter.Err() != nil {
		log.Fatalf("error during all account group iteration: %v", allAccountGroupIter.Err())
	}
	fmt.Printf("Total account groups found: %d\n", allAccountGroupCount)

	// List account groups for our specific account
	fmt.Println("\nAccountGroups for our test account:")
	accountGroupIter, err := client.AccountsGroups().List(ctx, fmt.Sprintf("accountid:%s", accountID))
	if err != nil {
		log.Fatalf("failed to list account groups: %v", err)
	}

	accountGroupCount := 0
	for accountGroupIter.Next() {
		accountGroup := accountGroupIter.Value()
		accountGroupCount++
		fmt.Printf("  %d. ID: %s, Account ID: %s, Group ID: %s\n", accountGroupCount, accountGroup.ID, accountGroup.AccountID, accountGroup.GroupID)
	}

	if accountGroupIter.Err() != nil {
		log.Fatalf("error during account group iteration: %v", accountGroupIter.Err())
	}
	fmt.Printf("Account groups for test account: %d\n", accountGroupCount)

	// === DELETE ===
	fmt.Println("\n3. DELETE - Removing account group relationships...")

	_, err = client.AccountsGroups().Delete(ctx, createdAccountGroup1.ID)
	if err != nil {
		log.Fatalf("failed to delete account group 1: %v", err)
	}
	fmt.Printf("Deleted AccountGroup 1 (ID: %s)\n", createdAccountGroup1.ID)

	_, err = client.AccountsGroups().Delete(ctx, createdAccountGroup2.ID)
	if err != nil {
		log.Fatalf("failed to delete account group 2: %v", err)
	}
	fmt.Printf("Deleted AccountGroup 2 (ID: %s)\n", createdAccountGroup2.ID)

	// Verify deletion
	fmt.Println("\nVerifying deletion by listing account groups for our test account...")
	verifyIter, err := client.AccountsGroups().List(ctx, fmt.Sprintf("accountid:%s", accountID))
	if err != nil {
		log.Fatalf("failed to verify account group deletion: %v", err)
	}

	remainingCount := 0
	for verifyIter.Next() {
		remainingCount++
	}

	if verifyIter.Err() != nil {
		log.Fatalf("error during deletion verification: %v", verifyIter.Err())
	}

	fmt.Printf("Account groups remaining after deletion: %d\n", remainingCount)

	// === CLEANUP ===
	fmt.Println("\nCleanup: Removing prerequisite resources...")

	// Delete the account
	_, err = client.Accounts().Delete(ctx, accountID)
	if err != nil {
		log.Printf("Warning: failed to delete account %s: %v", accountID, err)
	} else {
		fmt.Printf("Deleted Account: %s\n", account.Email)
	}

	// Delete the groups
	_, err = client.Groups().Delete(ctx, createdGroup1.ID)
	if err != nil {
		log.Printf("Warning: failed to delete group %s: %v", createdGroup1.ID, err)
	} else {
		fmt.Printf("Deleted Group 1: %s\n", createdGroup1.Name)
	}

	_, err = client.Groups().Delete(ctx, createdGroup2.ID)
	if err != nil {
		log.Printf("Warning: failed to delete group %s: %v", createdGroup2.ID, err)
	} else {
		fmt.Printf("Deleted Group 2: %s\n", createdGroup2.Name)
	}

	fmt.Println("\n=== AccountsGroups CRUD Example Completed Successfully! ===")
}
