package main

import (
	"log"
)

func info(f string, v ...interface{}) {
	log.Printf(f, v...)
}

func fatal(f string, v ...interface{}) {
	log.Fatalf(f, v...)
}
