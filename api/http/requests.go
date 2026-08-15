package http

import (
	"Code-compilation-system/repository"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type getRequest struct {
	key string
}

func parseGetRequest(r *http.Request) (*getRequest, error) {
	key := r.URL.Query().Get("task_id")
	if key == "" {
		return nil, fmt.Errorf("missing key")
	}
	return &getRequest{key: key}, nil
}

func errorHandler(w http.ResponseWriter, err error, resp any) {
	if errors.Is(err, repository.NotFound) {
		http.Error(w, "Key not found", http.StatusNotFound)
		return
	} else if errors.Is(err, repository.KeyExists) {
		http.Error(w, "Key already exists", http.StatusConflict)
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
