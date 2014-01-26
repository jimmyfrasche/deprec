package main

import (
	"os/exec"
	"path/filepath"
	"strings"

	"code.google.com/p/go.tools/go/vcs"
	"github.com/jimmyfrasche/goutil"
)

type Root struct {
	root, path string
	cmd        *vcs.Cmd
}

func repoRoot(p *goutil.Package) (*Root, error) {
	rdir := filepath.Join(p.Build.Root, "src")
	cmd, path, err := vcs.FromDir(p.Build.Dir, rdir)
	if err != nil {
		return nil, err
	}
	dir := filepath.Join(rdir, path)
	return &Root{root: dir, path: path, cmd: cmd}, nil
}

var repo = map[string][]string{
	"bzr": []string{"revno"},
	"git": []string{"rev-parse", "HEAD"},
	"hg":  []string{"id", "-i"},
	"svn": nil, //special cased
}

func (r *Root) Rev() (string, error) {

	var cmd *exec.Cmd
	nm := r.cmd.Cmd

	if nm == "svn" {
		cmd = exec.Command("svnversion")
	} else {
		cmd = exec.Command(nm, repo[nm]...)
	}
	cmd.Dir = r.root

	bs, err := cmd.Output()
	return strings.TrimSpace(string(bs)), err
}
