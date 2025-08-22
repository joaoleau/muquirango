package main

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/go-chi/chi/v5"

	"github.com/joaoleau/muquirango/internal/config"
	"github.com/joaoleau/muquirango/internal/router"
	"github.com/joaoleau/muquirango/internal/adapter"
)

var chiLambda *adapter.ChiLambda

func main() {
    db, table := config.DynamoClient(context.Background())
    r := chi.NewRouter()
    router.RegisterRoutes(r, db, table)
    chiLambda = adapter.New(r)
    lambda.Start(handler)
}

func handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
    return chiLambda.ProxyWithContext(ctx, request)
}
