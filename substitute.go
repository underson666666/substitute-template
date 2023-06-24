package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v2"
)

// Substitution represents the substitution data for a template file.
type Substitution map[string]string

// Project represents the YAML configuration file structure.
type Project struct {
	Dir          string       `yaml:"dir"`
	Template     string       `yaml:"template"`
	Substitution Substitution `yaml:"substitution"`
}

func main() {
	// go run substitute.go -t from -o to -s substitution.yaml

	// Parse command-line options
	templateDir := flag.String("t", "", "Directory containing template files (*.tmpl)")
	outputDir := flag.String("o", "", "Output directory")
	substitutionFile := flag.String("s", "", "Substitution file (substitution.yaml)")
	flag.Parse()

	substitutionData, err := ioutil.ReadFile(*substitutionFile)
	if err != nil {
		log.Fatalf("Failed to read substitution file: %v", err)
	}

	// Parse YAML data into Project struct
	var data map[string][]Project
	err = yaml.Unmarshal(substitutionData, &data)
	if err != nil {
		log.Fatalf("Failed to parse substitution file: %v", err)
	}

	// Process each project
	for projectName, projects := range data {
		projectOutputDir := filepath.Join(*outputDir, projectName)
		fmt.Printf("outputdir: %s\n", projectOutputDir)

		// Process each substitution
		for _, project := range projects {

			templateFile := filepath.Join(*templateDir, project.Dir, project.Template)
			outputFile := filepath.Join(projectOutputDir, project.Dir, strings.TrimSuffix(project.Template, ".tmpl"))

			// // Read template file
			templateData, err := ioutil.ReadFile(templateFile)
			if err != nil {
				log.Fatalf("Failed to read template file: %v", err)
			}

			// // Apply substitutions
			outputData := applySubstitutions(templateData, project.Substitution)

			// Create directories if they don't exist
			err = os.MkdirAll(filepath.Dir(outputFile), 0755)
			if err != nil {
				log.Fatalf("Failed to create output directory: %v", err)
			}

			// Write output file
			err = ioutil.WriteFile(outputFile, outputData, 0644)
			if err != nil {
				log.Fatalf("Failed to write output file: %v", err)
			}
		}
	}

	fmt.Println("Template substitution completed successfully.")
}

// applySubstitutions applies the given substitutions to the template data.
func applySubstitutions(templateData []byte, substitutions Substitution) []byte {
	outputData := templateData

	// for _, substitution := range substitutions {
	for key, value := range substitutions {
		placeholder := "{{ " + key + " }}"
		outputData = bytesReplace(outputData, []byte(placeholder), []byte(value))
	}
	// }

	return outputData
}

// bytesReplace replaces all occurrences of old with new in the given data.
func bytesReplace(data []byte, old []byte, new []byte) []byte {
	return bytes.Replace(data, old, new, -1)
}
