package os

import (
	"sync"
	"time"
)

var (
	root *directory

	chdir sync.Mutex
	cwd   *directory
)

func init() {
	root = directory{
		metadata: *metadata{
			"",
			0,
			ModeDir | 0755,
			time.Now(),
		},
		contents: make([]FileInfo, 0),
	}
	root.self = root
	root.parent = root
	cwd = &root
	Mkdir("temp", 0755)
	Chmod("/", 0644)
	ChDir("temp")
}
