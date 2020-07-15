package main

import (
	"context"
	"log"
	"time"

	sdm "github.com/strongdm/strongdm-sdk-go"
)

const (
	eksClusterExampleName                 = "Example AWS EKS Cluster"
	eksClusterExampleCertificateAuthority = `-----BEGIN CERTIFICATE-----
MIICpjCCAY4CCQCYJT6s+JVzSTANBgkqhkiG9w0BAQsFADAVMRMwEQYDVQQDDApr
dWJlcm5ldGVzMB4XDTIwMDcxNTE0MjgzN1oXDTIxMDcxNTE0MjgzN1owFTETMBEG
A1UEAwwKa3ViZXJuZXRlczCCASIwDQYJKoZIhvcNAQEBBQADggEPADCCAQoCggEB
AJtI1pfbqy65FihJ6SadnrdDw6IjGJo7icoxcDR9Tn0Ljz7a7CO4VgpDfs/X4ljG
LkGTqDqLXZ61+lssfaUwMFA61McthTZfd6rYLBcxWFmaVqvUL0tguTrrUPuegHXv
IBs827JSH43BXqLgvZCaWYb5PtD+CI9F9bOBm+M+BUufrdS6gUkTqipZdgC8sl8E
SvixPjKPRu4EnBE/cPEMvYzkSpjixs87WKGPR0FM+6SQVr6o14Fs3QNlcElBAi27
U7XL+an/Fj0osEZGDhJ1u/TmmWlW7RopE1YS8gpVxBzQkBmeUU05a9l1f4L8j45E
TFuF5daWkNLZFO08u1GxnlsCAwEAATANBgkqhkiG9w0BAQsFAAOCAQEANyPDqSUZ
nLiOVGg4OWPmXJy3tk7+Mb6j/xOFFKoKfrXJVUB1F5IDMD673ozkhKyNcqfFOEeZ
+E3WC2/CxxwkJfEUrtij8qWMnafvDnaPan86jNkZsz9zvxphqdeA0hsYZF5tPLWT
Sk8uIHuRA36mYhzCrXQ7dhLn4mC147LRcZ73CTi4j4bNyGtCYgYE+Ta1pcrREIHp
PMiZH+tzwXAWeVKh3foHTjeXKAgXhg3Lbqxn6Uq3cejraUMi9b489KKPOlcaQ7wX
FPkubmy3vrhgJySlrfBDtCgFDwSosLniZU479S3oZBsKgPgLe3ELzAw1vLcuIgmd
JrXnKV7Z4r9uWg==
-----END CERTIFICATE-----
`
)

// CreateEKSClusterExample will create, find, and delete an EKS Cluster
// as an example of the StrongDM Go SDK.
func CreateEKSClusterExample(client *sdm.Client) {
	exampleEKSCluster := &sdm.AmazonEKS{
		Name:                 eksClusterExampleName,
		Endpoint:             "https://A1ADBDD0AE833267869C6ED0476D6B41.gr7.us-east-2.eks.amazonaws.com",
		AccessKey:            "AKIAIOSFODNN7EXAMPLE",
		SecretAccessKey:      "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY",
		CertificateAuthority: eksClusterExampleCertificateAuthority,
		Region:               "us-east-1",
		ClusterName:          "example",
		RoleArn:              "arn:aws:iam::000000000000:role/RoleName",
		HealthcheckNamespace: "default",
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	createResponse, err := client.Resources().Create(ctx, exampleEKSCluster)
	if err != nil {
		if _, ok := err.(*sdm.AlreadyExistsError); ok {
			log.Println("Resource already exists, continuing to allow for cleanup.")
			return
		}
		log.Fatalf("Could not create EKS Cluster: %v", err)

	}
	id := createResponse.Resource.GetID()
	name := createResponse.Resource.GetName()

	log.Printf("Successfully created EKS Cluster.\n\tName: %v\n\tID: %v\n", name, id)
}
