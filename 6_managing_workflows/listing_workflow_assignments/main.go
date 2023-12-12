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

	// Create the resources - used for workflow assignments
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

	// Create a Workflow and assign the resources via a static access rule
	workflow := &sdm.Workflow{
		Name:        "Example Create Workflow",
		Description: "Example Workflow Description",
		AccessRules: sdm.AccessRules{
			sdm.AccessRule{
				Type: "mysql",
				Tags: tags,
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
