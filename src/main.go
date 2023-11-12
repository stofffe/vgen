package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	// flags
	var recursive_flag bool
	flag.BoolVar(&recursive_flag, "r", false, "recursively parse files")
	flag.Parse()

	args := flag.Args()
	if len(args) == 0 {
		log.Fatal("must supply file path")
	}
	path := args[0]
	path_info, err := os.Stat(path)
	if err != nil {
		log.Fatalf("could not open file %s", path)
	}

	// parse single file
	if !path_info.IsDir() {
		handleFile(path)
		return
	}

	// parse directory
	err = filepath.Walk(path, func(current_path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// recursive check
		if info.IsDir() {
			if !recursive_flag && current_path != path {
				return filepath.SkipDir
			} else {
				return nil
			}
		}

		// skip generated files
		if strings.HasSuffix(info.Name(), ".vgen.go") {
			return nil
		}

		err = handleFile(current_path)
		if err != nil {
			fmt.Printf("%s: %v\n", current_path, err)
		}
		return nil
	})
	if err != nil {
		log.Fatalf("error walking file tree: %v", err)
	}
}

func handleFile(path string) error {
	// parse file
	info, err := parseFile(path)
	if err != nil {
		return fmt.Errorf("could not parse file: %v", err)
	}

	// generate vgen file from info
	buffer, err := generateFile(info)
	if err != nil {
		return fmt.Errorf("could not generate file: %v", err)
	}

	// write new file
	file_name := strings.Replace(path, ".go", ".vgen.go", 1)
	file, err := os.Create(file_name)
	if err != nil {
		return fmt.Errorf("could not create file %v", file_name)
	}
	_, err = file.Write(buffer)
	if err != nil {
		return fmt.Errorf("could not write to file %v", file_name)
	}

	return nil
}
