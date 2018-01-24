package main

import (
	"fmt"
	"os"

	"github.com/mgutz/str"
)

func p(f string, v ...interface{}) {
	fmt.Printf(f, v...)
}

func sane(s string) string {
	return str.Clean(s)
}

func fexists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
