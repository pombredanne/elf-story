package main

import (
	"log"
	"os"
)

func main() {
	if len(os.Args) <= 1 {
		log.Fatalln("usage: ./lddtree [file]")
	}
	New(os.Args[1]).ResolveIndent("0")
}
