package cmd

import (
	"fmt"
	"os"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/vqhuy/kindle-manga/bot"
)

var offlineCmd = &cobra.Command{
	Use:   "offline",
	Short: "offline",
	Long:  `offline`,
	Run: func(cmd *cobra.Command, args []string) {
		url := cmd.Flag("url").Value.String()
		name := cmd.Flag("name").Value.String()
		out := cmd.Flag("out").Value.String()
		chap, err := strconv.Atoi(cmd.Flag("chap").Value.String())
		if err != nil {
			chap = 1
		}
		bot.RunOffline([]string{url}, name, out, chap)
	},
}

func init() {
	offlineCmd.Flags().StringP("out", "o", ".", "Location of directory for storing generated files")
	offlineCmd.Flags().StringP("name", "n", "manga", "Name of the manga")
	offlineCmd.Flags().StringP("url", "i", "", "URL")
	offlineCmd.Flags().IntP("chap", "c", 1, "First chap")
}

func Execute() {
	if err := offlineCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
