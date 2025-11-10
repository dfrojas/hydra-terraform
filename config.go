package main

import (
	"fmt"
	"os"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/zclconf/go-cty/cty"
)

type MockConfig struct {
	Source        string
	Output        string
	KeepResources []string
}

func writeHCLConfig(config *MockConfig) error {
	f := hclwrite.NewEmptyFile()
	rootBody := f.Body()

	// Create transform block
	transformBlock := rootBody.AppendNewBlock("transform", nil)
	transformBody := transformBlock.Body()

	// Add source
	transformBody.SetAttributeValue("source", cty.StringVal(config.Source))

	// Add output
	transformBody.SetAttributeValue("output", cty.StringVal(config.Output))

	// Add keep_resources as a list
	if len(config.KeepResources) > 0 {
		resourceValues := make([]cty.Value, len(config.KeepResources))
		for i, r := range config.KeepResources {
			resourceValues[i] = cty.StringVal(r)
		}
		transformBody.SetAttributeValue("keep_resources", cty.ListVal(resourceValues))
	} else {
		transformBody.SetAttributeValue("keep_resources", cty.ListValEmpty(cty.String))
	}

	// Write to file
	return os.WriteFile("terraform-mock.hcl", f.Bytes(), 0644)
}

func readHCLConfig(filename string) (*MockConfig, error) {
	parser := hclparse.NewParser()

	file, diags := parser.ParseHCLFile(filename)
	if diags.HasErrors() {
		return nil, fmt.Errorf("failed to parse HCL: %s", diags.Error())
	}

	config := &MockConfig{}

	// Parse the transform block
	content, _, diags := file.Body.PartialContent(&hcl.BodySchema{
		Blocks: []hcl.BlockHeaderSchema{
			{
				Type: "transform",
			},
		},
	})

	if diags.HasErrors() {
		return nil, fmt.Errorf("failed to parse content: %s", diags.Error())
	}

	if len(content.Blocks) == 0 {
		return nil, fmt.Errorf("no transform block found")
	}

	transformBlock := content.Blocks[0]
	attrs, diags := transformBlock.Body.JustAttributes()
	if diags.HasErrors() {
		return nil, fmt.Errorf("failed to parse attributes: %s", diags.Error())
	}

	// Extract source
	if sourceAttr, ok := attrs["source"]; ok {
		val, diags := sourceAttr.Expr.Value(nil)
		if !diags.HasErrors() {
			config.Source = val.AsString()
		}
	}

	// Extract output
	if outputAttr, ok := attrs["output"]; ok {
		val, diags := outputAttr.Expr.Value(nil)
		if !diags.HasErrors() {
			config.Output = val.AsString()
		}
	}

	// Extract keep_resources
	if keepAttr, ok := attrs["keep_resources"]; ok {
		val, diags := keepAttr.Expr.Value(nil)
		if !diags.HasErrors() && val.Type().IsListType() {
			config.KeepResources = make([]string, 0)
			for _, v := range val.AsValueSlice() {
				config.KeepResources = append(config.KeepResources, v.AsString())
			}
		}
	}

	return config, nil
}
