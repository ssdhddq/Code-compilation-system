package main

import (
	"Code-compilation-system/api/http"
	"Code-compilation-system/repository/ram_storage"
	"Code-compilation-system/session"
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

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
	session.RegisterProvider("ram_storage", ram_storage.Pvr)
	manager, err := session.NewManager("ram_storage", "sessionID", 86400)
	if err != nil {
		panic("manager not started")
	}

	go manager.GC()

	u := ram_storage.NewObjectUser()
	t := ram_storage.NewObjectTask()
	repo := ram_storage.NewObject(u, t)

	handler := http.NewObject(repo, manager)

	host := flag.String("host", "0.0.0.0", "host addr")
	port := flag.Int("port", 8080, "port addr")

	r := chi.NewRouter()
	handler.WrapHandlers(r)
	srv := &httpLib.Server{
		Addr:    fmt.Sprintf("%s:%d", *host, *port),
		Handler: r,
	}
	go func() {
		err = srv.ListenAndServe()
		if err != nil {
			log.Fatal(err)
		}
	}()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh

	manager.StopGC()

	log.Println("wait active tasks to finish...")
	if handler.ShutdownTasks(30 * time.Second) {
		log.Println("all tasks finished")
	} else {
		log.Println("timeout waiting for tasks")
	}

	shdCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(shdCtx); err != nil {
		log.Printf("server shutdown err: %v", err)
	} else {
		log.Println("server shutdown success")
	}

	log.Println("App exit")
	os.Exit(0)
}
