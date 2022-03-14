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
	"math/rand"
	"os"
	"time"

	sdm "github.com/strongdm/strongdm-sdk-go/v2"
)

func main() {
	// The Role Grants API has been deprecated in favor of Access Rules.
	// When using Access Rules, the best practice is to give Roles access to Resources based on type and tags.
	// If it is _necessary_ to grant access to specific Resources in the same way as Role Grants did,
	// you can use Resource IDs directly in Access Rules as shown in the following examples.

	log.SetFlags(0)
	//	Load the SDM API keys from the environment.
	//	If these values are not set in your environment,
	//	please follow the documentation here:
	//	https://www.strongdm.com/docs/admin-guide/api-credentials/
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

	err = createRoleGrantViaAccessRulesExample(ctx, client)
	if err != nil {
		log.Fatalf("error in createRoleGrantViaAccessRulesExample: %v", err)
	}
	err = deleteRoleGrantViaAccessRulesExample(ctx, client)
	if err != nil {
		log.Fatalf("error in deleteRoleGrantViaAccessRulesExample: %v", err)
	}
	err = listRoleGrantsViaAccessRulesExample(ctx, client)
	if err != nil {
		log.Fatalf("error in listRoleGrantsViaAccessRulesExample: %v", err)
	}
}

// Example: Create a Role with empty Access Rules and return the ID
func createExampleRole(ctx context.Context, client *sdm.Client, ar sdm.AccessRules) string {
	role := &sdm.Role{
		Name:        "exampleRole-" + fmt.Sprint(rand.Int()),
		AccessRules: ar,
	}
	roleResp, err := client.Roles().Create(ctx, role)
	if err != nil {
		log.Fatalf("error creating role: %v", err)
	}
	return roleResp.Role.ID
}

// Example: Create a sample Resource and return the ID
func createExampleResource(ctx context.Context, client *sdm.Client) string {
	redis := &sdm.Redis{
		Name:         "exampleResource-" + fmt.Sprint(rand.Int()),
		Hostname:     "example.com",
		Port:         6379,
		PortOverride: int32(rand.Intn(20000) + 3000),
	}
	resp, err := client.Resources().Create(ctx, redis)
	if err != nil {
		log.Fatalf("error creating resource: %v", err)
	}
	return resp.Resource.GetID()
}

// Example: Create a Role grant via Access Rules
func createRoleGrantViaAccessRulesExample(ctx context.Context, client *sdm.Client) error {
	// Create example Resources
	resourceID1 := createExampleResource(ctx, client)
	resourceID2 := createExampleResource(ctx, client)
	roleID := createExampleRole(ctx, client, sdm.AccessRules{
		sdm.AccessRule{
			IDs: []string{resourceID1},
		},
	})

	// Get the Role
	getResp, err := client.Roles().Get(ctx, roleID)
	if err != nil {
		return fmt.Errorf("error getting role: %v", err)
	}
	role := getResp.Role

	if len(role.AccessRules) == 0 {
		role.AccessRules = sdm.AccessRules{sdm.AccessRule{}}
	}

	// Append the ID to an existing static Access Rule
	if len(role.AccessRules) != 1 || len(role.AccessRules[0].IDs) == 0 {
		return fmt.Errorf("unexpected access rules in role")
	}
	role.AccessRules[0].IDs = append(role.AccessRules[0].IDs, resourceID2)

	// Update the Role
	_, err = client.Roles().Update(ctx, role)
	if err != nil {
		return fmt.Errorf("error updating role: %v", err)
	}
	return nil
}

// Example: Delete a Role grant via Access Rules
func deleteRoleGrantViaAccessRulesExample(ctx context.Context, client *sdm.Client) error {
	// Create example Resources
	resourceID1 := createExampleResource(ctx, client)
	resourceID2 := createExampleResource(ctx, client)
	roleID := createExampleRole(ctx, client, sdm.AccessRules{
		sdm.AccessRule{
			IDs: []string{resourceID1},
		},
	})

	// Get the Role
	getResp, err := client.Roles().Get(ctx, roleID)
	if err != nil {
		return fmt.Errorf("error getting role: %v", err)
	}
	role := getResp.Role

	// Append the ID to an existing static Access Rule
	if len(role.AccessRules) != 1 || len(role.AccessRules[0].IDs) == 0 {
		return fmt.Errorf("unexpected access rules in role")
	}
	role.AccessRules[0].IDs = append(role.AccessRules[0].IDs, resourceID2)
	if len(role.AccessRules[0].IDs) == 0 {
		role.AccessRules = nil
	}

	// Update the Role
	_, err = client.Roles().Update(ctx, role)
	if err != nil {
		return fmt.Errorf("error updating role: %v", err)
	}
	return nil
}

// Example: List Role grants via Access Rules
func listRoleGrantsViaAccessRulesExample(ctx context.Context, client *sdm.Client) error {
	// Create example Resources
	resourceID := createExampleResource(ctx, client)
	roleID := createExampleRole(ctx, client, sdm.AccessRules{
		sdm.AccessRule{
			IDs: []string{resourceID},
		},
	})

	// Get the Role
	getResp, err := client.Roles().Get(ctx, roleID)
	if err != nil {
		return fmt.Errorf("error getting role: %v", err)
	}
	role := getResp.Role

	// role.AccessRules contains each AccessRule associated with the Role
	for _, resourceID := range role.AccessRules[0].IDs {
		fmt.Println(resourceID)
	}

	return nil
}
