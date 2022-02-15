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
	client, err := sdm.New(
		accessKey,
		secretKey,
	)
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

	// Update Access Rules
	role.AccessRules = sdm.AccessRules{
		sdm.AccessRule{
			Tags: sdm.Tags{
				"env": "staging",
			},
		},
		sdm.AccessRule{
			Type: "postgres",
		},
	}
	_, err = client.Roles().Update(ctx, role)
	if err != nil {
		log.Fatalf("failed to update role: %v", err)
	}

	// Create a User
	user := &sdm.User{
		Email:     "example@example.com",
		FirstName: "example",
		LastName:  "example",
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	accountResponse, err := client.Accounts().Create(ctx, user)
	if err != nil {
		log.Fatalf("Could not create user: %v", err)
	}

	accountID := accountResponse.Account.GetID()
	fmt.Println("Successfully created user.")
	fmt.Println("\tID:", accountID)

	// Assign the User or Service Account to Role
	attachment := &sdm.AccountAttachment{
		AccountID: accountID,
		RoleID:    roleID,
	}

	attachmentResponse, err := client.AccountAttachments().Create(ctx, attachment)
	if err != nil {
		log.Fatalf("Could not create account attachment: %v", err)
	}

	attachmentID := attachmentResponse.AccountAttachment.ID
	fmt.Println("Successfully created account attachment.")
	fmt.Println("\tID:", attachmentID)
}
