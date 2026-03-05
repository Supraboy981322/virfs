package virfs

import (
	"fmt"
	"errors"
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
		Content map[string]Entry
		Name string
		Parent *Dir
	}

	File struct {
		Content []byte
	}

	Entry struct {
		Entry_type Entry_type
		Dir *Dir //if file, will be nil
		File *File //if dir, will be nil
		Name string
	}
)

type Fs struct {
	Root Dir
}

var (
	EmptyPath = errors.New("cannot use empty path")
	DirNotExist = errors.New("path does not exist")
	FileExists = errors.New("file exists")
	InvalidPath = errors.New("malformed path")
	Type_mismatch = errors.New("missmatched entry type")
	DirNotEmpty = errors.New("directory is not empty")
	FileNotExist = errors.New("file does not exist")
)

func Init() Fs {
	fs := Fs {
		Root: Dir {
			Content: map[string]Entry{},
			Name: "/",
		},
	}
	fs.Root.Parent = &fs.Root
	return fs
}

func Init_UNIX() Fs {
	fs := Fs {
		Root: Unix_root_dir(),
	}
	fs.Root.Parent = &fs.Root
	return fs
}

func (fs Fs) Mkdir(path string) error {
	if len(path) < 1 { return EmptyPath }
	if path[0] != '/' { return InvalidPath }

	current, e := fs.goto_path(path)
	if e != nil { return e } 

	target := Get_name(path)
	if current.Contains(target) || fs.Is_root(path) { return FileExists }

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

func (fs Fs) MkFile(path string, content []byte) error {
	if !Is_valid_path(path) { return InvalidPath }

	current, e := fs.goto_path(path)
	if e != nil { return e }

	target := Get_name(path)
	if current.Contains(target) || fs.Is_root(path) { return FileExists }
	if len(target) == 0 { return InvalidPath }

	file := File{
		Content: content,
	}

	(*current).Content[target] = Entry { 
		Entry_type: File_entry,
		Dir: nil,
		File: &file,
		Name: target, 
	}
	
	return nil
}

//so I can just willy-nilly insert a print statement without importing again 
func _(){fmt.Print()}

func Unix_root_dir() Dir {
	var root = Dir {
		Content: map[string]Entry{},
		Name: "/",
	}
	paths := []string {
		"usr", "etc", "var", "bin", "lib", "media", "mnt", "tmp",
	}
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

func (fs Fs) RmDir(path string, force bool) error {
	p, e := fs.goto_path(path)
	if e != nil { return e }
	name := Get_name(path)
	if !p.Contains(name) { return DirNotExist }
	if p.Content[name].Entry_type != Dir_entry { return Type_mismatch }

	if len(p.Content[name].Dir.Content) != 0 && !force { return DirNotEmpty }
	delete(p.Content, name)
	return nil
}

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
