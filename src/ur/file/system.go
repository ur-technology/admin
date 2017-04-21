package file

import (
	"fmt"
	"io"
	"os"
	"time"
)

type SystemReader interface {
	List(dirs ...string) ([]string, error)
	Glob(pattern string) ([]string, error)
	Read(path string) ([]byte, error)
	ReadCloser(path string) (io.ReadCloser, time.Time, error)
	Stat(path string) (os.FileInfo, error)
}

type SystemWriter interface {
	Write(path string, reader io.Reader) error
	Chmod(path string, mode os.FileMode) error
}

type System interface {
	SystemReader
	SystemWriter
}

func Copy(src SystemReader, dest SystemWriter) error {

	files, err := src.List()

	if err != nil {
		return fmt.Errorf("List failed: %s", err)
	}

	for _, file := range files {

		rc, _, err := src.ReadCloser(file)

		if err != nil {
			return fmt.Errorf("Read '%s': %s", file, err)
		}

		if err := dest.Write(file, rc); err != nil {
			return fmt.Errorf("Write '%s': %s", file, err)
		}
	}

	return nil
}

type info struct {
	name    string      // base name of the file
	size    int64       // length in bytes for regular files; system-dependent for others
	mode    os.FileMode // file mode bits
	modTime time.Time   // modification time
	isDir   bool        // abbreviation for Mode().IsDir()
	sys     interface{} // underlying data source (can return nil)
}

func newInfo() os.FileInfo {
	return &info{}
}

func (i *info) Name() string {
	return i.name
}

func (i *info) Size() int64 {
	return i.size
}

func (i *info) Mode() os.FileMode {
	return i.mode
}

func (i *info) ModTime() time.Time {
	return i.modTime
}

func (i *info) IsDir() bool {
	return i.isDir
}

func (i *info) Sys() interface{} {
	return i.sys
}
