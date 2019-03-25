package main

import (
	"fmt"
	"log"
	"os"

	"github.com/ghodss/yaml"
)

func main() {
	if len(os.Args) <= 1 {
		log.Fatalln("usage: ./lddtree [file]")
	}
	root := New(os.Args[1])
	root.Resolve()

	b, err := yaml.Marshal(root)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Print(string(b))
}
