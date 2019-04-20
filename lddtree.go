package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	if len(os.Args) <= 1 {
		log.Fatal(`usage: ./lddtree [file]
[file] is searched in such order:
- /etc/ld.so.cache
- [file], interpreted as relative path
- PATH
`)
	}
	root := New(os.Args[1])
	//root.Resolve()
	fmt.Println("============================= .ResolvIndent() =============================")
	root.ResolveIndent(".")

	fmt.Println("============================= .String() =============================")
	fmt.Print(root.String())

	fmt.Println("============================= .StringIndent() =============================")
	fmt.Print(root.StringIndent("    "))

	fmt.Println("============================= .PrintIndent() =============================")
	root.PrintIndent("    ")

	fmt.Println("============================= .HTML() =============================")
	fmt.Println(root.HTML())
}
