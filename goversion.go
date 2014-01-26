package main

import (
	"log"
	"os/exec"
	"regexp"
	"runtime"
)

var verreg = regexp.MustCompile("version (?P<verstr>[^ ]+) ")

func goVersion() (version string) {
	bs, err := exec.Command("go", "version").Output()
	version = runtime.Version()
	if err != nil || !verreg.Match(bs) {
		log.Println("Warning: could not run go command, falling back on go version this command was built with", version)
	} else {
		version = string(verreg.FindSubmatch(bs)[1])
	}
	return
}
