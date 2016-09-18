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

// TODO: Find some criteria for cleanup of old versions.
/*func initJanitor(s *Scrap) {
	delay := s.Expiry - time.Now().Unix()
	ticker := time.NewTicker(time.Duration(delay) * time.Second)
	go func() {
		for {
			select {
			case <-ticker.C:
				ticker.Stop()
				info("Janitor: Deleting %s", s.Hash)
				s.delete()
				return
			}
		}
	}()
}*/

func fexists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
