package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func generateCmd() *cobra.Command {
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
	// Read terraform-mock.hcl
	config, err := readHCLConfig("terraform-mock.hcl")
	if err != nil {
		return fmt.Errorf("failed to read terraform-mock.hcl: %w", err)
	}

	// Parse the source module
	module, err := parseModuleDetailed(config.Source)
	if err != nil {
		return fmt.Errorf("failed to parse source module: %w", err)
	}

	// Filter resources based on keep_resources
	filtered := filterResources(module, config.KeepResources)

	// Create output directory
	if err := os.MkdirAll(config.Output, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Generate new Terraform files
	if err := generateTerraformFiles(config.Output, filtered); err != nil {
		return fmt.Errorf("failed to generate Terraform files: %w", err)
	}

	// Generate LocalStack provider configuration
	if err := generateLocalStackProvider(config.Output); err != nil {
		return fmt.Errorf("failed to generate provider config: %w", err)
	}

	fmt.Printf("Generated mocked Terraform files in: %s\n", config.Output)

	return nil
}
