package main

import (
	"fmt"
	"log"

	"github.com/stofffe/vgen/examples"
)

func main() {

	person, err := examples.PersonVgen{
		Name:  P(""),
		Age:   P(123),
		Vibes: nil,
	}.Validate()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("person: %v\n", person)
}

func P[T any](t T) *T {
	return &t
}
