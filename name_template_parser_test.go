package main

import (
	"reflect"
	. "testing"
)

func assertEquals(t *T, expected, actual interface{}) bool {
	if reflect.DeepEqual(expected, actual) {
		return true
	} else {
		t.Errorf("Expected %v <%T>, got %v <%T>",
			expected, expected, actual, actual)
		return false
	}
}

func TestParseTag(t *T) {
	// "Foo" -> (Tag "Foo")
	assertEquals(t,
		Tag("Foo"),
		parseNameTemplate("Foo"))

	// "Hello World" -> (Tag "Hello World")
	assertEquals(t,
		Tag("Hello World"),
		parseNameTemplate("Hello World"))
}

func TestAnd(t *T) {
	// "Foo + Bar"
	// -> (Or (And (Tag "Foo") (Tag "Bar"))
	assertEquals(t,
		Or([]And{
			And([]Chunk{
				Tag("Foo"),
				Tag("Bar"),
			}),
		}),
		parseNameTemplate("Foo + Bar"))
}

func TestOr(t *T) {
	// "Foo | Bar"
	// -> (Or (And (Tag "Foo")) (And (Tag "Bar")))
	assertEquals(t,
		Or([]And{
			And([]Chunk{
				Tag("Foo"),
			}),
			And([]Chunk{
				Tag("Bar"),
			}),
		}),
		parseNameTemplate("Foo | Bar"))
}

func TestNot(t *T) {
	// "Foo - Bar"
	// -> (Or (And (Tag "Foo") (Not "Bar")))
	assertEquals(t,
		Or([]And{
			And([]Chunk{
				Tag("Foo"),
				Not("Bar"),
			}),
		}),
		parseNameTemplate("Foo - Bar"))

	// "- Bar"
	// -> (Or (And (Not "Bar")))
	assertEquals(t,
		Or([]And{
			And([]Chunk{
				Not("Bar"),
			}),
		}),
		parseNameTemplate("- Bar"))

	// "-Bar"
	// -> (Or (And (Not "Bar")))
	assertEquals(t,
		Or([]And{
			And([]Chunk{
				Not("Bar"),
			}),
		}),
		parseNameTemplate("-Bar"))
}

func TestMultipleOperators(t *T) {
	// "Foo + Bar | Fizz - Buzz"
	// -> (Or
	//      (And (Tag "Foo") (Tag "Bar"))
	//      (And (Tag "Fizz") (Not "Buzz")))
	assertEquals(t,
		Or([]And{
			And([]Chunk{Tag("Foo"), Tag("Bar")}),
			And([]Chunk{Tag("Fizz"), Not("Buzz")}),
		}),
		parseNameTemplate("Foo + Bar | Fizz - Buzz"))
}

func TestOperatorAssociativity(t *T) {
	// "Alpha - Beta - Gamma - Delta"
	// -> (Or
	//      (And
	//        (Tag "Alpha")
	//        (Not "Beta")
	//        (Not "Gamma")
	//        (Not "Delta")))
	assertEquals(t,
		Or([]And{
			And([]Chunk{
				Tag("Alpha"),
				Not("Beta"),
				Not("Gamma"),
				Not("Delta"),
			}),
		}),
		parseNameTemplate("Alpha - Beta - Gamma - Delta"))

	// "Alpha + Beta + Gamma - Delta - Epsilon - Foxtrot | Omega"
	// -> (Or
	//      (And
	//        (Tag "Omega"))
	//      (And
	//        (Tag "Alpha")
	//        (Tag "Beta")
	//        (Tag "Gamma")
	//        (Not "Delta")
	//        (Not "Epsilon")
	//        (Not "Foxtrot")))
	assertEquals(t,
		Or([]And{
			And([]Chunk{
				Tag("Omega"),
			}),
			And([]Chunk{
				Tag("Alpha"),
				Tag("Beta"),
				Tag("Gamma"),
				Not("Delta"),
				Not("Epsilon"),
				Not("Foxtrot"),
			}),
		}),
		parseNameTemplate("Alpha + Beta + Gamma - Delta - Epsilon - Foxtrot | Omega"))
}

func TestOperatorPrecedence(t *T) {
	// "A + B | C - D"
	// -> (Or (And (Tag A) (Tag B)) (And (Tag C) (Not D)))
	assertEquals(t,
		Or([]And{
			And([]Chunk{Tag("A"), Tag("B")}),
			And{Tag("C"), Not("D")},
		}),
		parseNameTemplate("A + B | C - D"))

	// "A - B + C | D"
	// -> (Or (And (Tag A) (Not B) (Tag C)) (Tag D))
	assertEquals(t,
		Or([]And{
			And([]Chunk{Tag("A"), Not("B"), Tag("C")}),
			And([]Chunk{Tag("D")}),
		}),
		parseNameTemplate("A - B + C | D"))

	// "A | B - C + D"
	// -> (Or (And (Tag A))
	//        (And (Tag B) (Not C) (Tag D)))
	assertEquals(t,
		Or([]And{
			And([]Chunk{Tag("A")}),
			And([]Chunk{Tag("B"), Not("C"), Tag("D")}),
		}),
		parseNameTemplate("A | B - C + D"))
}

func TestFilters(t *T) {
	// ":foo" -> :foo
	assertEquals(t,
		Filter("foo"),
		parseNameTemplate(":foo"))

	// "A | :foo" -> (Or A :foo)
	assertEquals(t,
		Or([]And{
			And([]Chunk{Tag("A")}),
			And([]Chunk{Filter("foo")}),
		}),
		parseNameTemplate("A | :foo"))

	// "A:foo" -> (And A :foo)
	assertEquals(t,
		Or([]And{
			And([]Chunk{
				Tag("A"),
				Filter("foo"),
			}),
		}),
		parseNameTemplate("A:foo"))

	// "A + B:foo" -> (And A B:foo)
	assertEquals(t,
		Or([]And{
			And([]Chunk{
				Tag("A"),
				Filtered{"B", "foo"},
			}),
		}),
		parseNameTemplate("A + B:foo"))

	// "A:foo + B" -> (And A:foo B)
	assertEquals(t,
		Or([]And{
			And([]Chunk{
				Filtered{"A", "foo"},
				Tag("B"),
			}),
		}),
		parseNameTemplate("A:foo + B"))
}

func TestMaybe(t *T) {
	// "[A]" -> (Maybe A)
	assertEquals(t,
		Maybe{
			Or([]And{
				And([]Chunk{
					Tag("A"),
				}),
			}),
		},
		parseNameTemplate("[A]"))

	// "[A + B | C:foo]" -> (Maybe (Or (And A B) (And C:foo)))
	assertEquals(t,
		Maybe{
			Or([]And{
				And([]Chunk{
					Tag("A"),
					Tag("B"),
				}),
				And([]Chunk{
					Filtered{"C", "foo"},
				}),
			}),
		},
		parseNameTemplate("[A + B | C:foo]"))
}
