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

	// Create a manual Workflow with initial Access Rule
	workflow := &sdm.Workflow{
		Name:        "Full Example Create Manual Worfklow",
		Description: "Example Workflow Description",
		AccessRules: sdm.AccessRules{
			sdm.AccessRule{
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
	workflowID := wf.ID
	workflowName := wf.Name

	fmt.Println("Successfully created Workflow.")
	fmt.Println("\tID:", workflowID)
	fmt.Println("\tName:", workflowName)

	// To allow users access to the resources managed by this workflow, you must
	// add workflow roles to the workflow.
	// Two steps are needed to add a workflow role:
	// Step 1: create a Role
	// Step 2: create a WorkflowRole

	// Create a Role - used for creating a workflow role
	roleCreateResponse, err := client.Roles().Create(ctx, &sdm.Role{
		Name: "Example Role for Manual Workflow",
	})
	if err != nil {
		log.Fatalf("Could not create role: %v", err)
	}
	roleID := roleCreateResponse.Role.ID

	// Create a WorkflowRole
	_, err = client.WorkflowRoles().Create(ctx, &sdm.WorkflowRole{
		WorkflowID: workflowID,
		RoleID:     roleID,
	})
	if err != nil {
		log.Fatalf("Could not create workflow role: %v", err)
	}

	fmt.Println("Successfully created WorkflowRole.")
	fmt.Println("\tWorkflow ID:", workflowID)
	fmt.Println("\tRole ID:", roleID)

	// To manually enable this workflow, you must add workflow approvers
	// to this workflow.
	// Two steps are needed to add a workflow approver:
	// Step 1: create an Account
	// Step 2: create a WorkflowApprover

	// Create a approver - used for creating a workflow approver
	approverCreateResponse, err := client.Accounts().Create(ctx, &sdm.User{
		Email:     "create-workflow-full-example@example.com",
		FirstName: "example",
		LastName:  "example",
	})
	if err != nil {
		log.Fatalf("Could not create approver: %v", err)
	}
	approverID := approverCreateResponse.Account.GetID()

	// Create a WorkflowApprover
	_, err = client.WorkflowApprovers().Create(ctx, &sdm.WorkflowApprover{
		WorkflowID: workflowID,
		ApproverID: approverID,
	})
	if err != nil {
		log.Fatalf("Could not create workflow approver: %v", err)
	}

	fmt.Println("Successfully created WorkflowApprover.")
	fmt.Println("\tWorkflow ID:", workflowID)
	fmt.Println("\tApprover ID:", approverID)

	// You can enable this workflow after adding workflow approvers.
	// Update Workflow Enabled
	wf.Enabled = true
	workflowUpdateResponse, err := client.Workflows().Update(ctx, wf)
	if err != nil {
		log.Fatalf("Could not update workflow: %v", err)
	}
	wf = workflowUpdateResponse.Workflow

	fmt.Println("Successfully updated Workflow Enabled.")
	fmt.Println("\tWorkflow ID:", workflowID)
	fmt.Println("\tEnabled:", wf.Enabled)
}
