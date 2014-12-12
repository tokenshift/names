package main

import (
	"fmt"

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

func (f Filter) String() string {
	return fmt.Sprintf(":%s", string(f))
}

// A filter applied to a tag.
type Filtered struct {
	Tag
	Filter
}

func (f Filtered) String() string {
	return fmt.Sprintf("%s%s", f.Tag, f.Filter)
}

// A conjunction of tags/chunks that must all must match.
type And []Matcher

func (a And) String() string {
	return fmt.Sprintf("(And %v)", []Matcher(a))
}

// A disjunction of conjunctions of which at least one must match.
type Or []And

func (o Or) String() string {
	return fmt.Sprintf("(Or %v)", []And(o))
}

// A single, negated term.
type Not struct {
	Matcher
}

func (n Not) String() string {
	return fmt.Sprintf("(Not %s)", n.Matcher)
}


// Entry Point and Non-Terminals

func parseNameTemplate(template string) (Template, error) {
	scanner := p.NewScanner([]byte(template))
	r, _ := parseMaybe(scanner)
	if result, ok := r.(Template); ok {
		return result, nil
	} else {
		return nil, nil
	}
}

// Maybe ("[template]")
func parseMaybe(s p.Scanner) (p.ParsecNode, p.Scanner) {
	maybe := p.And(func(ns []p.ParsecNode) p.ParsecNode {
		return Maybe{
			ns[1].(Matcher),
		}
	}, lbracket, parseDisj, rbracket)

	return p.OrdChoice(func(ns []p.ParsecNode) p.ParsecNode {
		return ns[0].(Matcher)
	}, maybe, parseDisj)(s)
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
	}, parseAndTag, parseNotTag, parseTerm)

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
		return Not{ns[1].(Matcher)}
	}, minus, parseTerm)(s)
}

// An added tag (+ A)
func parseAndTag(s p.Scanner) (p.ParsecNode, p.Scanner) {
	return p.And(func(ns []p.ParsecNode) p.ParsecNode {
		return ns[1].(Matcher)
	}, plus, parseTerm)(s)
}

// A term (tag, filter, or both).
func parseTerm(s p.Scanner) (p.ParsecNode, p.Scanner) {
	return p.OrdChoice(func(ns []p.ParsecNode) p.ParsecNode {
		return ns[0].(Matcher)
	}, parseFiltered, parseFilter, tag)(s)
}

// A filtered term (Term:filter)
func parseFiltered(s p.Scanner) (p.ParsecNode, p.Scanner) {
	return p.And(func(ns []p.ParsecNode) p.ParsecNode {
		return Filtered{
			ns[0].(Tag),
			ns[1].(Filter),
		}
	}, tag, parseFilter)(s)
}

// A filter (:filter)
func parseFilter(s p.Scanner) (p.ParsecNode, p.Scanner) {
	return p.And(func(ns []p.ParsecNode) p.ParsecNode {
		return ns[1].(Filter)
	}, colon, filter)(s)
}
