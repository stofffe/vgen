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
	errs := make(map[string][]string)
	// Name string
	if t.Name != nil {
		name := *t.Name

		if !(len(name) > 0) {
			errs["name"] = append(errs["name"], fmt.Sprintf(`can not be empty`))
		}

		if !(len(name) < 20) {
			errs["name"] = append(errs["name"], fmt.Sprintf(`len must be < 20`))
		}

		if err := isBob(name); err != nil {
			errs["name"] = append(errs["name"], fmt.Sprintf(`failed custom isBob: %v`, err))
		}

		res.Name = name
	} else {
		errs["name"] = append(errs["name"], fmt.Sprintf("required"))
	}

	// Age int
	if t.Age != nil {
		age := *t.Age

		if !(age >= 18) {
			errs["age"] = append(errs["age"], fmt.Sprintf(`must be >= 18`))
		}

		if err := driveAge(age); err != nil {
			errs["age"] = append(errs["age"], fmt.Sprintf(`failed custom driveAge: %v`, err))
		}

		res.Age = age
	} else {
		errs["age"] = append(errs["age"], fmt.Sprintf("required"))
	}

	// Vibes bool
	if t.Vibes != nil {
		vibes := *t.Vibes

		res.Vibes = vibes
	} else {
		errs["vibes"] = append(errs["vibes"], fmt.Sprintf("required"))
	}

	if len(errs) > 0 {
		j, _ := json.Marshal(errs)
		return Person{}, fmt.Errorf("%s", j)
	}

	return res, nil
}
