package main

import (
	"strings"

	p "github.com/prataprc/goparsec"
)

func trimmedTerminal(pattern, name string) p.Parser {
	term := p.Token(pattern, name)
	return func(s p.Scanner) (p.ParsecNode, p.Scanner) {
		tkn, s2 := term(s)
		if t, ok := tkn.(*p.Terminal); ok {
			return strings.TrimSpace(t.Value), s2
		} else {
			return nil, s
		}
	}
}

var comment = p.Token(`^//.*\n`, "COMMENT")
var name = trimmedTerminal(`^[0-9a-zA-Z\.\-_ "']+`, "NAME")

func filter(s p.Scanner) (p.ParsecNode, p.Scanner) {
	n, s2 := p.Token(`^[a-z]+`, "FILTER")(s)
	if filter, ok := n.(*p.Terminal); ok {
		return Filter(strings.TrimSpace(filter.Value)), s2
	} else {
		return nil, s
	}
}

func tag(s p.Scanner) (p.ParsecNode, p.Scanner) {
	n, s2 := p.Token(`^[0-9a-zA-Z_ ]+`, "TAG")(s)
	if tag, ok := n.(*p.Terminal); ok {
		return Tag(strings.TrimSpace(tag.Value)), s2
	} else {
		return nil, s
	}
}

// Punctuation

var lbrace = p.Token(`^{`, "LBRACE")
var rbrace = p.Token(`^}`, "RBRACE")
var lbracket = p.Token(`^\[`, "LBRACKET")
var rbracket = p.Token(`^\]`, "RBRACKET")
var colon = p.Token(`^:`, "COLON")
var comma = p.Token(`^,`, "COMMA")
var minus = p.Token(`^\-`, "MINUS")
var pipe = p.Token(`^\|`, "PIPE")
var plus = p.Token(`^\+`, "PLUS")
