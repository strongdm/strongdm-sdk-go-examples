package main

import (
	"context"
	"log"
	"os"
	"time"

	sdm "github.com/strongdm/strongdm-sdk-go"
)

func main() {
	//	Load the SDM API keys from the environment.
	//	If these values are not set in your environment,
	//	please follow the documentation here:
	//	https://www.strongdm.com/docs/admin-guide/api-credentials/
	accessKey := os.Getenv("SDM_API_ACCESS_KEY")
	secretKey := os.Getenv("SDM_API_SECRET_KEY")
	if accessKey == "" || secretKey == "" {
		log.Fatal("SDM_API_ACCESS_KEY and SDM_API_SECRET_KEY must be provided")
	}

	// Create the client
	client, err := sdm.New(
		accessKey,
		secretKey,
	)
	if err != nil {
		log.Fatalf("could not create client: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Create a user
	user := &sdm.User{
		Email:     "example@example.com",
		FirstName: "example",
		LastName:  "example",
	}

	accountResponse, err := client.Accounts().Create(ctx, user)
	if err != nil {
		log.Fatalf("Could not create user: %v", err)
	}

	accountID := accountResponse.Account.GetID()
	log.Printf("Successfully created user.\n\tID: %v\n", accountID)

	// Create a role
	role := &sdm.Role{
		Name: "example role",
	}

	roleResponse, err := client.Roles().Create(ctx, role)
	if err != nil {
		log.Fatalf("Could not create role: %v", err)
	}

	roleID := roleResponse.Role.ID
	log.Printf("Successfully created role.\n\tID: %v\n", roleID)

	// Assign account to role
	attachment := &sdm.AccountAttachment{
		AccountID: accountID,
		RoleID:    roleID,
	}

	attachmentResponse, err := client.AccountAttachments().Create(ctx, attachment)
	if err != nil {
		log.Fatalf("Could not create account attachment: %v", err)
	}

	attachmentID := attachmentResponse.AccountAttachment.ID
	log.Printf("Successfully created account attachment.\n\tID: %v\n", attachmentID)
}
