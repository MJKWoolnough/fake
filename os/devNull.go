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

func (nullDevice) Read([]byte) (int, error) {
	return 0, io.EOF
}

func (nullDevice) Seek(int64, int) (int64, error) {
	return 0, ErrInvalid
}

func (nullDevice) Close() error {
	return nil
}

func (nullDevice) ReaderAt([]byte, int64) (int, error) {
	return 0, ErrInvalid
}

func (nullDevice) WriterAt([]byte, int64) (int, error) {
	return 0, ErrInvalid
}

func (nullDevice) WriteTo(w io.Writer) (int64, error) {
	return 0, io.EOF
}

func (nullDevice) Sync() error {
	return nil
}

func (nullDevice) Truncate(int64) error {
	return ErrInvalid
}
