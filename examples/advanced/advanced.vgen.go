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
	address, errs := t.InnerValidation()
	if len(errs) > 0 {
		j, _ := json.Marshal(errs)
		return Address{}, fmt.Errorf("%s", j)
	}
	return address, nil
}
func (t AddressVgen) InnerValidation() (Address, map[string][]string) {
	res := Address{}
	errs := make(map[string][]string)
	if t.Street != nil {
		street := *t.Street
		{
			if !(len(street) > 0) {
				errs[fmt.Sprintf("street")] = append(errs[fmt.Sprintf("street")], fmt.Sprintf("can not be empty"))
			}
			_ = street
		}
	} else {
		errs["street"] = append(errs["street"], fmt.Sprintf("required"))
	}
	if t.Number != nil {
		number := *t.Number
		{
			if !(number < 5) {
				errs[fmt.Sprintf("number")] = append(errs[fmt.Sprintf("number")], fmt.Sprintf("must be < 5"))
			}
			_ = number
		}
	} else {
		errs["number"] = append(errs["number"], fmt.Sprintf("required"))
	}
	if len(errs) > 0 {
		return Address{}, errs
	}
	return res, nil
}
type PersonVgen struct {
	Name       *string
	Address1   *Address
	Address2   *AddressVgen
	Addresses  *[]AddressVgen
	Addresses2 *[][]AddressVgen
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
		name := *t.Name
		{
			if !(len(name) > 0) {
				errs[fmt.Sprintf("name")] = append(errs[fmt.Sprintf("name")], fmt.Sprintf("can not be empty"))
			}
			if !(len(name) < 20) {
				errs[fmt.Sprintf("name")] = append(errs[fmt.Sprintf("name")], fmt.Sprintf("len must be < 20"))
			}
			_ = name
		}
	} else {
		errs["name"] = append(errs["name"], fmt.Sprintf("required"))
	}
	if t.Address1 != nil {
		address1 := *t.Address1
		{
			if err := valAddr(address1); err != nil {
				errs[fmt.Sprintf("address1")] = append(errs[fmt.Sprintf("address1")], err.Error())
			}
			_ = address1
		}
	}
	if t.Address2 != nil {
		address2 := *t.Address2
		{
			address2, struct_errs := address2.InnerValidation()
			for path, err_list := range struct_errs {
				for _, err := range err_list {
					errs[fmt.Sprintf("address2")+"."+path] = append(errs[fmt.Sprintf("address2")+"."+path], err)
				}
			}
			_ = address2
			if err := valAddr(address2); err != nil {
				errs[fmt.Sprintf("address2")] = append(errs[fmt.Sprintf("address2")], err.Error())
			}
		}
	} else {
		errs["address2"] = append(errs["address2"], fmt.Sprintf("required"))
	}
	if t.Addresses != nil {
		addresses := *t.Addresses
		{
			for i0, addresses := range addresses {
				addresses, struct_errs := addresses.InnerValidation()
				for path, err_list := range struct_errs {
					for _, err := range err_list {
						errs[fmt.Sprintf("addresses[%d]", i0)+"."+path] = append(errs[fmt.Sprintf("addresses[%d]", i0)+"."+path], err)
					}
				}
				_ = addresses
				_ = i0
			}
		}
	} else {
		errs["addresses"] = append(errs["addresses"], fmt.Sprintf("required"))
	}
	if t.Addresses2 != nil {
		addresses2 := *t.Addresses2
		{
			for i0, addresses2 := range addresses2 {
				for i1, addresses2 := range addresses2 {
					addresses2, struct_errs := addresses2.InnerValidation()
					for path, err_list := range struct_errs {
						for _, err := range err_list {
							errs[fmt.Sprintf("addresses2[%d][%d]", i0, i1)+"."+path] = append(errs[fmt.Sprintf("addresses2[%d][%d]", i0, i1)+"."+path], err)
						}
					}
					_ = addresses2
					_ = i1
				}
				_ = i0
			}
		}
	} else {
		errs["addresses2"] = append(errs["addresses2"], fmt.Sprintf("required"))
	}
	if len(errs) > 0 {
		return Person{}, errs
	}
	return res, nil
}
