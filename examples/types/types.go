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
	Name       string      // vgen:[ req, not_empty, len_lt=20 ]
	Address1   Address     // vgen:[ custom=valAddr ]
	Address2   Address     // vgen:[ i, req, custom=valAddr ]
	Addresses  []Address   // vgen:[ req ][ i ]
	Addresses2 [][]Address // vgen:[ req ][ ][ i ]
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
			Street: "address1",
			Number: 1,
		},
		Address2: &AddressVgen{
			Street: util.InitP("address2"),
			Number: util.InitP(10),
		},
		Addresses: &[]AddressVgen{
			{
				Street: util.InitP("address2"),
				Number: util.InitP(10),
			},
		},
		Addresses2: &[][]AddressVgen{
			{
				{
					Street: util.InitP("addressA"),
					Number: util.InitP(10),
				},
				{
					Street: util.InitP(""),
					Number: util.InitP(10),
				},
			},
			{
				{
					Street: util.InitP("addressC"),
				},
			},
		},
	}.Validate()
	if err != nil {
		util.DebugPrintString("err", err.Error())
	} else {
		fmt.Printf("person: %v\n", person)
	}
}
