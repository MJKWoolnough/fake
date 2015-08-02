package os

import "os"

func canExecute(f os.FileMode) bool {
	return f&0111 > 0
}

func canWrite(f os.FileMode) bool {
	return f&0222 > 0
}

func canRead(f os.FileMode) bool {
	return f&0444 > 0
}
