package os

import (
	"path"
	"strings"
	"time"
)

func navigateTo(p string) (*directory, error) {
	if len(p) == 0 {
		return cwd, nil
	}
	d := cwd
	if p[0] == '/' {
		d = root
		p = p[1:]
	}
	for _, dir := range strings.Split(p, "/") {
		switch dir {
		case "":
		case ".":
		case "..":
			d = d.parent
		default:
			fi, err := d.get(dir)
			if err != nil {
				return nil, err
			}
			if !fi.IsDir() {
				return nil, ErrInvalid
			}
			d = fi.(*directory)
		}
	}
	return d, nil
}

func getFile(p string) (FileInfo, error) {
	dir, file := path.Split(path.Clean(p))
	d, err := navigateTo(dir)
	if err != nil {
		return nil, err
	}
	return d.get(file)
}

func Chdir(p string) error {
	cwdmu.Lock()
	defer cwdmu.Unlock()
	c, err := navigateTo(path.Clean(p))
	if err != nil {
		return &PathError{
			"chdir",
			p,
			err,
		}
	}
	cwd = c
	return nil
}

func Chmod(p string, mode FileMode) error {
	f, err := getFile(p)
	if err != nil {
		return &PathError{
			"chmod",
			p,
			err,
		}
	}
	type i interface {
		chmod(FileMode)
	}
	f.(i).chmod(mode)
	return nil
}

func Chown(p string, _, _ int) error {
	return &PathError{
		"chown",
		p,
		ErrUnsupported,
	}
}

func Chtimes(p string, _, mtime time.Time) error {
	f, err := getFile(p)
	if err != nil {
		return &PathError{
			"chtimes",
			p,
			err,
		}
	}
	type i interface {
		setModTime(time.Time)
	}
	f.(i).setModTime(mtime)
	return nil
}

func Clearenv() {

}

func Environ() []string {
	return []string{}
}

func Exit(_ int) {

}

func Expand(s string, mapping func(string) string) string {
	return s
}

func ExpandEnv(s string) string {
	return s
}

func Getegid() int {
	return 0
}

func Getenv(_ string) string {
	return ""
}

func Geteuid() int {
	return 0
}

func Getgid() int {
	return 0
}

func Getgroups() ([]int, error) {
	return []int{}, nil
}

func Getpagesize() int {
	return 0
}

func Getpid() int {
	return 0
}

func Getppid() int {
	return 0
}

func Getuid() int {
	return 0
}

func Getwd() (string, error) {
	d := cwd
	names := make([]string, 1, 32)
	names[0] = d.Name()
	for d != d.parent {
		d = d.parent
		names = append(names, d.Name())
	}
	l := len(names) - 1
	for i := 0; i < l>>1-1; i++ {
		names[i], names[l-i] = names[l-i], names[i]
	}
	return strings.Join(names, "/"), nil
}

func Hostname() (string, error) {
	return "", nil
}

func IsPathSeparator(c uint8) bool {
	return c == '/'
}

func Lchown(p string, _, _ int) error {
	return &PathError{
		"lchown",
		p,
		ErrUnsupported,
	}
}

func Link(oldname, newname string) error {
	return &LinkError{
		"link",
		oldname,
		newname,
		ErrUnsupported,
	}
}

func Mkdir(p string, fileMode FileMode) error {
	dir, toMake := path.Split(path.Clean(p))
	d, err := navigateTo(dir)
	if err == nil {
		_, err = d.mkdir(toMake, fileMode)
	}
	if err != nil {
		return &PathError{
			"mkdir",
			p,
			err,
		}
	}
	return nil
}

func MkdirAll(p string, fileMode FileMode) error {
	p = path.Clean(p)
	d := cwd
	if p[0] == '/' {
		d = root
		p = p[1:]
	}
	var err error
	for _, dir := range strings.Split(p, "/") {
		d, err = d.mkdir(dir, fileMode)
		if IsPermission(err) {
			return &PathError{
				"mkdirall",
				p,
				err,
			}
		}
	}
	return nil
}

func NewSyscallError(_, string, _ error) error {
	return nil
}

func Readlink(name string) (string, error) {
	return "", &PathError{
		"readlink",
		name,
		ErrUnsupported,
	}
}

func Remove(name string) error {
	dir, file := path.Split(path.Clean(name))
	d, err := navigateTo(dir)
	if err == nil {
		err = d.remove(file, false)
	}
	if err != nil {
		return &PathError{
			"remove",
			name,
			err,
		}
	}
	return nil
}

func RemoveAll(name string) error {
	dir, file := path.Split(path.Clean(name))
	d, err := navigateTo(dir)
	if err == nil {
		err = d.remove(file, true)
	}
	if err != nil {
		return &PathError{
			"remove",
			name,
			err,
		}
	}
	return nil
}

func Rename(oldpath, newpath string) error {
	olddir, oldfile := path.Split(path.Clean(oldpath))
	newdir, newfile := path.Split(path.Clean(newpath))
	oldd, err := navigateTo(olddir)
	if err == nil {
		var newd *directory
		newd, err = navigateTo(newdir)
		f, err := oldd.get(oldfile)
		if err == nil {
			type i interface {
				move(string, *directory) error
			}
			err = f.(i).move(newfile, newd)
		}
	}
	if err != nil {
		return &LinkError{
			"rename",
			oldpath,
			newpath,
			err,
		}
	}
	return nil
}

func SameFile(f, g FileInfo) bool {
	return f == g
}

func Setenv(key, value string) error {
	return ErrUnsupported
}

func Symlink(oldname, newname string) error {
	return &LinkError{
		"symlink",
		oldname,
		newname,
		ErrUnsupported,
	}
}

func TempDir() string {
	return "/tmp"
}

func Truncate(name string, size int64) error {
	f, err := getFile(name)
	if err == nil {
		if f, ok := f.(*file); ok {
			if f.canWrite() {
				f.Contents = f.Contents[:size]
			} else {
				err = ErrPermission
			}
		} else {
			err = ErrInvalid
		}
	}
	if err != nil {
		return &PathError{
			"truncate",
			name,
			err,
		}
	}
	return nil
}

func Unsetenv(_ string) error {
	return ErrUnsupported
}
