package main

import (
	"Code-compilation-system/api/http"
	"Code-compilation-system/repository/ram_storage"
	"log"

	httpLib "net/http"

	"github.com/go-chi/chi/v5"
)

func main() {
	repo := ram_storage.NewObject()
	handler := http.NewObject(repo)

	r := chi.NewRouter()
	handler.WrapHandlers(r)
	err := httpLib.ListenAndServe("localhost:8080", r)
	if err != nil {
		log.Fatal(err)
	}
}
