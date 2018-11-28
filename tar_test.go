package sqlfsutil

import (
	"fmt"
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
	ioutil.WriteFile(fn, bytes('a'), 0600)

	a.NoError(os.Mkdir(path.Join(dir, "b"), 0600))
	fn = filepath.Join(dir, "b", "b")
	ioutil.WriteFile(fn, bytes('b'), 0600)

	f, e := ioutil.TempFile("", "*")
	a.NoError(e)

	a.NoError(Tar(f, dir, nil, true))

	a.NoError(exec.Command("tar", "-C", dir, "-xzf", f.Name()).Run())
	fmt.Println(dir)
}
