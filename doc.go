//Command deprec(1) records the dependencies of a repository in a simple file
//and allows easy comparison to help track down dependency differences.
//
//The file contains the Go version of the current system along with all
//dependencies' repositories and their revision id in the root of the
//repository.
//The Go version always goes on the first line.
//After that each dependency is listed in alphabetical order.
//
//REVISION IDS
//
//The commands used to extract the revision ids are:
//	bzr revno
//	git rev-parse HEAD
//	hg id -i
//	svnversion
//
//EXAMPLES
//
//Compare the current repositories dep.log with the current state of the system
//	deprec -check
//
//To create dep.log in current repository
//	deprec
//
//To create several dep.log's simultaneously
//	deprec path/one other/path a/third/path
//
//Use in a script
//	if deprec -s
//	then
//		echo different
//	fi
package main
