// THIS FILE IS GENERATED BY VGEN
// DO NOT EDIT
package main
import (
	"encoding/json"
	"fmt"
)
type PersonVgen struct {
	Nicknames *[]string `json:"nicknames"`
	A *[][]string `json:"a"`
}
func (t PersonVgen) Validate() (Person, error) {
	errs := t.InnerValidation()
	if len(errs) > 0 {
		j, _ := json.Marshal(errs)
		return Person{}, fmt.Errorf("%s", j)
	}
	return t.Convert(), nil
}
func (t PersonVgen) InnerValidation() map[string][]string {
	errs := make(map[string][]string)
	if t.Nicknames != nil {
		_Nicknames := *t.Nicknames
		{
			if !(len(_Nicknames) > 3) {
				errs[fmt.Sprintf("nicknames")] = append(errs[fmt.Sprintf("nicknames")], fmt.Sprintf("len must be > 3"))
			}
			for i0, _Nicknames := range _Nicknames {
				if !(len(_Nicknames) > 0) {
					errs[fmt.Sprintf("nicknames[%d]", i0)] = append(errs[fmt.Sprintf("nicknames[%d]", i0)], fmt.Sprintf("can not be empty"))
				}
				_ = i0
			}
		}
	}
	if t.A != nil {
		_A := *t.A
		{
			if !(len(_A) > 0) {
				errs[fmt.Sprintf("a")] = append(errs[fmt.Sprintf("a")], fmt.Sprintf("can not be empty"))
			}
			for i0, _A := range _A {
				if !(len(_A) > 1) {
					errs[fmt.Sprintf("a[%d]", i0)] = append(errs[fmt.Sprintf("a[%d]", i0)], fmt.Sprintf("len must be > 1"))
				}
				for i1, _A := range _A {
					if err := isBob(_A); err != nil {
						errs[fmt.Sprintf("a[%d][%d]", i0, i1)] = append(errs[fmt.Sprintf("a[%d][%d]", i0, i1)], err.Error())
					}
					if !(len(_A) > 0) {
						errs[fmt.Sprintf("a[%d][%d]", i0, i1)] = append(errs[fmt.Sprintf("a[%d][%d]", i0, i1)], fmt.Sprintf("can not be empty"))
					}
					_ = i1
				}
				_ = i0
			}
		}
	}
	return errs
}
func (t PersonVgen) Convert() Person {
	var res Person
	if t.Nicknames != nil {
		_Nicknames := *t.Nicknames
		res.Nicknames = make([]string, len(_Nicknames))
		for i0, _Nicknames := range _Nicknames {
			res.Nicknames[i0] = _Nicknames
		}
	}
	if t.A != nil {
		_A := *t.A
		res.A = make([][]string, len(_A))
		for i0, _A := range _A {
			res.A[i0] = make([]string, len(_A))
			for i1, _A := range _A {
				res.A[i0][i1] = _A
			}
		}
	}
	return res
}
func PersonFromJson(bytes []byte) (Person, error) {
	var v PersonVgen
	err := json.Unmarshal(bytes, &v)
	if err != nil {
		return Person{}, err
	}
	r, err := v.Validate()
	if err != nil {
		return Person{}, err
	}
	return r, nil
}
