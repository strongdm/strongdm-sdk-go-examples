// Copyright 2024 StrongDM Inc
//
// Licensed under the Apache License, Version 2.0 (the
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

var examplePolicies = []sdm.Policy{
	{
		Name:        "default-permit-policy",
		Description: "a default permit policy",
		Policy:      "permit (principal, action, resource);",
	},
	{
		Name:        "permit-sql-select-policy",
		Description: "a permit sql select policy",
		Policy:      `permit (principal, action == SQL::Action::"select", resource == Postgres::Database::"*");`,
	},
	{
		Name:        "default-forbid-policy",
		Description: "a default forbid policy",
		Policy:      "forbid (principal, action, resource);",
	},
	{
		Name:        "forbid-connect-policy",
		Description: "a forbid connect policy",
		Policy:      `forbid (principal, action == StrongDM::Action::"connect", resource);`,
	},
	{
		Name:        "forbid-sql-delete-policy",
		Description: "a forbid delete policy on all resources",
		Policy:      `forbid (principal, action == SQL::Action::"delete", resource == Postgres::Database::"*");`,
	},
}

func cleanup(client *sdm.Client, policyId string) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	_, err := client.Policies().Delete(ctx, policyId)
	if err != nil {
		log.Printf("Unable to cleanup policy with id=%q: %v", policyId, err)
	}
}

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
		log.Fatalf("Could not create client: %v", err)
	}

	// Create the example policies so we can find them.
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	for _, p := range examplePolicies {
		createdPolicy, err := client.Policies().Create(ctx, &p)
		if err != nil {
			log.Fatalf("Couldn't create example policy (%s): %v", p.Name, err)
		}
		defer cleanup(client, createdPolicy.Policy.ID)
		fmt.Println("Successfully created Policy.")
		fmt.Println("\tID:", createdPolicy.Policy.ID)
		fmt.Println("\tName:", createdPolicy.Policy.Name)
	}

	// Find policies that related to `sql` by Name
	fmt.Println("Finding all Policies with a name containing 'sql'")
	listResp, err := client.Policies().List(ctx, "name:*sql*")
	if err != nil {
		log.Fatalf("Could not list policies: %v", err)
	}
	for listResp.Next() {
		p := listResp.Value()
		fmt.Printf("\tID: %s\tName: %s\n", p.ID, p.Name)
	}

	// Find policies that forbid based on the Policy
	fmt.Println("Finding all Policies with a name containing 'sql'")
	listResp, err = client.Policies().List(ctx, "policy:forbid*")
	if err != nil {
		log.Fatalf("Could not list policies: %v", err)
	}
	for listResp.Next() {
		p := listResp.Value()
		fmt.Printf("\tID: %s\tName: %s\n", p.ID, p.Name)
	}

}
