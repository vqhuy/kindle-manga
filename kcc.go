package main

import (
	"log"
	"os/exec"
)

var kccCmd = "kcc-c2e"
var device = "KO"

func kcc(name, dir string) error {
	log.Printf("Creating MOBI file for " + name + "...")
	cmd := exec.Command(kccCmd, "-p", device, "-m", "-q", "-u", "-s", "-o", dir, dir)
	return cmd.Run()
}
