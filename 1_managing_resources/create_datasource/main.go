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

	examplePostgresDatasource := &sdm.Postgres{
		Name:             "Example Postgres Datasource",
		Hostname:         "example.strongdm.com",
		Port:             5432,
		Username:         "example",
		Password:         "example",
		Database:         "example",
		OverrideDatabase: true,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	createResponse, err := client.Resources().Create(ctx, examplePostgresDatasource)
	if err != nil {
		log.Fatalf("Could not create Postgres datasource: %v", err)
	}

	id := createResponse.Resource.GetID()
	name := createResponse.Resource.GetName()

	log.Printf("Successfully created Postgres datasource.\n\tName: %v\n\tID: %v\n", name, id)
}
