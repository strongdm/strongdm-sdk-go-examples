// Copyright 2023 StrongDM Inc
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

	sdm "github.com/strongdm/strongdm-sdk-go/v4"
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
	client, err := sdm.New(accessKey, secretKey)
	if err != nil {
		log.Fatal("failed to create strongDM client:", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Create a Workflow
	workflow := &sdm.Workflow{
		Name:        "Example Update Workflow",
		Description: "Example Workflow Description",
	}

	createResponse, err := client.Workflows().Create(ctx, workflow)
	if err != nil {
		log.Fatalf("Could not create workflow: %v", err)
	}

	wf := createResponse.Workflow

	// Update Workflow Name
	newName := "Example New Name"
	wf.Name = newName
	updated, err := client.Workflows().Update(ctx, wf)
	if err != nil {
		log.Fatalf("Could not update workflow: %v", err)
	}
	wf = updated.Workflow

	fmt.Println("Successfully update Workflow Name.")
	fmt.Println("\tNew Name:", wf.Name)

	// Update Workflow Description
	description := "Example New Description"
	wf.Description = description
	updated, err = client.Workflows().Update(ctx, wf)
	if err != nil {
		log.Fatalf("Could not update workflow: %v", err)
	}
	wf = updated.Workflow

	fmt.Println("Successfully update Workflow Description.")
	fmt.Println("\tNew Description:", wf.Description)

	// Update Workflow Weight
	oldWeight := wf.Weight
	wf.Weight = oldWeight + 20
	updated, err = client.Workflows().Update(ctx, wf)
	if err != nil {
		log.Fatalf("Could not update workflow: %v", err)
	}
	wf = updated.Workflow

	fmt.Println("Successfully update Workflow Weight.")
	fmt.Println("\tNew Weight:", wf.Weight)

	// Update Workflow AutoGrant
	auto := wf.AutoGrant
	wf.AutoGrant = !auto
	updated, err = client.Workflows().Update(ctx, wf)
	if err != nil {
		log.Fatalf("Could not update workflow: %v", err)
	}
	wf = updated.Workflow

	fmt.Println("Successfully update Workflow AutoGrant.")
	fmt.Println("\tAutoGrant:", wf.AutoGrant)

	// Update Workflow Enabled
	// The requirements to enable a workflow are that the workflow must be either set
	// up for with auto grant enabled or have one or more WorkflowApprovers created for
	// the workflow.
	wf.AutoGrant = true
	wf.Enabled = true
	updated, err = client.Workflows().Update(ctx, wf)
	if err != nil {
		log.Fatalf("Could not update workflow: %v", err)
	}
	wf = updated.Workflow

	fmt.Println("Successfully update Workflow Enabled.")
	fmt.Println("\tEnabled:", wf.Enabled)

}
