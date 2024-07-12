package main

import (
	"fmt"

	"github.com/stofffe/vgen"
)

// vgen
type Person struct {
	name     string
	nickname string
	age      int
}

func main() {
	age := 12
	person := PersonVgen{
		name:     nil,
		nickname: nil,
		age:      &age,
	}

	rules := PersonRules{
		name: vgen.NewRules( // required
			vgen.Required[string](),
			vgen.Eq("bob"),
		),
		nickname: vgen.NewRules( // not required
			vgen.Eq("bobsson"),
		),
		age: vgen.NewRules(
			vgen.Gte(18),
		),
	}

	err := person.Validate(rules)
	fmt.Println(err.Debug())
}
