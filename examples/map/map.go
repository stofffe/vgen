package main

import (
	"fmt"

	"github.com/stofffe/vgen"
)

// vgen
type Character struct {
	stats map[string]int
}

func main() {
	stats := make(map[string]int)
	stats["speed"] = 80
	stats["intelligence"] = 120
	character := CharacterVgen{
		stats: &stats,
	}

	rules := CharacterRules{
		stats: vgen.NewRules(true,
			vgen.MapHasKey[int]("speed"),
			vgen.MapHasKey[int]("strength"),
			vgen.MapHasKey[int]("intelligence"),
			vgen.MapValue(
				vgen.Lt(100),
			),
		),
	}

	err := character.Validate(rules)
	fmt.Println(err.Debug())
}
