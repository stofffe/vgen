package main

import (
	"fmt"

	"github.com/stofffe/vgen/util"
)

// vgen:[i]
type Person struct {
	Name      string   // vgen:[ req, not_empty, len_lt(20), custom(isBob) ]
	Age       int      // vgen:[ req, gte(18), custom(driveAge) ]
	Vibes     bool     // vgen:[ req ]
	Nicknames []string // vgen:[ req, len_gt(3), :len_gte(4), custom(totLen)]
}

func totLen(s []string) error {
	l := 0
	for _, v := range s {
		l += len(v)
	}
	if l > 10 {
		return fmt.Errorf("tot len must be < 10")
	}
	return nil
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
		Name:      util.InitP("helo"),
		Age:       util.InitP(17),
		Vibes:     nil,
		Nicknames: util.InitP([]string{"hello", "yoyo", "abc", "noooooooooooooooooo"}),
	}.Validate()
	if err != nil {
		util.DebugPrint("err", err.Error())
	} else {
		fmt.Printf("person: %v\n", person)
	}
}
