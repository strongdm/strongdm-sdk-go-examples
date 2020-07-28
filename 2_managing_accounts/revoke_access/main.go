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
	log.Printf("Successfully created user.\n\tID: %v\n", accountID)

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

	log.Printf("Successfully created Postgres datasource.\n\tID: %v\n", resourceID)

	// Create an account grant
	accountGrant := &sdm.AccountGrant{
		ResourceID: resourceID,
		AccountID:  accountID,
	}

	grantResponse, err := client.AccountGrants().Create(ctx, accountGrant)
	if err != nil {
		log.Fatalf("Could not create account grant: %v", err)
	}

	grantID := grantResponse.AccountGrant.ID

	log.Printf("Successfully created account grant.\n\tID: %v\n", grantID)

	// Delete an account grant
	_, err = client.AccountGrants().Delete(ctx, grantID)
	if err != nil {
		log.Fatalf("Could not delete account grant: %v", err)
	}
	log.Printf("Successfully deleted account grant.\n")
}
