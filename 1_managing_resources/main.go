package main

import (
	"log"
	"os"

	sdm "github.com/strongdm/strongdm-sdk-go"
)

const (
	// Set this to false if you want to leave the resources up for testing or
	// to view in the admin UI. The example resources will not be functional.
	cleanupResources = true
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

	//	Run the create_datasource example.
	CreateDatasourceExample(client)

	//	Run the create_eks_cluster example.
	CreateEKSClusterExample(client)

	//	Run the create_rdp_server example.
	//CreateRDPServer(client)

	//	Run the create_ssh_server example.
	//CreateSSHServer(client)
}
