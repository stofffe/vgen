package main

import (
	"fmt"

	"github.com/stofffe/vgen/util"
)

// vgen:[i]
type Person struct {
	Nicknames []string   `json:"nicknames"` // vgen:[ len_gt=3 ]
	A         [][]string `json:"a"`         // vgen:[ custom=isBob, not_empty ]
}

func isBob(t string) error {
	if t != "bob" {
		return fmt.Errorf("must be bob")
	}
	return nil
}

func main() {
	person, err := PersonVgen{
		Nicknames: util.InitP([]string{"hello", "", "abc", "noooooooooooooooooo"}),
		A:         util.InitP([][]string{{"abc"}, {"bob", "yo"}, {"bob", ""}}),
	}.ValidatedConvert()
	if err != nil {
		util.DebugPrintAny("err", err)
	} else {
		fmt.Printf("person: %v\n", person)
	}
}
