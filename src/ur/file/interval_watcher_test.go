package file

import (
	"bytes"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func mockFileReaderForWatch(t *testing.T) System {
	return NewMockSystem(
		MockSystemFile{
			Path:         "a/1.txt",
			Contents:     []byte("hello"),
			LastModified: time.Now(),
		},
		MockSystemFile{
			Path:         "a/2.txt",
			Contents:     []byte("hello"),
			LastModified: time.Now(),
		},
		MockSystemFile{
			Path:         "b/1.txt",
			Contents:     []byte("hello"),
			LastModified: time.Now(),
		},
		MockSystemFile{
			Path:         "b/2.txt",
			Contents:     []byte("hello"),
			LastModified: time.Now(),
		},
	)
}

func assertExpectedFileInfos(t *testing.T, input map[string]fileInfos, fis map[string][]string) {

	for key, set := range fis {
		for _, file := range set {
			require.Contains(t, input[key], file)
		}
	}
}

func Test_ATimeoutWatcherCanScanForAListOfFiles(t *testing.T) {

	w := NewIntervalWatcher(
		mockFileReaderForWatch(t),
		1*time.Millisecond,
		nil,
		"a/*.txt",
		"b/*.txt",
	).(*intervalWatcher)

	output, err := w.scan()
	require.Nil(t, err)

	assertExpectedFileInfos(
		t,
		output,
		map[string][]string{
			"a/*.txt": []string{"a/1.txt", "a/2.txt"},
			"b/*.txt": []string{"b/1.txt", "b/2.txt"},
		},
	)
}

func Test_ATimeoutWatcherCanDiffFileInfoSets(t *testing.T) {

	fs := mockFileReaderForWatch(t)

	w := NewIntervalWatcher(
		fs,
		1*time.Millisecond,
		nil,
		"a/*.txt",
		"b/*.txt",
	).(*intervalWatcher)

	scan, err := w.scan()
	require.Nil(t, err)
	w.infos = scan

	// Test create
	require.Nil(t, fs.Write("a/3.txt", bytes.NewBufferString("Hello!")))

	createScan, err := w.scan()
	require.Nil(t, err)

	diffs, err := w.diff(createScan)
	require.Equal(t,
		[]WatchDifferential{
			WatchDifferential{Path: "a/3.txt", Change: Created},
		},
		diffs,
	)

	// Test modify
	require.Nil(t, fs.Write("a/3.txt", bytes.NewBufferString("Hello!")))

	modifyScan, err := w.scan()
	require.Nil(t, err)

	diffs, err = w.diff(modifyScan)
	require.Equal(t,
		[]WatchDifferential{
			WatchDifferential{Path: "a/3.txt", Change: Modified},
		},
		diffs,
	)

	// Test delete
	fs.(*mockSystem).files = fs.(*mockSystem).files[:4]

	deleteScan, err := w.scan()
	require.Nil(t, err)

	diffs, err = w.diff(deleteScan)
	require.Equal(t,
		[]WatchDifferential{
			WatchDifferential{Path: "a/3.txt", Change: Deleted},
		},
		diffs,
	)
}

func Test_ATimeoutWatcherReturnsAChannelAndClosesItAfterUse(t *testing.T) {

	c := make(chan struct{})
	defer close(c)

	w := NewIntervalWatcher(
		mockFileReaderForWatch(t),
		1*time.Millisecond,
		c,
		"a/*.txt",
		"b/*.txt",
	)

	output, err := w.Watch()
	require.Nil(t, err)
	require.NotNil(t, output)

	c <- struct{}{}

	_, ok := <-output
	require.False(t, ok)

	assertExpectedFileInfos(t,
		w.(*intervalWatcher).infos,
		map[string][]string{
			"a/*.txt": []string{"a/1.txt", "a/2.txt"},
			"b/*.txt": []string{"b/1.txt", "b/2.txt"},
		},
	)
}

func Test_ATimeoutWatcherWatchesOnATimeout(t *testing.T) {

	cancel := make(chan struct{})
	defer func() { cancel <- struct{}{} }()

	fs := mockFileReaderForWatch(t)

	w := NewIntervalWatcher(
		fs,
		1*time.Millisecond,
		cancel,
		"a/*.txt",
		"b/*.txt",
	).(*intervalWatcher)

	output, err := w.Watch()
	require.Nil(t, err)
	require.NotNil(t, output)

	require.Nil(t, fs.Write("a/3.txt", bytes.NewBufferString("Hello!")))

	select {
	case diffs, ok := <-output:
		require.True(t, ok)
		require.Equal(t,
			[]WatchDifferential{
				WatchDifferential{Path: "a/3.txt", Change: Created},
			},
			diffs,
		)

	case <-time.After(10 * time.Millisecond):
		t.Fatalf("Timed out waiting for diff")
	}
}
