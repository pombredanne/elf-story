package main

import "fmt"
import "github.com/ghodss/yaml"
import "os"
import "log"
import "debug/elf"
import "github.com/rai-project/ldcache"

var ld_cache *ldcache.LDCache

func init() {
	var err error
	ld_cache, err = ldcache.Open()
	if err != nil {
		log.Fatalln(err)
	}
}

func lookup1(name string) (string, error) {
	_, result := ld_cache.Lookup(name)
	if len(result) == 0 {
		return "", fmt.Errorf("Not found: %s", name)
	}
	return result[0], nil
}

type A ELF

type ELF map[string][]A

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

var mem = map[string][]string{}

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
		e[e.Key()] = []A{} // ensure not nil after resolve
		return
	}
	for _, dep := range deps {
		path, err := lookup1(dep)
		if err != nil {
			log.Println(err)
			e[e.Key()] = append(e[e.Key()], A(New(dep)))
			continue
		}

		d := New(path)
		d.Resolve()
		e[e.Key()] = append(e[e.Key()], A(d))
	}
}

func main() {
	root := New(os.Args[1])
	root.Resolve()
	b, _ := yaml.Marshal(root)
	fmt.Print(string(b))
}
