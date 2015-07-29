package os

// Constants for FileMode
const (
	ModeDir FileMode = 1 << (32 - 1 - iota)
	ModeAppend
	ModeExclusive
	ModeTemporary
	ModeSymlink
	ModeDevice
	ModeNamedPipe
	ModeSocket
	ModeSetuid
	ModeSetgid
	ModeCharDevice
	ModeSticky
	ModeSpecial

	ModeType          = ModeDir | ModeSymlink | ModeNamedPipe | ModeSocket | ModeDevice
	ModePerm FileMode = 0777
)

type FileMode uint32

func (f FileMode) IsDir() bool {
	return f&ModeDir != 0
}

func (f FileMode) IsRegular() bool {
	return f&ModeType != 0
}

func (f FileMode) Perm() FileMode {
	return f & ModePerm
}

func (f FileMode) canExecute() bool {
	return f&0111 > 0
}

func (f FileMode) canWrite() bool {
	return f&0222 > 0
}

func (f FileMode) canRead() bool {
	return f&0444 > 0
}

func (f FileMode) isSpecial() bool {
	return f&ModeSpecial != 0
}

func (f FileMode) String() string {
	return ""
}
