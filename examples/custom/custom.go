package main

import (
	"fmt"

	"github.com/stofffe/vgen/util"
)

// vgen:[i]
type Person struct {
	Name  string `json:"name"`  // vgen:[ req, not_empty , custom=isBob]
	Age   int    `json:"age"`   // vgen:[ req, gt=18, lt=22 ]
	Vibes bool   `json:"vibes"` // vgen:[ req ]
}

func isBob(t string) error {
	if t != "bob" {
		return fmt.Errorf("must be Bob")
	}
	return nil
}

func main() {
	person, err := PersonVgen{
		Name:  util.InitP("helo"),
		Age:   util.InitP(17),
		Vibes: nil,
	}.Validate()
	if err != nil {
		util.DebugPrintString("err", err.Error())
	} else {
		fmt.Printf("person: %v\n", person)
	}
}
