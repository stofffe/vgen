package main

import (
	"fmt"

	"github.com/stofffe/vgen/util"
)

// vgen:[i]
type Person struct {
	Name      string     // vgen:[ len_lt=20]
	Age       int        // vgen:[ gte=18 ]
	Vibes     bool       // vgen:[ ]
	Nicknames []string   // vgen:[ len_gt=3 ][ not_empty ]
	A         [][]string // vgen:[ not_empty ][ len_gt=1 ][ len_gte=2 ]
}

func main() {
	person, err := PersonVgen{
		Name:      util.InitP("helo"),
		Age:       util.InitP(17),
		Vibes:     nil,
		Nicknames: util.InitP([]string{"hello", "", "abc", "noooooooooooooooooo"}),
		A:         util.InitP([][]string{{"abc"}, {"123", "yo"}, {"no", ""}}),
	}.Validate()
	if err != nil {
		util.DebugPrintString("err", err.Error())
	} else {
		fmt.Printf("person: %v\n", person)
	}
}
