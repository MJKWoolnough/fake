package os

import (
	"path"
	"reflect"
	"strings"
	"time"
	"unsafe"
)

func WriteBytes(p string, perm FileMode, data []byte) {
	var filename string
	p, filename = path.Split(path.Clean(p))
	if len(p) == 0 {
		return
	}
	d := cwd
	if p[0] == '/' {
		d = root
		p = p[1:]
	}
	for _, dir := range strings.Split(p, "/") {
		switch dir {
		case "", ".":
		case "..":
			d = d.parent
		default:
			fi, ok := d.Contents[dir]
			if !ok || !fi.IsDir() {
				e := &directory{
					node{
						ModeDir | 0777,
						time.Now(),
						dir,
						d,
					},
					make(map[string]FileInfo),
				}
				d.Contents[dir] = e
				d = e
			} else {
				d = fi.(*directory)
			}
		}
	}
	d.Contents[filename] = &file{
		node{
			perm,
			time.Now(),
			filename,
			d,
		},
		data,
	}
}

func WriteString(p, data string) {
	s := (*reflect.StringHeader)(unsafe.Pointer(&data))
	WriteBytes(p, ModeSpecial|0400, *(*[]byte)(unsafe.Pointer(&reflect.SliceHeader{s.Data, s.Len, s.Len})))
}
