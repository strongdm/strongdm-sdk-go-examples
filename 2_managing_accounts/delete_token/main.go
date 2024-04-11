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
		Permissions: []string{"role:create", "user:create_admin_token"},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	createResponse, err := client.Accounts().Create(ctx, apiKey)
	if err != nil {
		log.Fatalf("Could not create api key: %v", err)
	}

	id := createResponse.Account.GetID()

	fmt.Println("Successfully created api key.")
	fmt.Println("\tID:", id)

	// Delete the Token
	_, err = client.Accounts().Delete(ctx, id)
	if err != nil {
		log.Fatalf("Could not delete account: %v", err)
	}
	fmt.Println("Successfully deleted account.")
}
