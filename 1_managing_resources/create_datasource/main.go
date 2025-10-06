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
	// Load the SDM API keys from the environment.
	// If these values are not set in your environment,
	// please follow the documentation here:
	// https://docs.strongdm.com/references/api/api-keys
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
	datasource := &sdm.Postgres{
		Name:     "Example Postgres Datasource",
		Hostname: "example.strongdm.com",
		Port:     5432,
		Username: "example",
		Password: "example",
		Database: "example",
		// May be set to one of the ResourceIPAllocationMode constants to select between VNM,
		// loopback, or default allocation. If not set, will be behave as if configured for
		// 'default'.
		// For more details on Virtual Networking Mode see documentation here:
		// https://docs.strongdm.com/admin/clients/client-networking/virtual-networking-mode
		BindInterface: sdm.ResourceIPAllocationModeLoopback,
		// Set `PortOverride` to `-1` to auto-allocate an available port.
		PortOverride: 19999,
		Tags: sdm.Tags{
			"example": "example",
		},
	}

	// You can also specify an explicit loopback IP address to bind to if Loopback IP Ranges
	// are enabled as documented here:
	// https://docs.strongdm.com/admin/clients/client-networking/loopback-ip-ranges
	datasource.BindInterface = "127.0.0.2"

	// ...Or if your organization has Virtual Networking Mode enabled,
	// you may configure the resource's bind interface
	// to automatically get an IP allocated upon creation:
	datasource.BindInterface = sdm.ResourceIPAllocationModeVNM

	// ...Or specify an explicit VNM IP address to bind to...
	datasource.BindInterface = "100.64.0.1"

	// Your organization can default either to 'loopback' or 'vnm', and
	// ResourceIPAllocationModeDefault will honor that.
	datasource.BindInterface = sdm.ResourceIPAllocationModeDefault

	// Create the Datasource
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	createResponse, err := client.Resources().Create(ctx, datasource)
	if err != nil {
		log.Fatalf("Could not create Postgres datasource: %v", err)
	}

	fmt.Println("Successfully created Postgres datasource.")
	fmt.Println("\tID:", createResponse.Resource.GetID())
	fmt.Println("\tName:", createResponse.Resource.GetName())
}
