// Copyright 2024 StrongDM Inc
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

	sdm "github.com/strongdm/strongdm-sdk-go/v8"
)

func main() {
	log.SetFlags(0)
	//	Load the SDM API keys from the environment.
	//	If these values are not set in your environment,
	//	please follow the documentation here:
	//	https://www.strongdm.com/docs/api/api-keys/

	// For this token rotation script, the API key used here
	// must have permissions to create and delete tokens,
	// Any token created by this API Key can only have permissions
	// that are a subset of permissions possessed by this API Key.
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

	// Create an API Key
	apiKey := &sdm.Token{
		Name:        "example-token", // name of token must be unique
		AccountType: "api",
		Duration:    time.Hour,
		Permissions: []string{sdm.PermissionRoleCreate, sdm.PermissionUserCreateAdminToken},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	createResponse, err := client.Accounts().Create(ctx, apiKey)
	if err != nil {
		log.Fatalf("Could not create api key: %v", err)
	}

	id := createResponse.Account.GetID()
	tokenName := createResponse.Account.(*sdm.Token).Name
	createdAccessKey := createResponse.AccessKey
	createdSecretKey := createResponse.SecretKey

	fmt.Println("Successfully created api key.")
	fmt.Println("\tID:", id)
	fmt.Println("\tName:", tokenName)
	fmt.Println("\tAccessKey:", createdAccessKey)
	fmt.Println("\tSecretKey:", createdSecretKey)

	// Find the ID of the Token based on its unique name
	listResp, err := client.Accounts().List(ctx, "name:?", tokenName)
	if err != nil {
		log.Fatalf("Could not list token by name: %v", err)
	}
	var oldToken *sdm.Token
	for listResp.Next() {
		oldToken = listResp.Value().(*sdm.Token)
	}
	oldTokenID := oldToken.ID

	// Temporarily update the name of the old token
	// so new one can be created with its name
	deprecatedToken := &sdm.Token{
		ID:   oldTokenID,
		Name: tokenName + "-deprecated",
	}

	_, err = client.Accounts().Update(ctx, deprecatedToken)
	if err != nil {
		log.Fatalf("Failed to update name of old token to deprecated: %v", err)
	}

	fmt.Println("Successfully updated name of old token")

	// Create new token with same name and permissions as old token
	newApiKey := &sdm.Token{
		Name:        tokenName,
		AccountType: oldToken.AccountType,
		Duration:    oldToken.Duration,
		Permissions: oldToken.Permissions,
	}

	createResponse, err = client.Accounts().Create(ctx, newApiKey)
	if err != nil {
		log.Fatalf("Could not create api key: %v", err)
	}

	newId := createResponse.Account.GetID()
	tokenName = createResponse.Account.(*sdm.Token).Name
	newAccessKey := createResponse.AccessKey
	newSecretKey := createResponse.SecretKey

	fmt.Println("Successfully created new api key.")
	fmt.Println("\tID:", newId)
	fmt.Println("\tName:", tokenName)
	fmt.Println("\tAccessKey:", newAccessKey)
	fmt.Println("\tSecretKey:", newSecretKey)

	// Delete the old token once the new token is successfully created
	_, err = client.Accounts().Delete(ctx, oldTokenID)
	if err != nil {
		log.Fatalf("Could not delete old token: %v", err)
	}
	fmt.Println("Successfully deleted old token.")
}
