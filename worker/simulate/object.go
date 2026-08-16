package simulate

import (
	"Code-compilation-system/repository"
	"fmt"
	"time"
)

type Object struct {
}

func NewObject() *Object {
	return &Object{}
}

func (o *Object) GoWork(task *repository.Task) {
	fmt.Print("Начало обработки таски\n")
	time.Sleep(time.Second * 5)
	fmt.Print("Конец обработки таски\n")
}
