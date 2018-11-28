# `tar`

![](https://travis-ci.org/wangkuiyi/tar.svg?branch=develop)

To create a gziped tarball, simply call

```go
tarfile, _ := os.Create("/tmp/my.tar.gz")
Tar(tarfile, dir, nil, true)
```

To create a tarball in a SQL table by taking the table as a filesystem, we can use the [`sqlfs`](https://github.com/wangkuiyi/sqlfs) package:

```go
func Tar(db *sql.DB, dir string, inc include, compress bool) (e error) {
	// Create a file in the SQLFS.
	fn := strings.Replace(strings.Replace(dir, ".", "_"), "-", "_")
	sqlfn := fmt.Sprintf("sqlflow_models.%s", fn)
	sqlf, e := sqlfs.Create(db, sqlfn)
	if e != nil {
		return fmt.Errorf("Cannot create sqlfs file %s: %v", sqlfn, e)
	}
	defer func() { e = sqlf.Close() }()

	return TarDir(sqlf, dir, inc, compress)
}
```
