package main

import "github.com/stofffe/vgen/vgen"

type PersonVgen struct {
	name      *string
	age       *int
	nicknames *[]string
	matrix    *[][]int
}

func (p PersonVgen) Validate(rules PersonRules) vgen.ErrorMap {
	errors := make(vgen.ErrorMap)
	new_errors := make(vgen.ErrorMap)
	// name
	new_errors = rules.name.Validate("name", p.name)
	errors.AddErrors(new_errors)
	// age
	new_errors = rules.age.Validate("age", p.age)
	errors.AddErrors(new_errors)
	// nicknames
	errors.AddErrors(rules.nicknames.Validate("nicknames", p.nicknames))
	// matrix
	errors.AddErrors(rules.matrix.Validate("matrix", p.matrix))

	return errors
}

type PersonRules struct {
	name      vgen.Rules[string]
	age       vgen.Rules[int]
	nicknames vgen.Rules[[]string]
	matrix    vgen.Rules[[][]int]
}
