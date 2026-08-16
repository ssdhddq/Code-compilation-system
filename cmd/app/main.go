package main

import (
	"Code-compilation-system/api/http"
	"Code-compilation-system/repository/ram_storage"
	"log"

	httpLib "net/http"

	_ "Code-compilation-system/docs"

	"github.com/go-chi/chi/v5"
)

// @title           Swagger Example API
// @version         1.0
// @description     This is a sample server for a task service.
// @host            localhost:8080
// @BasePath        /
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
