package examples

import (
	"encoding/json"
	"fmt"
)

type PersonVgen struct {
	Name  *string
	Age   *int
	Vibes *bool
}

func (t PersonVgen) Validate() (Person, error) {
	res := Person{}
	errs := make(map[string]string)
	// Name string
	if t.Name != nil {
		name := *t.Name

		// req

		// len > 0

		res.Name = name
	} else {
		errs["Name"] = fmt.Sprintf("required")
	}

	// Age int
	if t.Age != nil {
		age := *t.Age

		// val >= 3

		res.Age = age
	}

	// Vibes bool
	if t.Vibes != nil {
		vibes := *t.Vibes

		// req

		res.Vibes = vibes
	} else {
		errs["Vibes"] = fmt.Sprintf("required")
	}

	if len(errs) > 0 {
		j, _ := json.Marshal(errs)
		return Person{}, fmt.Errorf("%s", j)
	}

	return res, nil
}
