package main

import (
	"Code-compilation-system/api/http"
	"Code-compilation-system/repository/ram_storage"
	"Code-compilation-system/session"
	"flag"
	"fmt"
	"log"

	httpLib "net/http"

	_ "Code-compilation-system/docs"

	"github.com/go-chi/chi/v5"
)

var globalSession *session.Manager

// @title           Swagger Example API
// @version         1.0
// @description     This is a sample server for a task service.
// @host            localhost:8080
// @BasePath        /
func main() {
	s := ram_storage.NewObjectSession()
	u := ram_storage.NewObjectUser()
	t := ram_storage.NewObjectTask()
	repo := ram_storage.NewObject(u, t, s)
	handler := http.NewObject(repo)

	host := flag.String("host", "0.0.0.0", "host addr")
	port := flag.Int("port", 8080, "port addr")


	r := chi.NewRouter()
	handler.WrapHandlers(r)
	err := httpLib.ListenAndServe(fmt.Sprintf("%s:%d", *host, *port), r)
	if err != nil {
		log.Fatal(err)
	}
}
