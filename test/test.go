// TODO: *more* thorough testing

package main

import (
	"fmt"
	"github.com/Supraboy981322/virfs"
)
//so I can just willy-nilly insert a print statement without importing again 
func _(){fmt.Print()}

func main() {
	unix_test()
}

func unix_test() {
	//initialize a UNIX root dir
	var fs = virfs.Init_UNIX()
	var e error

	//print the starting fs
	fmt.Printf("\n%#v\n", fs.Root.Content)

	//attempt to create a dir that should already exist
	task("attempt to create dir that already exists (/usr)")
	e = fs.Mkdir("/usr")
	if e != virfs.FileExists {
		failed("%v", e)
	} else {
		passed("file already exists")
	}

	//create a file in a sub dir 
	task("create file /usr/foo")
	e = fs.MkFile("/usr/foo", []byte("bar"))
	if e != nil && fs.Root.Content["usr"].Dir.Contains("foo") {
		failed("%v", e)
	} else {
		passed("created file in subdir to root")
}

	//verify file contents match expected
	task("verify /usr/foo contents")
	if string(fs.Root.Content["usr"].Dir.Content["foo"].File.Content) != "bar" {
		failed("file contents did not match")
	} else {
		passed("file contents matched")
	}

	//attempt to replace root dir with a file
	task("attempt to replace root dir with file")
	e = fs.MkFile("/", []byte("foo"))
	if e != virfs.FileExists && !fs.Root.Contains("/") && !fs.Root.Contains("") {
		failed("successfully (bad) repaced root dir with file (err: %v)", e)
	} else {
		passed("prevented replacing root dir with file")
	}

	//create a file in fs root
	task("make file /bar")
	e = fs.MkFile("/bar", []byte("foo"))
	if e != nil && fs.Root.Contains("bar") {
		failed("%v", e)
	} else {
		passed("created file in root dir")
	}

	//verify the file contents match expected
	task("verify /bar contents")
	if string(fs.Root.Content["bar"].File.Content) != "foo" {
		failed("file contents did not match")
	} else {
		passed("file contents matched")
	}

	//attempt to create a file that was already created
	task("attempt to create file that already exists (/bar)")
	e = fs.MkFile("/bar", []byte("bar"))
	if e == nil || !fs.Root.Contains("bar") {
		failed("successfully (bad) created file that should already exist")
	} else {
		passed("disallowed created file that already exists")
	}
	
	//create a subdir
	task("create dir /usr/keeper")
	e = fs.Mkdir("/usr/keeper")
	if e != nil || !fs.Root.Content["usr"].Dir.Contains("keeper") {
		failed("failed to  that should already exist")
	} else {
		passed("disallowed created file that already exists")
	}

	//remove file from fs root
	task("delete /bar")
	e = fs.RmFile("/bar", false)
	if e != nil {
		failed("failed to remove file from root dir (/bar)")
	} else {
		passed("removed file file root dir")
	}

	//remove dir
	task("delete /usr/keeper")
	e = fs.RmDir("/usr/keeper", false)
	if e != nil {
		failed("failed to delete dir")
	} else {
		passed("removed dir")
	}

	//sanity check for dir that was changed several times (make sure not empty) 
	if len(fs.Root.Content["usr"].Dir.Content) < 1 {
		panic("rearange tests, expected the /usr dir to not be empty but was")
	}

	//attempt to remove non-empty dir without -f
	task("attempt to delete /usr without force (has contents)")
	e = fs.RmDir("/usr", false)
	if e == nil {
		failed("seccessfully (bad) deleted non-empty dir that has contents")
	} else {
		passed("prevented removal of non-empty dir")
	}

	//attempt to delete fs root
	task("attempting to delete filesystem root")
	previous_len := len(fs.Root.Content)
	e = fs.RmDir("/", true) 
	if e != virfs.PermissionDenied || len(fs.Root.Content) != previous_len {
		//print different err depending on context
		if e != virfs.PermissionDenied {
			failed("failed to prevent root dir deletion (err: %v)", e)
		} else {
			failed("size changed (before{%d} ; after{%d})", previous_len, len(fs.Root.Content))
		}
	} else {
		passed("prevented deleting root dir")
	}
}

// TODO: mirror UNIX fs test here
//  (should be the same, but just to be sure)
func empty_fs() {
	var fs = virfs.Init()
	if e := fs.Mkdir("/usr"); e != nil {
		panic(e)
	}
	fmt.Printf("\n%#v\n", fs.Root.Content["usr"])
	if e := fs.MkFile("/usr/foo", []byte("bar")); e != nil {
		panic(e)
	}
	fmt.Printf("\n%#v\n", fs.Root.Content["usr"].Dir.Content["foo"].File)
}
