package os

import "errors"

var (
	ErrInvalid     = errors.New("invalid argument")
	ErrPermission  = errors.New("permission denied")
	ErrExist       = errors.New("file already exists")
	ErrNotExist    = errors.New("file does not exist")
	ErrUnsupported = errors.New("unsupported feature")
	ErrNotEmpty    = errors.New("directory not empty")
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
	if p, ok := err.(*PathError); ok {
		err = p.Err
	}
	return err == ErrExist
}

func IsNotExist(err error) bool {
	if p, ok := err.(*PathError); ok {
		err = p.Err
	}
	return err == ErrNotExist
}

func IsPermission(err error) bool {
	if p, ok := err.(*PathError); ok {
		err = p.Err
	}
	return err == ErrPermission
}
