package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func main() {
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
