// Copyright (c) 2018, The Goki Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"goki.dev/gi/v2/gi"
	"goki.dev/gi/v2/gimain"
	"goki.dev/gi/v2/texteditor"
	"goki.dev/girl/styles"
)

var samplefile gi.FileName = "../demo/demo.go"

// var samplefile gi.FileName = "../../Makefile"

// var samplefile gi.FileName = "../../README.md"

func main() { gimain.Run(app) }

func app() {
	b := gi.NewAppBody("texteditor").SetTitle("GoGi texteditor.Editor Test")
	b.App().About = `This is a demo of the texteditor.Editor in the <b>GoGi</b> graphical interface system, within the <b>Goki</b> tree framework.  See <a href="https://github.com/goki">Goki on GitHub</a>`

	splt := gi.NewSplits(b, "split-view")
	splt.SetSplits(.5, .5)
	// these are all inherited so we can put them at the top "editor panel" level
	splt.Style(func(s *styles.Style) {
		s.Text.WhiteSpace = styles.WhiteSpacePreWrap
		s.Text.TabSize = 4
		s.Font.Family = string(gi.GeneralSettings.MonoFont)
	})

	txed1 := texteditor.NewEditor(splt, "texteditor-1")
	txed1.Style(func(s *styles.Style) {
		s.Min.X.Ch(20)
		s.Min.Y.Ch(10)
	})
	txed2 := texteditor.NewEditor(splt, "texteditor-2")
	txed2.Style(func(s *styles.Style) {
		s.Min.X.Ch(20)
		s.Min.Y.Ch(10)
	})

	txbuf := texteditor.NewBuf()
	txed1.SetBuf(txbuf)
	txed2.SetBuf(txbuf)

	// txbuf.Hi.Lang = "Markdown" // "Makefile" // "Go" // "Markdown"
	txbuf.Hi.Lang = "Go"
	txbuf.Open(samplefile)
	// pr := txbuf.Hi.PiLang.Parser()
	// giv.InspectorDialog(&pr.Lexer)

	b.NewWindow().Run().Wait()
}
