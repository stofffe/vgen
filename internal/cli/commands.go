package cli

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/spf13/cobra"
)

const suffix = ".vgen.go"

func CreateCommands() {
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
	var genVerbose bool
	generateCmd := &cobra.Command{
		Use:   "generate",
		Short: "generate from exsisting go files",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			generate(args, genVerbose)
		},
	}
	generateCmd.Flags().BoolVarP(&genVerbose, "verbose", "v", false, "output verbose errors")
	rootCmd.AddCommand(generateCmd)

	// clean
	var cleanVerbose bool
	cleanCmd := &cobra.Command{
		Use:   "clean",
		Short: "clean exsisting vgen files",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			clean(args, cleanVerbose)
		},
	}
	cleanCmd.Flags().BoolVarP(&cleanVerbose, "verbose", "v", false, "output verbose errors")
	rootCmd.AddCommand(cleanCmd)

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func clean(args []string, verbose bool) {
	errors := []ErrorMessage{}
	warnings := []WarningMessage{}
	infos := []InfoMessage{}
	paths := []string{}

	// get files to be removed
	for _, path := range args {
		path := path
		pathInfo, err := os.Stat(path)
		if err != nil {
			errors = append(errors, ErrorMessage{
				path: path,
				err: DetailedError{
					msg: "could not open file",
					err: fmt.Errorf("open file info for %s", path),
				},
			})
			continue
		}

		if !pathInfo.IsDir() {
			paths = append(paths, path)
			continue
		}

		// parse directory
		filepath.Walk(path, func(current_path string, info os.FileInfo, err error) error {
			// file tree traversal errors
			if err != nil {
				errors = append(errors, ErrorMessage{
					path: path,
					err:  fmt.Errorf("error walking file tree: %v", err),
				})
				return nil
			}

			// dont parse folders
			if info.IsDir() {
				return nil
			}

			// add vgen files
			if strings.HasSuffix(info.Name(), suffix) {
				paths = append(paths, current_path)
			}

			return nil
		})
	}

	if len(paths) == 0 {
		warnings = append(warnings, WarningMessage{
			msg: "no files found",
		})
	}

	// remove files
	for _, path := range paths {
		path := path
		err := os.Remove(path)
		if err != nil {
			errors = append(errors, ErrorMessage{
				path: path,
				err:  fmt.Errorf("remove file: %v", err),
			})
		} else {
			infos = append(infos, InfoMessage{
				msg:  "removed file",
				path: path,
			})
		}
	}

	// logging
	for _, e := range errors {
		fmt.Println(e.Format(verbose))
	}
	for _, w := range warnings {
		fmt.Println(w.Format())
	}
	for _, r := range infos {
		fmt.Println(r.Format())
	}
}

func generate(args []string, verbose bool) {
	errors := []ErrorMessage{}
	warnings := []WarningMessage{}
	info := []InfoMessage{}

	// get files to be parsed
	paths := []string{}
	for _, path := range args {
		fileInfo, err := os.Stat(path)
		if err != nil {
			errors = append(errors, ErrorMessage{
				path: path,
				err: DetailedError{
					msg: "could not open file",
					err: fmt.Errorf("open file info"),
				},
			})
			continue
		}

		// parse single file
		if !fileInfo.IsDir() {
			paths = append(paths, path)
			continue
		}

		// parse directory
		filepath.Walk(path, func(current_path string, info os.FileInfo, err error) error {
			// file tree traversal errors
			if err != nil {
				errors = append(errors, ErrorMessage{
					path: path,
					err:  fmt.Errorf("walk file tree: %w", err),
				})
				return nil
			}

			// dont parse folders
			if info.IsDir() {
				return nil
			}

			// skip generated files
			if strings.HasSuffix(info.Name(), suffix) {
				return nil
			}

			paths = append(paths, current_path)
			return nil
		})
	}

	// parse files concurrently
	wg := sync.WaitGroup{}
	errorc := make(chan ErrorMessage, len(paths))
	warnc := make(chan WarningMessage, len(paths))
	infoc := make(chan InfoMessage, len(paths))
	for _, path := range paths {
		path := path // TODO fixed in 1.22?
		wg.Add(1)
		go func() {
			defer wg.Done()
			n, err := handleFile(path)
			if err != nil {
				errorc <- ErrorMessage{
					err:  fmt.Errorf("handle file: %w", err),
					path: path,
				}

				return
			}
			if n == 0 {
				warnc <- WarningMessage{
					msg:  "no parseable types",
					path: path,
				}
				return
			} else if n == 1 {
				infoc <- InfoMessage{
					msg:  fmt.Sprintf("parsed 1 type"),
					path: path,
				}
			} else {
				infoc <- InfoMessage{
					msg:  fmt.Sprintf("parsed %d types", n),
					path: path,
				}
			}
		}()
	}
	wg.Wait()
	close(errorc)
	close(infoc)
	close(warnc)

	for e := range errorc {
		errors = append(errors, e)
	}
	for w := range warnc {
		warnings = append(warnings, w)
	}
	for i := range infoc {
		info = append(info, i)
	}

	// logging
	for _, e := range errors {
		fmt.Println(e.Format(verbose))
	}
	for _, w := range warnings {
		fmt.Println(w.Format())
	}
	for _, info := range info {
		fmt.Println(info.Format())
	}
}

func handleFile(path string) (int, error) {
	// parse file
	info, err := parseFile(path)
	if err != nil {
		return 0, fmt.Errorf("parse file: %w", err)
	}

	// generate vgen file from info
	buffer, err := generateFile(info)
	if err != nil {
		return 0, fmt.Errorf("generate file: %w", err)
	}

	// write new file
	fileName := strings.Replace(path, ".go", suffix, 1)
	file, err := os.Create(fileName)
	if err != nil {
		return 0, fmt.Errorf("create file: %w", err)
	}
	_, err = file.Write(buffer)
	if err != nil {
		return 0, fmt.Errorf("write to file: %w", err)
	}

	return len(info.StructTypes), nil
}
