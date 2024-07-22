package main

import (
	"fmt"
	"log"

	"github.com/stofffe/vgen/pkg/vgen"
)

// vgen(include)
type Person struct {
	name string
	auth string
}

func main() {
	name := "bob"
	auth := "member"

	person := PersonVgen{
		name: &name,
		auth: &auth,
	}

	meberRules := PersonRules{
		name: vgen.RulesRequired(
			vgen.StringNotEmpty(),
		),
		auth: vgen.RulesRequired(
			vgen.StringNotEmpty(),
		),
	}

	adminRules := meberRules.Extend(PersonRules{
		auth: vgen.RulesRequired(
			vgen.Eq("admin"),
		),
	})

	result, verr := person.ValidatedConvert(adminRules)
	if verr != nil {
		log.Fatal(verr.Debug())
	}
	fmt.Println(result)

}
