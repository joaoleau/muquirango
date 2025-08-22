package handler

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/google/uuid"
	"github.com/joaoleau/muquirango/dto"
	"github.com/joaoleau/muquirango/interfaces"
	"github.com/joaoleau/muquirango/model"
)

type UserHandler struct {
	Repository interfaces.UserRepository
}

func NewUserHandler(repo interfaces.UserRepository) *UserHandler {
	return &UserHandler{Repository: repo}
}

func (e *UserHandler) Execute(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return ResponseWithError(http.StatusNotFound, fmt.Errorf("not found :)"))
}

func (e *UserHandler) Add(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	input, err := Deserialize[dto.CreateUserInput](request)
	if err != nil {
		return ResponseWithError(http.StatusBadRequest, err)
	}

	User := &model.User{
		ID:          uuid.NewString(),
		Type:        input.Type,
		Description: input.Description,
		CreatedAt:   time.Now().UTC(),
	}

	savedUser, err := e.Repository.Add(ctx, User)
	if err != nil {
		return ResponseWithError(http.StatusInternalServerError, err)
	}

	return ResponseWithSerialized(http.StatusCreated, savedUser)
}

func (e *UserHandler) GetAll(ctx context.Context) (events.APIGatewayProxyResponse, error) {
	entries, err := e.Repository.GetAll(ctx)
	if err != nil {
		return ResponseWithError(http.StatusInternalServerError, err)
	}

	return ResponseWithSerialized(http.StatusAccepted, entries)
}

func (e *UserHandler) UpdateByID(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	id := request.PathParameters["id"]
	updateUser, err := e.Repository.GetById(ctx, id)
	if err != nil {
		return ResponseWithError(http.StatusBadRequest, err)
	}

	input, err := Deserialize[dto.CreateUserInput](request)
	if err != nil {
		return ResponseWithError(http.StatusBadRequest, err)
	}

	newUser := &model.User{
		ID:          updateUser.ID,
		Type:        input.Type,
		Description: input.Description,
		CreatedAt:   updateUser.CreatedAt,
	}

	_, err = e.Repository.Update(ctx, newUser)
	if err != nil {
		return ResponseWithError(http.StatusBadRequest, err)
	}

	return ResponseWithSerialized(http.StatusAccepted, newUser)
}

func (e *UserHandler) DeleteByID(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	id := request.PathParameters["id"]
	deleteUser, err := e.Repository.GetById(ctx, id)
	if err != nil {
		return ResponseWithError(http.StatusBadRequest, err)
	}

	_, err = e.Repository.Delete(ctx, deleteUser)
	if err != nil {
		return ResponseWithError(http.StatusBadRequest, err)
	}

	return ResponseWithSerialized(http.StatusAccepted, deleteUser)
}

func (e *UserHandler) GetByID(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	id := request.PathParameters["id"]
	User, err := e.Repository.GetById(ctx, id)
	if err != nil {
		return ResponseWithError(http.StatusBadRequest, err)
	}

	_, err = e.Repository.GetById(ctx, id)
	if err != nil {
		return ResponseWithError(http.StatusBadRequest, err)
	}
	return ResponseWithSerialized(http.StatusAccepted, User)
}
