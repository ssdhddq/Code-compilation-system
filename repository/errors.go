package repository

import "errors"

var (
	NotFound   = errors.New("key not found")
	KeyExists  = errors.New("key already exists")
	NilTask    = errors.New("task is nil")
	NilSession = errors.New("session is nil")
	NilUser    = errors.New("user is nil")
)
