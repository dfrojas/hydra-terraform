<h1>
<p align="center">
  <img src="documentation/img/whale-7.png" width="150">
  <p align="center">Hydra Terraform</p>
</h1>
  <p align="center">
    Regenerate your Terraform code for LocalStack.
    <br /><br />
    <a href="#about">About</a>
    ·
    <a href="#documentation">Documentation</a>
    ·
    <a href="#roadmap">Roadmap</a>
  </p>
</p>

## About

Hydra Terraform is a CLI tool that regenerates Terraform code locally for use with LocalStack. It analyzes existing Terraform configurations, filters out resources that don't work well locally, and generates simplified versions pre-configured for LocalStack, keeping only what you need for fast local development.

The primary goal is to make LocalStack development practical by solving a common problem: many Terraform resources work fine in AWS but are problematic in LocalStack. Lambda functions with VPC configurations, complex networking setups, or layer dependencies which is not available in the Community Edition. Hydra lets you keep the core resources (S3, DynamoDB, IAM, etc) while removing the problematic ones, giving you a clean local development environment.

Hydra is written in Go and uses  <a href="https://pkg.go.dev/github.com/hashicorp/hcl/v2/hclparse" target="_blank">HashiCorp's official HCL parser</a> to ensure compatibility with Terraform syntax.

> [!IMPORTANT]
>
> This is a personal/side project that I'm building only for fun and learn in my free time.
> Please use carefully.


> [!NOTE]
>
> This project was developed from a design document outlining the original vision and architecture. Most of the initial codebase was generated with LLM assistance (Claude) as a way to learn how AI can help accelerate development from concept to implementation.  
>  
> You can find the original design doc <a href="desing-doc.md">here</a>.


### Use Cases

**Complex Lambda with VPC and Layers**: LocalStack doesn't support Lambda layers in the CE. Keep only S3 and IAM, skip Lambda.

**Private Module Sources**: Clone once locally, point Hydra to the local path, never need GitHub auth again.

**Simplified Testing**: Test core infrastructure (S3, DynamoDB) without dealing with complex networking (VPC, subnets).

## Installation & Usage

Build from source:
```bash
git clone https://github.com/dfrojas/hydra-terraform
cd hydra-terraform
make build
```

### Commands

**`hydratf init`**

Analyzes a Terraform module and creates a configuration file.

```bash
hydratf init --source <path-to-module> --name <output-name>
```

**Parameters:**
- `--source` - Path to your existing Terraform module
- `--name` - Name for the output directory

**Output:** Creates `terraform-localstack.hcl` with all discovered resources

**`hydratf generate`**

Generates adapted Terraform files based on your configuration. `geneate` command point automatically to the LocalStack default endpoints and default credentials.

```bash
hydratf generate
```

**Requirements:** Must have a `terraform-localstack.hcl` file in the current directory

**Output:** Creates adapted Terraform files in the configured output directory

### Configuration File

The `terraform-localstack.hcl` file controls what gets adapted:

```hcl
transform {
  source         = "/absolute/path/to/source/module"
  output         = "output-directory-name"
  keep_resources = ["resource_type_1", "resource_type_2"]
}
```

**Fields:**
- `source` - Absolute path to your source Terraform module
- `output` - Directory where adapted files will be generated
- `keep_resources` - Array of AWS resource types to include (e.g., `aws_s3_bucket`, `aws_iam_role`)
- `remove_blocks` - HashMap to remove attributes of a certain resource (e.g, `{"aws_lambda_function" = ["vpc_config" "layers"]}`)
- `skip_resources` - Array of AWS resources that we want to ignore.

## Roadmap

| Feature                  | Status   |
|--------------------------|----------|
| Generates HCL file       | ✅       |
| Support keep_resources   | ⏳       |
| Support remove_blocks    | ⏳       |
| Suport skip_resources    | ⏳       |


<!-- ### Complete Example

Let's adapt a complex Lambda module to work with LocalStack, keeping only S3 and IAM:

**1. Initialize from existing module**
```bash
hydratf init --source ./my-lambda-module --name localstack
```

**2. Edit `terraform-localstack.hcl`**
```hcl
transform {
  source         = "/Users/you/projects/my-lambda-module"
  output         = "localstack"
  keep_resources = ["aws_s3_bucket", "aws_iam_role"]
}
```

**3. Generate adapted module**
```bash
hydratf generate
```

**4. Deploy to LocalStack**
```bash
cd localstack
terraform init
terraform apply
```

Done! Your S3 and IAM resources are now running locally without Lambda/VPC complexity. -->