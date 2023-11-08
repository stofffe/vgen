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
	Name       string      `json:"name"`       // vgen:[ req, not_empty, len_lt=20 ]
	Address1   Address     `json:"address1"`   // vgen:[ custom=valAddr ]
	Address2   Address     `json:"address2"`   // vgen:[ i, req, custom=valAddr ]
	Addresses  []Address   `json:"addresses"`  // vgen:[ req ][ i ]
	Addresses2 [][]Address `json:"addresses2"` // vgen:[ req ][ ][ i ]
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
			Number: util.InitP(1),
		},
		Addresses: &[]AddressVgen{
			{
				Street: util.InitP("address2"),
				Number: util.InitP(1),
			},
		},
		Addresses2: &[][]AddressVgen{
			{
				{
					Street: util.InitP("addressA"),
					Number: util.InitP(1),
				},
				{
					Street: util.InitP(""),
					Number: util.InitP(1),
				},
			},
			{
				{
					Street: util.InitP("addressC"),
				},
			},
		},
	}.ValidatedConvert()
	if err != nil {
		util.DebugPrintAny("err", err)
	} else {
		fmt.Printf("person: %v\n", person)
	}
}

//	func main() {
//		person, err := PersonVgen{
//			Name: util.InitP("helo"),
//			Address1: &Address{
//				Street: "address1",
//				Number: 1,
//			},
//			Address2: &AddressVgen{
//				Street: util.InitP("address2"),
//				Number: util.InitP(10),
//			},
//			Addresses: &[]AddressVgen{
//				{
//					Street: util.InitP("address2"),
//					Number: util.InitP(10),
//				},
//			},
//			Addresses2: &[][]AddressVgen{
//				{
//					{
//						Street: util.InitP("addressA"),
//						Number: util.InitP(10),
//					},
//					{
//						Street: util.InitP(""),
//						Number: util.InitP(10),
//					},
//				},
//				{
//					{
//						Street: util.InitP("addressC"),
//					},
//				},
//			},
//		}.Validate()
//		if err != nil {
//			util.DebugPrintString("err", err.Error())
//		} else {
//			fmt.Printf("person: %v\n", person)
//		}
//	}
