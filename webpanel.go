package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"strings"
	"sync"

	"github.com/hoisie/web"
)

type IndexPage struct {
	Title    string
	Sitename string
	Content  template.HTML
}

var indexTemplate, indexerr = template.ParseFiles("tpl/index.tpl")
var share sync.Mutex

type DevNull struct{}

func (*DevNull) Write(p []byte) (int, error) {
	return len(p), nil
}

func initWeb() {
	info("Launching web service.")
	//	web.SetLogger(log.New(&DevNull{}, "", 0))
	web.Get("/(.*)", get_index)
	web.Run(fmt.Sprintf("%s:%d", cfg.Main.Address, cfg.Main.Port))
}

func endWeb() {
	info("Shutting down web service.")
	web.Close()
}

func get_index(ctx *web.Context, arg string) {
	ctx.SetHeader("Content-type", "text/html", true)
	ctx.SetHeader("Cache-Control", "no-cache", true)
	title := cfg.Main.Sitename
	name := cfg.Main.Sitename
	// If there's an argument we try to get that version's links
	message := ""
	if arg != "" {
		fn := sane(arg)
		s, _ := ioutil.ReadFile("versions/" + fn)
		data := string(s)
		if data != "" {
			lines := strings.Split(data, "\n")
			count := len(lines) - 1
			l := 0
			for l < count {
				message += "<p>" + lines[l] + ": "
				l++
				message += "<a href=\"" + lines[l] + "\">" + lines[l] + "</a><br />"
				l++
				message += "MD5: " + lines[l] + "</p>"
				l++
			}
		} else {
			message = "<p>Unknown version (or not updated - nudge the admin!)</p>"
		}
	} else {
		files, err := ioutil.ReadDir("versions")
		if err == nil {
			for _, fi := range files {
				if !fi.IsDir() {
					message += "<li><a href=\"" + fi.Name() + "\">" + fi.Name() + "</a></li>"
				}
			}
		}
	}
	msg := template.HTML(message)
	indexTemplate.Execute(ctx.ResponseWriter, &IndexPage{
		Title:    title,
		Sitename: name,
		Content:  msg,
	})
}
