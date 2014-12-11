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
type Chunk interface {
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
type And []Chunk

func (a And) String() string {
	return fmt.Sprintf("(And %v)", a)
}

// A disjunction of conjunctions of which at least one must match.
type Or []And

func (o Or) String() string {
	return fmt.Sprintf("(Or %v)", o)
}

// A single, negated tag.
type Not Tag

func (n Not) String() string {
	return fmt.Sprintf("(Not %s)", n)
}


// Entry Point

func parseNameTemplate(template string) Template {
	scanner := p.NewScanner([]byte(template))
	r, _ := parseDisj(scanner)
	if result, ok := r.(Template); ok {
		return result
	} else {
		return nil
	}
}

// Disjunction: A | B.
func parseDisj(s p.Scanner) (p.ParsecNode, p.Scanner) {
	disj := p.And(func (ns []p.ParsecNode) p.ParsecNode {
		return nil
	}, parseConj, pipe, parseDisj)

	return p.OrdChoice(func (ns []p.ParsecNode) p.ParsecNode {
		return ns[0].(Chunk)
	}, disj, parseConj)(s)
}

// Conjunction: A + B.
func parseConj(s p.Scanner) (p.ParsecNode, p.Scanner) {
	conj := p.And(func (ns []p.ParsecNode) p.ParsecNode {
		return nil
	}, parseNeg, plus, parseConj)

	return p.OrdChoice(func (ns []p.ParsecNode) p.ParsecNode {
		return nil
	}, conj, parseNeg)(s)
}

// Negation: A - B
func parseNeg(s p.Scanner) (p.ParsecNode, p.Scanner) {
	neg := p.And(func (ns []p.ParsecNode) p.ParsecNode {
		return nil
	}, parseTag, parseUnaryNeg)

	return p.OrdChoice(func (ns []p.ParsecNode) p.ParsecNode {
		return nil
	}, neg, parseUnaryNeg)(s)
}

// Unary Negation: - B
func parseUnaryNeg(s p.Scanner) (p.ParsecNode, p.Scanner) {
	neg := p.And(func (ns []p.ParsecNode) p.ParsecNode {
		return nil
	}, minus, parseNeg)

	return p.OrdChoice(func (ns []p.ParsecNode) p.ParsecNode {
		return nil
	}, neg, parseTag)(s)
}

// Negation: 

func parseTag(s p.Scanner) (p.ParsecNode, p.Scanner) {
	n, s2 := p.Token(`^[^,{}:\r\n\+\-\|]+`, "TAG")(s)
	if tag, ok := n.(*p.Terminal); ok {
		return Tag(strings.TrimSpace(tag.Value)), s2
	} else {
		return nil, s
	}
}


// Terminals

var plus = p.Token(`^\+`, "PLUS")
var minus = p.Token(`^\-`, "MINUS")
var pipe = p.Token(`^\|`, "PIPE")
