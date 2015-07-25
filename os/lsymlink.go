package os

import "time"

type symbolic struct {
	metadata
	link string
}

func newSymbolic(d *directory, name string, mode FileMode, modTime time.Time, contents string) {

}

func (s symbolic) Sys() interface{} {
	return s.link
}
