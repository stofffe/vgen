package main

import (
	"fmt"

	"github.com/stofffe/vgen"
)

// If json tag exists that name will be used as alias

// vgen(include)
type Person struct {
	Name  string   // vgen(alias=name)
	Pet   Pet      // vgen(n, alias=pet)
	Pets  []Pet    // vgen(n, alias=pets)
	Names []string // vgen(alias=names)
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
		Pets: &[]PetVgen{
			{Name: &petName},
			{Name: &petName},
		},
		Names: &[]string{},
	}

	petRules := PetRules{
		Name: vgen.NewRules(true,
			vgen.Eq("bobby"),
		),
	}

	personRules := PersonRules{
		Name: vgen.NewRules(true,
			vgen.Eq("bob"),
		),
		Pet: vgen.NewRules(true,
			vgen.Nested(petRules),
		),
		Pets: vgen.NewRules(true,
			vgen.LenGt[[]PetVgen](3),
			vgen.List(
				vgen.Nested(petRules),
			),
		),
	}

	err := person.Validate(personRules)
	fmt.Println(err.Debug())
}
