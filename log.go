package main

import (
	"encoding/json"
	"fmt"
	"github.com/jimmyfrasche/goutil"
	"io/ioutil"
	"path/filepath"
	"sort"
	"strings"
)

type (
	Entry struct {
		repo, old, new string
	}

	sortEntry []*Entry

	DepLog struct {
		root      string
		Go, oldGo string
		Revs      map[string]*Entry
		sorted    sortEntry
		fromGodep bool
	}

	Godeps struct {
		GoVersion string
		Deps      []struct {
			ImportPath string
			Rev        string
		}
	}
)

func (s sortEntry) Len() int {
	return len(s)
}

func (s sortEntry) Less(i, j int) bool {
	return s[i].repo < s[j].repo
}

func (s sortEntry) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func ReadDepLog(path string) (*DepLog, error) {
	out := &DepLog{
		root: path,
		Go:   goVersion(),
		Revs: map[string]*Entry{},
	}
	bs, err := ioutil.ReadFile(out.path())
	if err != nil {
		//if there wasn't an old dep file, see if there was a Godep file
		err = readGodeps(out)
		if err != nil {
			//no old dep info, nothing to do
			return out, nil
		}
		return out, nil
	}
	for i, record := range strings.Split(string(bs), "\n") {
		if record == "" {
			continue
		}
		fields := strings.SplitN(record, " ", 2)
		if len(fields) != 2 {
			return nil, fmt.Errorf("dep.log: malformed line %d: %s", i, record)
		}
		repo, ver := fields[0], fields[1]
		if i == 0 && repo != "Go" {
			return nil, fmt.Errorf("dep.log: malformed, first line must be go version")
		}
		if i == 0 {
			out.oldGo = ver
			continue
		}
		out.Revs[repo] = &Entry{repo: repo, old: ver}
	}
	return out, nil
}

func readGodeps(d *DepLog) error {
	bs, err := ioutil.ReadFile(d.godepPath())
	if err != nil {
		//no old dep info, nothing to do
		return nil
	}

	var godep Godeps
	err = json.Unmarshal(bs, &godep)
	if err != nil {
		return err
	}

	d.fromGodep = true
	d.oldGo = godep.GoVersion
	for _, rev := range godep.Deps {
		//XXX below nil may cause problems if we ever include -tags
		pkg, err := goutil.Import(nil, rev.ImportPath)
		if err != nil {
			return err
		}
		root, err := repoRoot(pkg)
		if err != nil {
			return err
		}
		p := root.path
		d.Revs[p] = &Entry{repo: p, old: rev.Rev}
	}
	return nil
}

func (d *DepLog) path() string {
	return filepath.Join(d.root, "dep.log")
}

func (d *DepLog) godepPath() string {
	return filepath.Join(d.root, "Godeps")
}

func (d *DepLog) Add(repo, rev string) {
	if _, ok := d.Revs[repo]; ok {
		d.Revs[repo].new = rev
	} else {
		d.Revs[repo] = &Entry{repo: repo, new: rev}
	}
}

func (d *DepLog) sort() {
	d.sorted = nil
	for _, r := range d.Revs {
		d.sorted = append(d.sorted, r)
	}
	sort.Sort(d.sorted)
}

func (d *DepLog) Diff() (out []string) {
	if d.oldGo == "" {
		return
	}
	if len(d.sorted) == 0 {
		d.sort()
	}
	pushf := func(s string, v ...interface{}) {
		out = append(out, fmt.Sprintf(s, v...))
	}
	if d.oldGo == "" && d.Go != d.oldGo {
		pushf("Package built with %s but you are using %s", d.Go, d.oldGo)
	}
	for _, e := range d.Revs {
		if e.old != e.new {
			switch {
			case e.new == "":
				//nothing to report
			case e.old == "":
				pushf("Package %s is a new dependency", e.repo)
			default:
				pushf("Package built with %s revision %s but you are using %s", e.repo, e.old, e.new)
			}
		}
	}
	return
}

func (d *DepLog) Write() error {
	if len(d.sorted) == 0 {
		d.sort()
	}
	out := []string{fmt.Sprint("Go ", d.Go)}
	for _, rec := range d.sorted {
		if rec.new == "" {
			//no longer dep
			continue
		}
		out = append(out, fmt.Sprint(rec.repo, " ", rec.new))
	}
	return ioutil.WriteFile(d.path(), []byte(strings.Join(out, "\n")+"\n"), 0644)
}
