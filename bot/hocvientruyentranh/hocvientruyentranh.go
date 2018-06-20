package hocvientruyentranh

import (
	"strconv"
	"strings"

	"github.com/c633/kindle-manga/bot"
	"github.com/gocolly/colly"
)

func init() {
	bot.Register(&collector{
		Bot: &bot.Bot{
			Colly: colly.NewCollector(),
		},
	})
}

type collector struct {
	*bot.Bot
}

var _ bot.Collector = (*collector)(nil)

func (b *collector) Page() string {
	return "http://hocvientruyentranh.com/"
}

func (b *collector) GetLink(base string, chap int) string {
	// should not use the base colly's collector
	c := colly.NewCollector()
	var link string
	c.OnHTML(`div.box-body tbody tr`, func(e *colly.HTMLElement) {
		s := e.DOM.Children().First().Children().First()
		href, _ := s.Attr("href")
		title, _ := s.Attr("title")
		if strings.Contains(title, strconv.Itoa(chap)) {
			link = href
		}
	})

	c.Visit(base)
	return link
}

func (b *collector) Collect(base string, chap int, outputDir string) {
	b.Colly.OnHTML(`div.manga-container`, func(e *colly.HTMLElement) {
		e.ForEach("img", func(i int, ee *colly.HTMLElement) {
			link := ee.Attr("src")
			b.Colly.Visit(link)
		})
	})

	b.Bot.Collect(base, chap, outputDir)

	link := b.GetLink(base, chap)
	b.Colly.Visit(link)
}
