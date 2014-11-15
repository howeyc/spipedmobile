// servebundle
package assets

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/alecthomas/gobundle"
)

type ServeBundle struct {
	Bundle *gobundle.Bundle
}

type BundleFile struct {
	name    string
	reader  *bytes.Reader
	modTime time.Time
}

func (srv *ServeBundle) Open(filename string) (http.File, error) {
	if r, err := srv.Bundle.Open(filename[1:]); err == nil {
		ball, errr := ioutil.ReadAll(r)
		return &BundleFile{filename, bytes.NewReader(ball), time.Now()}, errr
	} else {
		return nil, err
	}
}

func (bf *BundleFile) Read(p []byte) (n int, err error) {
	return bf.reader.Read(p)
}

func (bf *BundleFile) Seek(offset int64, whence int) (int64, error) {
	return bf.reader.Seek(offset, whence)
}

func (bf *BundleFile) Close() error {
	return nil
}

func (bf *BundleFile) Readdir(count int) ([]os.FileInfo, error) {
	return nil, errors.New("WTF?!?")
}

func (bf *BundleFile) Stat() (os.FileInfo, error) {
	return bf, nil
}

func (bf *BundleFile) Name() string {
	return bf.name
}

func (bf *BundleFile) Size() int64 {
	return int64(bf.reader.Len())
}

func (bf *BundleFile) Mode() os.FileMode {
	return os.FileMode(0777)
}

func (bf *BundleFile) ModTime() time.Time {
	return bf.modTime
}

func (bf *BundleFile) IsDir() bool {
	return false
}

func (bf *BundleFile) Sys() interface{} {
	return nil
}
