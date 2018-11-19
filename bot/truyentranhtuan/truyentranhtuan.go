package truyentranhtuan

import (
	"encoding/json"
	"fmt"
	"html"
	"log"
	"regexp"

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
	return "http://truyentranhtuan.com/"
}

func (b *collector) GetLink(base string, chap int) string {
	return fmt.Sprintf("%s-chuong-%d/", base, chap)
}

func (b *collector) Collect(base string, chap int, outputDir string) {
	b.Colly = colly.NewCollector()

	re := regexp.MustCompile(`(?m)var slides_page_url_path = \s*(.*).$`)

	b.Colly.OnHTML(`script`, func(e *colly.HTMLElement) {
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
				b.Visit(u)
			}
		}
	})

	b.Bot.Collect(base, chap, outputDir)

	link := b.GetLink(base, chap)
	b.Colly.Visit(link)
}
