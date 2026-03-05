package main

import (
	"fmt"
	"github.com/Supraboy981322/virfs"
)


func foo() { fmt.Print() }

func main() {
	unix_test()
}

func unix_test() {
	var fs = virfs.Init_UNIX()
	var e error
	fmt.Printf("\n%#v\n", fs.Root.Content)

	task("attempt to create dir that already exists (/usr)")
	e = fs.Mkdir("/usr")
	if e != virfs.FileExists {
		failed("%v", e)
	} else { passed("file already exists") }

	task("create file /usr/foo")
	e = fs.MkFile("/usr/foo", []byte("bar"))
	if e != nil && fs.Root.Content["usr"].Dir.Contains("foo") {
		failed("%v", e)
	} else { passed("created file in subdir to root") }

	task("verify /usr/foo contents")
	if string(fs.Root.Content["usr"].Dir.Content["foo"].File.Content) != "bar" {
		failed("file contents did not match")
	} else { passed("file contents matched") }

	task("attempt to replace root dir with file")
	e = fs.MkFile("/", []byte("foo"))
	if e != virfs.FileExists && !fs.Root.Contains("/") && !fs.Root.Contains("") {
		failed("successfully (bad) repaced root dir with file (err: %v)", e)
	} else { passed("prevented replacing root dir with file") }

	task("make file /bar")
	e = fs.MkFile("/bar", []byte("foo"))
	if e != nil && fs.Root.Contains("bar") {
		failed("%v", e)
	} else { passed("created file in root dir") }

	task("verify /bar contents")
	if string(fs.Root.Content["bar"].File.Content) != "foo" {
		failed("file contents did not match")
	} else { passed("file contents matched") }

	task("attempt to create file that already exists (/bar)")
	e = fs.MkFile("/bar", []byte("bar"))
	if e == nil || !fs.Root.Contains("bar") {
		failed("successfully (bad) created file that should already exist")
	} else { passed("disallowed created file that already exists") }
	
	task("create dir /usr/keeper")
	e = fs.Mkdir("/usr/keeper")
	if e != nil || !fs.Root.Content["usr"].Dir.Contains("keeper") {
		failed("failed to  that should already exist")
	} else { passed("disallowed created file that already exists") }

	task("delete /bar")
	e = fs.RmFile("/bar", false)
	if e != nil {
		failed("failed to remove file from root dir (/bar)")
	} else { passed("removed file file root dir") }

	task("delete /usr/keeper")
	e = fs.RmDir("/usr/keeper", false)
	if e != nil {
		failed("failed to delete dir")
	} else { passed("removed dir") }

	if len(fs.Root.Content["usr"].Dir.Content) < 1 {
		panic("rearange tests, expected the /usr dir to not be empty but was")
	}
	task("attempt to delete /usr without force (has contents)")
	e = fs.RmDir("/usr", false)
	if e == nil {
		failed("seccessfully (bad) deleted non-empty dir that has contents")
	} else { passed("prevented removal of non-empty dir") }
}

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

func passed(msg string, a ...any) {
	fmt.Printf("[\x1b[32mGOOD\x1b[0m]:  " + msg + "\n", a...)
}
func task(msg string, a ...any) {
	fmt.Printf("\n[\x1b[36mtask\x1b[0m]:  " + msg + "\n", a...)
}
func failed(msg string, a ...any) {
	fmt.Printf("[\x1b[31;1mFAILED\x1b[0m]:  " + msg + "\n", a...)
}
