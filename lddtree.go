package main

import (
	"fmt"
	"log"
	"os"

	"github.com/ghodss/yaml"
)

func main() {
	root := New(os.Args[1])
	root.Resolve()

	b, err := yaml.Marshal(root)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Print(string(b))
}
