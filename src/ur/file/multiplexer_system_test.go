package file

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_MultiplexerSystemIsASystem(t *testing.T) {
	func(s SystemReader) {}(NewMultiplexerSystem())
}

func Test_MultiplexerSystemCanList(t *testing.T) {

	for _, expected := range []struct {
		inputs  [][]string
		outputs []string
	}{
		{
			inputs: [][]string{
				{
					"one.js",
					"two.js",
				},
				{
					"three.js",
				},
			},
			outputs: []string{
				"one.js",
				"two.js",
				"three.js",
			},
		},
	} {

		var systems []SystemReader

		for _, inputset := range expected.inputs {

			system := NewMockSystem().(*mockSystem)

			for _, file := range inputset {
				system.files = append(system.files, MockSystemFile{
					Path: file,
				})
			}

			systems = append(systems, system)
		}

		multi := NewMultiplexerSystem(systems...)

		list, err := multi.List()
		assert.Nil(t, err)

		assert.Equal(t, expected.outputs, list)
	}
}

func Test_MultiplexerSystemCanGlob(t *testing.T) {

	for _, expected := range []struct {
		glob    string
		inputs  [][]string
		outputs []string
	}{
		{
			glob: "*.js",
			inputs: [][]string{
				{
					"one.js",
					"two.js",
				},
				{
					"three.js",
				},
			},
			outputs: []string{
				"one.js",
				"two.js",
				"three.js",
			},
		},
	} {

		var systems []SystemReader

		for _, inputset := range expected.inputs {

			system := NewMockSystem().(*mockSystem)

			for _, file := range inputset {
				system.files = append(system.files, MockSystemFile{
					Path: file,
				})
			}

			systems = append(systems, system)
		}

		multi := NewMultiplexerSystem(systems...)

		fileglob, err := multi.Glob(expected.glob)
		assert.Nil(t, err)

		assert.Equal(t, expected.outputs, fileglob)
	}
}

func Test_MultiplexerSystemCanRead(t *testing.T) {

	system1 := NewMockSystem(
		MockSystemFile{
			Path:     "something.txt",
			Contents: []byte("hello"),
		},
	)

	system2 := NewMockSystem(
		MockSystemFile{
			Path:     "something-else.txt",
			Contents: []byte("bye"),
		},
	)

	multi := NewMultiplexerSystem(system1, system2)

	read1, err1 := multi.Read("something.txt")
	assert.Nil(t, err1)
	assert.Equal(t, []byte("hello"), read1)

	read2, err2 := multi.Read("something-else.txt")
	assert.Nil(t, err2)
	assert.Equal(t, []byte("bye"), read2)

	_, err3 := multi.Read("nothing.txt")
	assert.NotNil(t, err3)
}
