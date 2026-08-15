package repository

type Object interface {
	Get(key string) (*string, error)
	Put(key, value string) error
	Post(key, value string) error
	Delete(key, value string) error
}
