package http

import (
	"Code-compilation-system/repository"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type getRequest struct {
	id uuid.UUID
}

func parseGetRequest(r *http.Request) (*getRequest, error) {
	idStr := chi.URLParam(r, "task_id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return nil, fmt.Errorf("получение UUID %v", err)
	}
	return &getRequest{id: id}, nil
}

type postRequest struct {
	Data string `json:"data"`
}

func parsePostRequest(r *http.Request) (*postRequest, error) {
	var req postRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, fmt.Errorf("получение тела запроса: %v", err)
	}
	if req.Data == "" {
		return nil, fmt.Errorf("тело запроса пустое")
	}
	return &req, nil
}

func errorHandler(w http.ResponseWriter, err error, resp any) {
	if errors.Is(err, repository.NotFound) {
		http.Error(w, "Id not found", http.StatusNotFound)
		return
	} else if errors.Is(err, repository.KeyExists) {
		http.Error(w, "Id already exists", http.StatusConflict)
		return
	} else if err != nil {
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}

	if resp != nil {
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			http.Error(w, "InternalError", http.StatusInternalServerError)
		}
	}
}
