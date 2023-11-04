package main

import (
	"encoding/json"
	"fmt"
	"log"

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

	var personVgen PersonVgen
	err := json.Unmarshal([]byte(body), &personVgen)
	if err != nil {
		log.Fatal(err)
	}
	person, err := personVgen.Validate()
	if err != nil {
		util.DebugPrintString("err", err.Error())
	} else {
		fmt.Printf("person: %v\n", person)
	}
}
