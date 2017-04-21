package file

import (
	"fmt"
	"os"
	"time"
)

type fileInfos map[string]os.FileInfo

type intervalWatcher struct {
	reader  SystemReader
	globs   []string
	timeout time.Duration
	cancel  <-chan struct{}
	infos   map[string]fileInfos
}

func NewIntervalWatcher(reader SystemReader, timeout time.Duration, cancel <-chan struct{}, globs ...string) Watcher {

	return &intervalWatcher{
		reader:  reader,
		timeout: timeout,
		cancel:  cancel,
		globs:   globs,
		infos:   make(map[string]fileInfos),
	}
}

func (w *intervalWatcher) scan() (map[string]fileInfos, error) {

	infos := make(map[string]fileInfos)

	for _, glob := range w.globs {

		files, err := w.reader.Glob(glob)

		if err != nil {
			return nil, fmt.Errorf("Glob on '%s' failed: %s", glob, err)
		}

		fim := make(fileInfos)

		for _, file := range files {

			finfo, err := w.reader.Stat(file)

			if err != nil {
				return nil, fmt.Errorf("Stat on '%s' failed: %s", file, err)
			}

			fim[file] = finfo
		}

		infos[glob] = fim
	}

	return infos, nil
}

func (w *intervalWatcher) diff(inputs map[string]fileInfos) ([]WatchDifferential, error) {

	diffs := []WatchDifferential{}

	for glob, input := range inputs {

		// Loop 0: deleted
		for name := range w.infos[glob] {
			if _, exists := input[name]; !exists {
				diffs = append(diffs, WatchDifferential{Path: name, Change: Deleted})
			}
		}

		// Loop 1: created or modified
		for name, infoNew := range input {

			infoExisting, exists := w.infos[glob][name]

			if !exists {
				diffs = append(diffs, WatchDifferential{Path: name, Change: Created})
				continue
			}

			if infoExisting.ModTime() != infoNew.ModTime() || infoExisting.Size() != infoNew.Size() {
				diffs = append(diffs, WatchDifferential{Path: name, Change: Modified})
				continue
			}
		}

	}

	// Save the new inputs
	w.infos = inputs

	return diffs, nil
}

func (w *intervalWatcher) Watch() (chan []WatchDifferential, error) {

	baseInfos, err := w.scan()

	if err != nil {
		return nil, err
	}

	w.infos = baseInfos

	c := make(chan []WatchDifferential, 1)

	go func() {

		ticker := time.NewTicker(w.timeout)
		defer ticker.Stop()

		for {
			select {
			case <-w.cancel:
				close(c)
				return

			case <-ticker.C:
				scan, err := w.scan()

				if err != nil {
					close(c)
					return
				}

				diffs, err := w.diff(scan)

				if err != nil {
					close(c)
					return
				}

				if len(diffs) > 0 {
					c <- diffs
				}
			}

		}
	}()

	return c, nil
}
