package file

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_MockSystemIsASystemReader(t *testing.T) {
	func(s SystemReader) {}(NewMockSystem())
}

func Test_MockSystemIsASystemWriter(t *testing.T) {
	func(s SystemWriter) {}(NewMockSystem())
}

func Test_MockSystemCanList(t *testing.T) {

	for _, expected := range []struct {
		inputs []string
	}{
		{
			inputs: []string{
				"one.txt",
				"two.js",
				"three.js",
				"four.txt",
				"five.go",
			},
		},
		{
			inputs: []string{
				"one.txt",
				"two.js",
				"three.js",
				"four.txt",
				"five.go",
			},
		},
		{
			inputs: []string{
				"one.txt",
				"two.js",
				"three.js",
				"four.txt",
				"five.go",
			},
		},
	} {

		system := NewMockSystem().(*mockSystem)

		for _, file := range expected.inputs {
			system.files = append(system.files, MockSystemFile{
				Path: file,
			})
		}

		list, err := system.List()
		require.Nil(t, err)

		require.Equal(t, expected.inputs, list)
	}
}

func Test_MockSystemCanGlob(t *testing.T) {

	for _, expected := range []struct {
		glob    string
		inputs  []string
		outputs []string
	}{
		{
			glob: "*.js",
			inputs: []string{
				"one.txt",
				"two.js",
				"three.js",
				"four.txt",
				"five.go",
			},
			outputs: []string{
				"two.js",
				"three.js",
			},
		},
		{
			glob: "*.txt",
			inputs: []string{
				"one.txt",
				"two.js",
				"three.js",
				"four.txt",
				"five.go",
			},
			outputs: []string{
				"one.txt",
				"four.txt",
			},
		},
		{
			glob: "*.dasdqw",
			inputs: []string{
				"one.txt",
				"two.js",
				"three.js",
				"four.txt",
				"five.go",
			},
		},
	} {

		system := NewMockSystem().(*mockSystem)

		for _, file := range expected.inputs {
			system.files = append(system.files, MockSystemFile{
				Path: file,
			})
		}

		fileglob, err := system.Glob(expected.glob)
		require.Nil(t, err)

		require.Equal(t, expected.outputs, fileglob)
	}
}

func Test_MockSystemCanRead(t *testing.T) {

	system := NewMockSystem(
		MockSystemFile{
			Path:     "something.txt",
			Contents: []byte("hello"),
		},
	)

	read1, err1 := system.Read("something.txt")
	require.Nil(t, err1)
	require.Equal(t, []byte("hello"), read1)

	_, err2 := system.Read("nothing.txt")
	require.NotNil(t, err2)
}

func Test_AMockSystemCanWrite(t *testing.T) {
	ASystemCanWrite(t, NewMockSystem())
}
