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

type EntryHandler struct {
	Repository interfaces.EntryRepository
}

func NewEntryHandler(repo interfaces.EntryRepository) *EntryHandler {
	return &EntryHandler{Repository: repo}
}

func (e *EntryHandler) Execute(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	switch request.HTTPMethod {
	case http.MethodPost:
		return e.Add(ctx, request)
	case http.MethodGet:
		if request.PathParameters["id"] != "" {
			return e.GetByID(ctx, request)
		}
		return e.GetAll(ctx)
	case http.MethodPut:
		if request.PathParameters["id"] != "" {
			return e.UpdateByID(ctx, request)
		}
	case http.MethodDelete:
		if request.PathParameters["id"] != "" {
			return e.DeleteByID(ctx, request)
		}
	default:
		return ResponseWithError(http.StatusMethodNotAllowed, fmt.Errorf("method not allowed"))
	}

	return ResponseWithError(http.StatusNotFound, fmt.Errorf("not found :)"))
}

func (e *EntryHandler) Add(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	input, err := Deserialize[dto.CreateEntryInput](request)
	if err != nil {
		return ResponseWithError(http.StatusBadRequest, err)
	}

	entry := &model.Entry{
		ID:          uuid.NewString(),
		Type:        input.Type,
		Description: input.Description,
		CreatedAt:   time.Now().UTC(),
	}

	savedEntry, err := e.Repository.Add(ctx, entry)
	if err != nil {
		return ResponseWithError(http.StatusInternalServerError, err)
	}

	return ResponseWithSerialized(http.StatusCreated, savedEntry)
}

func (e *EntryHandler) GetAll(ctx context.Context) (events.APIGatewayProxyResponse, error) {
	entries, err := e.Repository.GetAll(ctx)
	if err != nil {
		return ResponseWithError(http.StatusInternalServerError, err)
	}

	return ResponseWithSerialized(http.StatusAccepted, entries)
}

func (e *EntryHandler) UpdateByID(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	id := request.PathParameters["id"]
	updateEntry, err := e.Repository.GetById(ctx, id)
	if err != nil {
		return ResponseWithError(http.StatusBadRequest, err)
	}

	input, err := Deserialize[dto.CreateEntryInput](request)
	if err != nil {
		return ResponseWithError(http.StatusBadRequest, err)
	}

	newEntry := &model.Entry{
		ID:          updateEntry.ID,
		Type:        input.Type,
		Description: input.Description,
		CreatedAt:   updateEntry.CreatedAt,
	}

	_, err = e.Repository.Update(ctx, newEntry)
	if err != nil {
		return ResponseWithError(http.StatusBadRequest, err)
	}

	return ResponseWithSerialized(http.StatusAccepted, newEntry)
}

func (e *EntryHandler) DeleteByID(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	id := request.PathParameters["id"]
	deleteEntry, err := e.Repository.GetById(ctx, id)
	if err != nil {
		return ResponseWithError(http.StatusBadRequest, err)
	}

	_, err = e.Repository.Delete(ctx, deleteEntry)
	if err != nil {
		return ResponseWithError(http.StatusBadRequest, err)
	}

	return ResponseWithSerialized(http.StatusAccepted, deleteEntry)
}

func (e *EntryHandler) GetByID(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	id := request.PathParameters["id"]
	entry, err := e.Repository.GetById(ctx, id)
	if err != nil {
		return ResponseWithError(http.StatusBadRequest, err)
	}

	_, err = e.Repository.GetById(ctx, id)
	if err != nil {
		return ResponseWithError(http.StatusBadRequest, err)
	}
	return ResponseWithSerialized(http.StatusAccepted, entry)
}
