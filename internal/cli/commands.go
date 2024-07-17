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
	var genRecursive, genVerbose bool
	generateCmd := &cobra.Command{
		Use:   "generate",
		Short: "generate from exsisting go files",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			generate(args, genRecursive, genVerbose)
		},
	}
	generateCmd.Flags().BoolVarP(&genRecursive, "recursive", "r", false, "recursively parse specified paths")
	generateCmd.Flags().BoolVarP(&genVerbose, "verbose", "v", false, "output more detailed information")
	rootCmd.AddCommand(generateCmd)

	// clean
	var cleanRecursive, cleanVerbose bool
	cleanCmd := &cobra.Command{
		Use:   "clean",
		Short: "clean exsisting vgen files",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			clean(args, cleanRecursive, cleanVerbose)
		},
	}
	cleanCmd.Flags().BoolVarP(&cleanRecursive, "recursive", "r", false, "recursively clean specified paths")
	cleanCmd.Flags().BoolVarP(&cleanVerbose, "verbose", "v", false, "output more detailed information")
	rootCmd.AddCommand(cleanCmd)

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

type CleanedFileInfo struct {
	path string
}
type CleanedFileError struct {
	path string
	err  DetailedError
}

func (c CleanedFileInfo) Format() string {
	return fmt.Sprintf("%s: removed", c.path)
}
func (c CleanedFileError) Format() string {
	return fmt.Sprintf("%s: %s", c.path, c.err.detailed)
}

func clean(args []string, recursive, verbose bool) {
	errors := []CleanedFileError{}
	removed := []CleanedFileInfo{}
	paths := []string{}

	// get files to be removed
	for _, path := range args {
		path := path
		pathInfo, err := os.Stat(path)
		if err != nil {
			errors = append(errors, CleanedFileError{
				path: path,
				err: DetailedError{
					inner:    fmt.Errorf("could not open file %s", path),
					detailed: fmt.Errorf("could not open file info for %s", path),
				},
			})
			continue
		}

		// parse single file
		if !pathInfo.IsDir() {
			paths = append(paths, path)
			continue
		}

		// parse directory
		filepath.Walk(path, func(current_path string, info os.FileInfo, err error) error {
			// file tree traversal errors
			if err != nil {
				errors = append(errors, CleanedFileError{
					path: path,
					err:  NewInternalError(fmt.Errorf("error walking file tree: %v", err)),
				})
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
			if strings.HasSuffix(info.Name(), suffix) {
				paths = append(paths, current_path)
			}

			return nil
		})
	}

	// remove files
	for _, path := range paths {
		path := path
		err := os.Remove(path)
		if err != nil {
			errors = append(errors, CleanedFileError{
				path: path,
				err:  NewInternalError(fmt.Errorf("could not remove file %s: %v", path, err)),
			})
		} else {
			removed = append(removed, CleanedFileInfo{
				path: path,
			})
		}
	}

	// verbose logging
	if verbose {
		fmt.Printf("errors: %d\n", len(errors))
		for _, e := range errors {
			fmt.Printf("\t%s\n", e.Format())
		}

		fmt.Printf("info: %d\n", len(removed))
		for _, r := range removed {
			fmt.Printf("\t%s\n", r.Format())
		}
	}
}

type GeneratedFileInfo struct {
	path      string
	typeCount int
}
type GeneratedFileWarning struct {
	warning string
	path    string
}
type GeneratedFileError struct {
	path string
	err  DetailedError
}

func (g GeneratedFileInfo) Format() string {
	return fmt.Sprintf("\t%s: parsed %d types", g.path, g.typeCount)
}
func (g GeneratedFileWarning) Format() string {
	return fmt.Sprintf("\t%s: %s", g.path, g.warning)
}
func (g GeneratedFileError) Format(detailed bool) string {
	if detailed {
		return fmt.Sprintf("\t%s: %s", g.path, g.err.detailed)
	} else {
		return fmt.Sprintf("\t%s: %s", g.path, g.err.inner)
	}
}

func generate(args []string, recursive, verbose bool) {
	errors := []GeneratedFileError{}
	warnings := []GeneratedFileWarning{}
	info := []GeneratedFileInfo{}

	// get files to be parsed
	paths := []string{}
	for _, path := range args {
		fileInfo, err := os.Stat(path)
		if err != nil {
			errors = append(errors, GeneratedFileError{
				path: path,
				err: DetailedError{
					inner:    fmt.Errorf("could not open file"),
					detailed: fmt.Errorf("could not open file info"),
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
				errors = append(errors, GeneratedFileError{
					path: path,
					err:  NewInternalError(fmt.Errorf("could not walk file tree: %v", err)),
				})
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
			if strings.HasSuffix(info.Name(), suffix) {
				return nil
			}

			paths = append(paths, current_path)
			return nil
		})
	}

	// parse files concurrently
	wg := sync.WaitGroup{}
	errorc := make(chan GeneratedFileError, len(paths))
	warnc := make(chan GeneratedFileWarning, len(paths))
	infoc := make(chan GeneratedFileInfo, len(paths))
	for _, path := range paths {
		path := path // TODO fixed in 1.22?
		wg.Add(1)
		go func() {
			defer wg.Done()
			n, err := handleFile(path)
			if err != nil {
				errorc <- GeneratedFileError{
					path: path,
					err:  NewDetailedError(err, fmt.Sprintf("could not handle file")),
				}

				return
			}
			if n == 0 {
				warnc <- GeneratedFileWarning{
					warning: "no parseable types",
					path:    path,
				}
				return
			}
			infoc <- GeneratedFileInfo{
				path:      path,
				typeCount: n,
			}
		}()
	}
	wg.Wait()
	close(errorc)
	close(infoc)
	close(warnc)

	// verbose logging
	if verbose {
		for e := range errorc {
			errors = append(errors, e)
		}
		for w := range warnc {
			warnings = append(warnings, w)
		}
		for i := range infoc {
			info = append(info, i)
		}

		// log
		fmt.Printf("errors: %d\n", len(errors))
		for _, e := range errors {
			fmt.Println(e.Format(true))
		}
		fmt.Printf("warnings: %d\n", len(warnings))
		for _, w := range warnings {
			fmt.Println(w.Format())
		}
		fmt.Printf("info: %d\n", len(info))
		for _, info := range info {
			fmt.Println(info.Format())
		}
	}
}

func handleFile(path string) (int, error) {
	// parse file
	info, err := parseFile(path)
	if err != nil {
		return 0, NewDetailedError(err, "%s: could not parse file")
	}

	// generate vgen file from info
	buffer, err := generateFile(info)
	if err != nil {
		return 0, NewDetailedError(err, "could not generate file")
	}

	// write new file
	fileName := strings.Replace(path, ".go", suffix, 1)
	file, err := os.Create(fileName)
	if err != nil {
		return 0, NewInternalError(fmt.Errorf("could not create file %s: %v", fileName, err))
	}
	_, err = file.Write(buffer)
	if err != nil {
		return 0, NewInternalError(fmt.Errorf("could not write to file %v", fileName))

	}

	return len(info.StructTypes), nil
}
