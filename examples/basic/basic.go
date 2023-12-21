// Copyright (c) 2023, The Goki Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"goki.dev/gi/v2/gi"
	"goki.dev/gi/v2/gimain"
	"goki.dev/grr"
	"goki.dev/video"
)

func main() { gimain.Run(app) }

func app() {
	sc := gi.NewScene("basic-video").SetTitle("Basic Video Example")
	v := video.NewVideo(sc)
	grr.Log(v.Open("../videos/deer.mp4"))
	w := gi.NewWindow(sc).Run()
	grr.Log(v.Play(0, 0))
	w.Wait()
}
