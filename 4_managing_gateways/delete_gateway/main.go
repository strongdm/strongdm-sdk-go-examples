// Copyright 2020 StrongDM Inc
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	sdm "github.com/strongdm/strongdm-sdk-go/v2"
)

func main() {
	//	Load the SDM API keys from the environment.
	//	If these values are not set in your environment,
	//	please follow the documentation here:
	//	https://www.strongdm.com/docs/api/api-keys/
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

	// Create the Gateway
	create := &sdm.Gateway{
		Name:          "gateway-for-delete-example",
		ListenAddress: "gateway.example.com:5555",
	}

	createResponse, err := client.Nodes().Create(ctx, create)
	if err != nil {
		log.Fatalf("Could not create gateway: %v", err)
	}

	id := createResponse.Node.GetID()
	token := createResponse.Token
	fmt.Println("Successfully created gateway.")
	fmt.Println("\tID:", id)
	fmt.Println("\tToken:", token)

	// Delete the Gateway
	_, err = client.Nodes().Delete(ctx, id)
	if err != nil {
		log.Fatalf("Could not delete gateway: %v", err)
	}
	fmt.Println("Successfully deleted gateway.")
}
