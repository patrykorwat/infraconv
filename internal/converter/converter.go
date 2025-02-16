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
	"context"
	"github.com/patrykorwat/infraconv/internal/format"
	internalParser "github.com/patrykorwat/infraconv/internal/parser"
	transformer "github.com/patrykorwat/infraconv/internal/transformer"
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

	crossplane := transformer.NewCrossplaneTransformer()
	ctx := context.Background()
	err = crossplane.Transform(ctx, config, "test/manual/output")
	if err != nil {
		return errors.Wrap(err, "transform error")
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
