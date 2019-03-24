package main

import "fmt"
import "github.com/ghodss/yaml"
import "os"
import "log"
import "path/filepath"
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

type ELF struct {
	Path     string          `json:"file"`
	Children map[string]*ELF `json:"kids"`
}

func (e *ELF) Deps() []string {
	f, err := elf.Open(e.Path)
	if err != nil {
		log.Println(err)
		return nil
	}
	defer f.Close()
	libs, err := f.ImportedLibraries()
	if err != nil {
		log.Println(err)
	}
	return libs // possibly nil
}

func main() {
	path, err := filepath.Abs(os.Args[1])
	if err != nil {
		log.Fatalln(err)
	}

	root := &ELF{
		Path: path,
	}

	root.Resolve()

	b, _ := yaml.Marshal(root)
	fmt.Print(string(b))
}

func lookup1(name string) (string, error) {
	_, result := ld_cache.Lookup(name)
	if len(result) == 0 {
		return "", fmt.Errorf("None: %s", name)
	}
	return result[0], nil
}

func (e *ELF) Resolve() {
	deps := e.Deps()
	if len(deps) == 0 {
		return
	}
	for _, dep := range deps {
		path, err := lookup1(dep)
		if err != nil {
			log.Println(err)
			d := &ELF{
				Path:     dep,
				Children: nil,
			}
			e.Children = append(e.Children, d)
			continue
		}

		d := &ELF{
			Path:     path,
			Children: []*ELF{},
		}

		d.Resolve()

		e.Children = append(e.Children, d)
	}
}
