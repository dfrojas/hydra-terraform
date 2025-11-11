package config

import (
	"fmt"
	"os"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/zclconf/go-cty/cty"
)

func WriteHCLConfig(config *LocalstackConfig) error {
	f := hclwrite.NewEmptyFile()
	rootBody := f.Body()

	transformBlock := rootBody.AppendNewBlock("transform", nil)
	transformBody := transformBlock.Body()

	transformBody.SetAttributeValue("source", cty.StringVal(config.Source))

	transformBody.SetAttributeValue("output", cty.StringVal(config.Output))

	if len(config.KeepResources) > 0 {
		resourceValues := make([]cty.Value, len(config.KeepResources))
		for i, r := range config.KeepResources {
			resourceValues[i] = cty.StringVal(r)
		}
		transformBody.SetAttributeValue("keep_resources", cty.ListVal(resourceValues))
	} else {
		transformBody.SetAttributeValue("keep_resources", cty.ListValEmpty(cty.String))
	}

	transformBody.SetAttributeValue("remove_attributes", cty.MapValEmpty(cty.List(cty.String)))

	return os.WriteFile("terraform-hydra.hcl", f.Bytes(), 0644)
}

func ReadHCLConfig(filename string) (*LocalstackConfig, error) {
	parser := hclparse.NewParser()

	file, diags := parser.ParseHCLFile(filename)
	if diags.HasErrors() {
		return nil, fmt.Errorf("failed to parse HCL: %s", diags.Error())
	}

	config := &LocalstackConfig{}

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

	if sourceAttr, ok := attrs["source"]; ok {
		val, diags := sourceAttr.Expr.Value(nil)
		if !diags.HasErrors() {
			config.Source = val.AsString()
		}
	}

	if outputAttr, ok := attrs["output"]; ok {
		val, diags := outputAttr.Expr.Value(nil)
		if !diags.HasErrors() {
			config.Output = val.AsString()
		}
	}

	if keepAttr, ok := attrs["keep_resources"]; ok {
		val, diags := keepAttr.Expr.Value(nil)
		if !diags.HasErrors() && val.Type().IsListType() {
			config.KeepResources = make([]string, 0)
			for _, v := range val.AsValueSlice() {
				config.KeepResources = append(config.KeepResources, v.AsString())
			}
		}
	}

	if removeAttrs, ok := attrs["remove_attributes"]; ok {
		val, diags := removeAttrs.Expr.Value(nil)
		if !diags.HasErrors() && val.Type().IsObjectType() {
			config.RemoveAttributes = make(map[string][]string)
			for resource, attrList := range val.AsValueMap() {
				attributes := []string{}
				for _, attr := range attrList.AsValueSlice() {
					attributes = append(attributes, attr.AsString())
				}
				config.RemoveAttributes[resource] = attributes
			}
		}
	}

	return config, nil
}
