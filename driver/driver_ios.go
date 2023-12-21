// Copyright 2023 The Goki Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build ios && !offscreen

package driver

import "C"

import (
	"goki.dev/goosi"
	"goki.dev/goosi/driver/ios"
)

func driverMain(f func(goosi.App)) {
	ios.Main(f)
}
