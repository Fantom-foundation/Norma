package scenarios

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Fantom-foundation/Norma/driver/parser"
)

// TestCheckScenarious iterates through all scenarios in this directory
// and its sub-directories and checks whether the contained YAML files
// define valid scenarios.
func TestCheckScenarious(t *testing.T) {
	files, err := listAll()
	if err != nil {
		t.Fatalf("failed to get list of all scenario files: %v", err)
	}
	if len(files) == 0 {
		t.Fatalf("failed to locate any scenario files!")
	}
	for _, file := range files {
		t.Run(file, func(t *testing.T) {
			scenaro, err := parser.ParseFile(file)
			if err != nil {
				t.Fatalf("failed to parse file: %v", err)
			}
			if err = scenaro.Check(); err != nil {
				t.Fatalf("scenaro check failed: %v", err)
			}
		})
	}
}

func listAll() ([]string, error) {
	files := []string{}
	err := filepath.Walk(".",
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if strings.HasSuffix(path, ".yml") {
				files = append(files, path)
			}
			return nil
		})
	if err != nil {
		return nil, err
	}
	return files, nil
}
