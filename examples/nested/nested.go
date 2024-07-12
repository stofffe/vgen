package main

import (
	"fmt"

	"github.com/stofffe/vgen"
)

// If json tag exists that name will be used as alias

// vgen(include)
type Test struct {
	BestResult Result   // vgen(n, alias=pet)
	AllResults []Result // vgen(n, alias=pets)
}

// vgen(include)
type Result struct {
	Score string // vgen(alias=name)
}

var (
	A = "A"
	B = "B"
	C = "C"
	D = "D"
	E = "E"
	F = "F"

	X = "X"
)

func main() {
	test := TestVgen{
		BestResult: &ResultVgen{
			Score: &X,
		},
		AllResults: &[]ResultVgen{
			{Score: &B},
			{Score: &D},
			{Score: &X},
			{Score: &A},
			{},
			{Score: &C},
		},
	}

	resultRules := ResultRules{
		Score: vgen.NewRules(
			vgen.Required[string](),
			vgen.OneOf(A, B, C, D, E),
		),
	}

	testRules := TestRules{
		BestResult: vgen.NewRules(
			vgen.Nested(resultRules),
		),
		AllResults: vgen.NewRules(
			vgen.List(
				vgen.Nested(resultRules),
			),
		),
	}

	err := test.Validate(testRules)
	fmt.Println(err.Debug())
}
