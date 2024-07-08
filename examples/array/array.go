package main

import (
	"fmt"

	"github.com/stofffe/vgen"
)

// vgen
type Array struct {
	vector []int
	matrix [][]int
}

func main() {
	arrays := ArrayVgen{
		vector: &[]int{0, 2, 12},
		matrix: &[][]int{
			{1, 2},
			{4, 5, 6},
			{7, 88, 9},
		},
	}

	rules := ArrayRules{
		vector: vgen.NewRules(true,
			vgen.LenEq[[]int](3),
			vgen.List(
				vgen.Gt(0),
				vgen.Lte(10),
			),
		),
		matrix: vgen.NewRules(true,
			vgen.LenEq[[][]int](3),
			vgen.List(
				vgen.LenEq[[]int](3),
				vgen.List(
					vgen.Gt(0),
					vgen.Lte(10),
				),
			),
		),
	}

	err := arrays.Validate(rules)
	fmt.Println(err.Debug())

}
