package main

import (
	"context"
	"log"
	"time"

	sdm "github.com/strongdm/strongdm-sdk-go"
)

const (
	sshServerExampleName      = "Example SSH Server"
	sshServerExamplePublicKey = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQClgTkc" +
		"smqpGTlhFCSyT6xTUOSyAo4a66niRZXf/AjB3Cc6H/BE+jSQUjtEJySO5Ak/kjL37ojI" +
		"mNWZICy3tPWLJsWKb6mzWJmcIZulOoXX2wLnGaVYwNvoo5AKRc9phGwGvuMmKsS9D9Zo" +
		"X4LRnvw5ONAMATPu/mJ+nGJ03mEHwraYMExaBC6+MkKukZbgumFjAtW7V7zFE6pxSGa2" +
		"BEG0fXDSED+ZcvxcqIyB+HKdYXyA91HhvRF0jGxwrDDcbHgVek9JJyYvNAdUpCwuU67j" +
		"yhtRdnM13bPGt0zpd8tgNBmr+/Vvx95/ZFB6+qj0hNEygslHebm2S3jXdfrPH8KF+XxB" +
		"LcOyFop2bVg6SRIA503D175fEmrV/GdoR3uMhMAh/prhtH5Q1+0OCkbRHAaAdy3kBONV" +
		"3i3B0ZRWhsH0VbaGYjVNnQJLPkwqsTEWNVrQOq2796M9ko2UhpFCHd6SX1mIQ75lL6kj" +
		"xaH0iKA7EOaE1aoxFZLNH1MonYgHrHs= example@strongdm.com"
)

// CreateSSHServerExample will create, find, and delete an SSH Server
// as an example of the StrongDM Go SDK.
func CreateSSHServerExample(client *sdm.Client) {
	exampleSSHServer := &sdm.SSH{
		Name:      sshServerExampleName,
		Hostname:  "203.0.113.23",
		Username:  "example",
		Port:      22,
		PublicKey: sshServerExamplePublicKey,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if createResponse, err := client.Resources().Create(ctx, exampleSSHServer); err != nil {
		if _, ok := err.(*sdm.AlreadyExistsError); ok {
			log.Println("Resource already exists, continuing to allow for cleanup.")
		} else {
			log.Fatalf("Could not create SSH server: %v", err)
		}
	} else {
		id := createResponse.Resource.GetID()
		name := createResponse.Resource.GetName()

		log.Printf("Successfully created SSH server.\n\tName: %v\n\tID: %v\n", name, id)
	}

	if cleanupResources {
		resource := getResourceByName(client, sshServerExampleName)
		deleteResource(client, resource)
	}
}
