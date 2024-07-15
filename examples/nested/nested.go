package main

import (
	"fmt"
	"log"

	"github.com/stofffe/vgen"
)

// If json tag exists that name will be used as alias

// vgen(include)
type Test struct {
	BestResult  Result                // vgen(nested)
	AllResults  []Result              // vgen(nested)
	All2Results [][]Result            // vgen(nested)
	All3Results [][][]Result          // vgen(nested)
	MapResults  []map[string]Result   // vgen(nested)
	Map2Results []map[string][]Result // vgen(nested)
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
	ma1 := make(map[string]ResultVgen)
	ma1["bobC"] = ResultVgen{Score: &C}
	ma1["bobA"] = ResultVgen{Score: &A}
	ma2 := make(map[string]ResultVgen)
	ma2["aaa"] = ResultVgen{Score: &A}
	ma3 := make(map[string][]ResultVgen)
	ma3["ooo"] = []ResultVgen{
		{Score: &C},
		{Score: &B},
	}
	test := TestVgen{
		BestResult: &ResultVgen{
			Score: &A,
		},
		AllResults: &[]ResultVgen{
			{Score: &B},
			{Score: &D},
			{Score: &X},
			{Score: &A},
			{},
			{Score: &C},
		},
		All2Results: &[][]ResultVgen{
			{
				{Score: &B},
				{Score: &B},
				{Score: &B},
			},
			{
				{Score: &A},
				{Score: &A},
				{Score: &A},
			},
		},
		MapResults: &[]map[string]ResultVgen{
			ma1, ma2,
		},
		Map2Results: &[]map[string][]ResultVgen{
			ma3,
		},
	}

	resultRules := ResultRules{
		Score: vgen.RulesRequired(
			vgen.OneOf(A, B, C, D, E),
		),
	}

	testRules := TestRules{
		BestResult: vgen.RulesOptional(
			resultRules,
		),
		AllResults: vgen.RulesOptional(
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
