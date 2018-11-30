package tar

import (
	"archive/tar"
	"compress/gzip"
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
		if _, e := io.Copy(w, tr); e != nil {
			return e
		}
		w.Close()
	}
	return nil
}
