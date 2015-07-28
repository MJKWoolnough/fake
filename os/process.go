package os

type Process struct{}

func FindProcess(_ int) (*Process, error) {
	return nil, ErrUnsupported
}

func StartProcess(name string, _ []string, _ interface{}) (*Process, error) {
	return nil, &PathError{
		"fork/exec",
		name,
		ErrUnsupported,
	}
}
