/*
 * Copyright 2025 by Crossplaner authors
 *
 * This program is a free software product. You can redistribute it and/or
 * modify it under the terms of the GNU Affero General Public License (AGPL)
 * version 3 as published by the Free Software Foundation.
 *
 * For details, see the GNU AGPL at: http://www.gnu.org/licenses/agpl-3.0.html
 */

package tf

import (
	internalConverter "github.com/patrykorwat/infraconv/internal/converter"
	"github.com/patrykorwat/infraconv/internal/format"
)

type Command struct {
	Convert startCommand `cmd:"" help:"Convert a Terraform infrastructure definition"`
}

type startCommand struct{}

func (c *startCommand) Run() error {
	converter := internalConverter.NewConverter(format.Terraform, format.Crossplane)
	return converter.Convert()
}
