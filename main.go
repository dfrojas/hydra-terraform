package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func main() {
	//fmt.Println("Hello from tfmock!")
	rootCmd := &cobra.Command{
		Use:   "hydratf",
		Short: "Terraform regeneation tool for LocalStack",
		Long:  "A tool to regenerate Terraform modules for use with LocalStack",
	}

	rootCmd.AddCommand(initCmd())
	rootCmd.AddCommand(generateCmd())

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
