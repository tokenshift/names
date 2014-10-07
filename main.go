package main

import (
	"fmt"
)

func main() {
	nameFile := parseNameFile("test.names")
	fmt.Println(nameFile)
}
