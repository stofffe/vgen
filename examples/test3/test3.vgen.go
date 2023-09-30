// THIS FILE IS GENERATED BY VGEN
// DO NOT EDIT

package main

import (
	"encoding/json"

	"fmt"
)

type AddressVgen struct {
	Street *string
	Number *int
}

func (t AddressVgen) Validate() (Address, error) {
	// TODO add output formatting here
	address, errs := t.innerValidation("")
	if len(errs) > 0 {
		j, _ := json.Marshal(errs)
		return Address{}, fmt.Errorf("%s", j)
	}
	return address, nil
}

func (t AddressVgen) innerValidation(prefix string) (Address, map[string][]string) {
	res := Address{}
	errs := make(map[string][]string)

	if t.Street != nil {
		street := *t.Street

		if !(len(street) > 0) {
			errs["street"] = append(errs["street"], fmt.Sprintf(`can not be empty`))
		}

		res.Street = street
	} else {
		errs["street"] = append(errs["street"], fmt.Sprintf("required"))
	}

	if t.Number != nil {
		number := *t.Number

		if !(number >= 0) {
			errs["number"] = append(errs["number"], fmt.Sprintf(`must be greater than or equal to 0`))
		}

		res.Number = number
	} else {
		errs["number"] = append(errs["number"], fmt.Sprintf("required"))
	}

	if len(errs) > 0 {
		return Address{}, errs
	}

	return res, nil
}

type PersonVgen struct {
	Name     *string
	Address1 *Address
	Address2 *AddressVgen
}

func (t PersonVgen) Validate() (Person, error) {
	// TODO add output formatting here
	person, errs := t.innerValidation("")
	if len(errs) > 0 {
		j, _ := json.Marshal(errs)
		return Person{}, fmt.Errorf("%s", j)
	}
	return person, nil
}

func (t PersonVgen) innerValidation(prefix string) (Person, map[string][]string) {
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

		res.Name = name
	} else {
		errs["name"] = append(errs["name"], fmt.Sprintf("required"))
	}

	if t.Address1 != nil {
		address1 := *t.Address1

		if err := valAddr(address1); err != nil {
			errs["address1"] = append(errs["address1"], fmt.Sprintf(`%v`, err))
		}

		res.Address1 = address1
	} else {
		errs["address1"] = append(errs["address1"], fmt.Sprintf("required"))
	}

	if t.Address2 != nil {
		address2_ := *t.Address2

		address2, err := address2_.innerValidation("")
		if err != nil {
			for k, v := range err {
				errs["address2."+k] = append(errs["address2."+k], v...)
			}
		}

		if err := valAddr(address2); err != nil {
			errs["address2"] = append(errs["address2"], fmt.Sprintf(`%v`, err))
		}

		res.Address2 = address2
	} else {
		errs["address2"] = append(errs["address2"], fmt.Sprintf("required"))
	}

	if len(errs) > 0 {
		return Person{}, errs
	}

	return res, nil
}
