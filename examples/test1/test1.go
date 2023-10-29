package main

import (
	"fmt"

	"github.com/stofffe/vgen/util"
)

// vgen:[i]
type Person struct {
	Name  string // vgen:[ req, not_empty ]
	Age   int    // vgen:[ req, gt=18, lt=22 ]
	Vibes bool   // vgen:[ req ]
}

func driveAge(t int) error {
	if t < 18 {
		return fmt.Errorf("must be at least 18")
	}
	return nil
}

func isBob(t string) error {
	if t != "Bob" {
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
