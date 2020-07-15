package main

import (
	"context"
	"fmt"
	"log"
	"time"

	sdm "github.com/strongdm/strongdm-sdk-go"
)

// DeleteExampleResources will attempt to find and delete the resources
// which are created within the examples.
func DeleteExampleResources(client *sdm.Client) {
	resourceNames := []string{
		datasourceExampleName,
		eksClusterExampleName,
		rdpServerExampleName,
		sshServerExampleName,
	}

	for _, name := range resourceNames {
		if resource := getResourceByName(client, name); resource != nil {
			deleteResource(client, resource)
		}
	}
}

func getResourceByName(client *sdm.Client, name string) sdm.Resource {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	filter := fmt.Sprintf("name:%q", name)

	resources, err := client.Resources().List(ctx, filter)
	if err != nil {
		log.Fatalf("could not list resources: %v", err)
	}

	if !resources.Next() {
		log.Printf("could not find resource with name %q\n", name)
		return nil
	}

	resource := resources.Value()
	log.Printf("Successfully read resource by name.\n\tName: %v\n\tID: %v\n", resource.GetName(), resource.GetID())
	return resource
}

func deleteResource(client *sdm.Client, resource sdm.Resource) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	_, err := client.Resources().Delete(ctx, resource.GetID())
	if err != nil {
		log.Fatalf("Could not delete resource: %v", err)
	}

	log.Printf("Successfully deleted resource.\n\tName: %v\n\tID: %v\n", resource.GetName(), resource.GetID())

}
