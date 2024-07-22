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
	name := ""
	person := PersonVgen{
		name:     &name,
		nickname: nil,
		age:      &age,
	}

	rules := PersonRules{
		name: vgen.RulesRequired(
			vgen.StringNotEmpty(),
		),
		nickname: vgen.RulesRequired(
			vgen.Eq("bobby"),
		),
		age: vgen.RulesOptional(
			vgen.Gte(18),
		),
	}

	result, err := person.ValidatedConvert(rules)
	if err != nil {
		log.Fatal(err.Debug())
	}
	fmt.Println(result)
}
