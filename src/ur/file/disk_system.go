package file

import (
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type diskSystem struct {
	basePath string
}

func NewDiskSystem(basePath ...string) System {

	base := ""

	if len(basePath) > 0 {
		base = basePath[0]
	}

	return &diskSystem{base}
}

func (d *diskSystem) path(path string) string {
	return filepath.Join(d.basePath, strings.TrimPrefix(path, d.basePath))
}

func (d *diskSystem) List(dirs ...string) (out []string, err error) {

	if len(dirs) == 0 {
		dirs = append(dirs, "")
	}

	for _, dir := range dirs {

		err = filepath.Walk(d.path(dir), func(path string, info os.FileInfo, err error) error {

			if !info.IsDir() {
				out = append(out, path)
			}

			return nil
		})

		if err != nil {
			return
		}
	}

	return
}

func (d *diskSystem) Glob(pattern string) (out []string, err error) {
	return filepath.Glob(d.path(pattern))
}

func (d *diskSystem) Read(path string) ([]byte, error) {
	return ioutil.ReadFile(d.path(path))
}

func (d *diskSystem) ReadCloser(path string) (io.ReadCloser, time.Time, error) {

	reader, err := os.Open(d.path(path))

	if err != nil {
		return nil, time.Time{}, err
	}

	stat, err := reader.Stat()

	if err != nil {
		return nil, time.Time{}, err
	}

	return reader, stat.ModTime(), nil
}

func (d *diskSystem) Stat(path string) (os.FileInfo, error) {
	return os.Stat(d.path(path))
}

func (d *diskSystem) Write(path string, reader io.Reader) error {

	path = d.path(path)

	if err := os.MkdirAll(filepath.Dir(path), 0744); err != nil {
		return nil
	}

	b, err := ioutil.ReadAll(reader)

	if err != nil {
		return err
	}

	return ioutil.WriteFile(path, b, 0644)
}

func (d *diskSystem) Chmod(path string, mode os.FileMode) error {
	return os.Chmod(path, mode)
}
