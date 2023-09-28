package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
)

func main() {
	file_path := "examples/test.go"

	// parse
	info, err := parseFile(file_path)
	if err != nil {
		log.Fatal(err)
	}

	// generate
	buffer, err := generateFile(info)
	if err != nil {
		log.Fatal(err)
	}

	// file
	file_name := strings.Replace(file_path, ".go", ".vgen.go", 1)
	file, err := os.Create(file_name)
	if err != nil {
		log.Fatal(err)
	}
	_, err = file.Write(buffer)
	if err != nil {
		log.Fatal(err)
	}
}

func PrettyPrint(name string, val any) {
	j, err := json.MarshalIndent(val, "", "  ")
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

// func parseFieldStruct(node *ast.StructType) error {
// 	if node.Fields == nil || len(node.Fields.List) == 0 {
// 		return fmt.Errorf("empty structs not supported")
// 	}
//
// 	for _, field := range node.Fields.List {
// 		err := parseField(field)
// 		if err != nil {
// 			return fmt.Errorf("could not parse field: %v", err)
// 		}
// 	}
//
// 	return nil
// }
