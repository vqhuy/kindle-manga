package main

import (
	"encoding/json"
	"fmt"
	"html"
	"log"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/gocolly/colly"
)

// truyentranhtuan.com specific collector
func collectTruyenTranhTuan(url, outputDir string) {
	nameInd := 0

	c := colly.NewCollector()
	re := regexp.MustCompile(`(?m)var slides_page_url_path = \s*(.*).$`)

	c.OnHTML(`script`, func(e *colly.HTMLElement) {
		dom, _ := e.DOM.Html()
		if re.MatchString(dom) {
			urls := re.FindStringSubmatch(dom)[1]
			urls = html.UnescapeString(urls)
			var images []string
			if err := json.Unmarshal([]byte(urls), &images); err != nil {
				log.Println(err)
				return
			}
			for _, u := range images {
				c.Visit(u)
			}
		}
	})

	c.OnRequest(func(r *colly.Request) {
		log.Println("Visiting", r.URL.String())
	})

	c.OnResponse(func(r *colly.Response) {
		name := fmt.Sprintf("%03d%s", nameInd, filepath.Ext(r.FileName()))
		nameInd++
		if strings.Index(r.Headers.Get("Content-Type"), "image") > -1 {
			if err := r.Save(filepath.Join(outputDir, name)); err != nil {
				log.Println(err)
			}
		} else if strings.Index(r.Headers.Get("Content-Type"), "application/octet-stream") > -1 {
			if err := r.Save(filepath.Join(outputDir, name)); err != nil {
				log.Println(err)
			}
		}
	})

	c.Visit(url)
}

// truyentranh.net specific collector
func collectTruyenTranh(url, outputDir string) {
	nameInd := 0
	c := colly.NewCollector()

	c.OnHTML(`#viewer`, func(e *colly.HTMLElement) {
		e.ForEach("img", func(i int, ee *colly.HTMLElement) {
			link := ee.Attr("src")
			c.Visit(link)
		})
	})
	c.OnHTML(`div.each-page`, func(e *colly.HTMLElement) {
		e.ForEach("img", func(i int, ee *colly.HTMLElement) {
			link := ee.Attr("src")
			c.Visit(link)
		})
	})
	c.OnHTML(`div.OtherText`, func(e *colly.HTMLElement) {
		e.ForEach("img", func(i int, ee *colly.HTMLElement) {
			link := ee.Attr("src")
			c.Visit(link)
		})
	})

	c.OnRequest(func(r *colly.Request) {
		log.Println("Visiting", r.URL.String())
	})

	c.OnResponse(func(r *colly.Response) {
		name := fmt.Sprintf("%03d%s", nameInd, filepath.Ext(r.FileName()))
		nameInd++
		if strings.Index(r.Headers.Get("Content-Type"), "image") > -1 {
			if err := r.Save(filepath.Join(outputDir, name)); err != nil {
				log.Println(err)
			}
		} else if strings.Index(r.Headers.Get("Content-Type"), "application/octet-stream") > -1 {
			if err := r.Save(filepath.Join(outputDir, name)); err != nil {
				log.Println(err)
			}
		}
	})

	c.Visit(url)
}
