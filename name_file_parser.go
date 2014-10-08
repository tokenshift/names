package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	. "github.com/prataprc/goparsec"
)

type NameDB interface {
	// Add a name to the DB.
	Add(TaggedName)

	// Enumerates all names in the DB.
	Each() (<-chan TaggedName)

	// Applies a function to every name in the DB.
	ForEach(func (TaggedName))
}

type nameDB map[string]TaggedName

func (db nameDB) Add(n TaggedName) {
	if name, ok := db[n.Name]; ok {
		name = TaggedName {
			n.Name,
			append(name.Tags, n.Tags...),
		}
		db[n.Name] = name
	} else {
		db[n.Name] = n
	}
}

func (db nameDB) Each() (<-chan TaggedName) {
	out := make(chan TaggedName)

	go func() {
		for _, name := range(db) {
			out <- name
		}
		close(out)
	}()

	return out
}

func (db nameDB) ForEach(f func(TaggedName)) {
	for name := range(db.Each()) {
		f(name)
	}
}

type TaggedName struct {
	Name string
	Tags []string
}

func (n TaggedName) String() string {
	return fmt.Sprintf("%s: %s", n.Name, strings.Join(n.Tags, ", "))
}


/* Instead of constructing an AST, the individual parsers return nil and update
 * a global database of names+tags as a side effect.
 */

var tagStack [][]string
var currentDB NameDB

func parseNameFile(filename string) NameDB {
	buffer, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	db := nameDB(make(map[string]TaggedName))

	currentDB = db
	tagStack = make([][]string, 0)
	scanner := NewScanner(buffer)
	parseBlockContents(scanner)
	currentDB = nil

	return db
}

func parseBlock(s Scanner) (ParsecNode, Scanner) {
	return And(func (ns []ParsecNode) ParsecNode {
		tagStack = tagStack[:len(tagStack)-1]
		return struct{}{}
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
		tags := append(ns[0].([]string), ns[1].(string))
		tagStack = append(tagStack, tags)
		return struct{}{}
	}, inits, parseTagListEnd)(s)
}

func parseTagListStart(s Scanner) (ParsecNode, Scanner) {
	return And(func (ns []ParsecNode) ParsecNode {
		return strings.TrimSpace(ns[0].(*Terminal).Value)
	}, ident, comma)(s)
}

func parseTagListEnd(s Scanner) (ParsecNode, Scanner) {
	return And(func (ns []ParsecNode) ParsecNode {
		return strings.TrimSpace(ns[0].(*Terminal).Value)
	}, ident, lbrace)(s)
}

func parseBlockContents(s Scanner) (ParsecNode, Scanner) {
	return Kleene(func (ns []ParsecNode) ParsecNode {
		return struct{}{}
	}, parseBlockContent)(s)
}

func parseBlockContent(s Scanner) (ParsecNode, Scanner) {
	return OrdChoice(func (ns []ParsecNode) ParsecNode {
		if term, ok := ns[0].(*Terminal); ok && term.Name == "IDENT" {
			tagSet := make(map[string]bool)
			for _, tags := range(tagStack) {
				for _, tag := range(tags) {
					tagSet[tag] = true
				}
			}

			tags := make([]string, 0, len(tagSet))
			for tag, _ := range(tagSet) {
				tags = append(tags, tag)
			}
			currentDB.Add(TaggedName {
				Name: term.Value,
				Tags: tags,
			})
		}
		return struct{}{}
	}, comment, parseBlock, ident)(s)
}

func ident(s Scanner) (ParsecNode, Scanner) {
	return Token(`^[^,{}\r\n]+`, "IDENT")(s)
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
