package ram_storage

import "Code-compilation-system/repository"

type Object struct {
	data map[string]string
}

func NewObject() *Object {
	return &Object{
		data: make(map[string]string),
	}
}

func (o *Object) Get(key string) (*string, error) {
	value, exists := o.data[key]
	if !exists {
		return nil, repository.NotFound
	}
	return &value, nil
}

func (o *Object) Put(key, value string) error {
	o.data[key] = value
	return nil
}

func (o *Object) Post(key, value string) error {
	if _, exists := o.data[key]; exists {
		return repository.KeyExists
	}
	o.data[key] = value
	return nil
}

func (o *Object) Delete(key, value string) error {
	if _, exists := o.data[key]; !exists {
		return repository.NotFound
	}
	delete(o.data, key)
	return nil
}
