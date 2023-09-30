package main

import (
	"log"
	"os"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("must supply file path")
	}

	file_path := os.Args[1]

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
