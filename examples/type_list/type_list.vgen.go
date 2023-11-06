// THIS FILE IS GENERATED BY VGEN
// DO NOT EDIT
package main
import (
	"encoding/json"
	"fmt"
)
type NameVgen struct {
	value *string
}
func (t NameVgen) Validate() (Name, error) {
	errs := t.InnerValidation()
	if len(errs) > 0 {
		j, _ := json.Marshal(errs)
		return Name{}, fmt.Errorf("%s", j)
	}
	return t.Convert(), nil
}
func (t NameVgen) InnerValidation() map[string][]string {
	errs := make(map[string][]string)
	if t.value != nil {
		_value := *t.value
		{
			_ = _value
		}
	} else {
		errs["value"] = append(errs["value"], fmt.Sprintf("required"))
	}
	return errs
}
func (t NameVgen) Convert() Name {
	var res Name
	if t.value != nil {
		_value := *t.value
		res.value = _value
	}
	return res
}
func NameFromJson(bytes []byte) (Name, error) {
	var v NameVgen
	err := json.Unmarshal(bytes, &v)
	if err != nil {
		return Name{}, err
	}
	r, err := v.Validate()
	if err != nil {
		return Name{}, err
	}
	return r, nil
}
type PersonVgen struct {
	A *string
	B *[]string
	C *[][]string
	D *NameVgen
	E *[]NameVgen
	F *[][]NameVgen
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
	if t.A != nil {
		_A := *t.A
		{
			_ = _A
		}
	} else {
		errs["A"] = append(errs["A"], fmt.Sprintf("required"))
	}
	if t.B != nil {
		_B := *t.B
		{
			for i0, _B := range _B {
				_ = _B
				_ = i0
			}
		}
	}
	if t.C != nil {
		_C := *t.C
		{
			for i0, _C := range _C {
				for i1, _C := range _C {
					_ = _C
					_ = i1
				}
				_ = i0
			}
		}
	}
	if t.D != nil {
		_D := *t.D
		{
			struct_errs := _D.InnerValidation()
			for path, err_list := range struct_errs {
				for _, err := range err_list {
					errs[fmt.Sprintf("D")+"."+path] = append(errs[fmt.Sprintf("D")+"."+path], err)
				}
			}
			_ = _D
		}
	} else {
		errs["D"] = append(errs["D"], fmt.Sprintf("required"))
	}
	if t.E != nil {
		_E := *t.E
		{
			for i0, _E := range _E {
				struct_errs := _E.InnerValidation()
				for path, err_list := range struct_errs {
					for _, err := range err_list {
						errs[fmt.Sprintf("E[%d]", i0)+"."+path] = append(errs[fmt.Sprintf("E[%d]", i0)+"."+path], err)
					}
				}
				_ = _E
				_ = i0
			}
		}
	}
	if t.F != nil {
		_F := *t.F
		{
			for i0, _F := range _F {
				for i1, _F := range _F {
					struct_errs := _F.InnerValidation()
					for path, err_list := range struct_errs {
						for _, err := range err_list {
							errs[fmt.Sprintf("F[%d][%d]", i0, i1)+"."+path] = append(errs[fmt.Sprintf("F[%d][%d]", i0, i1)+"."+path], err)
						}
					}
					_ = _F
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
	if t.A != nil {
		_A := *t.A
		res.A = _A
	}
	if t.B != nil {
		_B := *t.B
		res.B = make([]string, len(_B))
		for i0, _B := range _B {
			res.B[i0] = _B
		}
	}
	if t.C != nil {
		_C := *t.C
		res.C = make([][]string, len(_C))
		for i0, _C := range _C {
			res.C[i0] = make([]string, len(_C))
			for i1, _C := range _C {
				res.C[i0][i1] = _C
			}
		}
	}
	if t.D != nil {
		_D := *t.D
		res.D = _D.Convert()
	}
	if t.E != nil {
		_E := *t.E
		res.E = make([]Name, len(_E))
		for i0, _E := range _E {
			res.E[i0] = _E.Convert()
		}
	}
	if t.F != nil {
		_F := *t.F
		res.F = make([][]Name, len(_F))
		for i0, _F := range _F {
			res.F[i0] = make([]Name, len(_F))
			for i1, _F := range _F {
				res.F[i0][i1] = _F.Convert()
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