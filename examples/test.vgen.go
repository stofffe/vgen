// THIS FILE IS GENERATED FROM VGEN
// DO NOT EDIT

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
			errs["name"] = fmt.Sprintf(`can not be empty`)
		}

		if !(len(name) > 20) {
			errs["name"] = fmt.Sprintf(`len must be > 20`)
		}

		res.Name = name
	} else {
		errs["Name"] = fmt.Sprintf("required")
	}

	// Age int
	if t.Age != nil {
		age := *t.Age

		if !(age >= 18) {
			errs["age"] = fmt.Sprintf(`must be >= 18`)
		}

		res.Age = age
	} else {
		errs["Age"] = fmt.Sprintf("required")
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
