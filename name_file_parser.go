package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	. "github.com/prataprc/goparsec"
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
	scanner := NewScanner(buffer)
	parseBlockContents(scanner)
}

// Block = TagList BlockContents '}'
func parseBlock(s Scanner) (ParsecNode, Scanner) {
	return And(func (ns []ParsecNode) ParsecNode {
		tagStack = tagStack[:len(tagStack)-1]
		return struct{}{}
	}, parseTagList, parseBlockContents, rbrace)(s)
}

// TagList = TagListStart* TagListEnd
func parseTagList(s Scanner) (ParsecNode, Scanner) {
	inits := Kleene(func (ns []ParsecNode) ParsecNode {
		tags := make([]string, len(ns))
		for i, n := range(ns) {
			tags[i] = n.(string)
		}
		return tags
	}, parseTagListStart)

	return And(func (ns []ParsecNode) ParsecNode {
		tags := append(ns[0].([]string), ns[1].(string))
		tagStack = append(tagStack, tags)
		return struct{}{}
	}, inits, parseTagListEnd)(s)
}

// TagListStart = ident ','
func parseTagListStart(s Scanner) (ParsecNode, Scanner) {
	return And(func (ns []ParsecNode) ParsecNode {
		return strings.TrimSpace(ns[0].(*Terminal).Value)
	}, ident, comma)(s)
}

// TagListEnd = ident '{'
func parseTagListEnd(s Scanner) (ParsecNode, Scanner) {
	return And(func (ns []ParsecNode) ParsecNode {
		return strings.TrimSpace(ns[0].(*Terminal).Value)
	}, ident, lbrace)(s)
}

// BlockContents = BlockContent*
func parseBlockContents(s Scanner) (ParsecNode, Scanner) {
	return Kleene(func (ns []ParsecNode) ParsecNode {
		return struct{}{}
	}, parseBlockContent)(s)
}

// BlockContent = comment | Block | Name
func parseBlockContent(s Scanner) (ParsecNode, Scanner) {
	return OrdChoice(func (ns []ParsecNode) ParsecNode {
		if name, ok := ns[0].(TaggedName); ok {
			name.Tags = append(name.Tags, tags()...)
			fmt.Println(name)
		}
		return struct{}{}
	}, comment, parseBlock, parseName)(s)
}

// Name = ident [InlineTags]
func parseName(s Scanner) (ParsecNode, Scanner) {
	inlineTags := Maybe(func (ns []ParsecNode) ParsecNode {
		if ns == nil {
			return []string {}
		} else {
			return ns[0].([]string)
		}
	}, parseInlineTags)

	return And(func (ns []ParsecNode) ParsecNode {
		if tags, ok := ns[1].([]string); ok {
			return TaggedName {
				Name: ns[0].(*Terminal).Value,
				Tags: tags,
			}
		} else {
			return TaggedName {
				Name: ns[0].(*Terminal).Value,
			}
		}
	}, ident, inlineTags)(s)
}

// InlineTags = ':' TagListStart* ident
func parseInlineTags(s Scanner) (ParsecNode, Scanner) {
	inits := Kleene(func (ns []ParsecNode) ParsecNode {
		tags := make([]string, len(ns))
		for i, n := range(ns) {
			tags[i] = n.(string)
		}
		return tags
	}, parseTagListStart)

	return And(func (ns []ParsecNode) ParsecNode {
		return append(ns[1].([]string), ns[2].(*Terminal).Value)
	}, colon, inits, ident)(s)
}

func ident(s Scanner) (ParsecNode, Scanner) {
	return Token(`^[^,{}:\r\n]+`, "IDENT")(s)
}

func comma(s Scanner) (ParsecNode, Scanner) {
	return Token(`^,`, "COMMA")(s)
}

func colon(s Scanner) (ParsecNode, Scanner) {
	return Token(`^:`, "COLON")(s)
}

func lbrace(s Scanner) (ParsecNode, Scanner) {
	return Token(`^{`, "LBRACE")(s)
}

func rbrace(s Scanner) (ParsecNode, Scanner) {
	return Token(`^}`, "RBRACE")(s)
}

func comment(s Scanner) (ParsecNode, Scanner) {
	return Token(`^//.*\n`, "COMMENT")(s)
}

func values(ns []ParsecNode) []string {
	values := make([]string, len(ns))
	for i, n := range(ns) {
		values[i] = strings.TrimSpace(n.(*Terminal).Value)
	}
	return values
}
