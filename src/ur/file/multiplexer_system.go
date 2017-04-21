package file

import (
	"fmt"
	"io"
	"os"
	"time"
)

type multiplexerSystem struct {
	systems []SystemReader
}

func NewMultiplexerSystem(systems ...SystemReader) SystemReader {
	return &multiplexerSystem{
		systems: systems,
	}
}

func (m *multiplexerSystem) List(dirs ...string) (out []string, err error) {

	seens := map[string]bool{}

	for _, system := range m.systems {

		if files, err := system.List(dirs...); err == nil {

			for _, file := range files {
				if _, seen := seens[file]; !seen {
					out = append(out, file)
					seens[file] = true
				}
			}
		}
	}

	return
}

func (m *multiplexerSystem) Glob(pattern string) (out []string, err error) {

	seens := map[string]bool{}

	for _, system := range m.systems {

		if files, err := system.Glob(pattern); err == nil {

			for _, file := range files {
				if _, seen := seens[file]; !seen {
					out = append(out, file)
					seens[file] = true
				}
			}
		}
	}

	return
}

func (m *multiplexerSystem) Read(path string) ([]byte, error) {

	for _, system := range m.systems {

		if content, err := system.Read(path); err == nil {
			return content, nil
		}
	}

	return nil, fmt.Errorf("Can't read file from any system")
}

func (m *multiplexerSystem) ReadCloser(path string) (io.ReadCloser, time.Time, error) {

	for _, system := range m.systems {

		if file, tm, err := system.ReadCloser(path); err == nil {
			return file, tm, nil
		}
	}

	return nil, time.Time{}, fmt.Errorf("Can't get reader from any system")
}

func (m *multiplexerSystem) Stat(path string) (os.FileInfo, error) {

	for _, system := range m.systems {

		if info, err := system.Stat(path); err == nil {
			return info, err
		}
	}

	return nil, fmt.Errorf("Can't get stat from any system")
}
