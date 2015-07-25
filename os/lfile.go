package os

import "time"

type file struct {
	metadata
	data []byte
}

func newFile(d *directory, name string, mode FileMode, modTime time.Time, contents []byte) {

}

func (f file) Sys() interface{} {
	return f.data
}
