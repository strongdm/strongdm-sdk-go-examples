package main

import (
	"context"
	"log"
	"time"

	sdm "github.com/strongdm/strongdm-sdk-go"
)

const (
	postgresName             = "Example Postgres Instance"
	postgresHostname         = "example.strongdm.com"
	postgresPort             = 5432
	postgresUsername         = "example"
	postgresPassword         = "example"
	postgresDatabase         = "example"
	postgresOverrideDatabase = true
	postgresCleanupResource  = true
)

// CreateDatasourceExample will create, find, and delete a Postgres datasource
// as an example of the StrongDM Go SDK.
func CreateDatasourceExample(client *sdm.Client) {

	createDatasource(client)

	if postgresCleanupResource {
		resource := getResourceByName(client, postgresName)
		deleteResource(client, resource)
	}
}

func createDatasource(client *sdm.Client) {
	examplePostgresDatasource := &sdm.Postgres{
		Name:             postgresName,
		Hostname:         postgresHostname,
		Port:             postgresPort,
		Username:         postgresUsername,
		Password:         postgresPassword,
		Database:         postgresDatabase,
		OverrideDatabase: postgresOverrideDatabase,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	createResponse, err := client.Resources().Create(ctx, examplePostgresDatasource)
	if err != nil {
		if _, ok := err.(*sdm.AlreadyExistsError); ok {
			log.Println("Resource already exists, continuing to allow for cleanup.")
			return
		}
		log.Fatalf("Could not create Postgres datasource: %v", err)
	}

	id := createResponse.Resource.GetID()
	name := createResponse.Resource.GetName()

	log.Printf("Successfully created Postgres resource.\n\tName: %v\n\tID: %v\n", name, id)
}
