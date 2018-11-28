# `tar`

To create a gziped tarball, simply call

```go
tarfile, _ := os.Create("/tmp/my.tar.gz")
Tar(tarfile, dir, nil, true)
```
