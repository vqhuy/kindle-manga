package cmd

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"path/filepath"

	"github.com/BurntSushi/toml"
	"github.com/spf13/cobra"
	"github.com/vqhuy/kindle-manga/bot"
	"github.com/vqhuy/kindle-manga/util"
)

var robotCmd = &cobra.Command{
	Use:   "bot",
	Short: "bot",
	Long:  `bot`,
	Run: func(cmd *cobra.Command, args []string) {
		initConfig()

		var coll collection
		configPath := filepath.Join(configDir, "config.toml")
		if _, err := toml.DecodeFile(configPath, &coll); err != nil {
			log.Fatalln(err)
		}

		bot := newRobot(coll, configPath)
		bot.run()
		bot.save()

	},
}

func init() {
	offlineCmd.AddCommand(robotCmd)
}

var configDir string

func initConfig() {
	usr, err := user.Current()
	if err != nil {
		panic(err)
	}

	configDir = filepath.Join(usr.HomeDir, ".config", "kindle-manga")
	err = os.MkdirAll(configDir, 0755)
	if err != nil {
		panic(err)
	}
}

type manga struct {
	Name string
	Chap int
	URL  []string
}
type collection struct {
	Manga []manga
}
type robot struct {
	configPath string
	fav        collection
}

func newRobot(coll collection, configPath string) *robot {
	return &robot{
		configPath: configPath,
		fav:        coll,
	}
}

func (b *robot) run() {
	mail, err := restoreMailSettings()
	if err != nil {
		panic(err)
	}

	util.Rmdir() // remove tmp dirs of the previous run

	for i := range b.fav.Manga {
		index := i
		manga := b.fav.Manga[i]
		name := fmt.Sprintf("%s_[%d]", manga.Name, manga.Chap)

		dir, err := util.Mkdir(name)
		if err != nil {
			logErr(err, "mkdir")
			continue
		}

		output := bot.Run(manga.URL, name, manga.Chap, dir)
		if len(output) == 0 {
			continue
		}
		for _, o := range output {
			if err := sendToKindle(mail, o); err != nil {
				logErr(err, "send-to-kindle")
				continue
			}
		}

		// update config
		b.fav.Manga[index].Chap++
	}
}

func (b *robot) save() error {
	var confBuf bytes.Buffer

	e := toml.NewEncoder(&confBuf)
	if err := e.Encode(b.fav); err != nil {
		return err
	}
	return ioutil.WriteFile(b.configPath, confBuf.Bytes(), 0755)
}

func logErr(e error, extra string) {
	log.Println("[bot] " + extra + ": " + e.Error())
}
