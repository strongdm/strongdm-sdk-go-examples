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
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	sdm "github.com/strongdm/strongdm-sdk-go/v11"
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
		log.Fatalf("Could not create client: %v", err)
	}

	// Define a Policy to add to the system.
	policy := &sdm.Policy{
		Name:        "forbid-everything",
		Description: "Forbid everything",
		Policy:      `forbid ( principal, action, resource );`,
	}

	// Create the Policy
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	createResponse, err := client.Policies().Create(ctx, policy)
	if err != nil {
		log.Fatalf("Could not create policy: %v", err)
	}

	fmt.Println("Successfully created a policy to forbid all actions.")
	fmt.Println("\tID:", createResponse.Policy.ID)
	fmt.Println("\tName:", createResponse.Policy.Name)

	// Update the Policy

	// Note: The `Policy` field in `createResponse` can also be used to
	// make an update. However, we'll load it from the API to
	// demonstrate `Get()`.
	getResponse, err := client.Policies().Get(ctx, createResponse.Policy.ID)
	if err != nil {
		log.Fatalf("could not read policy: %v", err)
	}

	fmt.Println("Successfully retrieved policy.")
	fmt.Println("\tID:", getResponse.Policy.ID)
	fmt.Println("\tName:", getResponse.Policy.Name)

	updatePolicy := getResponse.Policy
	updatePolicy.Name = "forbid-one-thing"
	updatePolicy.Description = "Forbid connecting to the bad resource"
	updatePolicy.Policy = `forbid (
     principal,
     action == StrongDM::Action::"connect",
     resource == StrongDM::Resource::"rs-123d456789"
);`

	updatedResponse, err := client.Policies().Update(ctx, updatePolicy)
	if err != nil {
		log.Fatalf("Could not update policy: %v", err)
	}

	fmt.Println("Successfully updated Policy.")
	fmt.Println("\tID:", updatedResponse.Policy.ID)
	fmt.Println("\tName:", updatedResponse.Policy.Name)
	fmt.Println("\tDescription:", updatedResponse.Policy.Description)
	fmt.Println("\tPolicy:", updatedResponse.Policy.Policy)

	// Get the updated policy.
	getResponse, err = client.Policies().Get(ctx, createResponse.Policy.ID)
	if err != nil {
		log.Fatalf("Could not read policy: %v", err)
	}

	fmt.Println("Successfully retrieved updated policy.")
	fmt.Println("\tID:", getResponse.Policy.ID)
	fmt.Println("\tName:", getResponse.Policy.Name)
	fmt.Println("\tDescription:", getResponse.Policy.Description)
	fmt.Println("\tPolicy:", getResponse.Policy.Policy)

	// Delete the policy
	_, err = client.Policies().Delete(ctx, createResponse.Policy.ID)
	if err != nil {
		log.Fatalf("Could not delete policy: %v", err)
	}
	fmt.Println("Successfully deleted the policy.")

	// Trying to retrieve the policy again should result in a
	// *sdm.NotFoundError
	var notFound *sdm.NotFoundError
	_, err = client.Policies().Get(ctx, createResponse.Policy.ID)
	if err == nil {
		log.Fatalf("Something went wrong. We found a deleted policy.")
	} else if errors.As(err, &notFound) {
		fmt.Println("Successfully could not retrieve deleted policy.")
	} else {
		log.Fatalf("Could not read policy: %v", err)
	}
}
