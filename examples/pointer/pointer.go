package main

import (
	"fmt"
	"log"

	"github.com/stofffe/vgen"
)

// vgen(include)
type Person struct {
	Name     *string `json:"name"`     // vgen(alias=name)
	Nickname *string `json:"nickname"` // vgen(alias=nickname)
}

func main() {
	nickname := "bob"
	nicknamePtr := &nickname
	person := PersonVgen{
		Nickname: &nicknamePtr,
	}

	rules := PersonRules{
		Name: vgen.RulesRequired[*string](),
		Nickname: vgen.RulesOptional(
			vgen.Deref(
				vgen.NotOneOf("bob", "bobby"),
			),
		),
	}

	result, verr := person.ValidatedConvert(rules)
	if verr != nil {
		log.Fatal(verr.Debug())
	}
	fmt.Println(result)
}
