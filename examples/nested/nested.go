package main

import (
	"fmt"

	"github.com/stofffe/vgen"
)

// If json tag exists that name will be used as alias

// vgen(include)
type Person struct {
	Name string // vgen( alias=name )
	// vgen(
	//   nested,
	//   alias=pet
	// )
	Pet Pet
}

// vgen
type Pet struct {
	Name string // vgen(alias=name)
}

func main() {
	name := "bo"
	petName := "boby"
	person := PersonVgen{
		Name: &name,
		Pet: &PetVgen{
			Name: &petName,
		},
	}

	rules := PersonRules{
		Name: vgen.NewRules(true,
			vgen.Eq("bob"),
		),
		Pet: PetRules{
			Name: vgen.NewRules(true,
				vgen.Eq("bobby"),
			),
		},
	}

	err := person.Validate(rules)
	fmt.Println(err.Debug())
}
