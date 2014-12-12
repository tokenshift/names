package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

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

// A single name, with tags and other properties.
type Entry struct {
	Name string
	Type string
	Tags []string
}

// Tracks sets of tags in a push/pop stack.
type TagStack struct {
	// A single slice of tags is kept to make it easy to return all of the
	// tags at once.
	tags []string

	// Indices into the slice of all tags are kept to identify different
	// segments of the stack. When a segment is popped, everything after the
	// last segment index is discarded from the list of tags.
	segments []int
}

// Push a new segment of tags onto the stack.
func (stack *TagStack) Push(tags...string) {
	// Record the start index of the new tags.
	stack.segments = append(stack.segments, len(stack.tags))
	stack.tags = append(stack.tags, tags...)
}

// Pop the topmost segment of tags from the stack.
func (stack *TagStack) Pop() []string {
	if len(stack.segments) == 0 {
		return []string{}
	}

	index := stack.segments[len(stack.segments)-1]
	stack.segments = stack.segments[0:len(stack.segments)-1]
	if index >= len(stack.tags) {
		return []string{}
	}

	popped := stack.tags[index:len(stack.tags)]
	stack.tags = stack.tags[0:index]
	return popped
}

// Returns all tags in the tag stack.
func (stack TagStack) Tags() []string {
	return stack.tags
}


// Entry Point

func parseNameFiles(filenames []string) (<-chan Entry) {
	entries := make(chan Entry)

	go func() {
		defer close(entries)
		for _, filename := range(filenames) {
			for entry := range(parseNameFile(filename)) {
				entries <- entry
			}
		}
	}()

	return entries
}

func parseNameFile(filename string) (<-chan Entry) {
	buffer, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	block := parseBuffer(buffer)

	entries := make(chan Entry)

	go func() {
		defer close(entries)

		var tags TagStack
		sendNamesInBlock(block, tags, entries)
	}()

	return entries
}

func parseBuffer(buffer []byte) Block {
	scanner := p.NewScanner(buffer)
	result, _ := parseBlockContents(scanner)
	return result.(Block)
}

// Recursively iterates through names in this block, splitting them into
// components and sending them to the output channel.
func sendNamesInBlock(b Block, tags TagStack, out chan<- Entry) {
	for _, fullName := range(b.Names) {
		for _, entry := range(fullNameToComponents(fullName)) {
			entry.Tags = tags.Tags()
			out <- entry
		}
	}

	for _, child := range(b.Children) {
		tags.Push(child.Tags...)
		sendNamesInBlock(child.Block, tags, out)
		tags.Pop()
	}
}

// Breaks a full name into individual components, with their component type
// (e.g. "first", "last", "nick").
func fullNameToComponents(full string) []Entry {
	var entries []Entry

	comps := strings.Split(full, " ")
	for i, c := range(comps) {
		// Nicknames are denoted by surrounding them with double quotes.
		if strings.HasPrefix(c, "\"") && strings.HasSuffix(c, "\"") {
			entries = append(entries, Entry{
				Name: c[1:len(c)-1],
				Type: "nick",
			})
			continue
		}

		// Hyphenated names are broken into individual components.
		cs := strings.Split(c, "-")
		for _, c2 := range(cs) {
			if i == 0 {
				entries = append(entries, Entry{
					Name: c2,
					Type: "first",
				})
			}

			if i+1 < len(comps) {
				entries = append(entries, Entry{
					Name: c2,
					Type: "given",
				})
			}

			if i+1 == len(comps) {
				entries = append(entries, Entry{
					Name: c2,
					Type: "last",
				})
			}
		}
	}

	return entries
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
