package virfs

import (
	"strings"
	"fmt"
	"path/filepath"
)
//so I can just willy-nilly insert a print statement without importing again 
func _(){fmt.Print()}

//helper to validate path
func Is_valid_path(path string) bool {
	var e error

	//cleanup path
	foo := filepath.Clean(path)

	//make sure it's absolute
	foo, e = filepath.Abs(path)
	if e != nil { return false }

	//ensure it's in UNIX format 
	foo = filepath.ToSlash(foo)

	//just to be absolutely sure
	return foo[0] == '/'
}

//get the dir of a path
func Get_basepath(path string) (string, error) {
	//resolve to absolute path 
	p, e := Resolve_path(path)
	if e != nil { return "", e }

	//get directory
	d := filepath.Dir(p)

	//be 100% sure it's in UNIX format
	d = filepath.ToSlash(d)
	if d[0] != '/' { panic("FAILED TO FORCE UNIX DIR") } 

	//return the directory
	return d, nil
}

//helper to get the name of entry from path 
func Get_name(path string) string {
	//resolve to abs
	p, _ := Resolve_path(path)

	//filepath package is good enough for this
	return filepath.Base(p)
}

//internal helper to traverse to a path 
func (fs Fs) goto_path(path string) (*Dir, error) {
	//validate path
	if !Is_valid_path(path) { return nil, InvalidPath }

	//resolve to abs
	p, e := Resolve_path(path)
	if e != nil { return nil, e }

	//get the parent dir
	base, e := Get_basepath(p)
	if e != nil { return nil, e }
	//if parent name is empty, assume it's the root dir and return it
	if len(base) == 0 { return &fs.Root, nil }

	//remove the leading forward slash (for splitting) 
	if base[0] == '/' { base = base[1:] }
	
	//split the path
	path_split := strings.Split(base, "/")

	//if the first entry is empty string, shift the slice
	if len(path_split[0]) == 0 {
		path_split = path_split[1:]
	}

	//start at fs root
	current := &(fs.Root)
	//iterate over split path
	for _, d := range path_split {

		//err if target entry isn't present in current dir
		if !current.Contains(d) { return nil, DirNotExist }

		//move current dir to the next entry in path
		current = (*current).Content[d].Dir
	}

	//return pointer to last dir traversed into 
	return current, nil
}

//helper to check if a dir contains an entry name
func (d Dir) Contains(name string) bool {
	for n := range d.Content {
		if name == n { return true }
	}
	return false
}

//helper to check if a path is root dir
func (fs Fs) Is_path_root(path string) bool {
	p, _ := Resolve_path(path)
	return p == filepath.Dir(p)
}
//helper to check if a path is root dir
func (fs Fs) Is_root(dir *Dir) bool {
	return dir == &fs.Root
}

//helper to resolve a path to valid (hopefully), abs path
func Resolve_path(path string) (string, error) {
	var e error

	//cleanup relative traversal
	foo := filepath.Clean(path)

	//ensure it's resolved to abs
	foo, e = filepath.Abs(path)
	if e != nil { return "", e }

	//be absolutely sure ensure it's in UNIX format
	foo = filepath.ToSlash(foo)
	if foo[0] != '/' { return "", InvalidPath }

	//return resulting path
	return foo, nil
}
