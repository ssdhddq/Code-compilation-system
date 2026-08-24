package simulate

import (
	"Code-compilation-system/repository"
	"context"
	"log"
	"time"
)

type Object struct {
}

func NewObject() *Object {
	return &Object{}
}

func (w *Object) GoWork(ctx context.Context, task *repository.Task, onComplete func(result string, err error)) {
	go func() {
		log.Printf("Starting work %s", task.ID)
		select {
		case <-time.After(5 * time.Second):
			log.Printf("Finished work %s", task.ID)
			onComplete("datamoc", nil) // результат
		case <-ctx.Done():
			log.Printf("Task %s cancelled", task.ID)
			onComplete("", ctx.Err())
		}
	}()
}
