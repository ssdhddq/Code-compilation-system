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

func (o *Object) GetHandler(w http.ResponseWriter, r *http.Request) {
	req, err := parseGetRequest(r)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	value, err := o.repo.Get(req.key)
	errorHandler(w, err, getRequest{*value})
}
