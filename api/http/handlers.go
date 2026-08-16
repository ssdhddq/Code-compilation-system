package http

import (
	"Code-compilation-system/repository"
	"Code-compilation-system/worker"
	"Code-compilation-system/worker/simulate"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"
)

type Object struct {
	repo   repository.Object
	worker worker.Object
}

func NewObject(object repository.Object) *Object {
	return &Object{
		repo:   object,
		worker: simulate.NewObject(),
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
		Result:    "",
		Status:    "",
		CreatedAT: time.Time{},
		UpdatedAT: time.Time{},
		Data:      req.Data,
	}

	err = o.repo.Create(task)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
	}

	o.workHandler(id, w, task)

	createPostResponse(w, id)
}

func (o *Object) WrapHandlers(r chi.Router) {
	r.Use(middleware.Logger)
	r.Get("/status/{task_id}", o.GetStatusHandler)
	r.Get("/result/{task_id}", o.GetResultHandler)
	r.Post("/task", o.PostHandler)
}

func (o *Object) workHandler(id uuid.UUID, w http.ResponseWriter, task *repository.Task) {
	err := o.repo.UpdateStatus(id, "in process")
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
	}
	o.worker.GoWork(task)
	err = o.repo.UpdateResult(id, "data_data_data_moc")
	err = o.repo.UpdateStatus(id, "done")
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
	}
}
