// +build map

package main

import "fmt"

type S map[string]int

func New() S {
	var s = S{}
	return s
}

func (s S) Change() {
	//s = New()
	s["b"] = 3
}

func main() {
	var s = S{
		"a": 1,
	}

	s.Change()

	fmt.Println(s)
}
