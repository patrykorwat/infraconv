/*
 * Copyright 2025 by Crossplaner authors
 *
 * This program is a free software product. You can redistribute it and/or
 * modify it under the terms of the GNU Affero General Public License (AGPL)
 * version 3 as published by the Free Software Foundation.
 *
 * For details, see the GNU AGPL at: http://www.gnu.org/licenses/agpl-3.0.html
 */

package main

import (
	"github.com/alecthomas/kong"

	"github.com/patrykorwat/infraconv/cmd/infraconv/tf"
)

type cli struct {
	Tf tf.Command `cmd:"" help:"Run a job to convert infrastructure definitions from Terrraform."`
}

func main() {
	ctx := kong.Parse(&cli{},
		kong.Name("crossplaner"),
		kong.Description("A handy tool to convert infra definitions"),
		kong.UsageOnError(),
	)
	ctx.FatalIfErrorf(ctx.Run())
}
