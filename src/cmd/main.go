package main

import (
	"context"


	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	
	adapter "github.com/awslabs/aws-lambda-go-api-proxy/chi"
	"github.com/go-chi/chi/v5"

	"github.com/joaoleau/muquirango/internal/config"
	"github.com/joaoleau/muquirango/internal/router"
	"github.com/joaoleau/muquirango/internal/repository"
)

var chiLambda *adapter.ChiLambda

func main() {
    db, table := config.DynamoClient(context.Background())
    chiRouter := chi.NewRouter()
    
	repositories := &router.Repositories{
		TransactionRepo: repository.NewTransactionRepository(db, table),
	}

	router.RegisterRoutes(chiRouter, repositories)
    chiLambda = adapter.New(chiRouter)
    lambda.Start(handler)
}

func handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
    return chiLambda.ProxyWithContext(ctx, request)
}
