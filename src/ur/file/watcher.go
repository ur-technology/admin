package file

//go:generate stringer -type=FileChange
type FileChange int

const (
	Created FileChange = iota
	Modified
	Deleted
)

type WatchDifferential struct {
	Path   string
	Change FileChange
}

type Watcher interface {
	Watch() (chan []WatchDifferential, error)
}
