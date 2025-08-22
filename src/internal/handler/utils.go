package handler

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"github.com/aws/aws-lambda-go/events"
)

func Serialize[T any](obj T) ([]byte, error) {
	var body []byte
	body, err := json.Marshal(obj)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func Deserialize[T any](request events.APIGatewayProxyRequest) (*T, error) {
	body := request.Body

	if request.IsBase64Encoded {
		decoded, err := base64.StdEncoding.DecodeString(body)
		if err != nil {
			return nil, fmt.Errorf("failed to decode base64: %w", err)
		}
		body = string(decoded)
	}

	var t T
	if body == "" {
		return nil, errors.New("input is empty")
	}
	err := json.Unmarshal([]byte(body), &t)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func ResponseWithError(status int, err error) (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{
		StatusCode: status,
		Body:       err.Error(),
	}, nil
}

func ResponseWithMessage(status int, message string) (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{
		StatusCode: status,
		Headers: map[string]string{"Content-Type": "application/json"},
		Body:       message,
	}, nil
}

func ResponseWithSerialized(status int, data any) (events.APIGatewayProxyResponse, error) {
	body, err := Serialize(data)
	if err != nil {
		return ResponseWithError(http.StatusInternalServerError, err)
	}
	return ResponseWithMessage(status, string(body))
}