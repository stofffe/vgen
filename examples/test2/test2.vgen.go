// THIS FILE IS GENERATED BY VGEN
// DO NOT EDIT
package main
import (
	"encoding/json"
	"fmt"
)
type PersonVgen struct {
	Name      *string
	Age       *int
	Vibes     *bool
	Nicknames *[]string
	A         *[][]string
}
func (t PersonVgen) Validate() (Person, error) {
	person, errs := t.InnerValidation()
	if len(errs) > 0 {
		j, _ := json.Marshal(errs)
		return Person{}, fmt.Errorf("%s", j)
	}
	return person, nil
}
func (t PersonVgen) InnerValidation() (Person, map[string][]string) {
	res := Person{}
	errs := make(map[string][]string)
	if t.Name != nil {
		name := *t.Name // TODO not working for 0 rules
		if !(len(name) < 20) {
			errs[fmt.Sprintf("name")] = append(errs[fmt.Sprintf("name")], fmt.Sprintf("len must be < 20"))
		}
		_ = name // No rules fix
	}
	if t.Age != nil {
		age := *t.Age // TODO not working for 0 rules
		if !(age >= 18) {
			errs[fmt.Sprintf("age")] = append(errs[fmt.Sprintf("age")], fmt.Sprintf("must be >= 18"))
		}
		_ = age // No rules fix
	}
	if t.Vibes != nil {
		vibes := *t.Vibes // TODO not working for 0 rules
		_ = vibes // No rules fix
	}
	if t.Nicknames != nil {
		nicknames := *t.Nicknames // TODO not working for 0 rules
		if !(len(nicknames) > 3) {
			errs[fmt.Sprintf("nicknames")] = append(errs[fmt.Sprintf("nicknames")], fmt.Sprintf("len must be > 3"))
		}
		for i0, nicknames := range nicknames {
			if !(len(nicknames) > 0) {
				errs[fmt.Sprintf("nicknames[%d]", i0)] = append(errs[fmt.Sprintf("nicknames[%d]", i0)], fmt.Sprintf("can not be empty"))
			}
			_ = nicknames // No rules fix
			_ = i0
		}
	}
	if t.A != nil {
		a := *t.A // TODO not working for 0 rules
		if !(len(a) > 0) {
			errs[fmt.Sprintf("a")] = append(errs[fmt.Sprintf("a")], fmt.Sprintf("can not be empty"))
		}
		for i0, a := range a {
			if !(len(a) > 1) {
				errs[fmt.Sprintf("a[%d]", i0)] = append(errs[fmt.Sprintf("a[%d]", i0)], fmt.Sprintf("len must be > 1"))
			}
			for i1, a := range a {
				if !(len(a) >= 2) {
					errs[fmt.Sprintf("a[%d][%d]", i0, i1)] = append(errs[fmt.Sprintf("a[%d][%d]", i0, i1)], fmt.Sprintf("len must be >= 2"))
				}
				_ = a // No rules fix
				_ = i1
			}
			_ = i0
		}
	}
	if len(errs) > 0 {
		return Person{}, errs
	}
	return res, nil
}
