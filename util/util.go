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
func DebugPrintAny(name string, val any) {
	j, err := json.MarshalIndent(val, "", "  ")
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

// prints any value by marshalling
// crashes if value is not marshallable
func DebugPrintString(name string, val string) {
	var unmarshalled any
	err := json.Unmarshal([]byte(val), &unmarshalled)
	if err != nil {
		log.Fatalf("could not unmarshal: %v", val)
	}
	DebugPrintAny(name, unmarshalled)
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
