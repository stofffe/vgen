package main

import (
	"fmt"

	"github.com/stofffe/vgen/util"
)

// vgen:[i]
type Person struct {
	Name      string   // vgen:[ req, not_empty, len_lt(20)]
	Age       int      // vgen:[ req, gte(18) ]
	Vibes     bool     // vgen:[ req ]
	Nicknames []string // vgen:[ req, len_gt(3), :len_gte(4)]
}

func main() {
	person, err := PersonVgen{
		Name:      util.InitP("helo"),
		Age:       util.InitP(17),
		Vibes:     nil,
		Nicknames: util.InitP([]string{"hello", "yoyo", "abc", "noooooooooooooooooo"}),
	}.Validate()
	if err != nil {
		util.DebugPrintString("err", err.Error())
	} else {
		fmt.Printf("person: %v\n", person)
	}
}
