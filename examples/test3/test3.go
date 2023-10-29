package main

import (
	"fmt"

	"github.com/stofffe/vgen/util"
)

// vgen:[i]
type Address struct {
	Street string // vgen:[req, not_empty]
	Number int    // vgen:[req, lt=5]
}

// vgen:[i]
type Person struct {
	Name     string  // vgen:[ req, not_empty, len_lt=20]
	Address1 Address // vgen:[custom=valAddr]
	Address2 Address // vgen:[i, req, custom=valAddr, custom=abc]
	A        []int   // vgen:[req]
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
			Street: util.InitP(""),
			Number: util.InitP(100),
		},
	}.Validate()
	if err != nil {
		util.DebugPrintString("err", err.Error())
	} else {
		fmt.Printf("person: %v\n", person)
	}
}
