package adapter

import (
	"context"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/awslabs/aws-lambda-go-api-proxy/core"
	"github.com/go-chi/chi/v5"
)

type ChiLambda struct {
	core.RequestAccessor
	chiMux *chi.Mux
}

func New(chi *chi.Mux) *ChiLambda {
	return &ChiLambda{chiMux: chi}
}

func (g *ChiLambda) Proxy(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	chiRequest, err := g.ProxyEventToHTTPRequest(req)
	return g.proxyInternal(chiRequest, err)
}

func (g *ChiLambda) ProxyWithContext(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	chiRequest, err := g.EventToRequestWithContext(ctx, req)
	return g.proxyInternal(chiRequest, err)
}

func (g *ChiLambda) proxyInternal(chiRequest *http.Request, err error) (events.APIGatewayProxyResponse, error) {
	if err != nil {
		return core.GatewayTimeout(), core.NewLoggedError("Could not convert proxy event to request: %v", err)
	}

	respWriter := core.NewProxyResponseWriter()
	g.chiMux.ServeHTTP(http.ResponseWriter(respWriter), chiRequest)

	proxyResponse, err := respWriter.GetProxyResponse()
	if err != nil {
		return core.GatewayTimeout(), core.NewLoggedError("Error while generating proxy response: %v", err)
	}

	return proxyResponse, nil
}