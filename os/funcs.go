package os

import (
	"path"
	"strings"
	"time"
)

func navigateTo(p string) (*directory, error) {
	d := cwd
	if p[0] == '/' {
		d = root
		p = p[1:]
	}
	for _, dir := range strings.Split(p, "/") {
		switch dir {
		case ".":
		case "..":
			d = d.Parent
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
	return d
}

func getFile(p string) (FileInfo, error) {
	dir, file := p.Split(path.Clean(p))
	d, err := navigateTo(dir)
	if err != nil {
		return err
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
	f.(i).chmod(mode)
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
	0
}

func GetGroups() ([]int, error) {
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
	for d != d.Parent {
		d = d.Parent
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
		err = dir.mkdir(toMake, fileMode)
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
	for _, dir := range strings.Split(p, "/") {
		err := d.mkdir(dir, fileMode)
		if IsPermissions(err) {
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
	dir, file := p.Split(path.Clean(p))
	d, err := navigateTo(dir)
	if err == nil {
		err = dir.remove(file, false)
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
	dir, file := p.Split(path.Clean(p))
	d, err := navigateTo(dir)
	if err == nil {
		err = dir.remove(file, true)
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
	olddir, oldfile := p.Split(path.Clean(oldpath))
	newdir, newfile := p.Split(path.Clean(newpath))
	oldd, err := navigateTo(olddir)
	var (
		newd *directory
		f    FileInfo
	)
	if err == nil {
		newd, err = navigateTo(newdir)
		f, err := oldd.get(oldfile)
		if err == nil {
			type i interface {
				move(string, *directory) error
			}
			err = f.(i).move(newfile, newdir)
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