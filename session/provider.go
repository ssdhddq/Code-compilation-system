package session

type Provider interface {
	SessionInit(sid string) (Session, error)
	SessionUpdate(sid string) error
	SessionRead(sid string) (Session, error)
	SessionDestroy(sid string) error
	SessionGC(maxLifeTime int64)
}
