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

	// Define the RDP server
	// Set `PortOverride` to `-1` to auto-generate a port if Port Overrides is enabled.
	server := &sdm.RDP{
		Name:         "Example RDP Server",
		Hostname:     "example.strongdm.com",
		Username:     "example",
		Password:     "example",
		Port:         3389,
		PortOverride: 19200,
		Tags: sdm.Tags{
			"example": "example",
		},
	}

	// Create the server
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	createResponse, err := client.Resources().Create(ctx, server)
	if err != nil {
		log.Fatalf("Could not create RDP server: %v", err)
	}

	fmt.Println("Successfully created RDP server.")
	fmt.Println("\tID:", createResponse.Resource.GetID())
	fmt.Println("\tName:", createResponse.Resource.GetName())
}
