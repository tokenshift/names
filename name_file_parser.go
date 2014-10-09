package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	p "github.com/prataprc/goparsec"
)

type TaggedName struct {
	Name string
	Tags []string
}

func (n TaggedName) String() string {
	if len(n.Tags) > 0 {
		return fmt.Sprintf("%s: %s", n.Name, strings.Join(n.Tags, ", "))
	} else {
		return n.Name
	}
}

var tagStack [][]string
func tags() []string {
	set := make(map[string]bool)
	for _, tags := range(tagStack) {
		for _, tag := range(tags) {
			set[tag] = true
		}
	}

	distinct := make([]string, 0, len(set))
	for tag, _ := range(set) {
		distinct = append(distinct, tag)
	}
	return distinct
}

func parseNameFiles(filenames []string) {
	for _, filename := range(filenames) {
		parseNameFile(filename)
	}
}

func parseNameFile(filename string) {
	buffer, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	tagStack = make([][]string, 0)
	scanner := p.NewScanner(buffer)
	parseBlockContents(scanner)
}

// Block = TagList BlockContents '}'
func parseBlock(s p.Scanner) (p.ParsecNode, p.Scanner) {
	return p.And(func (ns []p.ParsecNode) p.ParsecNode {
		tagStack = tagStack[:len(tagStack)-1]
		return struct{}{}
	}, parseTagList, parseBlockContents, rbrace)(s)
}

// TagList = TagListStart* TagListEnd
func parseTagList(s p.Scanner) (p.ParsecNode, p.Scanner) {
	inits := p.Kleene(func (ns []p.ParsecNode) p.ParsecNode {
		tags := make([]string, len(ns))
		for i, n := range(ns) {
			tags[i] = n.(string)
		}
		return tags
	}, parseTagListStart)

	return p.And(func (ns []p.ParsecNode) p.ParsecNode {
		tags := append(ns[0].([]string), ns[1].(string))
		tagStack = append(tagStack, tags)
		return struct{}{}
	}, inits, parseTagListEnd)(s)
}

// TagListStart = ident ','
func parseTagListStart(s p.Scanner) (p.ParsecNode, p.Scanner) {
	return p.And(func (ns []p.ParsecNode) p.ParsecNode {
		return strings.TrimSpace(ns[0].(*p.Terminal).Value)
	}, ident, comma)(s)
}

// TagListEnd = ident '{'
func parseTagListEnd(s p.Scanner) (p.ParsecNode, p.Scanner) {
	return p.And(func (ns []p.ParsecNode) p.ParsecNode {
		return strings.TrimSpace(ns[0].(*p.Terminal).Value)
	}, ident, lbrace)(s)
}

// BlockContents = BlockContent*
func parseBlockContents(s p.Scanner) (p.ParsecNode, p.Scanner) {
	return p.Kleene(func (ns []p.ParsecNode) p.ParsecNode {
		return struct{}{}
	}, parseBlockContent)(s)
}

// BlockContent = comment | Block | Name
func parseBlockContent(s p.Scanner) (p.ParsecNode, p.Scanner) {
	return p.OrdChoice(func (ns []p.ParsecNode) p.ParsecNode {
		if name, ok := ns[0].(TaggedName); ok {
			name.Tags = append(name.Tags, tags()...)
			fmt.Println(name)
		}
		return struct{}{}
	}, comment, parseBlock, parseName)(s)
}

// Name = ident [InlineTags]
func parseName(s p.Scanner) (p.ParsecNode, p.Scanner) {
	inlineTags := p.Maybe(func (ns []p.ParsecNode) p.ParsecNode {
		if ns == nil {
			return []string {}
		} else {
			return ns[0].([]string)
		}
	}, parseInlineTags)

	return p.And(func (ns []p.ParsecNode) p.ParsecNode {
		if tags, ok := ns[1].([]string); ok {
			return TaggedName {
				Name: ns[0].(*p.Terminal).Value,
				Tags: tags,
			}
		} else {
			return TaggedName {
				Name: ns[0].(*p.Terminal).Value,
			}
		}
	}, ident, inlineTags)(s)
}

// InlineTags = ':' TagListStart* ident
func parseInlineTags(s p.Scanner) (p.ParsecNode, p.Scanner) {
	inits := p.Kleene(func (ns []p.ParsecNode) p.ParsecNode {
		tags := make([]string, len(ns))
		for i, n := range(ns) {
			tags[i] = n.(string)
		}
		return tags
	}, parseTagListStart)

	return p.And(func (ns []p.ParsecNode) p.ParsecNode {
		return append(ns[1].([]string), ns[2].(*p.Terminal).Value)
	}, colon, inits, ident)(s)
}

func ident(s p.Scanner) (p.ParsecNode, p.Scanner) {
	return p.Token(`^[^,{}:\r\n]+`, "IDENT")(s)
}

func comma(s p.Scanner) (p.ParsecNode, p.Scanner) {
	return p.Token(`^,`, "COMMA")(s)
}

func colon(s p.Scanner) (p.ParsecNode, p.Scanner) {
	return p.Token(`^:`, "COLON")(s)
}

func lbrace(s p.Scanner) (p.ParsecNode, p.Scanner) {
	return p.Token(`^{`, "LBRACE")(s)
}

func rbrace(s p.Scanner) (p.ParsecNode, p.Scanner) {
	return p.Token(`^}`, "RBRACE")(s)
}

func comment(s p.Scanner) (p.ParsecNode, p.Scanner) {
	return p.Token(`^//.*\n`, "COMMENT")(s)
}

func values(ns []p.ParsecNode) []string {
	values := make([]string, len(ns))
	for i, n := range(ns) {
		values[i] = strings.TrimSpace(n.(*p.Terminal).Value)
	}
	return values
}
