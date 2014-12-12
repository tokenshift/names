package main

import (
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"time"
)

func main() {
	// Parse name templates on command line.
	matchers := make([]Matcher, len(os.Args) - 1)
	for i, arg := range(os.Args[1:len(os.Args)]) {
		matcher, err := parseNameTemplate(arg)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		matchers[i] = matcher
	}

	// Load name files.
	nameFiles, err := filepath.Glob("*.names")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	// Check all names against each matcher and store matches.
	matches := make([][]string, len(matchers))
	names := parseNameFiles(nameFiles)
	for name := range(names) {
		for i, matcher := range(matchers) {
			if matcher.Matches(name) {
				matches[i] = append(matches[i], name.Name)
			}
		}
	}

	// Check that every component was matched.
	allMatched := true
	for i, ms := range(matches) {
		if len(ms) == 0 {
			fmt.Fprintln(os.Stderr, "No match for", matchers[i])
			allMatched = false
		}
	}
	if !allMatched {
		os.Exit(1)
	}

	// Pick a random name for each component.
	rand.Seed(time.Now().UnixNano())
	for i, ms := range(matches) {
		if i > 0 {
			fmt.Print(" ")
		}
		fmt.Print(ms[rand.Intn(len(ms))])
	}
	fmt.Print("\n")
}
