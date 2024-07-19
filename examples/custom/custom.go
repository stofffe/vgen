package main

import (
	"fmt"
	"log"

	"github.com/stofffe/vgen/pkg/vgen"
)

// vgen(include)
type Vector struct {
	vector []float32
}

// Custom rule
// Components must sum to 1
func Normalized() vgen.RuleFunc[[]float32] {
	return func(field_name string, input []float32) vgen.ErrorMap {
		var sum float32
		for _, v := range input {
			sum += v
		}
		if sum != 1.0 {
			return vgen.NewErrorMap(field_name, fmt.Errorf("vector not normalized"))
		}
		return nil
	}
}

func main() {
	vector := VectorVgen{
		vector: &[]float32{1, 0, 1},
	}

	rules := VectorRules{
		vector: vgen.RulesOptional[[]float32](
			vgen.LenEq[[]float32](3),
			Normalized(),
		),
	}

	result, verr := vector.ValidatedConvert(rules)
	if verr != nil {
		log.Fatal(verr.Debug())
	}
	fmt.Println(result)
}
