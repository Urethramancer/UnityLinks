package main

import (
	"flag"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"
)

const (
	progname    = "UnityLinks"
	progversion = "0.5.4"
)

// Config holds the setup for the web server
type Config struct {
	Main struct {
		Address  string
		Port     int
		Sitename string
	}
}

var cfg Config

func init() {
	rand.Seed(time.Now().UnixNano())
	if !fexists("versions") {
		info("Creating data directory.")
		os.MkdirAll("versions", 0700)
	}
	if !fexists("updates") {
		info("Creating update directory.")
		os.MkdirAll("updates", 0700)
	}
}

func main() {
	var runupdates = flag.Bool("u", false, "Run updates instead of launching web server.")
	var address = flag.String("a", "0.0.0.0", "Address to bind to.")
	var port = flag.Int("p", 8000, "Port to bind to.")
	var display = flag.Bool("d", false, "Display versions and URLs available, then exit.")
	var version = flag.String("v", "", "Display a specific version. If not specified you get everything.")
	flag.Parse()

	if *runupdates {
		updateVersions()
		return
	}

	if *display {
		if *version == "" {
			displayVersions()
		} else {
			displayVersion(*version)
		}
		return
	}

	cfg.Main.Address = *address
	cfg.Main.Port = *port
	cfg.Main.Sitename = "Unity download links"
	initWeb()
	endWeb()
}

var baseurl1 = "http://download.unity3d.com/download_unity/"
var baseurl2 = "http://beta.unity3d.com/download/"

func updateVersions() {
	p("%s %s", progname, progversion)
	files, err := ioutil.ReadDir("./updates")
	if err != nil {
		p("Error: %s.", err.Error())
		return
	}
	count := 0
	max := 0
	for _, fi := range files {
		if !fi.IsDir() {
			max++
			s, _ := ioutil.ReadFile("updates/" + fi.Name())
			if s != nil {
				hash := strings.Replace(string(s), "\n", "", -1)
				getMacIni(hash, fi.Name())
				getWinIni(hash, fi.Name())
				count++
			}
		}
	}
	info("Updated %d versions out of %d.", count, max)
}

func getVar(line string) string {
	return strings.TrimSpace(line[strings.IndexByte(line, '=')+1:])
}

func getMacIni(hash string, filename string) {
	url := baseurl1 + hash + "/unity-" + filename + "-osx.ini"
	url2 := baseurl2 + hash + "/unity-" + filename + "-osx.ini"

	mybase := baseurl1
	info("Downloading %s", url)
	response, _ := http.Get(url)
	p("Status: %d\n", response.StatusCode)
	if response.StatusCode == 404 {
		p("Trying again with %s\n", url2)
		response, _ = http.Get(url2)
		if response.StatusCode == 404 {
			p("Error downloading %s\n", url2)
			return
		}
		mybase = baseurl2
	}
	defer response.Body.Close()
	bs, err := ioutil.ReadAll(response.Body)
	if err != nil {
		p("Error fetching macOS version: %s\n", err.Error())
		return
	}
	contents := string(bs)
	sections := strings.Split(contents, "[")
	data := ""
	for _, s := range sections {
		if len(s) > 0 {
			lines := strings.Split(s, "\n")
			data += getVar(lines[1]) + "\n"
			data += getVar(lines[2]) + "\n"
			data += mybase + hash + "/" + getVar(lines[3]) + "\n"
			data += getVar(lines[4]) + "\n"
		}
	}
	p("Saving %s (macOS)\n", filename)
	saveVersion(filename, data)
}

//TODO: gcfg can't read the -win.ini properly. Make custom hack.
func getWinIni(hash string, filename string) {
	url := baseurl1 + hash + "/unity-" + filename + "-win.ini"
	url2 := baseurl2 + hash + "/unity-" + filename + "-win.ini"
	mybase := baseurl1

	info("Downloading %s", url)
	response, _ := http.Get(url)
	if response.StatusCode == 404 {
		p("Trying again with %s\n", url2)
		response, _ = http.Get(url2)
		if response.StatusCode == 404 {
			p("Error downloading %s", url2)
			return
		}
		mybase = baseurl2
	}
	defer response.Body.Close()
	bs, err := ioutil.ReadAll(response.Body)
	if err != nil {
		p("Error fetching Windows version: %s\n", err.Error())
		return
	}
	contents := string(bs)
	sections := strings.Split(contents, "[")
	data := ""
	for _, s := range sections {
		if len(s) > 0 {
			lines := strings.Split(s, "\n")
			data += getVar(lines[1]) + "\n"
			data += getVar(lines[2]) + "\n"
			l3 := mybase + hash + "/" + getVar(lines[3])
			l4 := mybase + hash + "/" + getVar(lines[4])
			data += l3 + "\n" + l4 + "\n"
		}
	}
	p("Saving %s (Windows)\n", filename)
	saveVersion(filename+".win", data)
}

func saveVersion(version string, data string) {
	file, err := os.OpenFile("versions/"+version, os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		p("Error: %s.", err.Error())
		return
	}

	defer file.Close()
	file.WriteString(data)
	os.Remove("updates/" + version)
}

func displayVersions() {
	files, err := ioutil.ReadDir("versions")
	if err != nil {
		p("Error: %s.", err.Error())
		return
	}
	for _, fi := range files {
		if !fi.IsDir() {
			displayVersion(fi.Name())
		}
	}
}

func displayVersion(version string) {
	name := "versions/" + version
	s, _ := ioutil.ReadFile(name)
	p("%s\n", s)
}
