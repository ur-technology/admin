package file

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	stdlibpath "path"
	"strings"
	"time"

	"github.com/ryanuber/go-glob"
)

type MockSystemFile struct {
	Path         string
	Contents     []byte
	LastModified time.Time
}

type mockSystem struct {
	files []MockSystemFile
}

func NewMockSystem(files ...MockSystemFile) System {

	s := &mockSystem{
		files: files,
	}

	return s
}

func (m *mockSystem) List(dirs ...string) (out []string, err error) {

	if len(dirs) == 0 {
		dirs = append(dirs, "")
	}

	for _, dir := range dirs {
		for _, f := range m.files {
			if strings.HasPrefix(f.Path, dir) {
				out = append(out, f.Path)
			}
		}
	}

	return
}

func (m *mockSystem) Glob(pattern string) (out []string, err error) {

	for _, f := range m.files {

		if glob.Glob(pattern, f.Path) {
			out = append(out, f.Path)
		}
	}

	return
}

func (m *mockSystem) Read(path string) ([]byte, error) {

	for _, f := range m.files {

		if f.Path == path {
			return f.Contents, nil
		}
	}

	return []byte{}, fmt.Errorf("Path not found")
}

func (m *mockSystem) ReadCloser(path string) (io.ReadCloser, time.Time, error) {

	for _, f := range m.files {

		if f.Path == path {
			return ioutil.NopCloser(bytes.NewBuffer(f.Contents)), f.LastModified, nil
		}
	}

	return nil, time.Time{}, fmt.Errorf("Path not found")
}

func (m *mockSystem) Stat(path string) (os.FileInfo, error) {

	for _, f := range m.files {

		if f.Path == path {
			return &info{
				name:    stdlibpath.Base(path),
				size:    int64(len(f.Contents)),
				mode:    0755,
				modTime: f.LastModified,
			}, nil
		}
	}

	return nil, fmt.Errorf("Path not found")
}

func (m *mockSystem) Write(path string, reader io.Reader) error {

	buf := bytes.Buffer{}

	if _, err := io.Copy(&buf, reader); err != nil {
		return err
	}

	newFile := MockSystemFile{
		Path:         path,
		Contents:     buf.Bytes(),
		LastModified: time.Now(),
	}

	for index, f := range m.files {

		if f.Path == path {
			m.files[index] = newFile
			return nil
		}
	}

	m.files = append(m.files, newFile)
	return nil
}

func (m *mockSystem) Chmod(path string, mode os.FileMode) error {
	return nil
}
