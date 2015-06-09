package main

import (
	"fmt"
	"os"

	friff "github.com/odeke-em/friff/src"
)

func main() {
	rest := os.Args[1:]
	argc := len(rest)
	if argc < 1 {
		fmt.Fprintf(os.Stderr, "expecting <left> <right>\n")
		os.Exit(-1)
		return
	}

	var left, right string
	if argc < 2 {
		left = rest[0]
		right = left
	} else if argc < 3 {
		left = rest[0]
		right = rest[1]
	}

	pt := friff.MergePaths(left, right)
	diff := pt.Merge()
	deletions := diff.Deletions
	insertions := diff.Insertions

	fmt.Println("original", diff.Original)

	for _, del := range deletions {
		fmt.Println("deletions", del)
	}

	for _, ins := range insertions {
		fmt.Println("insertions", ins)
	}
}
