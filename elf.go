package main

import (
	"debug/elf"
	"fmt"
	"log"
	"strings"

	"github.com/ghodss/yaml"
)

var mem = map[string][]string{}

/*
   go doesn't allow self referencing types
   type ELF map[string][]ELF
   but you can use a workaround:
*/

type A ELF

type ELF map[string][]A

/*
   this is only recommended way to safely create
   a new ELF, ensuring the key exists, but not
   necessarily the value
*/
func New(key string) ELF {
	return ELF{key: nil}
}

func (e ELF) Key() string {
	var key string
	for a, _ := range e {
		key = a
		break
	}
	return key
}

func ELF2A(es []ELF) []A {
	var as []A
	if es != nil {
		as = []A{}
	}
	for _, e := range es {
		as = append(as, A(e))
	}
	return as
}

func A2ELF(as []A) []ELF {
	var es []ELF
	if as != nil {
		es = []ELF{}
	}
	for _, a := range as {
		es = append(es, ELF(a))
	}
	return es
}

func (e ELF) Val() []ELF {
	return A2ELF(e[e.Key()])
}

func (e ELF) Append(a ELF) {
	e[e.Key()] = append(ELF2A(e.Val()), A(a))
}

func (e ELF) Set(a []ELF) {
	e[e.Key()] = ELF2A(a)
}

func (e ELF) Deps() []string {
	if v, ok := mem[e.Key()]; ok {
		return v
	}
	f, err := elf.Open(e.Key())
	if err != nil {
		log.Println(err)
		return nil
	}
	defer f.Close()
	libs, err := f.ImportedLibraries()
	if err != nil {
		log.Println(err)
	}
	mem[e.Key()] = libs
	return libs // possibly nil
}

func (e ELF) Resolve() {
	deps := e.Deps()
	if len(deps) == 0 {
		e.Set([]ELF{}) // ensure not nil after resolve
		return
	}
	for _, dep := range deps {
		path, err := ldcacheLookup(dep)
		if err != nil {
			log.Println(err)
			e.Append(New(dep))
			continue
		}
		d := New(path)
		d.Resolve()
		e.Append(d)
	}
}

func (e ELF) String() string {
	b, err := yaml.Marshal(e)
	if err != nil {
		log.Fatalln(err)
	}
	return fmt.Sprintf("%s", b)
}

func (e ELF) StringIndent(indent string, lvl ...int) string {
	result := ""
	prefix := strings.Repeat(indent, len(lvl))
	suffix := ""
	if e.Val() == nil {
		suffix = " [NOT FOUND]"
	}
	result += fmt.Sprintf("%s%s%s\n", prefix, e.Key(), suffix)
	for _, v := range e.Val() {
		result += v.StringIndent(indent, append(lvl, 0)...)
	}
	return result
}
