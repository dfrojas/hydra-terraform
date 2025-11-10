package cli

import (
	"github.com/spf13/cobra"
)

func Execute() error {
	rootCmd := &cobra.Command{
		Use:   "hydratf",
		Short: "Terraform regeneration tool for LocalStack",
		Long:  "A tool to regenerate Terraform modules for use with LocalStack",
	}

	rootCmd.AddCommand(InitCmd())
	rootCmd.AddCommand(GenerateCmd())

	return rootCmd.Execute()
}
