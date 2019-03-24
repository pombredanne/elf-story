package main

import "fmt"
import "encoding/json"
import "os"
import "log"
import "path/filepath"
import "debug/elf"
import "github.com/rai-project/ldcache"

var ld_cache *ldcache.LDCache
var opencounter int

func init() {
	var err error
	ld_cache, err = ldcache.Open()
	if err != nil {
		panic(err)
	}
}

type ELF struct {
	Name     string
	Path     string
	Childs   []string
	Children []*ELF
}

func (e *ELF) IsLeaf() bool {
	return len(e.Children) == 0
}

func main() {
	// basename
	name := filepath.Base(os.Args[1])

	// fullpath
	path, err := filepath.Abs(os.Args[1])
	if err != nil {
		panic(err)
	}

	// lookup children
	f, err := elf.Open(path)
	if err != nil {
		panic(err)
	}
	opencounter += 1
	defer func() {
		f.Close()
		opencounter -= 1
		log.Println(opencounter)
	}()
	libs, err := f.ImportedLibraries()
	if err != nil {
		panic(err)
	}
	if libs == nil {
		libs = []string{}
	}

	root := &ELF{
		Name:     name,
		Path:     path,
		Childs:   libs,
		Children: []*ELF{},
	}

	if len(libs) != 0 {
		root.Children = resolve(root.Childs)
	}

	b, _ := json.MarshalIndent(root, "", "  ")
	fmt.Println(string(b))
}

func lookup(names []string) []string {
	_, b := ld_cache.Lookup(names...)
	return b
}

func lookup1(name string) (string, error) {
	result := lookup([]string{name})
	if len(result) == 0 {
		return "", fmt.Errorf("None: %s", name)
	}
	return result[0], nil
}

func resolve(names []string) []*ELF {
	//fmt.Println(names)
	results := []*ELF{}
	for _, name := range names {
		path, err := lookup1(name)
		if err != nil {
			log.Println(err)
			e := &ELF{
				Name: name,
			}
			results = append(results, e)
			continue
		}

		// lookup children
		f, err := elf.Open(path)
		if err != nil {
			panic(err)
		}
		opencounter += 1
		defer func() {
			f.Close()
			opencounter -= 1
			log.Println(opencounter)
		}()
		libs, err := f.ImportedLibraries()
		if err != nil {
			panic(err)
		}
		if libs == nil {
			libs = []string{}
		}

		e := &ELF{
			Name:     name,
			Path:     path,
			Childs:   libs,
			Children: []*ELF{},
		}

		if len(libs) != 0 {
			e.Children = resolve(libs)
		}

		results = append(results, e)
	}
	return results
}
