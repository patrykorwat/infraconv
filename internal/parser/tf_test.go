package parser

import (
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/rs/zerolog/log"
	"reflect"
	"testing"
)

var hclParser = hclparse.NewParser()

func parseHcl(content string) *hcl.File {
	parseHCL, diagnostics := hclParser.ParseHCL([]byte(content), "file")
	if diagnostics.HasErrors() {
		log.Fatal().Err(diagnostics.Errs()[0])
	}
	return parseHCL
}

func Test_parseConfig(t *testing.T) {
	type args struct {
		file *hcl.File
	}
	tests := []struct {
		name string
		args args
		want *Config
	}{
		{
			name: "",
			args: args{
				file: parseHcl(`
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
    Name = "example"
  }
}
`),
			},
			want: &Config{
				RequiredProviders: nil,
				Providers:         nil,
				Resources:         nil,
				Modules:           nil,
				Variables:         nil,
				Locals:            nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := parseConfig(tt.args.file); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}
