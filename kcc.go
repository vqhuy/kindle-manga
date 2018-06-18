package main

import (
	"io/ioutil"
	"log"
	"os/exec"
	"path/filepath"
	"strconv"
)

var kccCmd = "kcc-c2e"
var device = "KO"
var ext = ".mobi"

var limit = 30 // maximum 30 pages per file, to avoid email attachment limits.

func kcc(name, dir string) ([]string, error) {
	log.Printf("Creating MOBI file(s) for " + name + "...")
	// split if necessary
	files, _ := ioutil.ReadDir(dir)
	if len(files) <= limit {
		return []string{filepath.Join(dir, name+".mobi")}, kccExec(dir)
	}

	var output []string
	for i := 1; i <= len(files)/limit+1; i++ {
		part := name + "_" + strconv.Itoa(i)
		path := filepath.Join(dir, part)

		end := i * limit
		if end > len(files) {
			end = len(files)
		}

		if err := mv(files[(i-1)*limit:end], dir, part); err != nil {
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
