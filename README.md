# infraconv
Convert between infrastructure definition frameworks.

Tool for now only works for AWS provider.

Supported source formats:
* Terraform/OpenTofu

Supported target formats:
* Crossplane 

Example input TF file
```terraform
terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
  }

  required_version = ">= 1.2.0"
}

provider "aws" {
  region = "us-west-2"
}

resource "aws_instance" "test_server" {
  ami           = "ami-830c94e3"
  instance_type = "t2.micro"

  tags = {
    Name = "infraconv-test-ec2-instance"
  }
}

resource "aws_s3_bucket" "test-bucket" {
  bucket = "infraconv-test-bucket"
}
```

Sample output of command `infraconv tf convert`:
```yaml
apiVersion: ec2.aws.upbound.io/v1beta2
kind: Instance
spec:
  forProvider:
    ami: ami-830c94e3
    instanceType: t2.micro
    region: us-west-2
    tags:
      Name: infraconv-test-ec2-instance
---
apiVersion: s3.aws.upbound.io/v1beta2
kind: Bucket
metadata:
  name: infraconv-test-bucket
spec:
  forProvider:
    region: us-west-2
```

Things to consider when passing TF files:
1. Modules are not supported
2. References aren't supported
3. Boolean fields can only contain boolean values,`true` or `false`, string values like`"true"` aren't supported
