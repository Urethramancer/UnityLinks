package main

import (
	"log"
)

func info(f string, v ...interface{}) {
	log.Printf(f, v...)
}
