package main

import (
	"fmt"

	"github.com/stofffe/vgen/util"
)

// vgen:[i]
type Email struct {
	Title  string `json:"title"`  // vgen:[ req, not_empty, len_lt=50 ]
	Text   string `json:"text"`   // vgen:[ req, not_empty, len_gt=200 ]
	Sender string `json:"sender"` // vgen:[ req, not_empty, len_lt=20 ]
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
