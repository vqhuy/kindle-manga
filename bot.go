package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"sync"

	"github.com/BurntSushi/toml"
)

type bot struct {
	configPath string
	fav        collection
	waitGroup  sync.WaitGroup
}

func newBot(coll collection, configPath string) *bot {
	return &bot{
		configPath: configPath,
		fav:        coll,
	}
}

func (b *bot) run() {
	mail, err := restoreMailSettings()
	if err != nil {
		panic(err)
	}

	rmdir() // remove tmp dirs of the previous run

	for i := range b.fav.Manga {
		index := i
		manga := b.fav.Manga[i]
		name := fmt.Sprintf("%s_[%d]", manga.Name, manga.Chap)
		url1 := fmt.Sprintf("%s-chuong-%d/", manga.URL1, manga.Chap)
		url2 := fmt.Sprintf("%s/Chap-%03d/", manga.URL2, manga.Chap)

		dir, err := mkdir(name)
		if err != nil {
			botErr(err, "mkdir")
			continue
		}

		b.waitGroup.Add(1)
		go func() {
			defer b.waitGroup.Done()
			var output []string
			var err error

			collectTruyenTranh(url2, dir) // try truyentranh.net first
			if output, err = kcc(name, dir); err != nil {
				botErr(err, "[kcc-truyentranh] "+name)
				collectTruyenTranhTuan(url1, dir) // try truyentranhtuan.com
				if output, err = kcc(name, dir); err != nil {
					botErr(err, "[kcc-truyentranhtuan] "+name)
					return
				}
			}
			for _, o := range output {
				if err := sendToKindle(mail, o); err != nil {
					botErr(err, "send-to-kindle")
					return
				}
			}

			// update config
			b.fav.Manga[index].Chap++
		}()
	}

	b.waitGroup.Wait()
}

func (b *bot) save() error {
	var confBuf bytes.Buffer

	e := toml.NewEncoder(&confBuf)
	if err := e.Encode(b.fav); err != nil {
		return err
	}
	return ioutil.WriteFile(b.configPath, confBuf.Bytes(), 0755)
}

func botErr(e error, extra string) {
	log.Println("[bot] " + extra + ": " + e.Error())
}
