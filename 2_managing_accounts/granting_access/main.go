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
		sdm.WithHost("api.strongdmdev.com:443"),
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

	// Create a datasource
	examplePostgresDatasource := &sdm.Postgres{
		Name:             "Example Postgres Datasource",
		Hostname:         "example.strongdm.com",
		Port:             5432,
		Username:         "example",
		Password:         "example",
		Database:         "example",
		OverrideDatabase: true,
	}

	resourceResponse, err := client.Resources().Create(ctx, examplePostgresDatasource)
	if err != nil {
		log.Fatalf("Could not create Postgres datasource: %v", err)
	}

	resourceID := resourceResponse.Resource.GetID()

	log.Printf("Successfully created Postgres datasource.\n\tID: %v\n", resourceID)

	// Create an account grant
	accountGrant := &sdm.AccountGrant{
		ResourceID: resourceID,
		AccountID:  accountID,
		StartFrom:  time.Now().Add(time.Hour).UTC(),
		ValidUntil: time.Now().Add(time.Hour * 2).UTC(),
	}

	grantResponse, err := client.AccountGrants().Create(ctx, accountGrant)
	if err != nil {
		log.Fatalf("Could not create account grant: %v", err)
	}

	grantID := grantResponse.AccountGrant.ID

	log.Printf("Successfully created account grant.\n\tID: %v\n", grantID)
}
