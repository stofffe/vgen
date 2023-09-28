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

		if !(len(name) > 0) {
			errs[""] = fmt.Sprintf("len must be > 0")
		}
		res.Name = name
	} else {
		errs["Name"] = fmt.Sprintf("required")
	}

	// Age int
	if t.Age != nil {
		age := *t.Age
		// Rule not implemented for val>=3
		// Rule not implemented for gt(4)
		// Rule not implemented for email
		res.Age = age
	}

	// Vibes bool
	if t.Vibes != nil {
		vibes := *t.Vibes

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
