package main

import (
	"fmt"

	"github.com/stofffe/vgen"
)

// vgen
type Person struct {
	name      string
	age       int
	nicknames []string
	stats     map[string]int
}

func main() {
	age := 12
	stats := make(map[string]int)
	stats["stamina"] = 25
	person := PersonVgen{
		name: nil,
		age:  &age,
		nicknames: &[]string{
			"notbob1",
			"bob2",
		},
		stats: &stats,
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
		nicknames: vgen.NewRules(true,
			vgen.List(
				vgen.OneOf("bob1", "bob2"),
			),
		),
		matrix: vgen.NewRules(true,
			vgen.LenEq[[]int](3),
			vgen.List(
				vgen.LenEq[int](3),
				vgen.List(
					vgen.Eq(1),
				),
			),
		),
		stats: vgen.NewRules(true,
			vgen.MapValue(
				vgen.Eq(100),
			),
		),
	}

	err := person.Validate(rules)
	fmt.Println(err.Debug())
}
