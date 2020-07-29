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

	sdm "github.com/strongdm/strongdm-sdk-go"
)

func main() {
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
		log.Fatalf("could not create client: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Create a datasource
	examplePostgresDatasource := &sdm.Postgres{
		Name:             "Example Postgres Datasource",
		Hostname:         "example.strongdm.com",
		Port:             5432,
		Username:         "example",
		Password:         "example",
		Database:         "example",
		OverrideDatabase: true,
		PortOverride:     19999,
	}

	resourceResponse, err := client.Resources().Create(ctx, examplePostgresDatasource)
	if err != nil {
		log.Fatalf("Could not create Postgres datasource: %v", err)
	}

	resourceID := resourceResponse.Resource.GetID()
	fmt.Println("Successfully created Postgres datasource.")
	fmt.Println("\tID:", resourceID)

	// Create a role
	role := &sdm.Role{
		Name: "example role",
	}

	roleResponse, err := client.Roles().Create(ctx, role)
	if err != nil {
		log.Fatalf("Could not create role: %v", err)
	}

	roleID := roleResponse.Role.ID
	fmt.Println("Successfully created role.")
	fmt.Println("\tID:", roleID)

	// Create a role grant
	grant := &sdm.RoleGrant{
		ResourceID: resourceID,
		RoleID:     roleID,
	}

	grantResponse, err := client.RoleGrants().Create(ctx, grant)
	if err != nil {
		log.Fatalf("Could not create account grant: %v", err)
	}

	grantID := grantResponse.RoleGrant.ID
	fmt.Println("Successfully created role grant.")
	fmt.Println("\tID:", grantID)

	// Create a user
	user := &sdm.User{
		Email:     "example@example.com",
		FirstName: "example",
		LastName:  "example",
	}

	accountResponse, err := client.Accounts().Create(ctx, user)
	if err != nil {
		log.Fatalf("Could not create user: %v", err)
	}

	accountID := accountResponse.Account.GetID()
	fmt.Println("Successfully created user.")
	fmt.Println("\tID:", accountID)

	// Assign account to role
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

	// Detatch user from role
	_, err = client.AccountAttachments().Delete(ctx, attachmentID)
	if err != nil {
		log.Fatalf("Could not delete account attachment: %v", err)
	}
	fmt.Println("Successfully deleted account attachment.")

	// Delete grant from role
	_, err = client.RoleGrants().Delete(ctx, grantID)
	if err != nil {
		log.Fatalf("Could not delete role grant: %v", err)
	}
	fmt.Println("Successfully deleted role grant.")
}
