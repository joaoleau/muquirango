package interfaces

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
)

type Handler interface {
	Execute (ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error)
}
