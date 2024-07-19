package main

import (
	"fmt"
	"log"

	"github.com/stofffe/vgen/pkg/vgen"
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
		name: vgen.RulesRequired[string](
			vgen.Eq("bob"),
		),
		nickname: vgen.RulesOptional[string](
			vgen.Eq("bobby"),
		),
		age: vgen.RulesOptional[int](
			vgen.Gte(18),
		),
	}

	result, err := person.ValidatedConvert(rules)
	if err != nil {
		log.Fatal(err.Debug())
	}
	fmt.Println(result)
}
