package os

import "os"

const (
	O_RDONLY = os.O_RDONLY
	O_WRONLY = os.O_WRONLY
	O_RDWR   = os.RDWR
	O_APPEND = os.APPEND
	O_CREATE = os.O_CREATE
	O_EXCL   = os.O_EXCL
	O_SYNC   = os.O_SYNC
	O_TRUNC  = os.O_TRUNC
)

const (
	SEEK_SET = os.SEEK_SET
	SEEK_CUR = os.SEEK_CUR
	SEEK_END = os.SEEK_END
)

const (
	PathSeparator     = os.PathSeparator
	PathListSeparator = os.PathListSeparator
)

const DevNull = os.DevNull
