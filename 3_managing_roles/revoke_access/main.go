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

	sdm "github.com/strongdm/web/pkg/api/v1/generated/go"
)

func main() {
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
	client, err := sdm.New(accessKey, secretKey)
	if err != nil {
		log.Fatal("failed to create strongDM client:", err)
	}

	// Create a resource (e.g., Redis)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	
	redisID := createExampleResource(ctx, client)

	// Create a Role with initial Access Rule
	role := &sdm.Role{
		Name: "accessRulesTestRole",
		AccessRules: sdm.AccessRules{
			sdm.AccessRule{
				IDs: []string{redisID},
			},
		},
	}
	roleResp, err := client.Roles().Create(ctx, role)
	if err != nil {
		log.Fatalf("failed to create role: %v", err)
	}
	role = roleResp.Role

	// The RoleGrants API has been deprecated in favor of Access Rules.
	// When using Access Rules the best practice is to grant Resources access based on Type and Tags.
	// If it is _necessary_ to grant access to specific Resources in the same way as RoleGrants did,
	// you can use Resource IDs directly in Access Rules as shown in the following example.

	err = deleteRoleGrantViaAccessRulesExample(ctx, client)
	if err != nil {
		log.Fatalf("error in deleteRoleGrantViaAccessRulesExample: %v", err)
	}
}

// Example: Delete Role grant via Access Rules
func deleteRoleGrantViaAccessRulesExample(ctx context.Context, client *sdm.Client) error {
	// Create example resources
	resourceID1 := createExampleResource(ctx, client)
	resourceID2 := createExampleResource(ctx, client)
	roleID := createExampleRole(ctx, client, sdm.AccessRules{
		sdm.AccessRule{
			IDs: []string{resourceID1, resourceID2},
		},
	})

	// Get the Role
	getResp, err := client.Roles().Get(ctx, roleID)
	if err != nil {
		return fmt.Errorf("error getting role: %v", err)
	}
	role := getResp.Role

	if len(role.AccessRules) != 1 || len(role.AccessRules[0].IDs) == 0 {
		return fmt.Errorf("unexpected access rules in role")
	}

	// Remove the ID of the second resource
	for i, r := range role.AccessRules[0].IDs {
		if r == resourceID2 {
			role.AccessRules[0].IDs = removeElement(role.AccessRules[0].IDs, i)
			break
		}
	}

	// Update the Role
	_, err = client.Roles().Update(ctx, role)
	if err != nil {
		return fmt.Errorf("error updating role: %v", err)
	}
	return nil
}


