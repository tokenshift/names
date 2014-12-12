package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func main() {
	matchers := make([]Matcher, len(os.Args))
	for i, arg := range(os.Args[1:len(os.Args)]) {
		matcher, err := parseNameTemplate(arg)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		matchers[i] = matcher
	}

	fmt.Println(matchers)

	nameFiles, err := filepath.Glob("*.names")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	names := parseNameFiles(nameFiles)
	for name := range(names) {
		fmt.Println(name)
	}
}
