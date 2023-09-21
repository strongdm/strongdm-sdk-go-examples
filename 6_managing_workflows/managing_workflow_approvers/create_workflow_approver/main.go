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
		Name:        "Example Create Worfklow",
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

	// Create a approver - used for creating a workflow approver
	approverCreateResponse, err := client.Accounts().Create(ctx, &sdm.User{
		Email:     "create-workflow-approver-example@example.com",
		FirstName: "example",
		LastName:  "example",
	})
	if err != nil {
		log.Fatalf("Could not create approver: %v", err)
	}
	approverID := approverCreateResponse.Account.GetID()

	// Create a WorkflowApprover
	_, err = client.WorkflowApprovers().Create(ctx, &sdm.WorkflowApprover{
		WorkflowID: wf.ID,
		ApproverID: approverID,
	})
	if err != nil {
		log.Fatalf("Could not create workflow approver: %v", err)
	}

	fmt.Println("Successfully created WorkflowApprover.")
	fmt.Println("\tWorkflow ID:", wf.ID)
	fmt.Println("\tApprover ID:", approverID)
}
