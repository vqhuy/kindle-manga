package nettruyen

import (
	"fmt"

	"github.com/gocolly/colly"
	"github.com/vqhuy/kindle-manga/bot"
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
	return "http://www.nettruyen.com"
}

func (b *collector) GetLink(base string, chap int) string {
	// should not use the base colly's collector
	chapstr := fmt.Sprintf("Chapter %d", chap)
	c := colly.NewCollector()
	var link string
	c.OnHTML(`#nt_listchapter li div a`, func(e *colly.HTMLElement) {
		s := e.DOM
		href, _ := s.Attr("href")
		title, _ := s.Html()
		if title == chapstr {
			link = href
		}
	})

	c.Visit(base)
	return link
}

func (b *collector) Collect(base string, chap int, outputDir string) error {
	b.Colly = colly.NewCollector()

	b.Colly.OnHTML(`div.reading-detail`, func(e *colly.HTMLElement) {
		e.ForEach("img", func(i int, ee *colly.HTMLElement) {
			link := ee.Attr("src")
			b.Visit(link)
		})
	})

	b.Bot.Collect(base, chap, outputDir)

	link := b.GetLink(base, chap)
	if link == "" {
		return bot.ErrorChapNotFound
	}
	b.Colly.Visit(link)
	return nil
}
