package main

import (
	"fmt"
	"log"
	"os"
	"strings"
)

func main() {
	if len(os.Args) <= 1 {
		log.Fatalln("usage: ./lddtree [file]")
	}
	root := New(os.Args[1])
	root.Resolve()
	fmt.Print(root.String())
	fmt.Println("=============================")
	var pr func(e ELF, indent string, lvl ...int)
	pr = func(e ELF, indent string, lvl ...int) {
		prefix := strings.Repeat(indent, len(lvl))
		suffix := ""
		if e.Val() == nil {
			suffix = " [NOT FOUND]"
		}
		fmt.Printf("%s%s%s\n", prefix, e.Key(), suffix)
		for _, v := range e.Val() {
			pr(v, indent, append(lvl, 0)...)
		}
	}
	pr(root, "    ")
}
