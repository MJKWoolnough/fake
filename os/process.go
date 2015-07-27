package os

type Process struct{}

func FindProcess(_ int) (*Process, error) {
	return ErrUnsupported
}

func StartProcess(name string, _ []string, _ interface{}) (*Process, error) {
	return nil, &PathError{
		"fork/exec",
		name,
		ErrUnsupported,
	}
}
