// THIS FILE IS GENERATED BY VGEN
// DO NOT EDIT

package main

import (
	"encoding/json"
	"fmt"
	"strconv"
)

type PersonVgen struct {
	Name      *string
	Age       *int
	Vibes     *bool
	Nicknames *[]string
}

func (t PersonVgen) Validate() (Person, error) {
	res := Person{}
	errs := make(map[string][]string)

	if t.Name != nil {
		name := *t.Name

		if !(len(name) > 0) {
			errs["name"] = append(errs["name"], fmt.Sprintf(`can not be empty`))
		}

		if !(len(name) < 20) {
			errs["name"] = append(errs["name"], fmt.Sprintf(`len must be less than 20`))
		}

		if err := isBob(name); err != nil {
			errs["name"] = append(errs["name"], fmt.Sprintf(`failed custom isBob: %v`, err))
		}

		res.Name = name
	} else {
		errs["name"] = append(errs["name"], fmt.Sprintf("required"))
	}

	if t.Age != nil {
		age := *t.Age

		if !(age >= 18) {
			errs["age"] = append(errs["age"], fmt.Sprintf(`must be greater than or equal to 18`))
		}

		if err := driveAge(age); err != nil {
			errs["age"] = append(errs["age"], fmt.Sprintf(`failed custom driveAge: %v`, err))
		}

		res.Age = age
	} else {
		errs["age"] = append(errs["age"], fmt.Sprintf("required"))
	}

	if t.Vibes != nil {
		vibes := *t.Vibes

		res.Vibes = vibes
	} else {
		errs["vibes"] = append(errs["vibes"], fmt.Sprintf("required"))
	}

	if t.Nicknames != nil {
		nicknames := *t.Nicknames

		if !(len(nicknames) > 3) {
			errs["nicknames"] = append(errs["nicknames"], fmt.Sprintf(`len must be greater than 3`))
		}

		if err := totLen(nicknames); err != nil {
			errs["nicknames"] = append(errs["nicknames"], fmt.Sprintf(`failed custom totLen: %v`, err))
		}

		for i, nicknames := range nicknames {

			if !(len(nicknames) >= 4) {
				errs["nicknames"] = append(errs["nicknames"], fmt.Sprintf(strconv.Itoa(i)+":"+`len must be greater than or equal to 4`))
			}

		}

		res.Nicknames = nicknames
	} else {
		errs["nicknames"] = append(errs["nicknames"], fmt.Sprintf("required"))
	}

	if len(errs) > 0 {
		j, _ := json.Marshal(errs)
		return Person{}, fmt.Errorf("%s", j)
	}

	return res, nil
}
