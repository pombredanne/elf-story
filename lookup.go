// TODO:
// respect $LD_PRELOADED_PATH

package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/rai-project/ldcache"
)

var ldCache *ldcache.LDCache

func init() {
	// initialize global var
	var err error
	// avoid := shadowing here
	ldCache, err = ldcache.Open()
	if err != nil {
		log.Fatalln(err)
	}
}

func Lookup(name string) (string, error) {
	str, err := ldcacheLookup(name)
	if err == nil {
		return str, nil
	}

	str, err = relpathLookup(name)
	if err == nil {
		return str, nil
	}

	str, err = pathLookup(name)
	if err == nil {
		return str, nil
	}

	return name, err
}

func pathLookup(name string) (string, error) {
	if strings.Contains(name, "/") {
		return "", fmt.Errorf("file name should not contain /")
	}
	path, ok := os.LookupEnv("PATH")
	if !ok {
		return "", fmt.Errorf("PATH is empty")
	}
	paths := strings.Split(path, ":")
	mode := os.FileMode(0100)
	for _, path := range paths {
		fullpath, err := dirLookup(path, name, mode)
		if err == nil && name == filepath.Base(fullpath) {
			return fullpath, nil
		}
	}
	return "", fmt.Errorf("%s not found in PATH", name)
}

func dirLookup(path, name string, mode os.FileMode) (string, error) {
	fullpath := filepath.Join(path, name)
	fi, err := os.Stat(fullpath)
	if err != nil {
		return "", err
	}
	if fi.IsDir() {
		return "", fmt.Errorf("not a file: %s", fullpath)
	}
	if (fi.Mode()&mode)>>6 != 1 {
		return "", fmt.Errorf("not executable: %s", fullpath)
	}
	return fullpath, nil
}

func relpathLookup(name string) (string, error) {
	_, err := os.Stat(name)
	if err == nil {
		return name, nil
	}
	if os.IsNotExist(err) {
		return "", err
	}
	return "", err
}

func ldcacheLookup(name string) (string, error) {
	_, result64 := ldCache.Lookup(name)
	if len(result64) == 0 {
		return "", fmt.Errorf("Not found: %s", name)
	}
	return result64[0], nil
}
