package main

import (
	"fmt"

	"github.com/stofffe/vgen/util"
)

// vgen:[i]
type Name struct {
	value string // vgen:[ req ]
}

// vgen:[i]
type Person struct {
	A string     // vgen:[ req ]
	B []string   // vgen:[ ][ ]
	C [][]string // vgen:[ ][ ][ ]
	D Name       // vgen:[ i, req ]
	E []Name     // vgen:[ ][ i ]
	F [][]Name   // vgen:[ ][ ][ i ]
}

func main() {
	person, err := PersonVgen{
		A: util.InitP("hello"),
		B: &[]string{"1", "2"},
		C: &[][]string{{"1", "2"}, {"3", "4"}},
		D: &NameVgen{value: util.InitP("hello")},
		E: &[]NameVgen{
			{value: util.InitP("hello")},
		},
		F: &[][]NameVgen{{
			{value: util.InitP("1")}, {value: util.InitP("2")},
		}, {
			{value: util.InitP("3")},
		}},
	}.Validate()
	if err != nil {
		util.DebugPrintString("err", err.Error())
	} else {
		fmt.Printf("person: %v\n", person)
	}
}
