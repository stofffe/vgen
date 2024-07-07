package main

import (
	"fmt"
	"log"
)

func main() {
	path := "examples/person/person.go"
	count, err := HandleFile(path)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(count)
}
