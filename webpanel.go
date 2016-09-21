package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"strings"

	"github.com/hoisie/web"
)

// IndexPage holds basic variables for the index template
type IndexPage struct {
	Title    string
	Sitename string
	Content  template.HTML
}

var indexTemplate, _ = template.ParseFiles("tpl/index.tpl")

func initWeb() {
	info("Launching web service.")
	web.Get("/robots.txt", robots)
	web.Get("/(.*)", getIndex)
	web.Run(fmt.Sprintf("%s:%d", cfg.Main.Address, cfg.Main.Port))
}

func endWeb() {
	info("Shutting down web service.")
	web.Close()
}

var robotstxt = []byte("User-agent: *\nDisallow: /\n")

func robots(ctx *web.Context) {
	ctx.SetHeader("X-Robots-Tag", "noindex", true)
	ctx.WriteHeader(200)
	ctx.ResponseWriter.Write(robotstxt)
}

func getIndex(ctx *web.Context, arg string) {
	ctx.SetHeader("Content-type", "text/html", true)
	ctx.SetHeader("Cache-Control", "no-cache", true)
	title := cfg.Main.Sitename
	name := cfg.Main.Sitename
	// If there's an argument we try to open that as a file
	message := ""
	windows := "<h2>Windows</h2>\n<ul>\n"
	p("Real IP: %s\n", ctx.Request.Header.Get("X-Real-IP"))
	if arg != "" {
		fn := sane(arg)
		s, _ := ioutil.ReadFile("versions/" + fn)
		data := string(s)
		if data != "" {
			lines := strings.Split(data, "\n")
			count := len(lines) - 1
			l := 0
			for l < count {
				if len(fn) > 4 && fn[len(fn)-4:] == ".win" {
					message += "<p>" + lines[l] + ": "
					l++
					message += lines[l] + "\n<br />\n"
					l++
					message += "<a href=\"" + lines[l] + "\">64-bit</a>\n"
					l++
					message += "<a href=\"" + lines[l] + "\">32-bit</a>\n<p>\n"
					l++
				} else {
					message += "<p>" + lines[l] + ": "
					l++
					message += lines[l] + "\n<br />\n"
					l++
					message += "<a href=\"" + lines[l] + "\">" + lines[l] + "</a><br />\n"
					l++
					message += "MD5: " + lines[l] + "</p>"
					l++
				}
			}
		} else {
			message = "<p>Unknown version (or not updated - nudge the admin!)</p>"
		}
	} else {
		files, err := ioutil.ReadDir("versions")
		if err == nil {
			message = "<h2>OS X</h2>\n<ul>\n"
			for _, fi := range files {
				if !fi.IsDir() {
					if fi.Name()[len(fi.Name())-4:] == ".win" {
						lt := fi.Name()[:len(fi.Name())-4]
						windows += "<li><a href=\"" + fi.Name() + "\">" + lt + "</a></li>\n"
					} else {
						message += "<li><a href=\"" + fi.Name() + "\">" + fi.Name() + "</a></li>\n"
					}
				}
			}
			message += "</ul>\n"
			message += windows + "</ul>\n"
		}
	}
	msg := template.HTML(message)
	indexTemplate.Execute(ctx.ResponseWriter, &IndexPage{
		Title:    title,
		Sitename: name,
		Content:  msg,
	})
}
