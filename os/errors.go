package os

import "errors"

var (
	ErrInvalid     = errors.New("invalid argument")
	ErrPermission  = errors.New("permission denied")
	ErrExist       = errors.New("file already exists")
	ErrNotExist    = errors.New("file does not exist")
	ErrUnsupported = errors.New("unsupported feature")
	ErrNotEmpty    = errors.New("directory not empty")
	ErrClosed      = errors.New("file closed")
)

type PathError struct {
	Op, Path string
	Err      error
}

func (p *PathError) Error() string {
	return p.Op + " " + p.Path + ": " + p.Err.Error()
}

type LinkError struct {
	Op, Old, New string
	Err          error
}

func (l *LinkError) Error() string {
	return l.Op + " " + l.Old + " " + l.New + ": " + l.Err.Error()
}

func IsExist(err error) bool {
	switch e := err.(type) {
	case *PathError:
		err = e.Err
	case *LinkError:
		err = e.Err
	}
	return err == ErrExist
}

func IsNotExist(err error) bool {
	switch e := err.(type) {
	case *PathError:
		err = e.Err
	case *LinkError:
		err = e.Err
	}
	return err == ErrNotExist
}

func IsPermission(err error) bool {
	switch e := err.(type) {
	case *PathError:
		err = e.Err
	case *LinkError:
		err = e.Err
	}
	return err == ErrPermission
}
