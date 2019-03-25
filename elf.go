package main

import (
	"debug/elf"
	"fmt"
	"log"

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

func (e ELF) Val() []A {
	return e[e.Key()]
}

func (e ELF) Append(a A) {
	e[e.Key()] = append(e.Val(), a)
}

func (e ELF) Set(a []A) {
	e[e.Key()] = a
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
		e.Set([]A{}) // ensure not nil after resolve
		return
	}
	for _, dep := range deps {
		path, err := ldcacheLookup(dep)
		if err != nil {
			log.Println(err)
			e.Append(A(New(dep)))
			continue
		}
		d := New(path)
		d.Resolve()
		e.Append(A(d))
	}
}

func (e ELF) String() string {
	b, err := yaml.Marshal(e)
	if err != nil {
		log.Fatalln(err)
	}
	return fmt.Sprintf("%s", b)
}
