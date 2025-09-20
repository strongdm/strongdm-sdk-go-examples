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

	sdm "github.com/strongdm/strongdm-sdk-go/v15"
)

// Example showing how to create approval workflows using groups as approvers
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

	// Create approver groups - these groups will be designated as approvers
	securityGroupResponse, err := client.Groups().Create(ctx, &sdm.Group{
		Name: "Security Team",
	})
	if err != nil {
		log.Fatalf("Could not create security group: %v", err)
	}
	securityGroupID := securityGroupResponse.Group.ID
	fmt.Println("Created Security Team group:", securityGroupID)

	adminGroupResponse, err := client.Groups().Create(ctx, &sdm.Group{
		Name: "Administrators",
	})
	if err != nil {
		log.Fatalf("Could not create admin group: %v", err)
	}
	adminGroupID := adminGroupResponse.Group.ID
	fmt.Println("Created Administrators group:", adminGroupID)

	devOpsGroupResponse, err := client.Groups().Create(ctx, &sdm.Group{
		Name: "DevOps Team",
	})
	if err != nil {
		log.Fatalf("Could not create devops group: %v", err)
	}
	devOpsGroupID := devOpsGroupResponse.Group.ID
	fmt.Println("Created DevOps Team group:", devOpsGroupID)

	// Create some users to add to groups (demonstrating group membership)
	securityUserResponse, err := client.Accounts().Create(ctx, &sdm.User{
		Email:     "security-lead@example.com",
		FirstName: "Security",
		LastName:  "Lead",
	})
	if err != nil {
		log.Fatalf("Could not create security user: %v", err)
	}
	securityUserID := securityUserResponse.Account.GetID()

	adminUserResponse, err := client.Accounts().Create(ctx, &sdm.User{
		Email:     "admin-user@example.com",
		FirstName: "Admin",
		LastName:  "User",
	})
	if err != nil {
		log.Fatalf("Could not create admin user: %v", err)
	}
	adminUserID := adminUserResponse.Account.GetID()

	// Add users to their respective groups
	_, err = client.AccountsGroups().Create(ctx, &sdm.AccountGroup{
		AccountID: securityUserID,
		GroupID:   securityGroupID,
	})
	if err != nil {
		log.Fatalf("Could not add security user to group: %v", err)
	}
	fmt.Println("Added security user to Security Team group")

	_, err = client.AccountsGroups().Create(ctx, &sdm.AccountGroup{
		AccountID: adminUserID,
		GroupID:   adminGroupID,
	})
	if err != nil {
		log.Fatalf("Could not add admin user to group: %v", err)
	}
	fmt.Println("Added admin user to Administrators group")

	// Create a manual approval workflow with groups as approvers
	af := &sdm.ApprovalWorkflow{
		Name:         "Group-Based Approval Workflow",
		Description:  "A workflow demonstrating group-based approvers",
		ApprovalMode: "manual",
		ApprovalWorkflowSteps: []*sdm.ApprovalFlowStep{
			{
				// First step: Any member of the Security Team can approve
				Quantifier: "any",
				Approvers: []*sdm.ApprovalFlowApprover{
					{GroupID: securityGroupID},
				},
			},
			{
				// Second step: All specified groups/roles must approve
				Quantifier: "all",
				SkipAfter:  time.Hour * 2,
				Approvers: []*sdm.ApprovalFlowApprover{
					{GroupID: adminGroupID},                              // Administrators group
					{GroupID: devOpsGroupID},                             // DevOps Team group
					{Reference: sdm.ApproverReferenceManagerOfRequester}, // Plus manager
				},
			},
			{
				// Third step: Mixed approvers - combination of groups and references
				Quantifier: "any",
				SkipAfter:  time.Hour,
				Approvers: []*sdm.ApprovalFlowApprover{
					{GroupID: securityGroupID}, // Security Team
					{GroupID: adminGroupID},    // Administrators
					{Reference: sdm.ApproverReferenceManagerOfManagerOfRequester},
				},
			},
		},
	}

	resp, err := client.ApprovalWorkflows().Create(ctx, af)
	if err != nil {
		log.Fatalf("Could not create approval workflow: %v", err)
	}

	fmt.Println("\nSuccessfully created group-based approval workflow.")
	fmt.Println("\tID:", resp.ApprovalWorkflow.ID)
	fmt.Println("\tName:", resp.ApprovalWorkflow.Name)
	fmt.Println("\tDescription:", resp.ApprovalWorkflow.Description)
	fmt.Println("\tNumber of Approval Steps:", len(resp.ApprovalWorkflow.ApprovalWorkflowSteps))

	for i, step := range resp.ApprovalWorkflow.ApprovalWorkflowSteps {
		fmt.Printf("\nStep %d:\n", i+1)
		fmt.Printf("\tQuantifier: %s\n", step.Quantifier)
		if step.SkipAfter > 0 {
			fmt.Printf("\tSkip After: %v\n", step.SkipAfter)
		}
		fmt.Printf("\tApprovers:\n")
		for _, approver := range step.Approvers {
			if approver.AccountID != "" {
				fmt.Printf("\t\t- Account ID: %s\n", approver.AccountID)
			} else if approver.RoleID != "" {
				fmt.Printf("\t\t- Role ID: %s\n", approver.RoleID)
			} else if approver.GroupID != "" {
				fmt.Printf("\t\t- Group ID: %s\n", approver.GroupID)
			} else if approver.Reference != "" {
				fmt.Printf("\t\t- Reference: %s\n", approver.Reference)
			}
		}
	}

	// Demonstrate updating workflow to use different group combinations
	flow := resp.ApprovalWorkflow
	updatedFlow := &sdm.ApprovalWorkflow{
		ID:           flow.ID,
		Name:         "Updated Group-Based Approval Workflow",
		Description:  "Updated workflow with different group approver combinations",
		ApprovalMode: "manual",
		ApprovalWorkflowSteps: []*sdm.ApprovalFlowStep{
			{
				// Single step with multiple group options
				Quantifier: "any",
				SkipAfter:  time.Hour * 24, // 24 hour timeout
				Approvers: []*sdm.ApprovalFlowApprover{
					{GroupID: securityGroupID},
					{GroupID: adminGroupID},
					{GroupID: devOpsGroupID},
				},
			},
		},
	}

	updated, err := client.ApprovalWorkflows().Update(ctx, updatedFlow)
	if err != nil {
		log.Fatalf("Could not update approval workflow: %v", err)
	}

	fmt.Println("\nSuccessfully updated approval workflow:")
	fmt.Println("\tNew Name:", updated.ApprovalWorkflow.Name)
	fmt.Println("\tNew Description:", updated.ApprovalWorkflow.Description)
	fmt.Println("\tSteps after update:", len(updated.ApprovalWorkflow.ApprovalWorkflowSteps))

	step := updated.ApprovalWorkflow.ApprovalWorkflowSteps[0]
	fmt.Printf("\tApprovers in updated step (any of these groups can approve):\n")
	for _, approver := range step.Approvers {
		if approver.GroupID != "" {
			fmt.Printf("\t\t- Group ID: %s\n", approver.GroupID)
		}
	}

	fmt.Println("\nExample demonstrates:")
	fmt.Println("  • Creating groups to act as approvers")
	fmt.Println("  • Adding users to groups (group membership)")
	fmt.Println("  • Using GroupID in approval workflow steps")
	fmt.Println("  • Combining group approvers with other approver types")
	fmt.Println("  • Different quantifiers (any/all) for group-based approval")

	// Clean up - delete the approval workflow
	_, err = client.ApprovalWorkflows().Delete(ctx, updated.ApprovalWorkflow.ID)
	if err != nil {
		log.Fatalf("Could not delete approval workflow: %v", err)
	}
	fmt.Println("\nCleaned up approval workflow.")
}
