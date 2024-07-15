package main

import (
	"fmt"
	"log"

	"github.com/stofffe/vgen"
)

// vgen(include)
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
		stats: vgen.RulesOptional(
			vgen.MapHasKey[map[string]int]("speed"),
			vgen.MapHasKey[map[string]int]("strength"),
			vgen.MapHasKey[map[string]int]("intelligence"),
			vgen.MapValue(
				vgen.Lt(100),
			),
		),
	}

	result, verr := character.ValidatedConvert(rules)
	if verr != nil {
		log.Fatal(verr.Debug())
	}
	fmt.Println(result)
}
