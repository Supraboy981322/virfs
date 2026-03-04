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
	fmt.Printf("\n%#v\n", fs.Root.Content)
	if e := fs.Mkdir("/usr"); e != nil && e != virfs.FileExists {
		panic(e)
	}
	fmt.Printf("\n%#v\n", fs.Root.Content["usr"])
	if e := fs.MkFile("/usr/foo", []byte("bar")); e != nil {
		panic(e)
	}
	fmt.Printf("\n%#v\n", fs.Root.Content["usr"].Dir.Content["foo"].File)
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
