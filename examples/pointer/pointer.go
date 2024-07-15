package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/stofffe/vgen"
)

// vgen(include)
type Pointer struct {
	Name *string `json:"name"` // vgen(alias=name)
	Age  *int    `json:"age"`  // vgen(alias=name)
}

func main() {
	js := `{
        "name": null,
        "age": 12
    }`

	var pointer PointerVgen

	err := json.Unmarshal([]byte(js), &pointer)
	if err != nil {
		log.Fatal(err)
	}

	rules := PointerRules{
		Name: vgen.RulesRequired(
			vgen.Deref(
				vgen.Eq("bobby"),
			),
		),
		Age: vgen.RulesRequired(
			vgen.Deref(
				vgen.Gte(18),
			),
		),
	}

	errs := pointer.Validate(rules)
	fmt.Println(errs.Debug())
}
