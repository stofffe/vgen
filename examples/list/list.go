package main

import (
	"fmt"

	"github.com/stofffe/vgen/util"
)

// vgen:[i]
type Person struct {
	Nicknames []string   `json:"nicknames"` // vgen:[ len_gt=3 ][ not_empty ]
	A         [][]string `json:"a"`         // vgen:[ not_empty ][ len_gt=1 ][ custom=isBob, not_empty ]
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
	}.Validate()
	if err != nil {
		util.DebugPrintString("err", err.Error())
	} else {
		fmt.Printf("person: %v\n", person)
	}
}
