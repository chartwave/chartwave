package chart

import (
	"fmt"
	"github.com/chartwave/chartwave/pkg/template"
	"io"
	"os"
	"path/filepath"
)

func (c *Chart) Template(out io.Writer) error {
	files, err := ListRecursiveFiles(c.Path)
	if err != nil {
		return err
	}

	tmpl := template.New()
	err = tmpl.AddFiles(files)
	if err != nil {
		return fmt.Errorf("failed to template chart: %w", err)
	}

	return tmpl.Run(out)
}

func ListRecursiveFiles(path string) ([]string, error) {
	result := make([]string, 0)

	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read %q directory: %w", path, err)
	}

	for _, entry := range entries {
		p := filepath.Join(path, entry.Name())
		if entry.IsDir() {
			subresult, err := ListRecursiveFiles(p)
			if err != nil {
				return nil, err
			}

			result = append(result, subresult...)
		} else {
			result = append(result, p)
		}
	}

	return result, nil
}
