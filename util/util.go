package util

import (
	"encoding/json"
	"fmt"
	"log"
)

// initialze pointer
// useful when creating pointers with value
func InitP[T any](t T) *T {
	return &t
}

// prints any value by marshalling
// crashes if value is not marshallable
func DebugPrint(name string, val string) {
	var unmarshalled any
	err := json.Unmarshal([]byte(val), &unmarshalled)
	if err != nil {
		log.Fatalf("could not unmarshal: %v", val)
	}
	j, err := json.MarshalIndent(unmarshalled, "", "  ")
	if err != nil {
		log.Fatalf("could not debug print: %v", err)
	}
	fmt.Printf(`
----------------------------------
Debug print %s
%s
----------------------------------
`, name, string(j))
}

// check if a list contains a specific value
func ListContains[T comparable](list []T, val T) bool {
	for _, v := range list {
		if v == val {
			return true
		}
	}
	return false
}
