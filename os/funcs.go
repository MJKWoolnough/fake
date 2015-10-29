package os

import (
	"os"
	"path"
	"strings"
	"time"
)

func (fs *filesystem) getDirectory(p string) (*breadcrumbs, error) {
	fs.RLock()
	d := fs.cwd
	fs.RUnlock()
	return fs.getDirectoryWithCwd(p, d)
}

func (fs *filesystem) getDirectoryWithCwd(p string, d *breadcrumbs) (*breadcrumbs, error) {
	if len(p) == 0 {
		return d, nil
	}
	if p[0] == '/' {
		d = fs.root
		p = p[1:]
	}
	parts := strings.Split(p, "/")
	for len(parts) > 0 && parts[0] == ".." {
		d = d.previous
	}
	for len(parts) > 0 {
		dir := parts[0]
		dir = dir[1:]
		fi, err := d.get(dir)
		if err != nil {
			return nil, err
		}
		switch f := fi.(type) {
		case *directory:
			d = &breadcrumbs{
				name:      dir,
				depth:     d.depth + 1,
				previous:  d,
				parent:    d,
				directory: f,
			}
		case *symlink:
			l, err := fs.getDirectoryWithCwd(f.link, d)
			if err != nil {
				return nil, err
			}
			d = &breadcrumbs{
				name:      dir,
				depth:     d.depth + 1,
				previous:  d,
				parent:    l.parent,
				directory: l.directory,
			}
		default:
			return nil, ErrIsNotDir
		}
		if d.FileMode&0111 > 0 {
			return nil, ErrPermission
		}
	}
	return d, nil
}

func (fs *filesystem) getNode(p string) (node, error) {
	dir, file := path.Split(p)
	d, err := fs.getDirectory(dir)
	if err != nil {
		return nil, err
	}
	return d.get(file)
}

func Chdir(p string) error {
	d, err := fs.getDirectory(path.Clean(p))
	if err != nil {
		return err
	}
	fs.Lock()
	fs.cwd = d
	if p[0] == '/' {
		fs.cwdPath = p
	} else {
		fs.cwdPath = path.Join(fs.cwdPath, p)
	}
	fs.Unlock()
	return nil
}

func Chmod(p string, mode os.FileMode) error {
	n, err := fs.getNode(path.Clean(p))
	if err != nil {
		return err
	}
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
	return nil
}

func Clearenv() {

}

func Environ() []string {
	return []string{}
}

func Exit(code int) {
	os.Exit(code)
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
	fs.RLock()
	cwd := fs.cwd
	fs.RUnlock()
	parts := make([]string, cwd.depth)
	for i := cwd.depth - 1; i >= 0; i-- {
		parts[i] = cwd.name
		cwd = cwd.previous
	}
	return path.Join("/", parts...)
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
}

func Mkdir(p string, fileMode os.FileMode) error {
}

func MkdirAll(p string, fileMode os.FileMode) error {
}

func NewSyscallError(_, string, _ error) error {
	return nil
}

func Readlink(name string) (string, error) {
}

func Remove(name string) error {
}

func RemoveAll(name string) error {
}

func Rename(oldpath, newpath string) error {
}

func SameFile(f, g os.FileInfo) bool {
	return f == g
}

func Setenv(key, value string) error {
	return ErrUnsupported
}

func Symlink(oldname, newname string) error {
}

func TempDir() string {
	return "/tmp"
}

func Truncate(name string, size int64) error {
}

func Unsetenv(_ string) error {
	return ErrUnsupported
}
