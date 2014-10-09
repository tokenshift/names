package main

import (
	p "github.com/prataprc/goparsec"
)

// Data/AST Definition

type Template interface {
}

type Maybe Template

type Matcher interface {
}

type Tag string
type Filter string

type And struct {
	L, R Matcher
}

type Or struct {
	L, R Matcher
}

type Not Matcher


// Entry Point and Non-Terminals

func parseNameTemplate(template string) Template {
	scanner := p.NewScanner([]byte(template))
	t, _ := tag(scanner)
	if term, ok := t.(*p.Terminal); ok {
		return Tag(term.Value)
	} else {
		return nil
	}
}


// Terminals

var tag = p.Token(`^[^,{}:\r\n\+\-\|]+`, "TAG")
