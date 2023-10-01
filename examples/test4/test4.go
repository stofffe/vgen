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
	Name      string    // vgen:[req, not_empty, len_lt=20]
	Address1  Address   // vgen:[req, custom=valAddr]
	Address2  Address   // vgen:[i, req, custom=valAddr]
	Addresses []Address // vgen:[req, len_gt=0, :custom=valAddr]
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
		Address2: &AddressVgen{},
		Addresses: &[]Address{
			Address{
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
