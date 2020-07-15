package main

import (
	"context"
	"log"
	"time"

	sdm "github.com/strongdm/strongdm-sdk-go"
)

const (
	rdpServerExampleName = "Example RDP Server"
)

// CreateRDPServerExample will create, find, and delete a RDP Server resource
// as an example of the StrongDM Go SDK.
func CreateRDPServerExample(client *sdm.Client) {
	exampleRDPServer := &sdm.RDP{
		Name:     "Example RDP Server",
		Hostname: "example.strongdm.com",
		Username: "example",
		Password: "example",
		Port:     3389,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	createResponse, err := client.Resources().Create(ctx, exampleRDPServer)
	if err != nil {
		if _, ok := err.(*sdm.AlreadyExistsError); ok {
			log.Println("Resource already exists, continuing to allow for cleanup.")
			return
		}
		log.Fatalf("Could not create RDP server: %v", err)

	}
	id := createResponse.Resource.GetID()
	name := createResponse.Resource.GetName()

	log.Printf("Successfully created RDP server.\n\tName: %v\n\tID: %v\n", name, id)
}
