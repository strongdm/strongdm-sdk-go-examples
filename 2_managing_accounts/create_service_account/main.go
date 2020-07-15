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

	client, err := sdm.New(
		accessKey,
		secretKey,
		sdm.WithHost("api.strongdmdev.com:443"),
	)
	if err != nil {
		log.Fatalf("could not create client: %v", err)
	}

	service := &sdm.Service{
		Name: "Example Service Account",
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	createResponse, err := client.Accounts().Create(ctx, service)
	if err != nil {
		log.Fatalf("Could not create service account: %v", err)
	}

	id := createResponse.Account.GetID()

	log.Printf("Successfully created service account.\n\tID: %v\n", id)
}
