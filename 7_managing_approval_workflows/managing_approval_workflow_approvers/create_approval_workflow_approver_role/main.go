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

	sdm "github.com/strongdm/strongdm-sdk-go/v6"
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

	// Create a manual approval workflow.
	approvalWorkflow := &sdm.ApprovalWorkflow{
		Name:         "Example Approval Workflow",
		ApprovalMode: "manual",
	}

	createResponse, err := client.ApprovalWorkflows().Create(ctx, approvalWorkflow)
	if err != nil {
		log.Fatalf("Could not create approval workflow: %v", err)
	}

	flow := createResponse.ApprovalWorkflow

	// Create an approval workflow step
	stepCreateResponse, err := client.ApprovalWorkflowSteps().Create(ctx, &sdm.ApprovalWorkflowStep{
		ApprovalFlowID: flow.ID,
	})
	if err != nil {
		log.Fatalf("Could not create approval workflow step: %v", err)
	}
	step := stepCreateResponse.ApprovalWorkflowStep

	// Create an approver role - used for creating an approval workflow approver
	roleCreateResponse, err := client.Roles().Create(ctx, &sdm.Role{
		Name: "example role for approval workflow approver role",
	})
	if err != nil {
		log.Fatalf("Could not create role: %v", err)
	}
	roleID := roleCreateResponse.Role.ID

	// Create an approval workflow approver
	approverCreateResponse, err := client.ApprovalWorkflowApprovers().Create(ctx, &sdm.ApprovalWorkflowApprover{
		ApprovalFlowID: flow.ID,
		ApprovalStepID: step.ID,
		RoleID:         roleID,
	})
	if err != nil {
		log.Fatalf("Could not create approval workflow approver: %v", err)
	}

	fmt.Println("Successfully created approval workflow approver.")
	fmt.Println("\tApproval Workflow ID:", flow.ID)
	fmt.Println("\tApproval Workflow Step ID:", stepCreateResponse.ApprovalWorkflowStep.ID)
	fmt.Println("\tRole ID:", roleID)
	fmt.Println("\tApproval Workflow Approver ID:", approverCreateResponse.ApprovalWorkflowApprover.ID)
}
