package blogtruyen

import (
	"strconv"
	"strings"

	"github.com/vqhuy/kindle-manga/bot"
	"github.com/gocolly/colly"
)

func init() {
	bot.Register(&collector{
		Bot: &bot.Bot{},
	})
}

type collector struct {
	*bot.Bot
}

var _ bot.Collector = (*collector)(nil)

func (b *collector) Page() string {
	return "http://blogtruyen.com"
}

func (b *collector) GetLink(base string, chap int) string {
	// should not use the base colly's collector
	c := colly.NewCollector()
	var link string
	c.OnHTML(`#list-chapters p`, func(e *colly.HTMLElement) {
		s := e.DOM.Children().First().Children().First()
		href, _ := s.Attr("href")
		title, _ := s.Attr("title")
		if strings.Contains(title, strconv.Itoa(chap)) {
			link = b.Page() + href
		}
	})

	c.Visit(base)
	return link
}

func (b *collector) Collect(base string, chap int, outputDir string) {
	b.Colly = colly.NewCollector()

	b.Colly.OnHTML(`#content`, func(e *colly.HTMLElement) {
		e.ForEach("img", func(i int, ee *colly.HTMLElement) {
			link := ee.Attr("src")
			b.Colly.Visit(link)
		})
	})

	b.Bot.Collect(base, chap, outputDir)

	link := b.GetLink(base, chap)
	b.Colly.Visit(link)
}
