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
	A string // vgen:[ req, not_empty ]
	B []string
	C [][]string
	D Name // vgen:[ i, req ]
	E []Name
	F [][]Name // vgen:[ i ]
}

func main() {
	person, err := PersonVgen{
		A: util.InitP("hello"),
		B: &[]string{"1", "2"},
		C: &[][]string{{"1", "2"}, {"3", "4"}},
		D: &NameVgen{value: util.InitP("hello")},
		E: &[]Name{{value: "hello"}},
		F: &[][]NameVgen{{
			{value: util.InitP("1")}, {value: util.InitP("2")},
		}, {
			{value: util.InitP("3")},
		}},
	}.ValidatedConvert()
	if err != nil {
		util.DebugPrintAny("err", err)
	} else {
		fmt.Printf("person: %v\n", person)
	}
}
