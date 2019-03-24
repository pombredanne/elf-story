package main

// this doesn't
// type ELF = map[string][]*ELF

// this compiles
type FOO map[string][]FOO
type BAR FOO
