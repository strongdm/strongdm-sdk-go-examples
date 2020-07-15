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

	// Create a composite role
	composite := &sdm.Role{
		Name:      "example composite role",
		Composite: true,
	}

	roleResponse, err = client.Roles().Create(ctx, composite)
	if err != nil {
		log.Fatalf("Could not create role: %v", err)
	}

	compositeID := roleResponse.Role.ID
	log.Printf("Successfully created composite role.\n\tID: %v\n", compositeID)

	// Attach role to composite
	attachment := &sdm.RoleAttachment{
		CompositeRoleID: compositeID,
		AttachedRoleID:  roleID,
	}

	attachmentResponse, err := client.RoleAttachments().Create(ctx, attachment)
	if err != nil {
		log.Fatalf("Could not create attachment: %v", err)
	}

	attachmentID := attachmentResponse.RoleAttachment.ID
	log.Printf("Successfully created attached role to composite role.\n\tID: %v\n", attachmentID)
}
