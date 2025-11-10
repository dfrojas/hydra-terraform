package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

func initCmd() *cobra.Command {
	var source string
	var name string

	cmd := &cobra.Command{
		Use:   "init",
		Short: "Initialize a mock configuration for a Terraform module",
		RunE: func(cmd *cobra.Command, args []string) error {
			if source == "" {
				return fmt.Errorf("--source is required")
			}
			if name == "" {
				return fmt.Errorf("--name is required")
			}

			return initMock(source, name)
		},
	}

	cmd.Flags().StringVar(&source, "source", "", "Path to the original Terraform module")
	cmd.Flags().StringVar(&name, "name", "", "Name for the mock configuration")

	return cmd
}

func initMock(source, name string) error {
	// Verify source exists
	absSource, err := filepath.Abs(source)
	if err != nil {
		return fmt.Errorf("failed to resolve source path: %w", err)
	}

	if _, err := os.Stat(absSource); os.IsNotExist(err) {
		return fmt.Errorf("source path does not exist: %s", absSource)
	}

	// Parse the Terraform module to find resources
	resources, err := parseModule(absSource)
	if err != nil {
		return fmt.Errorf("failed to parse module: %w", err)
	}

	// Generate HCL configuration
	outputDir := filepath.Join("mocks", name)
	config := &MockConfig{
		Source:        absSource,
		Output:        outputDir,
		KeepResources: resources,
	}

	// Write terraform-mock.hcl
	if err := writeHCLConfig(config); err != nil {
		return fmt.Errorf("failed to write HCL config: %w", err)
	}

	fmt.Printf("Created terraform-mock.hcl\n")
	fmt.Printf("Edit the file to customize which resources to keep, then run: tfmock generate\n")

	return nil
}
