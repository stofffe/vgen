package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/spf13/cobra"
)

const SUFFIX = ".vgen.go"

func createCommands() {
	// root
	rootCmd := &cobra.Command{
		Use:   "vgen",
		Short: "cli tool to generate validation logic",
		Long:  "generate validation logic from exsiting go struct",
		CompletionOptions: cobra.CompletionOptions{
			DisableDefaultCmd: true,
		},
	}

	// generate
	var genRecursive, genVerbose bool
	generateCmd := &cobra.Command{
		Use:   "generate",
		Short: "generate from exsisting go files",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			generateCommand(args, genRecursive, genVerbose)
		},
	}
	generateCmd.Flags().BoolVarP(&genRecursive, "recursive", "r", false, "recursively parse specified paths")
	generateCmd.Flags().BoolVarP(&genVerbose, "verbose", "v", false, "output more detailed information")

	// clean
	var cleanRecursive, cleanVerbose bool
	cleanCmd := &cobra.Command{
		Use:   "clean",
		Short: "clean exsisting vgen files",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			cleanCommand(args, cleanRecursive, cleanVerbose)
		},
	}
	cleanCmd.Flags().BoolVarP(&cleanRecursive, "recursive", "r", false, "recursively clean specified paths")
	cleanCmd.Flags().BoolVarP(&cleanVerbose, "verbose", "v", false, "output more detailed information")

	rootCmd.AddCommand(generateCmd, cleanCmd)
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}

}

func cleanCommand(args []string, recursive, verbose bool) {
	errors := []string{}
	removed := []string{}
	paths := []string{}
	for _, path := range args {
		path_info, err := os.Stat(path)
		if err != nil {
			errors = append(errors, fmt.Sprintf("could not open file info for %s", path))
			continue
		}

		// parse single file
		if !path_info.IsDir() {
			paths = append(paths, path)
			continue
		}

		// parse directory
		filepath.Walk(path, func(current_path string, info os.FileInfo, err error) error {
			// file tree traversal errors
			if err != nil {
				errors = append(errors, fmt.Sprintf("error walking file tree: %v", err))
				return nil
			}

			// recursive check
			if info.IsDir() {
				if !recursive && current_path != path {
					return filepath.SkipDir
				} else {
					return nil
				}
			}

			// add vgen files
			if strings.HasSuffix(info.Name(), SUFFIX) {
				paths = append(paths, current_path)
			}

			return nil
		})
	}

	for _, path := range paths {
		err := os.Remove(path)
		if err != nil {
			errors = append(errors, fmt.Sprintf("could not remove file %s: %v", path, err))
		} else {
			removed = append(removed, path)
		}
	}

	if verbose {
		fmt.Printf("errors: %d\n", len(errors))
		for _, e := range errors {
			fmt.Println(e)
		}
	}

	if verbose {
		fmt.Printf("removed files: %d\n", len(removed))
		for _, path := range removed {
			fmt.Printf("%s: removed\n", path)
		}
	}
}

func generateCommand(args []string, recursive, verbose bool) {
	errors := []string{}
	warnings := []string{}

	// get files to be parsed
	paths := []string{}
	for _, path := range args {
		path_info, err := os.Stat(path)
		if err != nil {
			errors = append(errors, fmt.Sprintf("could not open file info for %s", path))
			continue
		}

		// parse single file
		if !path_info.IsDir() {
			paths = append(paths, path)
			continue
		}

		// parse directory
		filepath.Walk(path, func(current_path string, info os.FileInfo, err error) error {
			// file tree traversal errors
			if err != nil {
				errors = append(errors, fmt.Sprintf("error walking file tree: %v", err))
				return nil
			}

			// recursive check
			if info.IsDir() {
				if !recursive && current_path != path {
					return filepath.SkipDir
				} else {
					return nil
				}
			}

			// skip generated files
			if strings.HasSuffix(info.Name(), SUFFIX) {
				return nil
			}

			paths = append(paths, current_path)
			return nil
		})
	}

	// parse files concurrently
	wg := sync.WaitGroup{}
	errc := make(chan string, len(paths))
	warnc := make(chan string, len(paths))
	genc := make(chan GeneratedFile, len(paths))
	for _, path := range paths {
		p := path
		wg.Add(1)
		go func() {
			defer wg.Done()
			n, err := handleFile(p)
			if err != nil {
				errc <- fmt.Sprintf("could not generate for file %s: %v", p, err)
				return
			}
			if n == 0 {
				warnc <- fmt.Sprintf("%s: no parseable types", p)
				return
			}
			genc <- GeneratedFile{
				path:      p,
				typeCount: n,
			}
		}()
	}
	wg.Wait()
	close(errc)
	close(genc)
	close(warnc)

	// output errors
	for e := range errc {
		errors = append(errors, e)
	}
	if verbose {
		fmt.Printf("errors: %d\n", len(errors))
		for _, e := range errors {
			fmt.Println(e)
		}
	}

	// output warnings
	for w := range warnc {
		warnings = append(warnings, w)
	}
	if verbose {
		fmt.Printf("warnings: %d\n", len(warnings))
		for _, w := range warnings {
			fmt.Println(w)
		}
	}

	// ouput generated files
	files := []GeneratedFile{}
	for path := range genc {
		files = append(files, path)
	}
	if verbose {
		fmt.Printf("generated files: %d\n", len(files))
		for _, file := range files {
			fmt.Printf("%s: parsed %d types\n", file.path, file.typeCount)
		}
	}
}

type GeneratedFile struct {
	path      string
	typeCount int
}

func handleFile(path string) (int, error) {
	// parse file
	info, err := parseFile(path)
	if err != nil {
		return 0, fmt.Errorf("could not parse file: %v", err)
	}

	// generate vgen file from info
	buffer, err := generateFile(info)
	if err != nil {
		return 0, fmt.Errorf("could not generate file: %v", err)
	}

	// write new file
	file_name := strings.Replace(path, ".go", SUFFIX, 1)
	file, err := os.Create(file_name)
	if err != nil {
		return 0, fmt.Errorf("could not create file %v", file_name)
	}
	_, err = file.Write(buffer)
	if err != nil {
		return 0, fmt.Errorf("could not write to file %v", file_name)
	}

	return len(info.StructTypes), nil
}
