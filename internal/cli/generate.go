package cli

import (
	"fmt"
	"os"

	"hydratf/internal/config"
	"hydratf/internal/generator"
	"hydratf/internal/parser"

	"github.com/spf13/cobra"
)

func GenerateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "generate",
		Short: "Generate mocked Terraform files from configuration",
		RunE: func(cmd *cobra.Command, args []string) error {
			return generate()
		},
	}

	return cmd
}

func generate() error {
	config, err := config.ReadHCLConfig("terraform-mock.hcl")
	if err != nil {
		return fmt.Errorf("failed to read terraform-mock.hcl: %w", err)
	}

	module, err := parser.ParseModuleDetailed(config.Source)
	if err != nil {
		return fmt.Errorf("failed to parse source module: %w", err)
	}

	filtered := generator.FilterResources(module, config.KeepResources)

	if err := os.MkdirAll(config.Output, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	if err := generator.GenerateTerraformFiles(config.Output, filtered); err != nil {
		return fmt.Errorf("failed to generate Terraform files: %w", err)
	}

	if err := generator.GenerateLocalStackProvider(config.Output); err != nil {
		return fmt.Errorf("failed to generate provider config: %w", err)
	}

	fmt.Printf("Generated mocked Terraform files in: %s\n", config.Output)

	return nil
}
