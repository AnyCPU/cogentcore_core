// Copyright 2023 The Goki Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build js && !offscreen

package driver

import (
	"goki.dev/goosi"
	"goki.dev/goosi/driver/web"
)

func driverMain(f func(goosi.App)) {
	web.Main(f)
}
