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

	sdm "github.com/strongdm/strongdm-sdk-go/v15"
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

	// Create an approver account that will participate in the manual approval flow.
	approverResp, err := client.Accounts().Create(ctx, &sdm.User{
		Email:     "manual-workflow-approver@example.com",
		FirstName: "Example",
		LastName:  "Approver",
	})
	if err != nil {
		log.Fatalf("Could not create approver account: %v", err)
	}
	approverID := approverResp.Account.GetID()

	// Create a role so we can demonstrate using role-based approvers in the approval workflow.
	roleCreateResponse, err := client.Roles().Create(ctx, &sdm.Role{
		Name: "Example Role for Manual Approval Workflow",
	})
	if err != nil {
		log.Fatalf("Could not create role: %v", err)
	}
	roleID := roleCreateResponse.Role.ID

	// Build a manual ApprovalWorkflow that references the approver created above.
	approvalWorkflow := &sdm.ApprovalWorkflow{
		Name:         "Example Manual Approval Workflow",
		Description:  "Manual approval flow to back the access workflow example",
		ApprovalMode: "manual",
		ApprovalWorkflowSteps: []*sdm.ApprovalFlowStep{
			{
				Quantifier: "any",
				Approvers: []*sdm.ApprovalFlowApprover{
					{AccountID: approverID},
				},
			},
		},
	}

	approvalCreateResponse, err := client.ApprovalWorkflows().Create(ctx, approvalWorkflow)
	if err != nil {
		log.Fatalf("Could not create approval workflow: %v", err)
	}
	createdApprovalWorkflow := approvalCreateResponse.ApprovalWorkflow

	fmt.Println("Successfully created manual ApprovalWorkflow.")
	fmt.Println("\tID:", createdApprovalWorkflow.ID)
	fmt.Println("\tName:", createdApprovalWorkflow.Name)

	// Create the Workflow and bind it to the ApprovalWorkflow by setting ApprovalFlowID.
	workflow := &sdm.Workflow{
		Name:           "Full Example Create Manual Workflow",
		Description:    "Example Workflow Description",
		ApprovalFlowID: createdApprovalWorkflow.ID,
		Enabled:        true,
		AccessRules: sdm.AccessRules{
			{
				Tags: sdm.Tags{
					"env": "dev",
				},
			},
		},
	}

	createResponse, err := client.Workflows().Create(ctx, workflow)
	if err != nil {
		log.Fatalf("Could not create workflow: %v", err)
	}

	wf := createResponse.Workflow

	fmt.Println("Successfully created Workflow.")
	fmt.Println("\tID:", wf.ID)
	fmt.Println("\tName:", wf.Name)
	fmt.Println("\tApproval Flow ID:", wf.ApprovalFlowID)

	// Add a workflow role so members of the role can request access through this workflow.
	_, err = client.WorkflowRoles().Create(ctx, &sdm.WorkflowRole{
		WorkflowID: wf.ID,
		RoleID:     roleID,
	})
	if err != nil {
		log.Fatalf("Could not create workflow role: %v", err)
	}

	fmt.Println("Successfully created WorkflowRole.")
	fmt.Println("\tWorkflow ID:", wf.ID)
	fmt.Println("\tRole ID:", roleID)
}
