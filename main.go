package main

import (
	"flag"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"

	"gopkg.in/gcfg.v1"
)

const PROGVERSION string = "0.2.0"

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
	var runupdates *bool = flag.Bool("u", false, "Run updates instead of launching web server.")
	var address *string = flag.String("a", "0.0.0.0", "Address to bind to.")
	var port *int = flag.Int("p", 8000, "Port to bind to.")
	var display *bool = flag.Bool("d", false, "Display versions and URLs available, then exit.")
	var version *string = flag.String("v", "", "Display a specific version. If not specified you get everything.")
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

func updateVersions() {
	p("UnityLinks %s", PROGVERSION)
	baseurl1 := "http://download.unity3d.com/download_unity/"
	//	baseurl2 := "http://beta.unity3d.com/download/"
	//	url2 := "http://beta.unity3d.com/download/b6c1a63227dc/unity-5.3.2p3-osx.ini"
	files, err := ioutil.ReadDir("./updates")
	if err != nil {
		fatal("Error: %s.", err.Error())
	}
	count := 0
	max := 0
	for _, fi := range files {
		if !fi.IsDir() {
			max++
			s, _ := ioutil.ReadFile("updates/" + fi.Name())
			if s != nil {
				hash := strings.Replace(string(s), "\n", "", -1)
				url := baseurl1 + hash + "/unity-" + fi.Name() + "-osx.ini"
				info("Downloading %s", url)
				response, err := http.Get(url)
				if err == nil {
					defer response.Body.Close()
					contents, err := ioutil.ReadAll(response.Body)
					var ini UnityIni
					gcfg.ReadStringInto(&ini, string(contents))
					if err == nil {
						data := ""
						data += ini.Unity.Title + "\n"
						data += baseurl1 + hash + "/" + ini.Unity.Url + "\n"
						data += ini.Unity.Md5 + "\n"

						data += ini.Documentation.Title + "\n"
						data += baseurl1 + hash + "/" + ini.Documentation.Url + "\n"
						data += ini.Documentation.Md5 + "\n"

						data += ini.Example.Title + "\n"
						data += baseurl1 + hash + "/" + ini.Example.Url + "\n"
						data += ini.Example.Md5 + "\n"

						data += ini.WebPlayer.Title + "\n"
						data += baseurl1 + hash + "/" + ini.WebPlayer.Url + "\n"
						data += ini.WebPlayer.Md5 + "\n"

						data += ini.Android.Title + "\n"
						data += baseurl1 + hash + "/" + ini.Android.Url + "\n"
						data += ini.Android.Md5 + "\n"

						data += ini.AppleTV.Title + "\n"
						data += baseurl1 + hash + "/" + ini.AppleTV.Url + "\n"
						data += ini.AppleTV.Md5 + "\n"

						data += ini.IOS.Title + "\n"
						data += baseurl1 + hash + "/" + ini.IOS.Url + "\n"
						data += ini.IOS.Md5 + "\n"

						data += ini.Linux.Title + "\n"
						data += baseurl1 + hash + "/" + ini.Linux.Url + "\n"
						data += ini.Linux.Md5 + "\n"

						data += ini.Mac.Title + "\n"
						data += baseurl1 + hash + "/" + ini.Mac.Url + "\n"
						data += ini.Mac.Md5 + "\n"

						data += ini.Samsung_TV.Title + "\n"
						data += baseurl1 + hash + "/" + ini.Samsung_TV.Url + "\n"
						data += ini.Samsung_TV.Md5 + "\n"

						data += ini.Tizen.Title + "\n"
						data += baseurl1 + hash + "/" + ini.Tizen.Url + "\n"
						data += ini.Tizen.Md5 + "\n"

						data += ini.WebGL.Title + "\n"
						data += baseurl1 + hash + "/" + ini.WebGL.Url + "\n"
						data += ini.WebGL.Md5 + "\n"

						data += ini.Windows.Title + "\n"
						data += baseurl1 + hash + "/" + ini.Windows.Url + "\n"
						data += ini.Windows.Md5 + "\n"

						saveVersion(fi.Name(), data)
						count++
					}
				}
			}
		}
	}
	info("Updated %d versions out of %d.", count, max)
}

func saveVersion(version string, data string) {
	file, err := os.OpenFile("versions/"+version, os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		fatal("Error: %s.", err.Error())
	}

	defer file.Close()
	file.WriteString(data)
	os.Remove("updates/" + version)
}

func displayVersions() {
	files, err := ioutil.ReadDir("versions")
	if err != nil {
		fatal("Error: %s.", err.Error())
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
