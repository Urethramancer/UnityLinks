package main

import (
	"strings"

	"github.com/PuerkitoBio/goquery"
)

const (
	betaurl     = "https://beta.unity3d.com/download/"
	downloadurl = "https://download.unity3d.com/download_unity/"
)

type Patch struct {
	Name   string
	Hash   string
	Mac    string
	MacAlt string
	Win    string
	WinAlt string
}

func GetHash(url string) string {
	var s string

	if strings.Contains(url, betaurl) {
		s = strings.TrimPrefix(url, betaurl)

	} else if strings.Contains(url, downloadurl) {
		s = strings.TrimPrefix(url, downloadurl)
	} else {
		return ""
	}

	a := strings.Split(s, "/")
	if len(a) == 0 {
		return ""
	}

	return a[0]
}

func GetPatches(url string) ([]Patch, error) {
	doc, err := goquery.NewDocument(url)
	if err != nil {
		return nil, err
	}

	var patches []Patch

	doc.Find(".patch").Each(func(i int, s *goquery.Selection) {
		title := s.Find(".release-notes").Find("h4").Text()
		title = strings.Replace(title, "Patch ", "", 1)
		p := Patch{Name: title}
		s.Find("a").Each(func(i int, s *goquery.Selection) {
			t := s.Text()
			if strings.Contains(t, "Download") {
				link, _ := s.Attr("href")
				var alt string
				if strings.Contains(link, "beta.unity3d.com") {
					alt = strings.Replace(link, "download/", "download_unity/", 1)
					alt = strings.Replace(alt, "beta.unity3d", "download.unity3d", 1)
				}
				if strings.Contains(t, "Mac") {
					p.Mac = link
					p.MacAlt = alt
				} else {
					p.Win = link
					p.WinAlt = alt
				}
			}
		})
		p.Hash = GetHash(p.Mac)
		patches = append(patches, p)
	})

	return patches, nil
}
