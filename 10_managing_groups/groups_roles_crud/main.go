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

	fmt.Println("=== GroupsRoles CRUD Example ===")

	// Setup: Create prerequisite resources
	fmt.Println("\nSetup: Creating prerequisite groups and roles...")

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

	// Create test roles
	role1 := &sdm.Role{
		Name: fmt.Sprintf("Test Role 1 - %d", time.Now().Unix()),
	}

	role2 := &sdm.Role{
		Name: fmt.Sprintf("Test Role 2 - %d", time.Now().Unix()),
	}

	createRole1Resp, err := client.Roles().Create(ctx, role1)
	if err != nil {
		log.Fatalf("failed to create role 1: %v", err)
	}
	createdRole1 := createRole1Resp.Role
	fmt.Printf("Created Role 1: %s (ID: %s)\n", createdRole1.Name, createdRole1.ID)

	createRole2Resp, err := client.Roles().Create(ctx, role2)
	if err != nil {
		log.Fatalf("failed to create role 2: %v", err)
	}
	createdRole2 := createRole2Resp.Role
	fmt.Printf("Created Role 2: %s (ID: %s)\n", createdRole2.Name, createdRole2.ID)

	// === CREATE ===
	fmt.Println("\n1. CREATE - Linking groups to roles...")

	// Link Group1 to Role1
	groupRole1 := &sdm.GroupRole{
		GroupID: createdGroup1.ID,
		RoleID:  createdRole1.ID,
	}

	// Link Group1 to Role2
	groupRole2 := &sdm.GroupRole{
		GroupID: createdGroup1.ID,
		RoleID:  createdRole2.ID,
	}

	// Link Group2 to Role1
	groupRole3 := &sdm.GroupRole{
		GroupID: createdGroup2.ID,
		RoleID:  createdRole1.ID,
	}

	createGroupRole1Resp, err := client.GroupsRoles().Create(ctx, groupRole1)
	if err != nil {
		log.Fatalf("failed to create group role link 1: %v", err)
	}
	createdGroupRole1 := createGroupRole1Resp.GroupRole
	fmt.Printf("Created GroupRole 1: Group %s linked to Role %s (ID: %s)\n", createdGroup1.Name, createdRole1.Name, createdGroupRole1.ID)

	createGroupRole2Resp, err := client.GroupsRoles().Create(ctx, groupRole2)
	if err != nil {
		log.Fatalf("failed to create group role link 2: %v", err)
	}
	createdGroupRole2 := createGroupRole2Resp.GroupRole
	fmt.Printf("Created GroupRole 2: Group %s linked to Role %s (ID: %s)\n", createdGroup1.Name, createdRole2.Name, createdGroupRole2.ID)

	createGroupRole3Resp, err := client.GroupsRoles().Create(ctx, groupRole3)
	if err != nil {
		log.Fatalf("failed to create group role link 3: %v", err)
	}
	createdGroupRole3 := createGroupRole3Resp.GroupRole
	fmt.Printf("Created GroupRole 3: Group %s linked to Role %s (ID: %s)\n", createdGroup2.Name, createdRole1.Name, createdGroupRole3.ID)

	// === READ (List) ===
	fmt.Println("\n2. READ - Listing group role relationships...")

	// List all group roles
	fmt.Println("\nAll GroupRoles:")
	allGroupRoleIter, err := client.GroupsRoles().List(ctx, "")
	if err != nil {
		log.Fatalf("failed to list all group roles: %v", err)
	}

	allGroupRoleCount := 0
	for allGroupRoleIter.Next() {
		groupRole := allGroupRoleIter.Value()
		allGroupRoleCount++
		fmt.Printf("  %d. ID: %s, Group ID: %s, Role ID: %s\n",
			allGroupRoleCount, groupRole.ID, groupRole.GroupID, groupRole.RoleID)
	}

	if allGroupRoleIter.Err() != nil {
		log.Fatalf("error during all group role iteration: %v", allGroupRoleIter.Err())
	}
	fmt.Printf("Total group roles found: %d\n", allGroupRoleCount)

	// List group roles for a specific group
	fmt.Printf("\nGroupRoles for Group 1 (%s):\n", createdGroup1.Name)
	group1RoleIter, err := client.GroupsRoles().List(ctx, fmt.Sprintf("groupid:%s", createdGroup1.ID))
	if err != nil {
		log.Fatalf("failed to list group 1 roles: %v", err)
	}

	group1RoleCount := 0
	for group1RoleIter.Next() {
		groupRole := group1RoleIter.Value()
		group1RoleCount++
		fmt.Printf("  %d. ID: %s, Group ID: %s, Role ID: %s\n", group1RoleCount, groupRole.ID, groupRole.GroupID, groupRole.RoleID)
	}

	if group1RoleIter.Err() != nil {
		log.Fatalf("error during group 1 role iteration: %v", group1RoleIter.Err())
	}
	fmt.Printf("Group roles for Group 1: %d\n", group1RoleCount)

	// List group roles for a specific role
	fmt.Printf("\nGroupRoles for Role 1 (%s):\n", createdRole1.Name)
	role1GroupIter, err := client.GroupsRoles().List(ctx, fmt.Sprintf("roleid:%s", createdRole1.ID))
	if err != nil {
		log.Fatalf("failed to list role 1 groups: %v", err)
	}

	role1GroupCount := 0
	for role1GroupIter.Next() {
		groupRole := role1GroupIter.Value()
		role1GroupCount++
		fmt.Printf("  %d. ID: %s, Group ID: %s, Role ID: %s\n", role1GroupCount, groupRole.ID, groupRole.GroupID, groupRole.RoleID)
	}

	if role1GroupIter.Err() != nil {
		log.Fatalf("error during role 1 group iteration: %v", role1GroupIter.Err())
	}
	fmt.Printf("Group roles for Role 1: %d\n", role1GroupCount)

	// === DELETE ===
	fmt.Println("\n3. DELETE - Removing group role relationships...")

	_, err = client.GroupsRoles().Delete(ctx, createdGroupRole1.ID)
	if err != nil {
		log.Fatalf("failed to delete group role 1: %v", err)
	}
	fmt.Printf("Deleted GroupRole 1 (ID: %s)\n", createdGroupRole1.ID)

	_, err = client.GroupsRoles().Delete(ctx, createdGroupRole2.ID)
	if err != nil {
		log.Fatalf("failed to delete group role 2: %v", err)
	}
	fmt.Printf("Deleted GroupRole 2 (ID: %s)\n", createdGroupRole2.ID)

	_, err = client.GroupsRoles().Delete(ctx, createdGroupRole3.ID)
	if err != nil {
		log.Fatalf("failed to delete group role 3: %v", err)
	}
	fmt.Printf("Deleted GroupRole 3 (ID: %s)\n", createdGroupRole3.ID)

	// Verify deletion
	fmt.Println("\nVerifying deletion by listing group roles for our test group...")
	verifyIter, err := client.GroupsRoles().List(ctx, fmt.Sprintf("groupid:%s", createdGroup1.ID))
	if err != nil {
		log.Fatalf("failed to verify group role deletion: %v", err)
	}

	remainingCount := 0
	for verifyIter.Next() {
		remainingCount++
	}

	if verifyIter.Err() != nil {
		log.Fatalf("error during deletion verification: %v", verifyIter.Err())
	}

	fmt.Printf("Group roles remaining after deletion for Group 1: %d\n", remainingCount)

	// === CLEANUP ===
	fmt.Println("\nCleanup: Removing prerequisite resources...")

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

	// Delete the roles
	_, err = client.Roles().Delete(ctx, createdRole1.ID)
	if err != nil {
		log.Printf("Warning: failed to delete role %s: %v", createdRole1.ID, err)
	} else {
		fmt.Printf("Deleted Role 1: %s\n", createdRole1.Name)
	}

	_, err = client.Roles().Delete(ctx, createdRole2.ID)
	if err != nil {
		log.Printf("Warning: failed to delete role %s: %v", createdRole2.ID, err)
	} else {
		fmt.Printf("Deleted Role 2: %s\n", createdRole2.Name)
	}

	fmt.Println("\n=== GroupsRoles CRUD Example Completed Successfully! ===")
}
