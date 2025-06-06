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

	sdm "github.com/strongdm/strongdm-sdk-go/v14"
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

	// Create an auto grant approval workflow.
	approvalWorkflow := &sdm.ApprovalWorkflow{
		Name:         "Example Update Approval Workflow",
		ApprovalMode: "automatic",
	}

	createResponse, err := client.ApprovalWorkflows().Create(ctx, approvalWorkflow)
	if err != nil {
		log.Fatalf("Could not create approval workflow: %v", err)
	}

	flow := createResponse.ApprovalWorkflow

	fmt.Println("Successfully created approval workflow.")
	fmt.Println("\tID:", flow.ID)
	fmt.Println("\tName:", flow.Name)

	// Get an approval workflow by id
	getResp, err := client.ApprovalWorkflows().Get(ctx, flow.ID)
	if err != nil {
		log.Fatalf("Could not get approval workflow: %v", err)
	}
	fmt.Println("Successfully got approval workflow.")
	fmt.Println("\tID:", getResp.ApprovalWorkflow.ID)
	fmt.Println("\tName:", getResp.ApprovalWorkflow.Name)

	// Update approval workflow Name
	newName := "Example New Name"
	flow.Name = newName
	updated, err := client.ApprovalWorkflows().Update(ctx, flow)
	if err != nil {
		log.Fatalf("Could not update approval workflow: %v", err)
	}
	flow = updated.ApprovalWorkflow

	fmt.Println("Successfully update approval workflow name.")
	fmt.Println("\tNew Name:", flow.Name)

	// Update approval workflow Description
	description := "Example New Description"
	flow.Description = description
	updated, err = client.ApprovalWorkflows().Update(ctx, flow)
	if err != nil {
		log.Fatalf("Could not update approval workflow: %v", err)
	}
	flow = updated.ApprovalWorkflow

	fmt.Println("Successfully update approval workflow description.")
	fmt.Println("\tNew Description:", flow.Description)

	// Update approval workflow approval mode
	newMode := "manual"
	flow.ApprovalMode = newMode
	updated, err = client.ApprovalWorkflows().Update(ctx, flow)
	if err != nil {
		log.Fatalf("Could not update approval workflow: %v", err)
	}
	flow = updated.ApprovalWorkflow

	fmt.Println("Successfully update approval workflow approval mode.")
	fmt.Println("\tNew Approval Mode:", flow.ApprovalMode)

	// Delete the approval workflow
	_, err = client.ApprovalWorkflows().Delete(ctx, flow.ID)
	if err != nil {
		log.Fatalf("Could not delete approval workflow: %v", err)
	}
	fmt.Println("Successfully deleted approval workflow.")
}
