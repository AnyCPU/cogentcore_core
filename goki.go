// Copyright (c) 2023, The GoKi Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

//go:generate goki generate ./...

import (
	"goki.dev/goki/config"
	"goki.dev/goki/generate"
	"goki.dev/goki/packman"
	"goki.dev/goki/tools"
	"goki.dev/grease"
	// "goki.dev/greasi"
)

func main() {
	opts := grease.DefaultOptions("goki", "GoKi", "Command line and GUI tools for developing apps and libraries using the GoKi framework.")
	opts.DefaultFiles = []string{".goki/config.toml"}
	opts.SearchUp = true
	grease.Run(opts, &config.Config{}, packman.Build, packman.Install, packman.Run, generate.Generate, tools.Init, tools.Setup, packman.Log, packman.Release, packman.GetVersion, packman.SetVersion, packman.UpdateVersion)
}
