package repository

import "errors"

var (
	NotFound  = errors.New("key not found")
	KeyExists = errors.New("key already exists")
)
