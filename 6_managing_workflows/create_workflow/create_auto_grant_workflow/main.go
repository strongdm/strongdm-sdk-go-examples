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

	// Create an auto grant Workflow with initial Access Rule
	workflow := &sdm.Workflow{
		Name:        "Example Create Auto Grant Worfklow",
		Description: "Example Workflow Description",
		AutoGrant:   true,
		Enabled:     true,
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

	workflowID := createResponse.Workflow.ID
	workflowName := createResponse.Workflow.Name

	fmt.Println("Successfully created Workflow.")
	fmt.Println("\tID:", workflowID)
	fmt.Println("\tName:", workflowName)
}
