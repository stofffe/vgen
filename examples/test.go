package examples

import (
	"encoding/json"
	"fmt"
)

// vgen:"i"
type Person struct {
	Name  string // vgen:[ req, len>0 ]
	Age   int    // invalidvgen
	Vibes bool   // vgen:[]
}

// vgen:[ val<=100, val>=0 ]
// vgen:"i"
type Person2 struct {
	Name  string `vgen:[ req, len>0 ]`
	Age   int    `vgen:[ val<=100, val>=0 ]`
	Vibes bool   `vgen:[ req ]`
}

// Person ouput

type PersonVgen struct {
	Name  *string
	Age   *int
	Vibes *bool
}

func (t PersonVgen) Validate() (Person, error) {
	errs := make(map[string]string)

	if t.Name != nil {
		name := *t.Name

		// len>0
		if !(len(name) > 0) {
			errs["name"] = fmt.Sprintf("len must be > 0")
		}
	} else {
		// req
		errs["name"] = fmt.Sprintf("required")
	}

	if t.Age != nil {
		age := *t.Age

		// val<=100
		if !(age <= 100) {
			errs["age"] = fmt.Sprintf("value must be <= 100")
		}

		// val>=0
		if !(age > 0) {
			errs["age"] = fmt.Sprintf("value must be > 0")
		}
	}

	if t.Vibes != nil {

	} else {
		errs["vibes"] = fmt.Sprintf("required")
	}

	if len(errs) > 0 {
		fmt.Println("err > 0", errs)
		j, err := json.Marshal(errs)
		if err != nil {
			return Person{}, fmt.Errorf("ERROR MARSHALLING SHOULD NOT HAPPEN")
		}
		return Person{}, fmt.Errorf("%s", j)
	}

	return Person{
		Name:  *t.Name,
		Age:   *t.Age,
		Vibes: *t.Vibes,
	}, nil

}

type EmptyType struct{}

type (
	A string
	B int
)
