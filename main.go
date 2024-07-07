package main

import (
	"fmt"
	"log"

	"github.com/stofffe/vgen/cli"
)

func main() {
	path := "examples/person/person.go"
	count, err := cli.HandleFile(path)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(count)
}
