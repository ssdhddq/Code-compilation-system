package http

import (
	"Code-compilation-system/repository"
	"net/http"
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
