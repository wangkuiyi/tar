package tar

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path"
)

// Untar from a reader into a directory
func Untar(r io.Reader, dir string) (e error) {
	gr, e := gzip.NewReader(r)
	if e != nil {
		return e
	}
	defer gr.Close()
	tr := tar.NewReader(gr)

	for {
		hdr, e := tr.Next()
		if e == io.EOF {
			break // End of archive
		}
		if e != nil {
			return e
		}
		fn := path.Join(dir, hdr.Name)
		dr, _ := path.Split(fn)
		os.MkdirAll(dr, 0777)
		w, e := os.Create(fn)
		if e != nil {
			return e
		}
		n, e := io.Copy(w, tr)
		if e != nil {
			return e
		}
		if n != hdr.Size {
			return fmt.Errorf("Read %d of %d bytes of %s", n, hdr.Size, hdr.Name)
		}
		w.Close()

		_, e = os.Stat(fn)
	}
	return nil
}
