package main

import (
	"fmt"
	"log"

	"github.com/stofffe/vgen"
)

// If json tag exists that name will be used as alias

// vgen(include)
type Test struct {
	bestResult Result   // vgen(nested)
	allResults []Result // vgen(nested)
}

// vgen(include)
type Result struct {
	score string // vgen(alias=score)
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
		bestResult: &ResultVgen{
			score: &X,
		},
		allResults: &[]ResultVgen{
			{score: &B},
			{score: &D},
			{score: &X},
			{score: &A},
			{},
			{score: &C},
		},
	}

	resultRules := ResultRules{
		score: vgen.RulesRequired(
			vgen.OneOf(A, B, C, D, E),
		),
	}

	testRules := TestRules{
		bestResult: vgen.RulesOptional(
			resultRules,
		),
		allResults: vgen.RulesOptional(
			vgen.List(
				resultRules,
			),
		),
	}

	result, verr := test.ValidatedConvert(testRules)
	if verr != nil {
		log.Fatal(verr.Debug())
	}
	fmt.Println(result)
}
