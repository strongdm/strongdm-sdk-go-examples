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

	client, err := sdm.New(
		accessKey,
		secretKey,
	)
	if err != nil {
		log.Fatalf("could not create client: %v", err)
	}

	user := &sdm.User{
		Email:     "example@example.com",
		FirstName: "example",
		LastName:  "example",
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	createResponse, err := client.Accounts().Create(ctx, user)
	if err != nil {
		log.Fatalf("Could not create user: %v", err)
	}

	id := createResponse.Account.GetID()
	fmt.Println("Successfully created user.")
	fmt.Println("\tID:", id)

	// Load the account to update
	getResponse, err := client.Accounts().Get(ctx, id)
	if err != nil {
		log.Fatalf("Could not get account: %v", err)
	}
	account := getResponse.Account

	// Set the fields to change
	account.SetSuspended(true)

	// Update the account
	updateResponse, err := client.Accounts().Update(ctx, account)
	if err != nil {
		log.Fatalf("Could not update account: %v", err)
	}

	fmt.Println("Successfully suspended account.")
	fmt.Println("\tID:", updateResponse.Account.GetID())
	fmt.Println("\tSuspended:", updateResponse.Account.IsSuspended())
}
