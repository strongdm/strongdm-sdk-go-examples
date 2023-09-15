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
		Name:        "Example Create Worfklow",
		Description: "Example Workflow Description",
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

	// Create a resource - used for workflow assignments
	tags := sdm.Tags{"example": "example"}
	resourceCreateResponse, err := client.Resources().Create(ctx, &sdm.Mysql{
		Name:         "Example MySQL Resource for Workflow Assignments",
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
	_, err = client.Workflows().Update(ctx, wf)
	if err != nil {
		log.Fatalf("Could not update workflow: %v", err)
	}

	fmt.Println("Successfully updated WorkflowAssignment.")
	fmt.Println("\tWorkflow ID:", workflowID)
	fmt.Println("\tResource ID:", resourceID)

	// List WorkflowAssignments
	filter := "resource:" + resourceID + " workflow:" + workflowID
	listResp, err := client.WorkflowAssignments().List(ctx, filter)
	if err != nil {
		log.Fatalf("Could not list workflow assignments: %v", err)
	}

	var gotWorkflowAssignments []*sdm.WorkflowAssignment
	for listResp.Next() {
		n := listResp.Value()
		fmt.Printf("WorkflowAssignment Resource ID: %s\n", n.ResourceID)
		fmt.Printf("WorkflowAssignment Workflow ID: %s\n", n.WorkflowID)
		gotWorkflowAssignments = append(gotWorkflowAssignments, n)
	}
	if listResp.Err() != nil {
		log.Fatalf("Could not list workflow assignments: %v", err)
	}
	if len(gotWorkflowAssignments) != 1 {
		log.Fatalf("list failed: expected %d workflows, got %d", 1, len(gotWorkflowAssignments))
	}
	fmt.Println("Successfully list WorkflowAssignment.")
}
