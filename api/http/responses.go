package http

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
)

type PostResponse struct {
	Id uuid.UUID `json:"task_id"`
}

func createPostResponse(w http.ResponseWriter, id uuid.UUID) {
	resp := &PostResponse{Id: id}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, "Internal error", http.StatusInternalServerError)
	}
}
