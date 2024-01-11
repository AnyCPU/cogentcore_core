// Copyright (c) 2023, The Goki Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"goki.dev/goki/grease"
	"goki.dev/goki/gti/gtigen"
)

func main() {
	opts := grease.DefaultOptions("gtigen", "GTIGen", "GTIGen provides the generation of general purpose type information for Go types, methods, functions and variables")
	grease.Run(opts, &gtigen.Config{}, gtigen.Generate)
}
