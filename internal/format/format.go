/*
 * Copyright 2025 by Crossplaner authors
 *
 * This program is a free software product. You can redistribute it and/or
 * modify it under the terms of the GNU Affero General Public License (AGPL)
 * version 3 as published by the Free Software Foundation.
 *
 * For details, see the GNU AGPL at: http://www.gnu.org/licenses/agpl-3.0.html
 */

package format

//go:generate go tool stringer -type=Format
type Format int

// GOLANG 1.24 - Feat 2: Generic type aliases
// type Format = int
// var nF int = 0
// var ww Format = nF

const (
	Terraform Format = iota
	Crossplane
)
