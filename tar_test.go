package tar

import (
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	dir string // to be set by TestMain
)

func TestRecursiveTar(t *testing.T) {
	a := assert.New(t)

	f, e := ioutil.TempFile("", "*.tar.gz")
	a.NoError(e)
	defer os.RemoveAll(f.Name())

	a.NoError(Tar(f, dir, nil))

	o, e := exec.Command("tar", "-ztf", f.Name()).Output()
	a.NoError(e)
	a.Equal("a\nb/b\n", string(o))
}

func TestRecursiveTarWithIncludeFileter(t *testing.T) {
	a := assert.New(t)

	f, e := ioutil.TempFile("", "*.tar.gz")
	a.NoError(e)
	defer os.RemoveAll(f.Name())

	inc := func(dir, relative string, fi os.FileInfo) bool {
		level := len(strings.Split(path.Join(relative, fi.Name()), "/")) - 1
		return level == 0 // only files in the top directory level.
	}
	a.NoError(Tar(f, dir, inc))

	o, e := exec.Command("tar", "-ztf", f.Name()).Output()
	a.NoError(e)
	a.Equal("a\n", string(o))
}

// Create a directory hierarchy in a temporary directory and returns
// this temporary directory.  The hierarchy includes two files:
//
// /a
// /b/b
//
// where /a is composed of 1024 'a's, and /b/b of 1024 'b's.
func createDirectoryHierarchy() (dir string, e error) {
	if dir, e = ioutil.TempDir("", ""); e != nil {
		return "", e
	}

	fill := func(a rune, n int) []byte {
		b := make([]byte, n)
		for i := range b {
			b[i] = byte(a)
		}
		return b
	}

	fn := filepath.Join(dir, "a")
	ioutil.WriteFile(fn, fill('a', 1024), 0700)

	if e = os.Mkdir(path.Join(dir, "b"), 0700); e != nil {
		return "", e
	}
	fn = filepath.Join(dir, "b", "b")
	ioutil.WriteFile(fn, fill('b', 1024), 0700)
	return dir, nil
}

func TestMain(m *testing.M) {
	d, e := createDirectoryHierarchy()
	if e != nil {
		log.Fatalf("Cannot create testing directory hierarchy: %v", e)
	}
	defer os.RemoveAll(dir)

	dir = d
	os.Exit(m.Run())
}
