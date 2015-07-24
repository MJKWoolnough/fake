package os

import "time"

var root = directory{
	metadata{
		"",
		0,
		0755,
		time.Now(),
	},
	&root, &root,
	make([]FileInfo, 0),
}

type metadata struct {
	name    string
	size    int64
	mode    FileMode
	modTime time.Time
}

func (m metadata) Name() string {
	return m.Name()
}

func (m metadata) Size() int64 {
	return m.size
}

func (m metadata) Mode() FileMode {
	return m.mode
}

func (m metadata) IsDir() bool {
	return m.mode.IsDir()
}

type file struct {
	metadata
	data []byte
}

func newFile(d *directory, name string, mode FileMode, modTime time.Time, contents []byte) {

}

func (f file) Sys() interface{} {
	return f.data
}

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

type symbolic struct {
	metadata
	link string
}

func newSymbolic(d *directory, name string, mode FileMode, modTime time.Time, contents string) {

}

func (s symbolic) Sys() interface{} {
	return s.link
}
