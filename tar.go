package sqlfsutil

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

type include func(dir string, fi os.FileInfo) bool

// Tar a directory into a writer.
func Tar(w io.Writer, dir string, inc include, compress bool) (e error) {
	var tw *tar.Writer
	if compress {
		gw := gzip.NewWriter(w)
		defer gw.Close()
		tw = tar.NewWriter(gw)
	} else {
		tw = tar.NewWriter(w)
	}
	defer tw.Close()

	return recursiveTar(tw, dir, inc, compress)
}

func recursiveTar(tw *tar.Writer, dir string, inc include, compress bool) (e error) {
	fis, e := ioutil.ReadDir(dir)
	if e != nil {
		return fmt.Errorf("Tar: ReadDir(%s) failed: %v", dir, e)
	}
	for _, fs := range fis {
		if inc != nil && inc(dir, fs) { // Include only certain files.
			fn := path.Join(dir, fs.Name())
			if fs.IsDir() {
				if e = recursiveTar(tw, fn, inc, compress); e != nil {
					return e
				}
			} else {
				if e = tw.WriteHeader(&tar.Header{
					Name: fn,
					Mode: 0600,
					Size: fs.Size()}); e != nil {
					return fmt.Errorf("Tar: WriteHeader(%s): %v", fn, e)
				}

				f, e := os.Open(fn)
				if e != nil {
					return fmt.Errorf("Tar: cannot open %s: %v", fn, e)
				}
				_, e = io.Copy(tw, f)
				if e != nil {
					return fmt.Errorf("Tar: failed to copy %s: %v", fn, e)
				}
				f.Close()
			}
		}
	}
	return nil
}
