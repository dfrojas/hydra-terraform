# Terraform for LocalStack

# Context

When I tried to implement LocalStack into the internal workflow development of the company where I work, I found that our Terraform resources relies heavily of Terraform modules. The source of these Terraform modules are in our private GitHub repositories. Although we could clone them manually or add an authentication to GitHub via CLI, it comes to my mind an idea to build an Open Source Project / Side Project where we could copy the source of a Terraform module without the need to re-write it in the current Terraform project, it means to be able to re-use the module without the need of have access to the source of the module. Also, in the cases where our real TF resources are not included as part of the premium version of Localstack.

# Current closest approaches

**override_module:** It is a native feature of Terraform but only works for “terraform tests” and for outputs. Take into account that we want to override the arguments, not the outputs.

**Terragrunt:** This tools is ideal for override Terraform resources and re-usem but it has not a native way to override a module. The module has to be re-written if we want to override it.

# Implementation Details

The tool will contain three types of filterings:

**keep_resources:** Whitelist approach, only specified resource types are included. Example: `["aws_s3_bucket", "aws_iam_role"]`


**skip_resources:** Blacklist approach specified resource types are excluded. Example: `["aws_cloudwatch_log_group"]`

**remove_attributes:** Block-level filtering, remove specific attributes within resources. Example: Remove `vpc_config` from Lambda functions

## LocalStack Configuration

Generated `provider.tf` automatically configures:

- AWS provider with LocalStack endpoints
- All services point to http://localstack:4566
- Credentials set to test/test
- Metadata checks disabled


# Design

Having this simple example:

```jsx
main.tf
-------

module "simple" {
  source = "./modules/simple"
  bucket_name = "example"
}

output "bucket_id" {
  value = module.aws_s3_bucket._.id
}

modules/simple/main.tf
----------------------
resource "aws_s3_bucket" "_" {
  bucket        = var.bucket_name
}

resource "aws_iam_role" "_" {
   ...
}

resource "aws_lambda_function" "_" {
   vpc_config = ...
   layers = ...
}
```

```jsx
hydratf init \
--source path-to-tf-code-in-file-system
--name simple-lambda
```

This outputs:

```jsx
# terraform-localstack.hcl
transform {
  source = "path-to-tf-code-in-file-system" // Path to the module
  output = "./localstack" // Where it should write the final source
  
  keep_resources = ["aws_s3_bucket", "aws_iam_role"]
  
  remove_attributes = {
    "aws_lambda_function" = ["vpc_config", "layers"]
  }
  
  skip_resources = ["aws_cloudwatch_log_group", "aws_lambda_layer_version"]
}
```

User can modify the generated HCL configuration file if needed. The init-mock should creates an initial skeleton only.

```jsx
hydratf generate
```

## Results

New TF files only with the resources set from the HCL file.  The outputs of the generated `hydratf` tool should point all the endpoints to 4566 which is the port of Localstack and with test default credentials.

# Tech Stack

Go

# Future Enhancements

- Smart Filtering: Analyze dependencies and warn about broken references
- Templates: Pre-defined configurations for common scenarios
- Validation: Verify generated modules before writing
- Dry Run: Preview changes before generating files
- Interactive Mode: Guide users through configuration
- Remote Sources: Direct GitHub/Registry integration
