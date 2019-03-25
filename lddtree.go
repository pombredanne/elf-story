package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	if len(os.Args) <= 1 {
		log.Fatalln("usage: ./lddtree [file]")
	}
	root := New(os.Args[1])
	root.Resolve()

	fmt.Println("============================= .String() =============================")
	fmt.Print(root.String())

	fmt.Println("============================= .StringIndent() =============================")
	fmt.Print(root.StringIndent("    "))

	fmt.Println("============================= .PrintIndent() =============================")
	root.PrintIndent("    ")
}
