package os

import (
	"errors"
	"os"
)

var (
	ErrInvalid     = os.ErrInvalid
	ErrPermission  = os.ErrPermission
	ErrExist       = os.ErrExist
	ErrNotExist    = os.ErrNotExist
	ErrUnsupported = errors.New("unsupported feature")
	ErrNotEmpty    = errors.New("directory not empty")
	ErrClosed      = errors.New("file closed")
	ErrIsDir       = errors.New("is directory")
	ErrIsNotDir    = errors.New("is not directory")
)

type PathError os.PathError

func (p *PathError) Error() string {
	return p.Op + " " + p.Path + ": " + p.Err.Error()
}

type LinkError os.LinkError

func (l *LinkError) Error() string {
	return l.Op + " " + l.Old + " " + l.New + ": " + l.Err.Error()
}

func IsExist(err error) bool {
	switch e := err.(type) {
	case *PathError:
		err = e.Err
	case *os.PathError:
		err = e.Err
	case *LinkError:
		err = e.Err
	case *os.LinkError:
		err = e.Err
	}
	return err == ErrExist
}

func IsNotExist(err error) bool {
	switch e := err.(type) {
	case *PathError:
		err = e.Err
	case *os.PathError:
		err = e.Err
	case *LinkError:
		err = e.Err
	case *os.LinkError:
		err = e.Err
	}
	return err == ErrNotExist
}

func IsPermission(err error) bool {
	switch e := err.(type) {
	case *PathError:
		err = e.Err
	case *os.PathError:
		err = e.Err
	case *LinkError:
		err = e.Err
	case *os.LinkError:
		err = e.Err
	}
	return err == ErrPermission
}
