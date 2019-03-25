// TODO:
// respect $LD_PRELOADED_PATH
// search $PATH
// check relative path

package main

import (
	"fmt"
	"log"

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

func ldcacheLookup(name string) (string, error) {
	_, result64 := ldCache.Lookup(name)
	if len(result64) == 0 {
		return "", fmt.Errorf("Not found: %s", name)
	}
	return result64[0], nil
}
