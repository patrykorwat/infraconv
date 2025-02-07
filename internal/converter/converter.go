/*
 * Copyright 2025 by Crossplaner authors
 *
 * This program is a free software product. You can redistribute it and/or
 * modify it under the terms of the GNU Affero General Public License (AGPL)
 * version 3 as published by the Free Software Foundation.
 *
 * For details, see the GNU AGPL at: http://www.gnu.org/licenses/agpl-3.0.html
 */

package converter

import (
	"fmt"
	"github.com/patrykorwat/infraconv/internal/format"
	internalParser "github.com/patrykorwat/infraconv/internal/parser"
	"github.com/pkg/errors"
)

type converter struct {
	source format.Format
	target format.Format
}

func (c converter) Convert() error {
	var parser internalParser.Parser
	switch c.source {
	case format.Terraform:
		parser = internalParser.NewTfParser()
	default:
		return errors.New("Unsupported source format: " + c.source.String())
	}
	config, err := parser.Parse("test/manual")
	if err != nil {
		return errors.Wrap(err, "parsing error")
	}

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

type Converter interface {
	Convert() error
}

func NewConverter(source, target format.Format) Converter {
	return converter{
		source: source,
		target: target,
	}
}
