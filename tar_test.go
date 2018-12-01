package tar

import (
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRecursiveTar(t *testing.T) {
	a := assert.New(t)

	// Create a directory
	dir, e := ioutil.TempDir("", "")
	a.NoError(e)
	defer os.RemoveAll(dir)

	bytes := func(a rune) []byte {
		b := make([]byte, 1024)
		for i := range b {
			b[i] = byte(a)
		}
		return b
	}

	fn := filepath.Join(dir, "a")
	ioutil.WriteFile(fn, bytes('a'), 0700)

	a.NoError(os.Mkdir(path.Join(dir, "b"), 0700))
	fn = filepath.Join(dir, "b", "b")
	ioutil.WriteFile(fn, bytes('b'), 0700)

	f, e := ioutil.TempFile("", "*.tar.gz")
	a.NoError(e)
	defer os.RemoveAll(f.Name())

	a.NoError(Tar(f, dir, nil, true))

	o, e := exec.Command("tar", "-ztf", f.Name()).Output()
	a.NoError(e)
	a.Equal("a\nb/b\n", string(o))
}
