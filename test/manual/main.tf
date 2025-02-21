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
