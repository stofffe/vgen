package main

// package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/stofffe/vgen/examples"
)

func main() {
	person, err := examples.PersonVgen{
		Name:  P("helo"),
		Age:   P(17),
		Vibes: nil,
	}.Validate()
	if err != nil {
		PrettyPrintJson("err", err.Error())
		log.Fatalf("failed to validate person: %v", err)
	}
	fmt.Printf("person: %v\n", person)
}

func P[T any](t T) *T {
	return &t
}
func PrettyPrintJson(name string, val string) {
	var unmarshalled any
	err := json.Unmarshal([]byte(val), &unmarshalled)
	if err != nil {
		log.Fatalf("could not unmarshal: %v", val)
	}
	j, err := json.MarshalIndent(unmarshalled, "", "  ")
	if err != nil {
		log.Fatalf("could not pretty print: %v", err)
	}
	fmt.Printf(`
----------------------------------
Pretty print %s
%s
----------------------------------
`, name, string(j))
}

// package main
//
// import (
// 	"fmt"
// 	"regexp"
// )

// func main() {
//
// 	// Define the regular expression pattern
// 	req := `(req)`
// 	len_gt := `(len>)(\d+)`
// 	len_lt := `(len<)(\d+)`
// 	len_gte := `(len>=)(\d+)`
// 	len_lte := `(len<=)(\d+)`
// 	pattern := req + "|" + len_gt + "|" + len_lt + "|" + len_gte + "|" + len_lte
//
// 	// Compile the regular expression
// 	regex := regexp.MustCompile(pattern)
//
// 	f, v := check("len<10", regex)
// 	fmt.Println(rules[f](v))
// 	f, v = check("len<=20", regex)
// 	fmt.Println(rules[f](v))
// 	f, v = check("req", regex)
// 	fmt.Println(rules[f](v))
// 	f, v = check("reeq", regex)
// 	fmt.Println(rules[f](v))
// }
//
// var rules = map[string]func(string) string{
// 	"req":   func(val string) string { return "REQ RULE" },
// 	"len<":  func(val string) string { return "LEN < RULE WITH " + val },
// 	"len<=": func(val string) string { return "LEN < RULE WITH " + val },
// 	"len>":  func(val string) string { return "LEN > RULE WITH " + val },
// 	"len>=": func(val string) string { return "LEN > RULE WITH " + val },
// }
//
// func check(str string, reg *regexp.Regexp) (string, string) {
// 	matches := reg.FindStringSubmatch(str)
// 	if len(matches) == 0 {
// 		fmt.Printf("invalid rule: %v\n", str)
// 		return "", ""
// 	}
// 	var filtered []string
// 	for i, v := range matches {
// 		if i != 0 && v != "" {
// 			filtered = append(filtered, v)
// 		}
// 	}
// 	// fmt.Printf("Matched: %s, len: %d\n", filtered, len(filtered))
// 	f := filtered[0]
// 	v := ""
// 	if len(filtered) > 1 {
// 		v = filtered[1]
// 	}
// 	return f, v
// }
