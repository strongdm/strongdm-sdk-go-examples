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

	// Create a role
	create := &sdm.Role{
		Name: "example role",
	}

	createResponse, err := client.Roles().Create(ctx, create)
	if err != nil {
		log.Fatalf("Could not create role: %v", err)
	}

	roleID := createResponse.Role.ID
	log.Printf("Successfully created role.\n\tID: %v\n", roleID)

	// Get the role
	getResponse, err := client.Roles().Get(ctx, roleID)
	if err != nil {
		log.Fatalf("Could not get role: %v", err)
	}

	// Modify fields
	role := getResponse.Role
	role.Name = "example role updated"

	// Update the role
	updateResponse, err := client.Roles().Update(ctx, role)
	if err != nil {
		log.Fatalf("Could not update role: %v", err)
	}

	log.Println("Successfully updated role.")
	log.Println("\tID:", updateResponse.Role.ID)
	log.Println("\tName:", updateResponse.Role.Name)
}
