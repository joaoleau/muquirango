package config

import (
	"context"
	"log"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

var client *dynamodb.Client
var tableName string = "GODynamo" 

// func init() {
// 	cfg, err := config.LoadDefaultConfig(context.Background())
// 	if err != nil {
// 		log.Fatalf("failed to load database: %v", err)
// 	}
// 	client = dynamodb.NewFromConfig(cfg)
// }

func DynamoClient(ctx context.Context) (*dynamodb.Client, string) {
	customResolver := aws.EndpointResolverWithOptionsFunc(
		func(service, region string, options ...interface{}) (aws.Endpoint, error) {
			return aws.Endpoint{
				URL:           "http://host.docker.internal:8000",
				SigningRegion: "sa-east-1",
			}, nil
		})

	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion("sa-east-1"),
		config.WithEndpointResolverWithOptions(customResolver),
	)
	if err != nil {
		log.Fatalf("failed to load database: %v", err)
	}

	client = dynamodb.NewFromConfig(cfg)
	return client, tableName
}
