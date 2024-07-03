package main

import (
	"fmt"

	"github.com/stofffe/vgen/vgen"
)

type Person struct {
	name      string
	age       int
	nicknames []string
	matrix    [][]int
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
			"notbob2",
		},
		matrix: &[][]int{
			{1, 0, 0},
			{0, 1, 0},
			{0, 0, 1},
		},
		stats: &stats,
	}

	rules := PersonRules{
		name: vgen.NewRules(true,
			vgen.Eq("bob"),
		),
		age: vgen.NewRules(true,
			vgen.Lte(5),
			vgen.Gte(120),
		),
		nicknames: vgen.NewRules(true,
			vgen.List(
				vgen.CustomMessage("test", vgen.Eq("bob1")),
				vgen.Eq("bob2"),
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
