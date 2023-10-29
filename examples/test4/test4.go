package main

import (
	"fmt"

	"github.com/stofffe/vgen/util"
)

// vgen:[i]
type Address struct {
	Street string // vgen:[req, not_empty]
	Number int    // vgen:[req, gte=0]
}

// vgen:[i]
type Person struct {
	Address   Address   // vgen:[req, i]
	Addresses []Address // vgen:[req, len_gt=0, :i, :custom=abc]
	Strings   []string  // vgen:[req, len_gt=0, :not_empty]
}

func abc(addr Address) error {
	return nil
}

func main() {
	person, err := PersonVgen{
		// Address: &AddressVgen{},
		Addresses: &[]Address{
			{
				Street: "123",
				Number: 123,
			},
		},
	}.Validate()
	if err != nil {
		util.DebugPrintString("err", err.Error())
	} else {
		fmt.Printf("person: %v\n", person)
	}
}
