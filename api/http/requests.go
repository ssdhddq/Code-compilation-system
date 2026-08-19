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

type PostRequestTask struct {
	Data string `json:"data,omitempty"`
}

func parsePostRequestTask(r *http.Request) (*PostRequestTask, error) {
	var req PostRequestTask
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, fmt.Errorf("получение тела запроса: %v", err)
	}
	if req.Data == "" {
		return nil, fmt.Errorf("тело запроса пустое")
	}
	return &req, nil
}

type PostRequestRegisterAndAuth struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func parsePostRequestReqisterAndAuth(r *http.Request) (*PostRequestRegisterAndAuth, error) {
	var req PostRequestRegisterAndAuth
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, fmt.Errorf("получение тела запроса: %v", err)
	}
	if req.Username == "" {
		return nil, fmt.Errorf("username не может быть пустым")
	}
	if req.Username == "" {
		return nil, fmt.Errorf("password не может быть пустым")
	}
	return &req, nil
}

func errorHandler(w http.ResponseWriter, err error) {
	if errors.Is(err, repository.NotFound) {
		http.Error(w, "Id not found", http.StatusNotFound)
		return
	} else if errors.Is(err, repository.KeyExists) {
		http.Error(w, "Id already exists", http.StatusConflict)
		return
	} else if errors.Is(err, repository.NilTask) {
		http.Error(w, "Nil task", http.StatusConflict)
		return
	} else if errors.Is(err, repository.NilSession) {
		http.Error(w, "Nil session", http.StatusConflict)
		return
	} else if errors.Is(err, repository.NilUser) {
		http.Error(w, "Nil user", http.StatusConflict)
		return
	} else if err != nil {
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}
}
