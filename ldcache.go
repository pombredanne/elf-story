// TODO:
// respect $LD_PRELOADED_PATH
// search $PATH
// check relative path

package main

import (
	"fmt"
	"log"
	"os"

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
	return name, err
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
