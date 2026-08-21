package ram_storage

type Object struct {
	*ObjectUser
	*ObjectTask
}

func NewObject(u *ObjectUser, t *ObjectTask) *Object {
	return &Object{
		ObjectUser: u,
		ObjectTask: t,
	}
}
