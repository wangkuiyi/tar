package tar

// The following example shows how to create a tarball as a table in the SQLFS.
/*
func Tar(db *sql.DB, dir string, inc include, compress bool) (e error) {
	// Create a file in the SQLFS.
	sqlfn := fmt.Sprintf("sqlflow_models.%s", dir)
	sqlf, e := sqlfs.Create(db, sqlfn)
	if e != nil {
		return fmt.Errorf("Cannot create sqlfs file %s: %v", sqlfn, e)
	}
	defer func() { e = sqlf.Close() }()

	return tar.Tar(sqlf, dir, inc, compress)
}
*/

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
)

type include func(dir, relative string, fi os.FileInfo) bool

// Tar a directory into a writer.
func Tar(w io.Writer, dir string, inc include) (e error) {
	gw := gzip.NewWriter(w)
	defer gw.Close()
	tw := tar.NewWriter(gw)
	defer tw.Close()

	return recursiveTar(tw, dir, "", inc)
}

func recursiveTar(tw *tar.Writer, dir, relative string, inc include) (e error) {
	abs := path.Join(dir, relative)
	fis, e := ioutil.ReadDir(abs)
	if e != nil {
		return fmt.Errorf("Tar: ReadDir(%s) failed: %v", abs, e)
	}

	for _, fs := range fis {
		if inc == nil || inc(dir, relative, fs) { // Include only certain files.
			fn := path.Join(relative, fs.Name())
			if fs.IsDir() {
				if e = recursiveTar(tw, dir, fn, inc); e != nil {
					return e
				}
			} else {
				if e = tw.WriteHeader(&tar.Header{
					Name: fn,
					Mode: 0600,
					Size: fs.Size()}); e != nil {
					return fmt.Errorf("Tar: WriteHeader(%s): %v", fn, e)
				}

				f, e := os.Open(path.Join(dir, fn))
				if e != nil {
					return fmt.Errorf("Tar: cannot open %s: %v", fn, e)
				}
				n, e := io.Copy(tw, f)
				if n != fs.Size() {
					return fmt.Errorf("Tar: copied %d of %d bytes of %s", n, fs.Size(), fn)
				}
				if e != nil {
					return fmt.Errorf("Tar: failed to copy %s: %v", fn, e)
				}
				f.Close()
			}
		}
	}
	return nil
}
