package cli

import (
	"fmt"
	"os"
	"path/filepath"

	"hydratf/internal/config"
	"hydratf/internal/parser"

	"github.com/spf13/cobra"
)

func InitCmd() *cobra.Command {
	var source string
	var name string

	cmd := &cobra.Command{
		Use:   "init",
		Short: "Initialize configuration for a Terraform module",
		RunE: func(cmd *cobra.Command, args []string) error {
			if source == "" {
				return fmt.Errorf("--source is required")
			}
			if name == "" {
				return fmt.Errorf("--name is required")
			}

			return initTransform(source, name)
		},
	}

	cmd.Flags().StringVar(&source, "source", "", "Path to the original Terraform module")
	cmd.Flags().StringVar(&name, "name", "", "Name for the configuration")

	return cmd
}

func initTransform(source, outputDir string) error {
	absSource, err := filepath.Abs(source)
	if err != nil {
		return fmt.Errorf("failed to resolve source path: %w", err)
	}

	if _, err := os.Stat(absSource); os.IsNotExist(err) {
		return fmt.Errorf("source path does not exist: %s", absSource)
	}

	resources, err := parser.ParseModule(absSource)
	if err != nil {
		return fmt.Errorf("failed to parse module: %w", err)
	}

	cfg := &config.LocalstackConfig{
		Source:        absSource,
		Output:        outputDir,
		KeepResources: resources,
	}

	if err := config.WriteHCLConfig(cfg); err != nil {
		return fmt.Errorf("failed to write HCL config: %w", err)
	}

	fmt.Printf("Created terraform-localstack.hcl\n")
	fmt.Printf("Edit the file to customize which resources to keep, then run: hydratf generate\n")

	return nil
}
