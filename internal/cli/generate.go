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
		Short: "Generate Terraform files for Localstack from configuration",
		RunE: func(cmd *cobra.Command, args []string) error {
			return generate()
		},
	}

	return cmd
}

func generate() error {
	config, err := config.ReadHCLConfig("terraform-hydra.hcl")
	if err != nil {
		return fmt.Errorf("failed to read terraform-hydra.hcl: %w", err)
	}

	module, err := parser.ParseModuleDetailed(config.Source)
	if err != nil {
		return fmt.Errorf("failed to parse source module: %w", err)
	}

	filtered := generator.FilterResources(module, config.KeepResources)

	if err := os.MkdirAll(config.Output, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	if err := generator.GenerateTerraformFiles(config.Output, filtered, config.RemoveAttributes); err != nil {
		return fmt.Errorf("failed to generate Terraform files: %w", err)
	}

	if err := generator.GenerateLocalStackProvider(config.Output); err != nil {
		return fmt.Errorf("failed to generate provider config: %w", err)
	}

	fmt.Printf("Generated Terraform files for Localstack in: %s\n", config.Output)

	return nil
}
