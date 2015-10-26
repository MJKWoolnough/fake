package os

import "os"

type FileInfo interface {
	os.FileInfo
}

type Signal interface {
	os.Signal
}
