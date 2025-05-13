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

	sdm "github.com/strongdm/strongdm-sdk-go/v2"
)

func main() {
	log.SetFlags(0)
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
		log.Fatal("failed to create strongDM client:", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Create a User
	user := &sdm.User{
		Email:     "revoke-access@example.com",
		FirstName: "example",
		LastName:  "example",
	}
	accountResponse, err := client.Accounts().Create(ctx, user)
	if err != nil {
		log.Fatalf("Could not create user: %v", err)
	}
	account := accountResponse.Account
	fmt.Println("Successfully created user.")
	fmt.Println("\tID:", account.GetID())

	// Create a Resource (e.g., Redis)
	redisResp, err := client.Resources().Create(ctx, &sdm.Redis{
		Name:         "Redis for Revoke Access Test",
		Hostname:     "example.com",
		PortOverride: 19204,
		Tags: sdm.Tags{
			"example": "revoke-access",
		},
	})
	if err != nil {
		log.Fatalf("failed to create resource: %v", err)
	}
	redis := redisResp.Resource
	fmt.Println("Successfully created resource.")
	fmt.Println("\tID:", redis.GetID())

	// Create a Role with initial Access Rule
	role := &sdm.Role{
		Name: "Role for Revoke Access Test",
		AccessRules: sdm.AccessRules{
			sdm.AccessRule{
				Tags: sdm.Tags{
					"example": "revoke-access",
				},
			},
		},
	}
	roleResp, err := client.Roles().Create(ctx, role)
	if err != nil {
		log.Fatalf("failed to create role: %v", err)
	}
	role = roleResp.Role
	fmt.Println("Successfully created role.")
	fmt.Println("\tID:", role.ID)

	// Assign role to account
	attachment := &sdm.AccountAttachment{
		AccountID: account.GetID(),
		RoleID:    role.ID,
	}
	attachmentResponse, err := client.AccountAttachments().Create(ctx, attachment)
	if err != nil {
		log.Fatalf("Could not create account attachment: %v", err)
	}
	attachment = attachmentResponse.AccountAttachment
	fmt.Println("Successfully created account attachment.")
	fmt.Println("\tID:", attachment.ID)

	// Option 1: delete access rules
	role.AccessRules = nil
	_, err = client.Roles().Update(ctx, role)
	if err != nil {
		log.Fatalf("failed to update role: %v", err)
	}
	fmt.Println("Successfully deleted access rules.")

	// Option 2: delete account attachment
	if _, err := client.AccountAttachments().Delete(ctx, attachment.ID); err != nil {
		log.Fatalf("failed to delete account attachment: %v", err)
	}
	fmt.Println("Successfully deleted account attachment.")

	// Option 3: remove tag from resource
	redis.SetTags(nil)
	if _, err := client.Resources().Update(ctx, redis); err != nil {
		log.Fatalf("failed to remove tag from resource: %v", err)
	}
	fmt.Println("Successfully removed tag from resource.")
}
