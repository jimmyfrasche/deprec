package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

var (
	check      = flag.Bool("check", false, "compare dep.log to current system and do not save a new dep.log")
	script     = flag.Bool("s", false, "exit with 1 if different instead of printing differences")
	writeGodep = flag.Bool("with-godep", false, "write dep.log even if repo contains a Godeps file from godep(1)")
)

func Usage() {
	log.Printf("usage: %s [flags] import-path*\n", os.Args[0])
	flag.PrintDefaults()
	os.Exit(2)
}

//Usage: %name %flags import-path*
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

	wasdiff, waserr := false, false

outer:
	for _, arg := range args {
		if multiple {
			fmt.Printf("%s:\n", arg)
		}

		root, roots, err := Import(nil, arg)
		if err != nil {
			logerr(err)
			waserr = true
			continue
		}

		deplog, err := ReadDepLog(root.root)
		if err != nil {
			logerr(err)
			waserr = true
			continue
		}

		for _, root := range roots {
			rev, err := root.Rev()
			if err != nil {
				logerr(err)
				waserr = true
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
			if deplog.fromGodep && !*writeGodep {
				f := filepath.Join(filepath.Base(deplog.root), "dep.log")
				log.Printf("Cannot write %s. Dependency information was from godep(1)\n", f)
				waserr = true
				continue
			}
			err = deplog.Write()
			if err != nil {
				waserr = true
				logerr(err)
			}
		}
	}

	if waserr || (*script && wasdiff) {
		os.Exit(1)
	}
}
