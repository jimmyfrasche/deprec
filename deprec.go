package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

var (
	check  = flag.Bool("check", false, "compare dep.log to current system and do not save a new dep.log")
	script = flag.Bool("s", false, "exit with 1 if different instead of printing differences")
)

func Usage() {
	log.Printf("usage: %s [flags] import-path*\n", os.Args[0])
	flag.PrintDefaults()
	os.Exit(2)
}

func main() {
	log.SetFlags(0)
	flag.Usage = Usage
	flag.Parse()

	args := flag.Args()
	if len(args) == 0 {
		args = append(args, ".")
	}
	multiple := len(args) > 1
	prefix := ""
	if multiple {
		prefix = "\t"
	}

	logerr := func(err error) {
		log.Printf("%s%s", prefix, err)
	}

	wasdiff := false

outer:
	for _, arg := range args {
		if multiple {
			fmt.Printf("%s:\n", arg)
		}

		root, roots, err := Import(nil, arg)
		if err != nil {
			logerr(err)
			continue
		}

		deplog, err := ReadDepLog(root.root)
		if err != nil {
			logerr(err)
			continue
		}

		for _, root := range roots {
			rev, err := root.Rev() //TODO
			if err != nil {
				logerr(err)
				continue outer
			}
			deplog.Add(root.path, rev)
		}

		diffs := deplog.Diff()

		for _, diff := range diffs {
			if !*script {
				fmt.Printf("%s%s\n", prefix, diff)
			}
			wasdiff = true
		}

		if !*check {
			err = deplog.Write()
			if err != nil {
				logerr(err)
			}
		}
	}

	if *script && wasdiff {
		os.Exit(1)
	}
}
