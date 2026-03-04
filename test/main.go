package main

import (
	"fmt"
	"github.com/Supraboy981322/virfs"
)

var fs = virfs.Init()

func main() {
	fmt.Printf("%#v\n", fs.Root.Content)
	if e := fs.Mkdir("/usr"); e != nil {
		panic(e)
	}
	fmt.Printf("%#v\n", fs.Root.Content)
}
