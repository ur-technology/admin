package file

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

func ASystemCanWrite(t *testing.T, system System) {

	for cycle, input := range []struct {
		path    string
		content string
	}{
		{
			path:    "test/one.txt",
			content: "some text",
		},
		{
			path:    "test/one/two/three/deep.txt",
			content: "some other text",
		},
		{
			path:    "root",
			content: "root",
		},
		{
			path:    "root",
			content: "root overwrite",
		},
	} {
		t.Logf("Cycle %d", cycle)

		buf := bytes.NewBufferString(input.content)
		err := system.Write(input.path, buf)
		require.Nil(t, err)

		content, err := system.Read(input.path)
		require.Nil(t, err)
		require.Equal(t, content, []byte(input.content))
	}
}

func Test_ASystemCanBeCopied(t *testing.T) {

	from := NewMockSystem(
		MockSystemFile{Path: "a.txt"},
		MockSystemFile{Path: "b.txt"},
		MockSystemFile{Path: "c.txt"},
	)

	to := NewMockSystem()

	require.Nil(t, Copy(from, to))

	fromList, err := from.List()
	require.Nil(t, err)

	toList, err := to.List()
	require.Nil(t, err)

	require.Equal(t, fromList, toList)
}
