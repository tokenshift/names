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
	// "Foo + Bar" -> (And (Tag "Foo") (Tag "Bar"))
	assertEquals(t,
		And{Tag("Foo"), Tag("Bar")},
		parseNameTemplate("Foo + Bar"))
}

func TestOr(t *T) {
	// "Foo | Bar" -> (Or (Tag "Foo") (Tag "Bar"))
	assertEquals(t,
		Or{Tag("Foo"), Tag("Bar")},
		parseNameTemplate("Foo | Bar"))
}

func TestNot(t *T) {
	// "Foo - Bar" -> (And (Tag "Foo") (Not (Tag "Bar")))
	assertEquals(t,
		And{Tag("Foo"), Not{Tag("Bar")}},
		parseNameTemplate("Foo - Bar"))

	// "- Bar" -> (Not (Tag "Bar"))
	assertEquals(t,
		Not{Tag("Bar")},
		parseNameTemplate("- Bar"))

	// "-Bar" -> (Not (Tag "Bar"))
	assertEquals(t,
		Not{Tag("Bar")},
		parseNameTemplate("-Bar"))
}

func TestMultipleOperators(t *T) {
	// "Foo + Bar | Fizz - Buzz"
	// -> (Or
	//      (And (Tag "Foo") (Tag "Bar"))
	//      (And (Tag "Fizz") (Not (Tag "Buzz"))))
	assertEquals(t,
		Or{
			And{Tag("Foo"), Tag("Bar")},
			And{Tag("Fizz"), Not{Tag("Buzz")}},
		},
		parseNameTemplate("Foo + Bar | Fizz - Buzz"))
}

func TestOperatorPrecedence(t *T) {
	// <Strong>
	// Not
	// And
	// Or
	// <Weak>

	// "A + B | C - D" -> (Or (And A B) (And (C (Not D))))
	assertEquals(t,
		Or{
			And{Tag("A"), Tag("B")},
			And{Tag("C"), Not{Tag("D")}},
		},
		parseNameTemplate("A + B | C - D"))

	// "A - B + C | D" -> (Or (And (And (A (Not B))) C) D)
	assertEquals(t,
		Or{
			And{
				And{Tag("A"), Not{Tag("B")}},
				Tag("C"),
			},
			Tag("D"),
		},
		parseNameTemplate("A - B + C | D"))

	// "A | B - C + D" -> (Or A (And (And B (Not C)) D))
	assertEquals(t,
		Or{
			Tag("A"),
			And{
				And{
					Tag("B"),
					Not{Tag("C")},
				},
				Tag("D"),
			},
		},
		parseNameTemplate("A | B - C + D"))
}

func TestGrouping(t *T) {
	// "A + B | C" -> (Or (And A B) C))
	assertEquals(t,
		Or{
			And{Tag("A"), Tag("B")},
			Tag("C"),
		},
		parseNameTemplate("A + B | C"))

	// "A + (B | C)" -> (And A (Or B C))
	assertEquals(t,
		And{
			Tag("A"),
			Or{Tag("B"), Tag("C")},
		},
		parseNameTemplate("A + (B | C)"))

	// "A - B | C" -> (Or (And A (Not B)) C)
	assertEquals(t,
		Or{
			And{Tag("A"), Not{Tag("B")}},
			Tag("C"),
		},
		parseNameTemplate("A - B | C"))

	// "A - (B | C)" -> (And A (Not (Or B C)))
	assertEquals(t,
		And{
			Tag("A"),
			Not{Or{Tag("B"), Tag("C")}},
		},
		parseNameTemplate("A + B | C"))

	// "A - (B + C)" -> (And A (Not (And B C)))
	assertEquals(t,
		Or{
			And{Tag("A"), Tag("B")},
			Tag("C"),
		},
		parseNameTemplate("A + B | C"))
}

func TestFilters(t *T) {
	// ":foo" -> :foo
	assertEquals(t,
		Filter("foo"),
		parseNameTemplate(":foo"))

	// "A | :foo" -> (Or A :foo)
	assertEquals(t,
		Or{Tag("A"), Filter("foo")},
		parseNameTemplate("A | :foo"))

	// "A:foo" -> (And A :foo)
	assertEquals(t,
		And{Tag("A"), Filter("foo")},
		parseNameTemplate("A:foo"))

	// "A + B:foo" -> (And A (And B :foo))
	assertEquals(t,
		And{
			Tag("A"),
			And{Tag("B"), Filter("foo")},
		},
		parseNameTemplate("A + B:foo"))

	// "A:foo + B" -> (And (And A :foo) B)
	assertEquals(t,
		And{
			And{Tag("A"), Filter("foo")},
			Tag("B"),
		},
		parseNameTemplate("A:foo + B"))

	// "(A + B):foo" -> (And (And A B) :foo)
	assertEquals(t,
		And{
			And{Tag("A"), Tag("B")},
			Filter("foo"),
		},
		parseNameTemplate("(A + B):foo"))
}

func TestMaybe(t *T) {
	// "[A]" -> (Maybe A)
	assertEquals(t,
		Maybe{Tag("A")},
		parseNameTemplate("[A]"))

	// "[A + B | C:foo]" -> (Maybe (Or (And A B) (And C :foo)))
	assertEquals(t,
		Maybe{Or{
			And{Tag("A"), Tag("B")},
			And{Tag("C"), Filter("foo")},
		}},
		parseNameTemplate("[A + B | C:foo]"))
}
