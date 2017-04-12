package database

type SessionFactory interface {
	NewSession() (Session, error)
	Close() error
}
