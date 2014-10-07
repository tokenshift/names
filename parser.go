package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	. "github.com/prataprc/goparsec"
)

type block struct {
	tags []string
	nameFile
}

type nameFile struct {
	names []string
	blocks []block
}

func (b block) String() string {
	var buffer bytes.Buffer
	fmt.Fprint(&buffer, strings.Join(b.tags, ", "))
	fmt.Fprint(&buffer, "{\n")
	fmt.Fprint(&buffer, b.nameFile.String())
	fmt.Fprint(&buffer, "}\n")
	return buffer.String()
}

func (bc nameFile) String() string {
	var buffer bytes.Buffer
	for _, name := range(bc.names) {
		fmt.Fprint(&buffer, name)
		fmt.Fprint(&buffer, "\n")
	}
	for _, child := range(bc.blocks) {
		fmt.Fprint(&buffer, child.String())
	}
	return buffer.String()
}

func parseNameFile(filename string) nameFile {
	buffer, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	scanner := NewScanner(buffer)
	result, _ := parseBlockContents(scanner)
	names, ok := result.(nameFile)
	if !ok {
		fmt.Fprintf(os.Stderr, "Unexpected return type from parsing: %T\n", result)
		os.Exit(1)
	}

	return names
}

func parseBlock(s Scanner) (ParsecNode, Scanner) {
	return And(func (ns []ParsecNode) ParsecNode {
		var b block
		b.tags = ns[0].([]string)
		b.names = ns[1].(nameFile).names
		b.blocks = ns[1].(nameFile).blocks
		return b
	}, parseTagList, parseBlockContents, rbrace)(s)
}

func parseTagList(s Scanner) (ParsecNode, Scanner) {
	inits := Kleene(func (ns []ParsecNode) ParsecNode {
		tags := make([]string, len(ns))
		for i, n := range(ns) {
			tags[i] = n.(string)
		}
		return tags
	}, parseTagListStart)

	return And(func (ns []ParsecNode) ParsecNode {
		return append(ns[0].([]string), ns[1].(string))
	}, inits, parseTagListEnd)(s)
}

func parseTagListStart(s Scanner) (ParsecNode, Scanner) {
	return And(func (ns []ParsecNode) ParsecNode {
		return ns[0].(*Terminal).Value
	}, ident, comma)(s)
}

func parseTagListEnd(s Scanner) (ParsecNode, Scanner) {
	return And(func (ns []ParsecNode) ParsecNode {
		return ns[0].(*Terminal).Value
	}, ident, lbrace)(s)
}

func parseBlockContents(s Scanner) (ParsecNode, Scanner) {
	return Kleene(func (ns []ParsecNode) ParsecNode {
		names := make([]string, 0)
		blocks := make([]block, 0)
		for _, n := range(ns) {
			if b, ok := n.(block); ok {
				names = append(names, b.names...)
				blocks = append(blocks, b.blocks...)
			} else {
				fmt.Fprintf(os.Stderr, "Expected a block")
				os.Exit(1)
			}
		}
		return nameFile {
			names: names,
			blocks: blocks,
		}
	}, parseBlockContent)(s)
}

func parseBlockContent(s Scanner) (ParsecNode, Scanner) {
	return OrdChoice(func (ns []ParsecNode) ParsecNode {
		var b block
		if name, ok := ns[0].(*Terminal); ok {
			b.names = []string{name.Value}
		} else if b2, ok := ns[0].(block); ok {
			b.blocks = []block{b2}
		} else {
			fmt.Fprintln(os.Stderr, "Expected a block or name.")
			os.Exit(1)
		}

		return b
	}, parseBlock, ident)(s)
}

func ident(s Scanner) (ParsecNode, Scanner) {
	return Token(`^[^,{}\n]+`, "NAME")(s)
}

func comma(s Scanner) (ParsecNode, Scanner) {
	return Token(`^,`, "COMMA")(s)
}

func lbrace(s Scanner) (ParsecNode, Scanner) {
	return Token(`^{`, "LBRACE")(s)
}

func rbrace(s Scanner) (ParsecNode, Scanner) {
	return Token(`^}`, "RBRACE")(s)
}

func values(ns []ParsecNode) []string {
	values := make([]string, len(ns))
	for i, n := range(ns) {
		values[i] = strings.TrimSpace(n.(*Terminal).Value)
	}
	return values
}
