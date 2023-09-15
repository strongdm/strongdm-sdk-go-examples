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

	// Create a Workflow
	workflow := &sdm.Workflow{
		Name:        "Example Create Worfklow",
		Description: "Example Workflow Description",
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	createResponse, err := client.Workflows().Create(ctx, workflow)
	if err != nil {
		log.Fatalf("Could not create workflow: %v", err)
	}

	wf := createResponse.Workflow
	workflowID := createResponse.Workflow.ID
	workflowName := createResponse.Workflow.Name

	fmt.Println("Successfully created Workflow.")
	fmt.Println("\tID:", workflowID)
	fmt.Println("\tName:", workflowName)

	// Create a Role - used for creating a workflow role
	roleCreateResponse, err := client.Roles().Create(ctx, &sdm.Role{
		Name: "Example Create Role",
	})
	if err != nil {
		log.Fatalf("Could not create role: %v", err)
	}
	roleID := roleCreateResponse.Role.ID

	// Create WorkflowRole
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

	// Create a approver - used for creating a workflow approver
	approverCreateResponse, err := client.Accounts().Create(ctx, &sdm.User{
		Email:     "create-user@example.com",
		FirstName: "example",
		LastName:  "example",
	})
	if err != nil {
		log.Fatalf("Could not create approver: %v", err)
	}
	approverID := approverCreateResponse.Account.GetID()

	// Create WorkflowApprover
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

	// Create a resource - used for workflow assignments
	tags := sdm.Tags{"example": "example"}
	resourceCreateResponse, err := client.Resources().Create(ctx, &sdm.Mysql{
		Name:         "Example resource",
		PortOverride: 19999,
		Hostname:     "example.strongdm.com",
		Username:     "example",
		Password:     "example",
		Database:     "example",
		Tags:         tags,
	})
	if err != nil {
		log.Fatalf("Could not create resource: %v", err)
	}
	resourceID := resourceCreateResponse.Resource.GetID()

	// Update workflow assignments
	wf.AccessRules = sdm.AccessRules{
		sdm.AccessRule{
			Type: "mysql",
			Tags: tags,
		},
	}

	// Update workflow enabled
	wf.Enabled = true

	// Update workflow
	workflowUpdateResponse, err := client.Workflows().Update(ctx, wf)
	if err != nil {
		log.Fatalf("Could not update workflow: %v", err)
	}

	fmt.Println("Successfully updated WorkflowAssignment.")
	fmt.Println("\tWorkflow ID:", workflowID)
	fmt.Println("\tResource ID:", resourceID)

	fmt.Println("Successfully updated Workflow Enabled.")
	fmt.Println("\tWorkflow ID:", workflowID)
	fmt.Println("\tEnabled:", workflowUpdateResponse.Workflow.Enabled)

}
