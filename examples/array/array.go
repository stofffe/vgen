package main

import (
	"fmt"
	"log"

	"github.com/stofffe/vgen/pkg/vgen"
)

// vgen(include)
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
		vector: vgen.RulesOptional[[]int](
			vgen.LenEq[[]int](3),
			vgen.List[int](
				vgen.Gt(0),
				vgen.Lte(10),
			),
		),
		matrix: vgen.RulesOptional[[][]int](
			vgen.LenEq[[][]int](3),
			vgen.List[[]int](
				vgen.LenEq[[]int](3),
				vgen.List[int](
					vgen.Gt(0),
					vgen.Lte(10),
				),
			),
		),
	}

	result, verr := arrays.ValidatedConvert(rules)
	if verr != nil {
		log.Fatal(verr.Debug())
	}
	fmt.Println(result)
}
