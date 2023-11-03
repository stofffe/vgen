package main

import (
	"fmt"

	"github.com/stofffe/vgen/util"
)

// vgen:[i]
type Email struct {
	Title  string // vgen:[ req, not_empty, len_lt=50 ]
	Text   string // vgen:[ req, not_empty, len_gt=200 ]
	Sender string // vgen:[ req, not_empty, len_lt=20 ]
}

func main() {
	person, err := EmailVgen{
		Text:   util.InitP("hello"),
		Sender: util.InitP(""),
	}.Validate()
	if err != nil {
		util.DebugPrintString("err", err.Error())
	} else {
		fmt.Printf("person: %v\n", person)
	}
}
