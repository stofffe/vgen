package main

import (
	"fmt"

	"github.com/stofffe/vgen/util"
)

// vgen:[i]
type Address struct {
	Street string // vgen:[req, not_empty]
	Number int    // vgen:[req, gt(0)]
}

// vgen:[i]
type Person struct {
	Name     string  // vgen:[ req, not_empty, len_lt(20)]
	Address1 Address // vgen:[req, custom(valAddr)]
	Address2 Address // vgen:[i, req, custom(valAddr), custom(abc)]
	// Address3 struct {
	// 	Street string
	// 	Number int
	// }
}

func abc(addr Address) error {
	return fmt.Errorf("abc")
}

func valAddr(addr Address) error {
	return fmt.Errorf("naa")
}

func main() {
	person, err := PersonVgen{
		Name: util.InitP("helo"),
		Address1: &Address{
			Street: "",
			Number: 123,
		},
		Address2: &AddressVgen{
			Street: util.InitP("st stree"),
			Number: util.InitP(0),
		},
	}.Validate()
	if err != nil {
		util.DebugPrintString("err", err.Error())
	} else {
		fmt.Printf("person: %v\n", person)
	}
}
