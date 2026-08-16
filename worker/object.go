package worker

import "Code-compilation-system/repository"

type Object interface {
	GoWork(task *repository.Task)
}
