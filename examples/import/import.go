package main

import (
	"fmt"
	"log"

	"github.com/stofffe/vgen/examples/import/timer"
	"github.com/stofffe/vgen/pkg/vgen"
)

// vgen(include)
type Import struct {
	clock timer.Clock
	timer timer.Timer // vgen(n)
}

func ClockIsZero() vgen.RuleFunc[timer.Clock] {
	return func(key string, input timer.Clock) vgen.ErrorMap {
		if input.Time != 0 {
			return vgen.NewErrorMap(key, fmt.Errorf("time must be zero"))
		}
		return nil
	}
}

func main() {
	currentTime := -120

	imp := ImportVgen{
		clock: &timer.Clock{
			Time: currentTime,
		},
		timer: &timer.TimerVgen{
			Time: &currentTime,
		},
	}

	rules := ImportRules{
		clock: vgen.RulesRequired(
			ClockIsZero(),
		),
		timer: vgen.RulesRequired(
			vgen.Nested(timer.TimerPositive),
		),
	}

	result, verr := imp.ValidatedConvert(rules)
	if verr != nil {
		log.Fatal(verr.Debug())
	}
	fmt.Println(result)

}
