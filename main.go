package main

import (
	"fmt"
	"errors"
	"strings"
)

type Entry_type int
const (
	File_entry Entry_type = iota
	Dir_entry
	Symlink_entry
	Fifo_entry
	Socket_entry
)

type (
	Dir struct {
		content map[string]Entry
	}

	File struct {
		content []byte
	}

	Entry struct {
		entry_type Entry_type
		dir *Dir //if file, will be nil
		file *File //if dir, will be nil
		name string
	}
)

type Fs struct {
	root Dir
}

var FS = Fs {
	root: Dir {
		content: map[string]Entry{},
	},
}

var (
	EmptyPath = errors.New("cannot use empty path")
	DirNotExist = errors.New("path does not exist")
	FileExists = errors.New("file exists")
	InvalidPath = errors.New("malformed path")
)

func main() {
	fmt.Printf("%#v\n", FS.root.content)
	if e := FS.Mkdir("/usr"); e != nil {
		panic(e)
	}
	fmt.Printf("%#v\n", FS.root.content)
}

func (fs Fs) Mkdir(path string) error {
	if len(path) < 1 { return EmptyPath }
	if path[0] != '/' { return InvalidPath }

	path_split := strings.Split(path[1:], "/")
	current, e := fs.goto_path(path)
	if e != nil { return e } 

	target := path_split[len(path_split)-1]
	if current.Contains(target) { return FileExists }

	(*current).content[target] = Entry {
		entry_type: Dir_entry,
		dir: &Dir{},
		file: nil,
		name: target,
	}

	return nil
}

func Is_valid_path(path string) bool { 
	if len(path) < 1 { return false }
	if path[0] != '/' { return false }
	return true
}

func Get_basepath(path string) (string, error) {
	if !Is_valid_path(path) { return "", InvalidPath }
	if len(path) < 2 { return "/", nil }
	split := strings.Split(path[1:], "/")
	return strings.Join(split[:len(split)-1], "/"), nil
}

func (fs Fs) goto_path(path string) (*Dir, error) {
	if !Is_valid_path(path) { return nil, InvalidPath }

	base, e := Get_basepath(path)
	if e != nil { return nil, e }
	path_split := strings.Split(base, "/")
	if len(path_split) < 2 { return &fs.root, nil }
	
	current := &(fs.root)
	for _, d := range path_split {
		if !current.Contains(d) { return nil, DirNotExist }
	}
	return current, nil
}

func (fs Fs) MkFile(path string, content []byte) error {
	current := fs.goto_path
	_=current

	return nil
}

func (d Dir) Contains(name string) bool {
	for n, _ := range d.content {
		if name == n { return true }
	}
	return false
}
