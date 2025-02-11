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
	"fmt"
	"github.com/alecthomas/kong"
	"runtime/debug"

	"github.com/patrykorwat/infraconv/cmd/infraconv/tf"
)

type (
	versionFlag bool
)

type cli struct {
	Version versionFlag `help:"Print version and quit." short:"v"`

	Tf tf.Command `cmd:"" help:"Run a job to convert infrastructure definitions from Terrraform."`
}

// GOLANG 1.24 - Feat 4: main module version
func (v versionFlag) BeforeApply(app *kong.Kong) error { //nolint:unparam // BeforeApply requires this signature.
	info, _ := debug.ReadBuildInfo()
	fmt.Println("Go version:", info.GoVersion)
	fmt.Println("App version:", info.Main.Version)
	app.Exit(0)
	return nil
}

func main() {
	ctx := kong.Parse(&cli{},
		kong.Name("infraconv"),
		kong.Description("A handy tool to convert infra definitions"),
		kong.UsageOnError(),
	)
	ctx.FatalIfErrorf(ctx.Run())
}
