package main

import (
	"context"
	"log"
	"time"

	sdm "github.com/strongdm/strongdm-sdk-go"
)

const (
	datasourceExampleName = "Example Postgres Datasource"
)

// CreateDatasourceExample will create, find, and delete a Postgres datasource
// as an example of the StrongDM Go SDK.
func CreateDatasourceExample(client *sdm.Client) {
	examplePostgresDatasource := &sdm.Postgres{
		Name:             datasourceExampleName,
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
		if _, ok := err.(*sdm.AlreadyExistsError); ok {
			log.Println("Resource already exists, continuing to allow for cleanup.")
			return
		}

		log.Fatalf("Could not create Postgres datasource: %v", err)
	}

	id := createResponse.Resource.GetID()
	name := createResponse.Resource.GetName()

	log.Printf("Successfully created Postgres datasource.\n\tName: %v\n\tID: %v\n", name, id)
}
