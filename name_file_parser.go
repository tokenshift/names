package main

import (
	"fmt"
	"io/ioutil"
	"os"

	p "github.com/prataprc/goparsec"
)

// Data Definitions and AST

// A block (delimited by curly braces) which may contain names and/or other
// blocks.
type Block struct {
	Names []string
	Children []TaggedBlock
}

func (b Block) String() string {
	return fmt.Sprintf("{Names: %v Children: %v}", b.Names, b.Children)
}

// A block with tags that apply to all of its contents.
type TaggedBlock struct {
	Tags []string
	Block
}

func (tb TaggedBlock) String() string {
	return fmt.Sprintf("Tags: %v %v", tb.Tags, tb.Block)
}

// Entry Point

func parseNameFiles(filenames []string) Block {
	block := Block{}

	for _, filename := range(filenames) {
		block = mergeBlocks(block, parseNameFile(filename))
	}

	return block
}

func parseNameFile(filename string) Block {
	buffer, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	return parseBuffer(buffer)
}

func parseBuffer(buffer []byte) Block {
	scanner := p.NewScanner(buffer)
	result, _ := parseBlockContents(scanner)
	return result.(Block)
}

// Non-Terminals

// Block Contents: the contents inside the curly braces of a block (not
// including the curly braces). This can consist of names and tagged blocks.
func parseBlockContents(s p.Scanner) (p.ParsecNode, p.Scanner) {
	entry := p.OrdChoice(func (ns []p.ParsecNode) p.ParsecNode {
		return ns[0]
	}, comment, parseTaggedBlock, parseName)

	return p.Kleene(func (ns []p.ParsecNode) p.ParsecNode {
		block := Block{}
		for _, n := range(ns) {
			if child, ok := n.(TaggedBlock); ok {
				block.Children = append(block.Children, child)
			} else if name, ok := n.(string); ok {
				block.Names = append(block.Names, name)
			}
		}
		return block
	}, entry)(s)
}

// A series of tags followed by a block, delimited by curly braces.
func parseTaggedBlock(s p.Scanner) (p.ParsecNode, p.Scanner) {
	return p.And(func (ns []p.ParsecNode) p.ParsecNode {
		return TaggedBlock{
			Tags: ns[0].([]string),
			Block: ns[2].(Block),
		}
	}, tags, lbrace, parseBlockContents, rbrace)(s)
}

// A single name, potentially tagged inline ("John Smith:tag1, tag2").
func parseName(s p.Scanner) (p.ParsecNode, p.Scanner) {
	withTags := p.And(func (ns []p.ParsecNode) p.ParsecNode {
		// Inline tags are just shorthand for a tagged block with a single
		// name.
		return TaggedBlock{
			Tags: ns[2].([]string),
			Block: Block{
				Names: []string{ns[0].(string)},
			},
		}
	}, name, colon, tags)

	return p.OrdChoice(func (ns []p.ParsecNode) p.ParsecNode {
		return ns[0]
	}, withTags, name)(s)
}

// A list of comma-delimited tags.
var tags = p.Kleene(func (ns []p.ParsecNode) p.ParsecNode {
	ts := make([]string, len(ns))
	for i, n := range(ns) {
		ts[i] = string(n.(Tag))
	}
	return ts
}, tag, comma)

func mergeBlocks(a, b Block) Block {
	return Block{
		Names: append(a.Names, b.Names...),
		Children: append(a.Children, b.Children...),
	}
}
