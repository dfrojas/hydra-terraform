package parser

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclwrite"
)

func ParseResources(modulePath string) ([]string, error) {
	resources := make(map[string]bool)

	// Find all .tf files
	err := filepath.Walk(modulePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && strings.HasSuffix(path, ".tf") {
			content, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			hclFile, diags := hclwrite.ParseConfig(content, path, hcl.Pos{Line: 1, Column: 1})
			if diags.HasErrors() {
				return fmt.Errorf("failed to parse %s: %s", path, diags.Error())
			}

			body := hclFile.Body()

			for _, block := range body.Blocks() {
				if block.Type() == "resource" {
					labels := block.Labels()
					if len(labels) >= 1 {
						resources[labels[0]] = true
					}
				}
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	result := make([]string, 0, len(resources))
	for resourceType := range resources {
		result = append(result, resourceType)
	}

	return result, nil
}

func ParseModuleDetailed(modulePath string) (*TerraformModule, error) {
	module := &TerraformModule{
		Resources: []*Resource{},
		Variables: []*Variable{},
		Outputs:   []*Output{},
		Locals:    []*Locals{},
	}

	err := filepath.Walk(modulePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && strings.HasSuffix(path, ".tf") {
			content, err := os.ReadFile(path)
			if err != nil {
				return err
			}

			hclFile, diags := hclwrite.ParseConfig(content, path, hcl.Pos{Line: 1, Column: 1})
			if diags.HasErrors() {
				return fmt.Errorf("failed to parse %s: %s", path, diags.Error())
			}

			body := hclFile.Body()

			for _, block := range body.Blocks() {
				switch block.Type() {
				case "resource":
					labels := block.Labels()
					/* labels are the name of the resource and the tag that we assign to it:
					resource "aws_iam_role" "lambda_role" {} where "aws_iam_role" is label[0]
					and lambda_role is label[1] */
					if len(labels) >= 2 {
						module.Resources = append(module.Resources, &Resource{
							Type: labels[0],
							Name: labels[1],
							Body: block,
						})
					}
				case "variable":
					labels := block.Labels()
					if len(labels) >= 1 {
						module.Variables = append(module.Variables, &Variable{
							Name: labels[0],
							Body: block,
						})
					}
				case "output":
					labels := block.Labels()
					if len(labels) >= 1 {
						module.Outputs = append(module.Outputs, &Output{
							Name: labels[0],
							Body: block,
						})
					}
				case "locals":
					module.Locals = append(module.Locals, &Locals{
						Body: block,
					})
				}
			}
		}

		return nil
	})

	return module, err
}
