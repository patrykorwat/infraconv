package parser

// Provider represents a Terraform provider.
type Provider struct {
	Name       string
	Attributes map[string]any
}

// Resource represents a Terraform resource.
type Resource struct {
	Type       string
	Name       string
	Labels     []string
	Attributes map[string]any
}

// Module represents a Terraform module.
type Module struct {
	Source     string
	Labels     []string
	Attributes map[string]any
}

// Local represents a Terraform local value.
type Local struct {
	Attributes map[string]any
}

// Variable represents a Terraform variable.
type Variable struct {
	Attributes map[string]any
}

// Config represents the Terraform configuration.
type Config struct {
	RequiredProviders []*string
	Providers         []*Provider
	Resources         []*Resource
	Modules           []*Module
	Variables         []*Variable
	Locals            []*Local
}

type Parser interface {
	Parse(directory string) (*Config, error)
}
