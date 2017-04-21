package file

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func newTmpDiskSystem(t *testing.T) (*diskSystem, func()) {

	dir, err := ioutil.TempDir("", "sepiaFilePackageTests")
	require.Nil(t, err)

	return NewDiskSystem(dir).(*diskSystem), func() { os.RemoveAll(dir) }
}

func Test_DiskSystemIsASystemReader(t *testing.T) {
	func(s SystemReader) {}(NewDiskSystem())
}

func Test_DiskSystemIsASystemWriter(t *testing.T) {
	func(s SystemWriter) {}(NewDiskSystem())
}

func Test_ADiskSystemCanWrite(t *testing.T) {

	system, cleanup := newTmpDiskSystem(t)
	defer cleanup()

	ASystemCanWrite(t, system)
}
