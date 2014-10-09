package main

import (
	//"github.com/prataprc/goparsec"
)

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

func parseNameTemplate(template string) Template {
	return nil
}
