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

import "github.com/patrykorwat/crossplaner/internal/format"

type converter struct {
	source format.Format
	target format.Format
}

func (c converter) Convert() error {
	//TODO implement me
	panic("implement me")
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
