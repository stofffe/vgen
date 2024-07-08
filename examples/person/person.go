package main

import (
	"fmt"

	"github.com/stofffe/vgen"
)

// vgen
type Person struct {
	name string
	age  int
}

func main() {
	age := 12
	person := PersonVgen{
		name: nil,
		age:  &age,
	}

	rules := PersonRules{
		name: vgen.NewRules(true,
			vgen.Eq("bob"),
		),
		age: vgen.NewRules(true,
			vgen.CustomMessage("not allowed to drive", vgen.Gte(18)),
			vgen.OneOf(1, 2, 3, 4, 5),
			vgen.NotOneOf(10, 11, 12),
		),
	}

	err := person.Validate(rules)
	fmt.Println(err.Debug())
}
