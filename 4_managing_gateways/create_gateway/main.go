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

	gateway := &sdm.Gateway{
		Name:          "example-gateway",
		ListenAddress: "gateway.example.com:5555",
	}

	gatewayResponse, err := client.Nodes().Create(ctx, gateway)
	if err != nil {
		log.Fatalf("Could not create gateway: %v", err)
	}

	gatewayID := gatewayResponse.Node.GetID()

	log.Printf("Successfully created gateway.\n\tID: %v\n", gatewayID)
}
