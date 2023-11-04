package main

import (
	"fmt"

	"github.com/stofffe/vgen/util"
)

// vgen:[i]
type Address struct {
	Street string `json:"street"` // vgen:[ req, not_empty ]
	Number int    `json:"number"` // vgen:[ req, lt=5 ]
}

// vgen:[i]
type Person struct {
	Name     string  `json:"name"`     // vgen:[ req, not_empty, len_lt=20 ]
	Address1 Address `json:"address1"` // vgen:[ custom=abc ]
	Address2 Address `json:"address2"` // vgen:[ i, req, custom=abc ]
}

func abc(addr Address) error {
	return fmt.Errorf("abc")
}

func main() {
	person, err := PersonVgen{
		Name: util.InitP("helo"),
		Address1: &Address{
			Street: "address1",
			Number: 1,
		},
		Address2: &AddressVgen{
			Street: util.InitP("address2"),
			Number: util.InitP(10),
		},
	}.Validate()
	if err != nil {
		util.DebugPrintString("err", err.Error())
	} else {
		fmt.Printf("person: %v\n", person)
	}
}
