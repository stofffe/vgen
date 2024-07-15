package main

import (
	"fmt"
	"log"

	"github.com/stofffe/vgen"
)

// vgen(include)
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
		name: vgen.RulesRequired(
			vgen.Eq("bob"),
		),
		nickname: vgen.RulesOptional(
			vgen.Eq("bobsson"),
		),
		age: vgen.RulesRequired(
			vgen.Gte(18),
		),
	}

	p, err := person.ValidatedConvert(rules)
	if err.HasError() {
		log.Fatal(err.Debug())
	}
	fmt.Println(p)
}
