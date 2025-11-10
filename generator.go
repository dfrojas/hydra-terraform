package main

import (
	//"fmt"
	"os"
	"path/filepath"

	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/zclconf/go-cty/cty"
)

func filterResources(module *TerraformModule, keepResources []string) *TerraformModule {
	if len(keepResources) == 0 {
		return module
	}

	// Create a set for faster lookup
	keepSet := make(map[string]bool)
	for _, r := range keepResources {
		keepSet[r] = true
	}

	filtered := &TerraformModule{
		Resources: []*Resource{},
		Variables: module.Variables,
		Outputs:   module.Outputs,
	}

	for _, resource := range module.Resources {
		if keepSet[resource.Type] {
			filtered.Resources = append(filtered.Resources, resource)
		}
	}

	return filtered
}

func generateTerraformFiles(outputDir string, module *TerraformModule) error {
	// Generate main.tf with resources
	if len(module.Resources) > 0 {
		if err := generateMainTF(outputDir, module.Resources, module.Locals); err != nil {
			return err
		}
	}

	// Generate variables.tf
	if len(module.Variables) > 0 {
		if err := generateVariablesTF(outputDir, module.Variables); err != nil {
			return err
		}
	}

	// Generate outputs.tf
	if len(module.Outputs) > 0 {
		if err := generateOutputsTF(outputDir, module.Outputs); err != nil {
			return err
		}
	}

	return nil
}

func generateMainTF(outputDir string, resources []*Resource, locals []*Locals) error {
	f := hclwrite.NewEmptyFile()
	rootBody := f.Body()

	for _, locals := range locals {
		// Copy the block
		newBlock := rootBody.AppendNewBlock("locals", []string{})
		newBody := newBlock.Body()

		// Copy attributes from original block
		for name, attr := range locals.Body.Body().Attributes() {
			newBody.SetAttributeRaw(name, attr.Expr().BuildTokens(nil))
		}

		rootBody.AppendNewline()
	}

	for _, resource := range resources {
		// Copy the block
		newBlock := rootBody.AppendNewBlock("resource", []string{resource.Type, resource.Name})
		newBody := newBlock.Body()

		// Copy attributes from original block
		for name, attr := range resource.Body.Body().Attributes() {
			newBody.SetAttributeRaw(name, attr.Expr().BuildTokens(nil))
		}

		// Copy nested blocks
		for _, block := range resource.Body.Body().Blocks() {
			copyBlock(newBody, block)
		}

		rootBody.AppendNewline()
	}

	return os.WriteFile(filepath.Join(outputDir, "main.tf"), f.Bytes(), 0644)
}

func generateVariablesTF(outputDir string, variables []*Variable) error {
	f := hclwrite.NewEmptyFile()
	rootBody := f.Body()

	for _, variable := range variables {
		newBlock := rootBody.AppendNewBlock("variable", []string{variable.Name})
		newBody := newBlock.Body()

		// Copy attributes from original block
		for name, attr := range variable.Body.Body().Attributes() {
			newBody.SetAttributeRaw(name, attr.Expr().BuildTokens(nil))
		}

		rootBody.AppendNewline()
	}

	return os.WriteFile(filepath.Join(outputDir, "variables.tf"), f.Bytes(), 0644)
}

func generateOutputsTF(outputDir string, outputs []*Output) error {
	f := hclwrite.NewEmptyFile()
	rootBody := f.Body()

	for _, output := range outputs {
		newBlock := rootBody.AppendNewBlock("output", []string{output.Name})
		newBody := newBlock.Body()

		// Copy attributes from original block
		for name, attr := range output.Body.Body().Attributes() {
			newBody.SetAttributeRaw(name, attr.Expr().BuildTokens(nil))
		}

		rootBody.AppendNewline()
	}

	return os.WriteFile(filepath.Join(outputDir, "outputs.tf"), f.Bytes(), 0644)
}

func copyBlock(parent *hclwrite.Body, block *hclwrite.Block) {
	newBlock := parent.AppendNewBlock(block.Type(), block.Labels())
	newBody := newBlock.Body()

	// Copy attributes
	for name, attr := range block.Body().Attributes() {
		newBody.SetAttributeRaw(name, attr.Expr().BuildTokens(nil))
	}

	// Recursively copy nested blocks
	for _, nestedBlock := range block.Body().Blocks() {
		copyBlock(newBody, nestedBlock)
	}
}

func generateLocalStackProvider(outputDir string) error {
	f := hclwrite.NewEmptyFile()
	rootBody := f.Body()

	// Add terraform block
	terraformBlock := rootBody.AppendNewBlock("terraform", nil)
	terraformBody := terraformBlock.Body()

	requiredProvidersBlock := terraformBody.AppendNewBlock("required_providers", nil)
	providersBody := requiredProvidersBlock.Body()

	// Add aws provider requirement
	awsBlock := providersBody.AppendNewBlock("aws", nil)
	awsBody := awsBlock.Body()
	awsBody.SetAttributeValue("source", cty.StringVal("hashicorp/aws"))
	awsBody.SetAttributeValue("version", cty.StringVal("~> 5.0")) // TODO: Could comes from a YAML config.

	rootBody.AppendNewline()

	// Add provider configuration for LocalStack
	providerBlock := rootBody.AppendNewBlock("provider", []string{"aws"})
	providerBody := providerBlock.Body()

	providerBody.SetAttributeValue("region", cty.StringVal("us-east-1"))
	providerBody.SetAttributeValue("access_key", cty.StringVal("test"))
	providerBody.SetAttributeValue("secret_key", cty.StringVal("test"))
	providerBody.SetAttributeValue("skip_credentials_validation", cty.True)
	providerBody.SetAttributeValue("skip_metadata_api_check", cty.True)
	providerBody.SetAttributeValue("skip_requesting_account_id", cty.True)

	// Add endpoints block
	endpointsBlock := providerBody.AppendNewBlock("endpoints", nil)
	endpointsBody := endpointsBlock.Body()

	// Common AWS services that LocalStack supports
	services := []string{
		"s3", "lambda", "dynamodb", "apigateway", "sqs", "sns",
		"cloudformation", "cloudwatch", "iam", "sts", "ec2",
		"kinesis", "kms", "secretsmanager", "ssm",
	}

	for _, service := range services {
		endpointsBody.SetAttributeValue(service, cty.StringVal("http://localstack:4566"))
	}

	return os.WriteFile(filepath.Join(outputDir, "provider.tf"), f.Bytes(), 0644)
}
