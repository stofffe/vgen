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
	email, err := EmailVgen{
		Title:  nil,
		Text:   util.InitP("hello"),
		Sender: util.InitP(""),
	}.ValidatedConvert()
	if err != nil {
		util.DebugPrintAny("err", err)
	} else {
		fmt.Printf("person: %v\n", email)
	}
}
