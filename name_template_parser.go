package main

import (
	"fmt"
	"strings"

	p "github.com/prataprc/goparsec"
)

// Data/AST Definition

type Template interface {
}

type Chunk interface {
}

type Maybe struct {
	Template
}

func (m Maybe) String() string {
	return fmt.Sprintf("[%s]", m.Template)
}

type Matcher interface {
}

type Tag string

type Filter string

func (f Filter) String() string {
	return fmt.Sprintf(":%s", string(f))
}

type And struct {
	L, R Matcher
}

func (a And) String() string {
	return fmt.Sprintf("(And %s %s)", a.L, a.R)
}

type Or struct {
	L, R Matcher
}

func (o Or) String() string {
	return fmt.Sprintf("(Or %s %s)", o.L, o.R)
}

type Not struct {
	Matcher
}

func (n Not) String() string {
	return fmt.Sprintf("(Not %s)", n.Matcher)
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

// Disjunction: A | B.
func parseDisj(s p.Scanner) (p.ParsecNode, p.Scanner) {
	disj := p.And(func (ns []p.ParsecNode) p.ParsecNode {
		return Or{
			ns[0].(Chunk),
			ns[2].(Chunk),
		}
	}, parseConj, pipe, parseDisj)

	return p.OrdChoice(func (ns []p.ParsecNode) p.ParsecNode {
		return ns[0].(Chunk)
	}, disj, parseConj)(s)
}

// Conjunction: A + B.
func parseConj(s p.Scanner) (p.ParsecNode, p.Scanner) {
	conj := p.And(func (ns []p.ParsecNode) p.ParsecNode {
		return And{
			ns[0].(Chunk),
			ns[2].(Chunk),
		}
	}, parseNeg, plus, parseConj)

	return p.OrdChoice(func (ns []p.ParsecNode) p.ParsecNode {
		return ns[0].(Chunk)
	}, conj, parseNeg)(s)
}

// Negation: A - B
func parseNeg(s p.Scanner) (p.ParsecNode, p.Scanner) {
	neg := p.And(func (ns []p.ParsecNode) p.ParsecNode {
		return And{
			ns[0].(Tag),
			ns[1].(Not),
		}
	}, parseTag, parseUnaryNeg)

	return p.OrdChoice(func (ns []p.ParsecNode) p.ParsecNode {
		return ns[0].(Chunk)
	}, neg, parseUnaryNeg)(s)
}

// Unary Negation: - B
func parseUnaryNeg(s p.Scanner) (p.ParsecNode, p.Scanner) {
	neg := p.And(func (ns []p.ParsecNode) p.ParsecNode {
		return Not{
			ns[1].(Chunk),
		}
	}, minus, parseTag)

	return p.OrdChoice(func (ns []p.ParsecNode) p.ParsecNode {
		return ns[0].(Chunk)
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
