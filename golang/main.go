package main

import (
	"context"
	"flag"
	"log"

	"cloud.google.com/go/storage"
	"golang.org/x/oauth2/google/externalaccount"

	//"github.com/aws/aws-sdk-go-v2/aws/session"

	"github.com/aws/aws-sdk-go-v2/aws"
	ac "github.com/aws/aws-sdk-go-v2/config"

	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

var (
	projectId = flag.String("projectId", "core-eso", "ProjectID")
	bucket    = flag.String("bucket", "core-eso-bucket", "GCS Bucket")
)

type awsSupplier struct {
	awsRegion             string
	awsCredentialProvider aws.CredentialsProvider
}

func NewAWSCredProvider(c aws.CredentialsProvider, r string) (awsSupplier, error) {
	return awsSupplier{
		awsCredentialProvider: c,
		awsRegion:             r,
	}, nil
}

func (supp awsSupplier) AwsRegion(ctx context.Context, options externalaccount.SupplierOptions) (string, error) {
	return supp.awsRegion, nil
}

func (supp awsSupplier) AwsSecurityCredentials(ctx context.Context, options externalaccount.SupplierOptions) (*externalaccount.AwsSecurityCredentials, error) {

	c, err := supp.awsCredentialProvider.Retrieve(ctx)
	if err != nil {
		return nil, err
	}
	return &externalaccount.AwsSecurityCredentials{
		AccessKeyID:     c.AccessKeyID,
		SecretAccessKey: c.SecretAccessKey,
		SessionToken:    c.SessionToken,
	}, nil
}

func main() {

	flag.Parse()

	log.Printf("======= Init  ========")

	ctx := context.Background()

	cfg, err := ac.LoadDefaultConfig(ctx, ac.WithRegion("us-east-1"))
	if err != nil {
		log.Fatal(err)
	}

	ac, err := NewAWSCredProvider(cfg.Credentials, "us-east-1")
	if err != nil {
		log.Fatal(err)
	}
	ts, err := externalaccount.NewTokenSource(ctx, externalaccount.Config{
		Audience:                       "//iam.googleapis.com/projects/995081019036/locations/global/workloadIdentityPools/aws-pool-1/providers/aws-provider-1",
		SubjectTokenType:               "urn:ietf:params:aws:token-type:aws4_request",
		Scopes:                         []string{"https://www.googleapis.com/auth/cloud-platform"},
		AwsSecurityCredentialsSupplier: ac,
	})

	if err != nil {
		log.Fatal(err)
	}

	// GCS does not support JWTAccessTokens, the following will only work if UseOauthToken is set to True
	storageClient, err := storage.NewClient(ctx, option.WithTokenSource(ts))
	if err != nil {
		log.Fatal(err)
	}
	it := storageClient.Bucket(*bucket).Objects(ctx, nil)
	for {
		attrs, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatalf("Bucket(%s).Objects: %v\n", *bucket, err)
		}
		log.Println(attrs.Name)
	}

}
