package main

import (
	"fmt"
	"strings"

	p "github.com/prataprc/goparsec"
)

// Data/AST Definition

type Template interface {
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
	r, _ := parseTemplateChunk(scanner)
	if result, ok := r.(Template); ok {
		return result
	} else {
		return nil
	}
}

func parseTemplateChunk(s p.Scanner) (p.ParsecNode, p.Scanner) {
	return p.OrdChoice(func (ns []p.ParsecNode) p.ParsecNode {
		return ns[0]
	}, parseUnaryNeg, parseNeg, parseDisj, parseConj, parseTag)(s)
}

// Conjunction: A + B.
func parseConj(s p.Scanner) (p.ParsecNode, p.Scanner) {
	return p.And(func (ns []p.ParsecNode) p.ParsecNode {
		return And{
			ns[0].(Tag),
			ns[2].(Tag),
		}
	}, parseTag, plus, parseTag)(s)
}

// Disjunction: A | B.
func parseDisj(s p.Scanner) (p.ParsecNode, p.Scanner) {
	return p.And(func (ns []p.ParsecNode) p.ParsecNode {
		return Or{
			ns[0].(Tag),
			ns[2].(Tag),
		}
	}, parseTag, pipe, parseTag)(s)
}

// Negation: A - B
func parseNeg(s p.Scanner) (p.ParsecNode, p.Scanner) {
	return p.And(func (ns []p.ParsecNode) p.ParsecNode {
		return And{
			ns[0].(Tag),
			ns[1].(Not),
		}
	}, parseTag, parseUnaryNeg)(s)
}

// Unary Negation: - B
func parseUnaryNeg(s p.Scanner) (p.ParsecNode, p.Scanner) {
	return p.And(func (ns []p.ParsecNode) p.ParsecNode {
		return Not{
			ns[1].(Tag),
		}
	}, minus, parseTag)(s)
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
