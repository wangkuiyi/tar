package tar

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUntar(t *testing.T) {
	a := assert.New(t)

	f, e := ioutil.TempFile("", "*.tar.gz")
	a.NoError(e)
	defer os.RemoveAll(f.Name())

	a.NoError(Tar(f, dir, nil))
	f.Close()

	r, e := os.Open(f.Name())
	a.NoError(e)

	untar, e := ioutil.TempDir("", "")
	a.NoError(e)
	defer os.RemoveAll(untar)

	a.NoError(Untar(r, untar))
}
