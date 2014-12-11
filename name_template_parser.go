package main

import (
	"fmt"
	"strings"

	p "github.com/prataprc/goparsec"
)

// Data/AST Definition

type Template interface {
}

// A single 'piece' of a template, like a Tag (possibly negated).
type Matcher interface {
}

type Maybe struct {
	Template
}

func (m Maybe) String() string {
	return fmt.Sprintf("[%s]", m.Template)
}

// A single tag to match.
type Tag string

// A standalone filter.
type Filter string

// A filter applied to a tag.
type Filtered struct {
	Tag
	Filter
}

func (f Filter) String() string {
	return fmt.Sprintf(":%s", string(f))
}

// A conjunction of tags/chunks that must all must match.
type And []Matcher

func (a And) String() string {
	return fmt.Sprintf("(And %s)", []Matcher(a))
}

// A disjunction of conjunctions of which at least one must match.
type Or []And

func (o Or) String() string {
	return fmt.Sprintf("(Or %s)", []And(o))
}

// A single, negated tag.
type Not Tag

func (n Not) String() string {
	return fmt.Sprintf("(Not %s)", string(n))
}


// Entry Point and Non-Terminals

func parseNameTemplate(template string) Template {
	scanner := p.NewScanner([]byte(template))
	r, _ := parseDisj(scanner)
	if result, ok := r.(Template); ok {
		return result
	} else {
		return nil
	}
}

// Disjunction: A (| B)*
func parseDisj(s p.Scanner) (p.ParsecNode, p.Scanner) {
	return p.Kleene(func(ns []p.ParsecNode) p.ParsecNode {
		terms := make([]And, len(ns))
		for i, n := range(ns) {
			terms[i] = n.(And)
		}
		return Or(terms)
	}, parseConj, pipe)(s)
}

// Conjunction: A (+ B)*
func parseConj(s p.Scanner) (p.ParsecNode, p.Scanner) {
	// A conjunction can be a tag (optionally with a + or -), followed by any
	// number of AndTags or NotTags ("+ A" or "- A").
	head := p.OrdChoice(func(ns []p.ParsecNode) p.ParsecNode {
		return ns[0].(Matcher)
	}, parseAndTag, parseNotTag, tag)

	tailEntry := p.OrdChoice(func(ns []p.ParsecNode) p.ParsecNode {
		return ns[0].(Matcher)
	}, parseAndTag, parseNotTag)

	tail := p.Kleene(func(ns []p.ParsecNode) p.ParsecNode {
		terms := make([]Matcher, len(ns))
		for i, n := range(ns) {
			terms[i] = n.(Matcher)
		}
		return terms
	}, tailEntry)

	return p.And(func(ns []p.ParsecNode) p.ParsecNode {
		t := ns[0].(Matcher)
		ts := ns[1].([]Matcher)
		return And(append([]Matcher{t}, ts...))
	}, head, tail)(s)
}

// A negated tag (- A)
func parseNotTag(s p.Scanner) (p.ParsecNode, p.Scanner) {
	return p.And(func(ns []p.ParsecNode) p.ParsecNode {
		return Not(ns[1].(Tag))
	}, minus, tag)(s)
}

// An added tag (+ A)
func parseAndTag(s p.Scanner) (p.ParsecNode, p.Scanner) {
	return p.And(func(ns []p.ParsecNode) p.ParsecNode {
		return ns[1].(Tag)
	}, plus, tag)(s)
}

func tag(s p.Scanner) (p.ParsecNode, p.Scanner) {
	n, s2 := p.Token(`^[^,{}:\r\n\+\-\|]+`, "TAG")(s)
	if tag, ok := n.(*p.Terminal); ok {
		return Tag(strings.TrimSpace(tag.Value)), s2
	} else {
		return nil, s
	}
}

var plus = p.Token(`^\+`, "PLUS")
var minus = p.Token(`^\-`, "MINUS")
var pipe = p.Token(`^\|`, "PIPE")
