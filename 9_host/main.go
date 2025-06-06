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

	sdm "github.com/strongdm/strongdm-sdk-go/v11"
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

	// Configure a client to communicate with the UK host.
	// If the WithHost option is not provided, it will default to the US control plane (api.strongdm.com:443)
	client, err := sdm.New(accessKey, secretKey, sdm.WithHost(sdm.APIHostUK))
	if err != nil {
		log.Fatal("failed to create strongDM client:", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Create a Postgres Datasource for example
	datasource := &sdm.Postgres{
		Name:         "Example Postgres Datasource",
		Hostname:     "example.strongdm.com",
		Port:         5432,
		Username:     "example",
		Password:     "example",
		Database:     "example",
		PortOverride: 19999,
		Tags: sdm.Tags{
			"example": "example",
		},
	}

	createResponse, err := client.Resources().Create(ctx, datasource)
	if err != nil {
		log.Fatalf("Could not create Postgres datasource: %v", err)
	}

	fmt.Println("Successfully created Postgres datasource.")
	fmt.Println("\tID:", createResponse.Resource.GetID())
}
