package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/georgetjose/openapi3gen/pkg/generator"
	"github.com/georgetjose/openapi3gen/pkg/parser"

	"github.com/spf13/cobra"
)

var (
	dir    string
	output string
)

func init() {
	generateCmd.Flags().StringVar(&dir, "dir", ".", "Directory of the Gin project")
	generateCmd.Flags().StringVar(&output, "output", "./openapi.json", "Output OpenAPI JSON file")
	rootCmd.AddCommand(generateCmd)
}

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate OpenAPI 3.0 spec from annotated routes",
	RunE: func(cmd *cobra.Command, args []string) error {
		// 1. Parse annotations
		routes, err := parser.ParseDirectory(dir)
		if err != nil {
			return fmt.Errorf("failed to parse: %w", err)
		}

		// 2. Register models
		registry := generator.NewModelRegistry()
		// TODO: optionally support JSON schema registry discovery later

		// 3. Generate spec
		spec := generator.GenerateSpec(routes, registry)

		// 4. Write to output
		specJson, err := json.MarshalIndent(spec, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal: %w", err)
		}

		if err := os.MkdirAll(filepath.Dir(output), os.ModePerm); err != nil {
			return err
		}

		if err := os.WriteFile(output, specJson, 0644); err != nil {
			return fmt.Errorf("failed to write: %w", err)
		}

		fmt.Printf("✅ OpenAPI 3.0 spec generated at: %s\n", output)
		return nil
	},
}
