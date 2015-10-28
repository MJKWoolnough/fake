package os

import (
	"io"
	"io/ioutil"
)

type discard interface {
	io.Writer
	io.ReaderFrom
	WriteString(string) (int, error)
}

type nullDevice struct {
	discard
}

func makeNull() nullDevice {
	return nullDevice{
		ioutil.Discard.(discard),
	}
}

func (devNull) Read([]byte) (int, error) {
	return 0, io.EOF
}

func (devNull) Seek(int64, int) (int64, error) {
	return 0, ErrInvalid
}

func (devNull) Close() error {
	return nil
}

func (devNull) ReaderAt([]byte, int64) (int, error) {
	return 0, ErrInvalid
}

func (devNull) WriterAt([]byte, int64) (int, error) {
	return 0, ErrInvalid
}

func (devNull) WriteTo(w io.Writer) (int64, error) {
	return 0, io.EOF
}

func (devNull) Sync() error {
	return nil
}

func (devNull) Truncate(int64) error {
	return ErrInvalid
}
