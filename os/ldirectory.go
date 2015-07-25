package os

import "time"

type directory struct {
	metadata
	parent, self *directory
	contents     []FileInfo
}

func newDirectory(d *directory, name string, mode FileMode, modTime time.Time, contents []FileInfo) {

}

func (d directory) Sys() interface{} {
	return d.contents
}
