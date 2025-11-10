package parser

import (
	"github.com/hashicorp/hcl/v2/hclwrite"
)

type TerraformModule struct {
	Resources []*Resource
	Variables []*Variable
	Outputs   []*Output
	Locals    []*Locals
}

type Resource struct {
	Type       string
	Name       string
	Body       *hclwrite.Block
	Attributes map[string]interface{}
}

type Variable struct {
	Name string
	Body *hclwrite.Block
}

type Output struct {
	Name       string
	Body       *hclwrite.Block
	Attributes map[string]interface{}
}

type Locals struct {
	Body *hclwrite.Block
}
