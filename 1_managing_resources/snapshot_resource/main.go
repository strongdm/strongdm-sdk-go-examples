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

	// Define the Postgres Datasource
	// Set `PortOverride` to `-1` to auto-generate a port if Port Overrides is enabled.
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

// Example setup to get snapshot of resource
func exampleResourceDeletion(client *sdm.Client) {
	start := time.Now()

	ctx := context.Background()

	// Set up some fake logs to query
	/* * * */
	postgresResp, err := client.Resources().Create(ctx, &sdm.Postgres{
		Name:     "example-postgres",
		Hostname: "example-postgres",
		Username: "example-username",
	})
	if err != nil {
		log.Fatalf("failed to create postgres: %v", err)
	}

	createdAt := time.Now()

	resourceID := postgresResp.Resource.GetID()

	postgresResp.Resource.SetName("example-postgres-renamed")
	_, err = client.Resources().Update(ctx, postgresResp.Resource)
	if err != nil {
		log.Fatalf("failed to rename postgres: %v", err)
	}

	renamedAt := time.Now()

	_, err = client.Resources().Delete(ctx, resourceID)
	if err != nil {
		log.Fatalf("failed to delete postgres: %v", err)
	}

	deletedAt := time.Now()
	/* * * */

	_, err = client.SnapshotAt(start).Resources().Get(ctx, resourceID)
	fmt.Println(err) // Does not exist

	getResp, err := client.SnapshotAt(createdAt).Resources().Get(ctx, resourceID)
	if err != nil {
		log.Fatalf("failed to retrieve created postgres: %v", err)
	}

	fmt.Println(getResp.Resource.GetName()) // example-postgres

	getResp, err = client.SnapshotAt(renamedAt).Resources().Get(ctx, resourceID)
	if err != nil {
		log.Fatalf("failed to retrieve created postgres: %v", err)
	}

	fmt.Println(getResp.Resource.GetName()) // example-postgres-renamed

	_, err = client.SnapshotAt(deletedAt).Resources().Get(ctx, resourceID)
	fmt.Println(err) // Does not exist

	history, err := client.ResourcesHistory().List(ctx, "id:?", resourceID)
	if err != nil {
		log.Fatalf("failed to list resource history: %v", err)
	}
	for history.Next() {
		v := history.Value()
		activity, err := client.Activities().Get(ctx, v.ActivityID)
		if err != nil {
			log.Fatalf("failed to look up history: %v", err)
		}
		fmt.Println(activity.Activity.Description) // Created, updated, deleted resource; in order
	}
	if err := history.Err(); err != nil {
		log.Fatalf("failed to finish resource history: %v", err)
	}

	//The example snapshot setup prints, as expected:

	//item does not exist: no resource found
	//example-postgres
	//example-postgres-renamed
	//item does not exist: no resource found
	//API Account [Testing] API Key (ddf09805-538a-4fe2-bbf9-415e6cd1d0b8) created a new postgres datasource named example-postgres.
	//API Account [Testing] API Key (ddf09805-538a-4fe2-bbf9-415e6cd1d0b8) updated datasource example-postgres-renamed.
	//API Account [Testing] API Key (ddf09805-538a-4fe2-bbf9-415e6cd1d0b8) deleted datasource example-postgres-renamed.
}