package timer

import "github.com/stofffe/vgen/pkg/vgen"

// vgen(i)
type Timer struct {
	Time int
}

var TimerPositive = TimerRules{
	Time: vgen.RulesRequired[int](
		vgen.Gte(0),
	),
}

type Clock struct {
	Time int
}
