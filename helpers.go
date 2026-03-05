package virfs

import (
	"strings"
	"fmt"
	"path/filepath"
)
func _(){fmt.Print()}

func Is_valid_path(path string) bool {
	var e error
	foo := filepath.Clean(path)
	foo, e = filepath.Abs(path)
	if e != nil { return false }
	foo = filepath.ToSlash(foo)
	return foo[0] == '/'
}

func Get_basepath(path string) (string, error) {
	p, e := Resolve_path(path)
	if e != nil { return "", e }
	return filepath.Dir(p), nil
}

func Get_name(path string) string {
	p, _ := Resolve_path(path)
	//split := strings.Split(p, "/") 
	return filepath.Base(p) //split[len(split)-1]
}

func (fs Fs) goto_path(path string) (*Dir, error) {
	if !Is_valid_path(path) { return nil, InvalidPath }

	p, e := Resolve_path(path)
	if e != nil { return nil, e }

	base, e := Get_basepath(p)
	if len(base) == 0 { return &fs.Root, nil }
	if e != nil { return nil, e }
	if base[0] == '/' { base = base[1:] }
	path_split := strings.Split(base, "/")
	if len(path_split[0]) == 0 {
		path_split = path_split[1:]
	}

	current := &(fs.Root)
	for _, d := range path_split {
		if !current.Contains(d) { return nil, DirNotExist }
		current = (*current).Content[d].Dir
	}
	return current, nil
}

func (d Dir) Contains(name string) bool {
	for n := range d.Content {
		if name == n { return true }
	}
	return false
}

func (fs Fs) Is_root(path string) bool {
	p, _ := Resolve_path(path)
	return p == filepath.Dir(p)
}

func Resolve_path(path string) (string, error) {
	var e error
	foo := filepath.Clean(path)
	foo, e = filepath.Abs(path)
	if e != nil { return "", e }
	foo = filepath.ToSlash(foo)
	if foo[0] != '/' { return "", InvalidPath }
	return foo, nil
}
