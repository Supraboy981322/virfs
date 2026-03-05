package virfs

import (
	"io"
	"fmt"
	"bytes"
	"errors"
)

//so I can just willy-nilly insert a print statement without importing again 
func _(){fmt.Print()}

//types
type (
	//why no dedicated enum type?
	Entry_type int

	//directory
	Dir struct {
		Content map[string]Entry
		Name string
		//points to self if root dir
		Parent *Dir

		// TODO:
		//  - properties
		//  - permissions
		//  - size
		//  - probably some other stuff
	}

	//file
	File struct {
		Content []byte
		Reader io.Reader
		Writer io.Writer
		Size uint

		// TODO:
		//  - properties
		//  - permissions
		//  - probably some other stuff
	}

	//item in a dir
	Entry struct {
		Entry_type Entry_type
		//name of file/dir (not path)
		Name string
		//if file, will be nil
		Dir *Dir
		//if dir, will be nil
		File *File
	}

	Fs struct {
		//root dir ('/')
		Root Dir
		// TODO:
		//  - properties
		//  - settings
		//  - size
		//  - probably some other stuff
	}
)

//why no dedicated enum type?
const (
	File_entry Entry_type = iota
	Dir_entry
	Symlink_entry
	Fifo_entry
	Socket_entry
)

//errors as values sure is a great idea
var (
	EmptyPath = errors.New("cannot use empty path")
	DirNotExist = errors.New("path does not exist")
	FileExists = errors.New("file exists")
	InvalidPath = errors.New("malformed path")
	Type_mismatch = errors.New("missmatched entry type")
	DirNotEmpty = errors.New("directory is not empty")
	FileNotExist = errors.New("file does not exist")
	PermissionDenied = errors.New("permission denied")
)

//initialize empty fs
func Init() Fs {
	fs := Fs {
		Root: Dir {
			Content: map[string]Entry{},
			Name: "/",
		},
	}

	//make root parent point to itself
	fs.Root.Parent = &fs.Root

	return fs
}

//initialize with UNIX dirs
func Init_UNIX() Fs {
	fs := Fs {
		Root: Unix_root_dir(),
	}

	//make root parent point to itself
	fs.Root.Parent = &fs.Root

	return fs
}

//make dir (takes absolute path only)
func (fs Fs) Mkdir(path string) error {
	//basic validity checks 
	if len(path) < 1 { return EmptyPath }
	if path[0] != '/' { return InvalidPath }

	//attempt to traverse to the parent of input path
	current, e := fs.goto_path(path)
	if e != nil { return e } 

	//get the name of the dir
	target := Get_name(path)
	//if already exists or is root, err
	if current.Contains(target) || fs.Is_root(path) { return FileExists }

	//add the entry to the path's parent dir 
	(*current).Content[target] = Entry {
		Entry_type: Dir_entry,
		Dir: &Dir{
			Name: target,
			Content: map[string]Entry{},
		},
		File: nil,
		Name: target,
	}

	return nil
}

//make file (takes absolute path only and content)
func (fs Fs) MkFile(path string, content []byte) error {
	//validate path
	if !Is_valid_path(path) { return InvalidPath }

	//traverse to path's parent
	current, e := fs.goto_path(path)
	if e != nil { return e }

	//get the name of the file
	target := Get_name(path)
	//make sure the file isn't already present and isn't root dir
	if current.Contains(target) || fs.Is_root(path) { return FileExists }
	//make sure the target isn't empty (can occur if trailing slash in path)
	if len(target) == 0 { return InvalidPath }

	buf := bytes.NewBuffer(content)

	//create the entry
	(*current).Content[target] = Entry { 
		Entry_type: File_entry,
		Dir: nil,
		File: &File{
			Content: content,
			Reader: buf,
			Writer: buf,
			// TODO: change this to prevent overflow
			Size: uint(len(content)),
		},
		Name: target, 
	}
	
	return nil
}

//remove a dir (takes absolute path only, set 'force' to true to
//  delete non-empty dir)
func (fs Fs) RmDir(path string, force bool) error {
	//traverse to the path's parent dir
	p, e := fs.goto_path(path)
	if e != nil { return e }

	//get the dir name
	name := Get_name(path)

	//make sure it's not the root dir
	if fs.Is_root(name) { return PermissionDenied }

	//make sure it's present
	if !p.Contains(name) { return DirNotExist }

	//err if not a dir
	if p.Content[name].Entry_type != Dir_entry { return Type_mismatch }

	//make sure the directory isn't empty if no force
	if !force && len(p.Content[name].Dir.Content) != 0 { return DirNotEmpty }

	//remove from fs
	delete(p.Content, name)

	return nil
}

//delete a file (takes absolute path only, set 'recurse' to true
//  if path may be dir)
func (fs Fs) RmFile(path string, recurse bool) error {
	p, e := fs.goto_path(path)
	if e != nil { return e }
	name := Get_name(path)
	if !p.Contains(name) { return FileNotExist }
	if p.Content[name].Entry_type != File_entry {
		if recurse { return fs.RmDir(path, true) }
		return Type_mismatch
	}
	delete(p.Content, name)
	return nil
}

//returns dir that mimiks basic UNIX root dir
func Unix_root_dir() Dir {
	//empty root dir
	var root = Dir {
		Content: map[string]Entry{},
		Name: "/",
	}

  //the most basics of basic ('/home' was a later addition to Unix)
	paths := []string {
		"usr", "etc", "var", "bin", "lib", "media", "mnt", "tmp",
	}

	//add each dir to fs root
	for _, p := range paths {
		root.Content[p] = Entry {
			Entry_type: Dir_entry,
			Dir: &Dir {
				Content: map[string]Entry{},
				Name: p,
				Parent: &root,
			},
			File: nil,
			Name: p,
		}
	}

	return root
}
