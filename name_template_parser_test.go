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
	// "Foo" -> (Or (And (Tag "Foo"))))
	result, _ := parseNameTemplate("Foo")
	assertEquals(t,
		Or([]And{And([]Matcher{Tag("Foo")})}),
		result)

	// "Hello World" -> (Or (And (Tag "Hello World")))
	result, _ = parseNameTemplate("Hello World")
	assertEquals(t,
		Or([]And{And([]Matcher{Tag("Hello World")})}),
		result)
}

func TestAnd(t *T) {
	// "Foo + Bar"
	// -> (Or (And (Tag "Foo") (Tag "Bar"))
	result, _ := parseNameTemplate("Foo + Bar")
	assertEquals(t,
		Or([]And{
			And([]Matcher{
				Tag("Foo"),
				Tag("Bar"),
			}),
		}),
		result)

	// "Foo + Bar + Fizz + Buzz"
	// -> (Or (And (Tag "Foo") (Tag "Bar") (Tag "Fizz") (Tag "Buzz"))
	result, _ = parseNameTemplate("Foo + Bar")
	assertEquals(t,
		Or([]And{
			And([]Matcher{
				Tag("Foo"),
				Tag("Bar"),
			}),
		}),
		result)
}

func TestOr(t *T) {
	// "Foo | Bar"
	// -> (Or (And (Tag "Foo")) (And (Tag "Bar")))
	result, _ := parseNameTemplate("Foo | Bar")
	assertEquals(t,
		Or([]And{
			And([]Matcher{
				Tag("Foo"),
			}),
			And([]Matcher{
				Tag("Bar"),
			}),
		}),
		result)
}

func TestNot(t *T) {
	// "Foo - Bar"
	// -> (Or (And (Tag "Foo") (Not "Bar")))
	result, _ := parseNameTemplate("Foo - Bar")
	assertEquals(t,
		Or([]And{
			And([]Matcher{
				Tag("Foo"),
				Not{Tag("Bar")},
			}),
		}),
		result)

	// "- Bar"
	// -> (Or (And (Not "Bar")))
	result, _ = parseNameTemplate("- Bar")
	assertEquals(t,
		Or([]And{
			And([]Matcher{
				Not{Tag("Bar")},
			}),
		}),
		result)

	// "-Bar"
	// -> (Or (And (Not "Bar")))
	result, _ = parseNameTemplate("-Bar")
	assertEquals(t,
		Or([]And{
			And([]Matcher{
				Not{Tag("Bar")},
			}),
		}),
		result)
}

func TestMultipleOperators(t *T) {
	// "Foo + Bar | Fizz - Buzz"
	// -> (Or
	//      (And (Tag "Foo") (Tag "Bar"))
	//      (And (Tag "Fizz") (Not "Buzz")))
	result, _ := parseNameTemplate("Foo + Bar | Fizz - Buzz")
	assertEquals(t,
		Or([]And{
			And([]Matcher{Tag("Foo"), Tag("Bar")}),
			And([]Matcher{Tag("Fizz"), Not{Tag("Buzz")}}),
		}),
		result)
}

func TestOperatorAssociativity(t *T) {
	// "Alpha - Beta - Gamma - Delta"
	// -> (Or
	//      (And
	//        (Tag "Alpha")
	//        (Not "Beta")
	//        (Not "Gamma")
	//        (Not "Delta")))
	result, _ := parseNameTemplate("Alpha - Beta - Gamma - Delta")
	assertEquals(t,
		Or([]And{
			And([]Matcher{
				Tag("Alpha"),
				Not{Tag("Beta")},
				Not{Tag("Gamma")},
				Not{Tag("Delta")},
			}),
		}),
		result)

	// "Alpha + Beta + Gamma - Delta - Epsilon - Foxtrot | Omega"
	// -> (Or
	//      (And
	//        (Tag "Alpha")
	//        (Tag "Beta")
	//        (Tag "Gamma")
	//        (Not "Delta")
	//        (Not "Epsilon")
	//        (Not "Foxtrot"))
	//      (And
	//        (Tag "Omega")))
	result, _ = parseNameTemplate("Alpha + Beta + Gamma - Delta - Epsilon - Foxtrot | Omega")
	assertEquals(t,
		Or([]And{
			And([]Matcher{
				Tag("Alpha"),
				Tag("Beta"),
				Tag("Gamma"),
				Not{Tag("Delta")},
				Not{Tag("Epsilon")},
				Not{Tag("Foxtrot")},
			}),
			And([]Matcher{
				Tag("Omega"),
			}),
		}),
		result)
}

func TestOperatorPrecedence(t *T) {
	// "A + B | C - D"
	// -> (Or (And (Tag A) (Tag B)) (And (Tag C) (Not D)))
	result, _ := parseNameTemplate("A + B | C - D")
	assertEquals(t,
		Or([]And{
			And([]Matcher{Tag("A"), Tag("B")}),
			And{Tag("C"), Not{Tag("D")}},
		}),
		result)

	// "A - B + C | D"
	// -> (Or (And (Tag A) (Not B) (Tag C)) (Tag D))
	result, _ = parseNameTemplate("A - B + C | D")
	assertEquals(t,
		Or([]And{
			And([]Matcher{Tag("A"), Not{Tag("B")}, Tag("C")}),
			And([]Matcher{Tag("D")}),
		}),
		result)

	// "A | B - C + D"
	// -> (Or (And (Tag A))
	//        (And (Tag B) (Not C) (Tag D)))
	result, _ = parseNameTemplate("A | B - C + D")
	assertEquals(t,
		Or([]And{
			And([]Matcher{Tag("A")}),
			And([]Matcher{Tag("B"), Not{Tag("C")}, Tag("D")}),
		}),
		result)
}

func TestFilters(t *T) {
	// ":foo" -> :foo
	result, _ := parseNameTemplate(":foo")
	assertEquals(t,
		Or([]And{
			And([]Matcher{
				Filter("foo"),
			}),
		}),
		result)

	// "A | :foo" -> (Or A :foo)
	result, _ = parseNameTemplate("A | :foo")
	assertEquals(t,
		Or([]And{
			And([]Matcher{Tag("A")}),
			And([]Matcher{Filter("foo")}),
		}),
		result)

	// "A:foo" -> (A:foo)
	result, _ = parseNameTemplate("A:foo")
	assertEquals(t,
		Or([]And{
			And([]Matcher{
				Filtered{"A", "foo"},
			}),
		}),
		result)

	// "A + B:foo" -> (And A B:foo)
	result, _ = parseNameTemplate("A + B:foo")
	assertEquals(t,
		Or([]And{
			And([]Matcher{
				Tag("A"),
				Filtered{"B", "foo"},
			}),
		}),
		result)

	// "A:foo + B" -> (And A:foo B)
	result, _ = parseNameTemplate("A:foo + B")
	assertEquals(t,
		Or([]And{
			And([]Matcher{
				Filtered{"A", "foo"},
				Tag("B"),
			}),
		}),
		result)
}

func TestMaybe(t *T) {
	// "[A]" -> (Maybe A)
	result, _ := parseNameTemplate("[A]")
	assertEquals(t,
		Maybe{
			Or([]And{
				And([]Matcher{
					Tag("A"),
				}),
			}),
		},
		result)

	// "[A + B | C:foo]" -> (Maybe (Or (And A B) (And C:foo)))
	result, _ = parseNameTemplate("[A + B | C:foo]")
	assertEquals(t,
		Maybe{
			Or([]And{
				And([]Matcher{
					Tag("A"),
					Tag("B"),
				}),
				And([]Matcher{
					Filtered{"C", "foo"},
				}),
			}),
		},
		result)
}

func TestParsingGarbage(t *T) {
	result, err := parseNameTemplate("% J#QOQ# ^#Q#")
	if err == nil {
		t.Errorf("Parsing should have failed.")
	}
	assertEquals(t, nil, result)
}
