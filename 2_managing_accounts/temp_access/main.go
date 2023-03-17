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

	// Create a resource
	resource := &sdm.Postgres{
		Name:         "Example Postgres Datasource for Go Temp Access",
		Hostname:     "example.strongdm.com",
		Port:         5432,
		Username:     "example",
		Password:     "example",
		Database:     "example",
		PortOverride: 18001,
		Tags: sdm.Tags{
			"example": "example",
		},
	}
	resourceResponse, err := client.Resources().Create(ctx, resource)
	if err != nil {
		log.Fatalf("Could not create Postgres datasource: %v", err)
	}
	fmt.Println("Successfully created Postgres datasource.")
	fmt.Println("\tID:", resourceResponse.Resource.GetID())
	fmt.Println("\tName:", resourceResponse.Resource.GetName())

	// Create a User
	user := &sdm.User{
		Email:     "go-temp-access@example.com",
		FirstName: "example",
		LastName:  "example",
	}
	userResponse, err := client.Accounts().Create(ctx, user)
	if err != nil {
		log.Fatalf("Could not create user: %v", err)
	}
	fmt.Println("Successfully created user.")
	fmt.Println("\tID:", userResponse.Account.GetID())

	// Grant temporary access
	grant := &sdm.AccountGrant{
		AccountID:  userResponse.Account.GetID(),
		ResourceID: resourceResponse.Resource.GetID(),
		ValidUntil: time.Now().Add(30 * time.Minute),
	}
	grantResponse, err := client.AccountGrants().Create(ctx, grant)
	if err != nil {
		log.Fatalf("Could not create temporary account grant: %v", err)
	}
	fmt.Println("Successfully created temporary account grant.")
	fmt.Println("\tID:", grantResponse.AccountGrant.ID)

	// Revoke temporary access
	_, err = client.AccountGrants().Delete(ctx, grantResponse.AccountGrant.ID)
	if err != nil {
		log.Fatalf("Could not delete temporary account grant: %v", err)
	}
	fmt.Println("Successfully deleted temporary account grant.")
}
