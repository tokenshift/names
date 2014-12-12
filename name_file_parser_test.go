package main

import (
	. "testing"
)

func TestParseSingleName(t *T) {
	assertEquals(t, Block{
		Names: []string{"William Wallace"},
	}, parseBuffer([]byte("William Wallace")))
}

func TestParseComment(t *T) {
	assertEquals(t, Block{}, parseBuffer([]byte("// This is a comment.")))
}

func TestParseCommentsAndNames(t *T) {
	assertEquals(t, Block{
		Names: []string{"William", "Wallace"},
	}, parseBuffer([]byte("William // This is a comment\n Wallace")))
}

func TestParseSingleTaggedBlock(t *T) {
	assertEquals(t, Block{
		Children: []TaggedBlock{
			TaggedBlock{
				Tags: []string{"Tag1"},
			},
		},
	}, parseBuffer([]byte("Tag1 {}")))
}

func TestParseMultiTaggedBlock(t *T) {
	assertEquals(t, Block{
		Children: []TaggedBlock{
			TaggedBlock{
				Tags: []string{"Tag1", "Tag2", "Tag3"},
			},
		},
	}, parseBuffer([]byte("Tag1, Tag2, Tag3 {}")))
}

func TestParseTaggedName(t *T) {
	assertEquals(t, Block{
		Children: []TaggedBlock{
			TaggedBlock{
				Tags: []string{"Braveheart", "Movie"},
				Block: Block{
					Names: []string{"William Wallace"},
				},
			},
		},
	}, parseBuffer([]byte("William Wallace:Braveheart,Movie")))
}

func TestParseNickName(t *T) {
	assertEquals(t, Block{
		Names: []string{"James \"Jimmy\" Douglas"},
	}, parseBuffer([]byte("James \"Jimmy\" Douglas")))
}

func TestParseHyphenatedeName(t *T) {
	assertEquals(t, Block{
		Names: []string{"James Clarence-Jones"},
	}, parseBuffer([]byte("James Clarence-Jones")))
}
