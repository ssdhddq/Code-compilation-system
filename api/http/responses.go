package http

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
)

type PostResponseTask struct {
	Id uuid.UUID `json:"task_id"`
}

func createPostResponseTask(w http.ResponseWriter, id uuid.UUID) {
	resp := &PostResponseTask{Id: id}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, "Internal error", http.StatusInternalServerError)
	}
}
