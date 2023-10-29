package main

import (
	"fmt"
	"log"
	"os"
	"strings"
)

func main() {
	// input := "asdaskdas vgen:[req, gt=3][len_gt=5, custom=abc][not_empty]   dssadassd"
	// test(input)
	// return

	// require file path
	if len(os.Args) < 2 {
		log.Fatal("must supply file path")
	}

	// parse file
	path := os.Args[1]
	info, err := parseFile(path)
	if err != nil {
		log.Fatalf("could not parse file: %v", err)
	}

	// generate vgen file from info
	buffer, err := generateFile(info)
	if err != nil {
		log.Fatalf("could not generate file: %v", err)
	}

	// write new file
	file_name := strings.Replace(path, ".go", ".vgen.go", 1)
	file, err := os.Create(file_name)
	if err != nil {
		log.Fatal(err)
	}
	_, err = file.Write(buffer)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%v\n", string(buffer))
}
