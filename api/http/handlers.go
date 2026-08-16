package http

import (
	"Code-compilation-system/repository"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"
)

type Object struct {
	repo repository.Object
}

func NewObject(object repository.Object) *Object {
	return &Object{
		repo: object,
	}
}

func (o *Object) GetStatusHandler(w http.ResponseWriter, r *http.Request) {
	req, err := parseGetRequest(r)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	task, err := o.repo.Get(req.id)
	errorHandler(w, err, task.Status)
}

func (o *Object) GetResultHandler(w http.ResponseWriter, r *http.Request) {
	req, err := parseGetRequest(r)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	task, err := o.repo.Get(req.id)
	errorHandler(w, err, task.Result)
}

func (o *Object) PostHandler(w http.ResponseWriter, r *http.Request) {
	req, err := parsePostRequest(r)
	if err != nil {
		http.Error(w, "Bad request parse", http.StatusBadRequest)
		return
	}

	id := uuid.New()

	task := &repository.Task{
		ID:        id,
		Result:    "testR",
		Status:    "testS",
		CreatedAT: time.Time{},
		UpdatedAT: time.Time{},
		Data:      req.Data,
	}

	err = o.repo.Create(task)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
	}

	createPostResponse(w, id)
}

func (o *Object) WrapHandlers(r chi.Router) {
	r.Use(middleware.Logger)
	r.Get("/status/{task_id}", o.GetStatusHandler)
	r.Get("/result/{task_id}", o.GetResultHandler)
	r.Post("/task", o.PostHandler)
}
