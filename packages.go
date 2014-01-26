package main

import (
	"go/build"
	"path/filepath"

	"github.com/jimmyfrasche/goutil"
)

//Import finds all the repo roots of all the dependencies of every
//package in the root of path, excluding standard library.
func Import(ctx *build.Context, path string) (*Root, []*Root, error) {
	pkg, err := goutil.Import(nil, path)
	if err != nil {
		return nil, nil, err
	}

	root, err := repoRoot(pkg)
	if err != nil {
		return nil, nil, err
	}

	tree := filepath.Join(pkg.Build.Root, "src", root.path)
	pkgs, err := goutil.ImportTree(ctx, tree)
	if err != nil {
		return nil, nil, err
	}

	var deps goutil.Packages
	for _, pkg := range pkgs {
		ps, err := pkg.ImportDeps()
		if err != nil {
			return nil, nil, err
		}
		deps = append(deps, ps...)
	}
	deps = deps.NoStdlib().Uniq()

	var roots []*Root
	seen := map[string]bool{root.path: true}
	for _, dep := range deps {
		root, err := repoRoot(dep)
		if err != nil {
			return nil, nil, err
		}
		if !seen[root.path] {
			seen[root.path] = true
			roots = append(roots, root)
		}
	}

	return root, roots, nil
}
