package worker

import (
	"Code-compilation-system/repository"
	"context"
)

type Object interface {
	GoWork(ctx context.Context, task *repository.Task, complete func(result string, err error))
}
