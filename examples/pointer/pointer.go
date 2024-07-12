package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/stofffe/vgen"
)

// vgen
type Pointer struct {
	Name *string `json:"name"` // vgen(alias=name)
}

func main() {
	js := `{
        "name": null
    }`

	var pointer PointerVgen

	err := json.Unmarshal([]byte(js), &pointer)
	if err != nil {
		log.Fatal(err)
	}

	rules := PointerRules{
		Name: vgen.NewRules(
			vgen.Required[*string](),
			vgen.Deref(
				vgen.Eq("bobby"),
			),
		),
	}

	errs := pointer.Validate(rules)
	fmt.Println(errs.Debug())
	// for _, v := range pointer.Name {
	// 	fmt.Println(*v)
	// }
}
