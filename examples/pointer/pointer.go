package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/stofffe/vgen"
)

// vgen
type Pointer struct {
	Name []*string `json:"name"`
	Age  []int
	Pet  *string
}

func main() {
	js := `{
        "name": [
            "bob", "bobby", null
        ]
    }`

	var pointer PointerVgen

	err := json.Unmarshal([]byte(js), &pointer)
	if err != nil {
		log.Fatal(err)
	}

	rules := PointerRules{
		Age: vgen.NewRules(
			vgen.Required[[]int](),
		),
		Name: vgen.NewRules(
			vgen.Required[[]*string](),
			vgen.List(
				vgen.PointerNotNil[*string](),
				vgen.Deref(
					vgen.Eq("bob"),
					vgen.Eq("aaa"),
				),
			),
		),
	}

	errs := pointer.Validate(rules)
	fmt.Println(errs.Debug())
	// for _, v := range pointer.Name {
	// 	fmt.Println(*v)
	// }
}
