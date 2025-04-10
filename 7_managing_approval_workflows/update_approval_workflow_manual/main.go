// Copyright 2024 StrongDM Inc
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

	// Create an approver account - used for creating an approval workflow approver
	accountCreateResponse, err := client.Accounts().Create(ctx, &sdm.User{
		Email:     "approval-workflow-approver-example@example.com",
		FirstName: "example",
		LastName:  "example",
	})
	if err != nil {
		log.Fatalf("Could not create approver: %v", err)
	}
	accountID := accountCreateResponse.Account.GetID()

	account2CreateResponse, err := client.Accounts().Create(ctx, &sdm.User{
		Email:     "approval-workflow-approver-example@example.com",
		FirstName: "example2",
		LastName:  "example2",
	})
	if err != nil {
		log.Fatalf("Could not create approver: %v", err)
	}
	account2ID := account2CreateResponse.Account.GetID()

	// Create an approver role - used for creating an approval workflow approver
	roleCreateResponse, err := client.Roles().Create(ctx, &sdm.Role{
		Name: "example role for approval workflow approver role",
	})
	if err != nil {
		log.Fatalf("Could not create role: %v", err)
	}
	roleID := roleCreateResponse.Role.ID

	// Create a manual grant approval workflow with multiple approval steps
	af := &sdm.ApprovalWorkflow{
		Name:         "Example approval workflow manual",
		Description:  "A manual grant approval workflow",
		ApprovalMode: "manual",
		ApprovalWorkflowSteps: []*sdm.ApprovalFlowStep{
			{
				Quantifier: "any",
				Approvers: []*sdm.ApprovalFlowApprover{
					{AccountID: accountID},
				},
			},
			{
				Quantifier: "all",
				SkipAfter:  time.Hour,
				Approvers: []*sdm.ApprovalFlowApprover{
					{AccountID: account2ID},
					{RoleID: roleID},
				},
			},
		},
	}

	resp, err := client.ApprovalWorkflows().Create(ctx, af)
	if err != nil {
		log.Fatalf("Could not create approval workflow: %v", err)
	}

	flow := resp.ApprovalWorkflow

	// Update approval workflow
	// Provide new configuration for approval workflow (Approval Workflow ID required)
	updatedFlow := &sdm.ApprovalWorkflow{
		ID:           flow.ID,
		Name:         "Example updated approval workflow",
		Description:  "An updated manual grant approval workflow",
		ApprovalMode: "manual",
		ApprovalWorkflowSteps: []*sdm.ApprovalFlowStep{
			{
				Quantifier: "all",
				SkipAfter:  time.Hour * 2,
				Approvers: []*sdm.ApprovalFlowApprover{
					{AccountID: accountID},
				},
			},
			{
				Quantifier: "any",
				Approvers: []*sdm.ApprovalFlowApprover{
					{AccountID: account2ID},
				},
			},
			{
				Quantifier: "any",
				SkipAfter:  time.Hour,
				Approvers: []*sdm.ApprovalFlowApprover{
					{RoleID: roleID},
				},
			},
		},
	}

	updated, err := client.ApprovalWorkflows().Update(ctx, updatedFlow)
	if err != nil {
		log.Fatalf("Could not update approval workflow: %v", err)
	}
	flow = updated.ApprovalWorkflow

	fmt.Println("Successfully update approval workflow.")
	fmt.Println("\tNew Name:", flow.Name)
	fmt.Println("\tNew Description:", flow.Description)
	fmt.Println("\tNew Approval Mode:", flow.ApprovalMode)
	fmt.Println("\tNew Approval Steps:", len(flow.ApprovalWorkflowSteps))

	// Updating an approval workflow from manual to autogrant deletes all approval steps and approvers
	flow.ApprovalMode = "automatic"
	flow.ApprovalWorkflowSteps = []*sdm.ApprovalFlowStep{}
	updated, err = client.ApprovalWorkflows().Update(ctx, flow)
	if err != nil {
		log.Fatalf("Could not update approval workflow: %v", err)
	}
	flow = updated.ApprovalWorkflow

	fmt.Println("Successfully update approval workflow.")
	fmt.Println("\tNew Name:", flow.Name)
	fmt.Println("\tNew Description:", flow.Description)
	fmt.Println("\tNew Approval Mode:", flow.ApprovalMode)
}
