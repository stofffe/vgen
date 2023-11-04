package main

import (
	"fmt"

	"github.com/stofffe/vgen/util"
)

// vgen:[i]
type Person struct {
	Name     string `json:"name"`     // vgen:[ req, not_empty, len_lt=10 ]
	Nickname string `json:"nickname"` // vgen:[ req ]
	Age      int    `json:"age"`      // vgen:[ req, gt=10 ]

}

func main() {
	body := `
		{
			"name": "very very long name",
			"age": 8
		}
	`
	person, err := PersonFromJson([]byte(body))
	if err != nil {
		util.DebugPrintString("err", err.Error())
	} else {
		fmt.Printf("person: %v\n", person)
	}
}
