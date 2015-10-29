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

func (fs *filesystem) getNode(p string, followSymLinks bool) (node, error) {
	fs.RLock()
	d := fs.cwd
	fs.RUnlock()
	return fs.getNodeWithCwd(p, followSymLinks, d)
}

func (fs *filesystem) getNodeWithCwd(p string, followSymLinks bool, cwd *breadcrumbs) (node, error) {
	dir, file := path.Split(p)
	d, err := fs.getDirectoryWithCwd(dir, cwd)
	if err != nil {
		return nil, err
	}
	f, err := d.get(file)
	if err != nil {
		return nil, err
	}
	s, ok := f.(*symlink)
	if ok && followSymLinks {
		return fs.getNodeWithCwd(s.link, true, d)
	}
	return f, nil
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
	n, err := fs.getNode(path.Clean(p), true)
	if err != nil {
		return err
	}
	n.SetMode(mode)
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
	n, err := fs.getNode(path.Clean(p), true)
	if err != nil {
		return err
	}
	n.SetModTime(mtime)
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
	return nil
}

func Mkdir(p string, fileMode os.FileMode) error {
}

func MkdirAll(p string, fileMode os.FileMode) error {
}

func NewSyscallError(_, string, _ error) error {
	return nil
}

func Readlink(name string) (string, error) {
	n, err := fs.getNode(path.Clean(p), false)
	if err != nil {
		return err
	}
	s, ok := n.(*symlink)
	if !ok {
		return "", &Path{
			Op:   "readlink",
			Path: name,
			Err:  ErrInvalid,
		}
	}
	if s.Mode&0444 == 0 {
		return "", &Path{
			Op:   "readlink",
			Path: name,
			Err:  ErrPermission,
		}
	}
	return s.link, nil
}

func Remove(name string) error {
	return nil
}

func RemoveAll(name string) error {
	return nil
}

func Rename(oldpath, newpath string) error {
	return nil
}

func SameFile(f, g os.FileInfo) bool {
	return f == g
}

func Setenv(key, value string) error {
	return ErrUnsupported
}

func Symlink(oldname, newname string) error {
	dir, file := path.Split(newname)
	d, err := fs.getDirectory(dir)
	if err != nil {
		return &LinkError{
			Op:  "symlink",
			Old: oldname,
			New: newname,
			Err: err,
		}
	}
	err = d.set(file, &symlink{
		modeTime: modeTime{},
		link:     oldname,
	})
	if err != nil {
		return &LinkError{
			Op:  "symlink",
			Old: oldname,
			New: newname,
			Err: err,
		}
	}
	return nil
}

func TempDir() string {
	return "/tmp"
}

func Truncate(name string, size int64) error {
	n, err := fs.getNode(path.Clean(p), true)
	if err != nil {
		return err
	}
	f, ok := n.Data().(file)
	if !ok {
		return &Path{
			Op:   "truncate",
			Path: name,
			Err:  ErrInvalid,
		}
	}
	if n.Mode()&0222 == 0 {
		return &Path{
			Op:   "truncate",
			Path: name,
			Err:  ErrPermission,
		}
	}
	return f.Truncate(size)
}

func Unsetenv(_ string) error {
	return ErrUnsupported
}
