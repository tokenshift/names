package main

import (
	"fmt"
)

func main() {
	db := parseNameFile("test.names")

	db.ForEach(func (name TaggedName) {
		fmt.Println(name)
	})
}
