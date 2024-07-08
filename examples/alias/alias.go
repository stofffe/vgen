package main

import (
	"fmt"

	"github.com/stofffe/vgen"
)

// If json tag exists that name will be used as alias

// vgen
type Person struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func main() {
	age := 12
	person := PersonVgen{
		Name: nil,
		Age:  &age,
	}

	rules := PersonRules{
		Name: vgen.NewRules(true,
			vgen.Eq("bob"),
		),
		Age: vgen.NewRules(true,
			vgen.CustomMessage("not allowed to drive", vgen.Gte(18)),
			vgen.OneOf(1, 2, 3, 4, 5),
			vgen.NotOneOf(10, 11, 12),
		),
	}

	err := person.Validate(rules)
	fmt.Println(err.Debug())
}
