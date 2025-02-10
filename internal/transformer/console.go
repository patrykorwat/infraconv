package transformer

import (
	"fmt"
	"github.com/patrykorwat/infraconv/internal/parser"
)

type consoleTransformer struct {
}

func (c consoleTransformer) Transform(config *parser.Config, directoryOutput string) error {
	// Print resources
	fmt.Println("Resources:")
	for _, resource := range config.Resources {
		fmt.Printf("Type: %s, Name: %s\n", resource.Type, resource.Name)
		for key, value := range resource.Attributes {
			fmt.Printf("    %s: %v\n", key, value)
		}
	}

	// Print modules
	fmt.Println("\nModules:")
	for _, module := range config.Modules {
		fmt.Printf("Source: %s\n", module.Source)
		for key, value := range module.Attributes {
			fmt.Printf("    %s: %v\n", key, value)
		}
	}

	// Print variables
	fmt.Println("\nVariables:")
	for _, variable := range config.Variables {
		for key, value := range variable.Attributes {
			fmt.Printf("    %s: %v\n", key, value)
		}
	}

	// Print locals
	fmt.Println("\nLocals:")
	for _, local := range config.Locals {
		for key, value := range local.Attributes {
			fmt.Printf("    %s: %v\n", key, value)
		}
	}

	return nil
}

func NewConsoleTransformer() Transformer {
	return &consoleTransformer{}
}
