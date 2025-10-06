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

	sdm "github.com/strongdm/strongdm-sdk-go/v15"
)

func main() {
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
		log.Fatalf("could not create client: %v", err)
	}

	// Define the Postgres Datasource
	examplePostgresDatasource := &sdm.Postgres{
		Name:         "Example Postgres Datasource for Update",
		Hostname:     "example.strongdm.com",
		Port:         5432,
		Username:     "example",
		Password:     "example",
		Database:     "example",
		PortOverride: 19202,
		Tags: sdm.Tags{
			"example": "example",
		},
	}

	// Create the Datasource
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	createResponse, err := client.Resources().Create(ctx, examplePostgresDatasource)
	if err != nil {
		log.Fatalf("Could not create Postgres datasource: %v", err)
	}

	id := createResponse.Resource.GetID()

	fmt.Println("Successfully created Postgres datasource.")
	fmt.Println("\tID:", id)
	fmt.Println("\tName:", createResponse.Resource.GetName())

	// Load the Datasource to update
	getResponse, err := client.Resources().Get(ctx, id)
	if err != nil {
		log.Fatalf("Could not read Postgres datasource: %v", err)
	}
	updatedPostgresDatasource := getResponse.Resource

	// Update the fields to change
	updatedPostgresDatasource.SetName("Example Name Updated")

	// If your organization has Virtual Networking Mode enabled,
	// you can automatically allocate an IP to that resource via the ResourceIPAllocationModeVNM constant...
	updatedPostgresDatasource.SetBindInterface(sdm.ResourceIPAllocationModeVNM)

	// ...Or fallback to whatever the default behavior is for your organization...
	updatedPostgresDatasource.SetBindInterface(sdm.ResourceIPAllocationModeDefault)

	// ...Or if there is a specific IP to bind to, you can specify it directly.
	// For more details on Virtual Networking Mode see documentation here:
	// https://docs.strongdm.com/admin/clients/client-networking/virtual-networking-mode
	updatedPostgresDatasource.SetBindInterface("127.0.0.1")

	if pg, ok := updatedPostgresDatasource.(*sdm.Postgres); ok {
		// Update `PortOverride` to `-1` to auto-allocate a different available port.
		pg.PortOverride = -1
	}

	// Update the Datasource
	updateResponse, err := client.Resources().Update(ctx, updatedPostgresDatasource)
	if err != nil {
		log.Fatalf("Could not update Postgres datasource: %v", err)
	}

	newPortOverride := 0
	if pg, ok := updateResponse.Resource.(*sdm.Postgres); ok {
		newPortOverride = int(pg.PortOverride)
	}

	fmt.Println("Successfully updated Postgres datasource.")
	fmt.Println("\tID:", updateResponse.Resource.GetID())
	fmt.Println("\tName:", updateResponse.Resource.GetName())
	fmt.Println("\tBindInterface:", updateResponse.Resource.GetBindInterface())
	fmt.Println("\tPortOverride:", newPortOverride)
}
