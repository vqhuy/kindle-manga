package kcc

import (
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"

	"github.com/vqhuy/kindle-manga/util"
)

var kccCmd = "/import/grecc/qhvu/.local/bin/kcc-c2e"
var device = "KO"
var ext = ".mobi"

var limit = 30 // maximum 30 pages per file, to avoid email attachment limits.

func Kcc(name, dir string) ([]string, error) {
	validate(dir)

	log.Printf("Creating MOBI file(s) for " + name + "...")
	// split if necessary
	files, _ := ioutil.ReadDir(dir)
	if len(files) <= limit {
		if err := kccExec(dir); err != nil {
			return nil, err
		}
		return []string{filepath.Join(dir, name+".mobi")}, nil
	}

	var output []string
	for i := 1; i <= len(files)/limit+1; i++ {
		part := name + "_" + strconv.Itoa(i)
		path := filepath.Join(dir, part)

		end := i * limit
		if end > len(files) {
			end = len(files)
		}

		if err := util.Mv(files[(i-1)*limit:end], dir, part); err != nil {
			return nil, err
		}

		if err := kccExec(path); err != nil {
			return nil, err
		}

		output = append(output, filepath.Join(path, part+ext))
	}
	return output, nil
}

func kccExec(io string) error {
	cmd := exec.Command(kccCmd, "-p", device, "-m", "-q", "-u", "-s", "-o", io, io)
	return cmd.Run()
}

func validate(dir string) error {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return err
	}

	for _, fp := range files {
		if fp.Size() == 0 {
			if err := os.Remove(filepath.Join(dir, fp.Name())); err != nil {
				return err
			}
		}
	}
	return nil
}
