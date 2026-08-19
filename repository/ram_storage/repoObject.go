package ram_storage

type Object struct {
	*ObjectUser
	*ObjectTask
	*ObjectSession
}

func NewObject(u *ObjectUser, t *ObjectTask, s *ObjectSession) *Object {
	return &Object{
		ObjectUser:    u,
		ObjectTask:    t,
		ObjectSession: s,
	}
}
