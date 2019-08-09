package bot

import (
	"fmt"
	"log"
	"net/url"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/gocolly/colly"
	"github.com/vqhuy/kindle-manga/kcc"
	"github.com/vqhuy/kindle-manga/util"
)

type Collector interface {
	Collect(base string, chap int, outputDir string)
	Page() string
	GetLink(base string, chap int) string
}

var collectors = make(map[string]Collector)

func Register(c Collector) {
	if _, ok := collectors[c.Page()]; ok {
		panic(fmt.Sprintf("%s is already registered", c))
	}
	collectors[c.Page()] = c
}

func Run(url []string, name string, chap int, dir string) []string {
	var output []string
	var err error

	for n, c := range collectors {
		link := find(n, url)
		if link == "" {
			continue
		}
		c.Collect(link, chap, dir)
		if output, err = kcc.Kcc(name, dir); err != nil {
			logErr(err, "["+c.Page()+"]-["+name+"]")
			continue
		}
		break
	}
	return output
}

func find(name string, url []string) string {
	for _, u := range url {
		if strings.Contains(u, name) {
			return u
		}
	}
	return ""
}

func logErr(e error, extra string) {
	log.Println("[collector]-" + extra + ": " + e.Error())
}

type Bot struct {
	Colly *colly.Collector
}

var _ Collector = (*Bot)(nil)

func (b *Bot) Page() string {
	return ""
}

func (b *Bot) GetLink(base string, chap int) string {
	return ""
}

func (b *Bot) Collect(base string, chap int, outputDir string) {
	nameInd := 0

	b.Colly.OnRequest(func(r *colly.Request) {
		log.Println("Visiting", r.URL.String())
	})

	b.Colly.OnResponse(func(r *colly.Response) {
		ext := util.GetExt(r.FileName())
		name := fmt.Sprintf("%03d%s", nameInd, ext)
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
}

func (b *Bot) Visit(url string) {
	org := strings.TrimSpace(url)
	org = b.filterGoogleCacheLink(org)
	b.Colly.Visit(org)
}

func (b *Bot) filterGoogleCacheLink(str string) string {
	google1 := "images1-focus-opensocial.googleusercontent.com/gadgets/proxy"
	google2 := "images2-focus-opensocial.googleusercontent.com/gadgets/proxy"

	re := regexp.MustCompile(`(?m)url=(.*)$`)

	if strings.Contains(str, google1) || strings.Contains(str, google2) {
		if re.MatchString(str) {
			org := re.FindStringSubmatch(str)[1]
			in, err := url.QueryUnescape(org)
			if err == nil {
				return in
			}
		}
	}
	return str
}
