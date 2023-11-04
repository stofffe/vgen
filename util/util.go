package util

import (
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"strings"
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

//
// func LowerFirstChar(str string) string {
// 	if str == "" {
// 		return str
// 	}
// 	firstchar := []rune(str)[0]
// 	firstchar = unicode.ToLower(firstchar)
// 	return string(firstchar) + str[1:]
// }

func ExtractJsonName(tag, backup string) string {
	reg := regexp.MustCompile(`json:".+\"`)
	match := reg.FindString(tag)
	match = strings.TrimPrefix(match, `json:"`)
	match = strings.TrimSuffix(match, `"`)
	match = strings.Split(match, ",")[0]
	if match == "" {
		return backup
	}
	return match
}

// func LowerFirstChar(str string) string {
// 	return str + "_"
// }
