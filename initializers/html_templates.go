package initializers

import (
	"log"
	"os"
	"path/filepath"
	"text/template"
)

var Templates *template.Template

func LoadHTMLTemplates() {
	root := "views"

	// Initialize a new template with a function map
	Templates = template.New("").Funcs(template.FuncMap{
		"add":               add,
		"selectDescription": selectDescription,
	})

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Only parse files with the .html extension
		if !info.IsDir() && filepath.Ext(path) == ".html" {
			_, err = Templates.ParseFiles(path)
			if err != nil {
				return err
			}
		}
		return nil
	})

	if err != nil {
		log.Fatalf("Failed to load HTML templates: %v", err)
	}
}

func add(a, b int) int {
	return a + b
}

func selectDescription(measurementSystem, description, descriptionMetric string) string {
	if measurementSystem == "metric" {
		return descriptionMetric
	}
	return description
}
