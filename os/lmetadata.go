package os

import "time"

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
