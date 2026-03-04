package main

import (
	"fmt"
	"github.com/Supraboy981322/virfs"
)

var fs = virfs.Init()

func foo() { fmt.Print() }

func main() {
	if e := fs.Mkdir("/usr"); e != nil {
		panic(e)
	}
	fmt.Printf("\n%#v\n", fs.Root.Content["usr"])
	if e := fs.MkFile("/usr/foo", []byte("bar")); e != nil {
		panic(e)
	}
	fmt.Printf("\n%#v\n", fs.Root.Content["usr"].Dir.Content["foo"])
}
