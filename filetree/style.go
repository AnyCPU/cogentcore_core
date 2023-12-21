// Copyright (c) 2023, The Goki Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package filetree

import (
	"goki.dev/colors"
	"goki.dev/gi/v2/gi"
	"goki.dev/girl/styles"
	"goki.dev/goosi/events"
	"goki.dev/grr"
	"goki.dev/icons"
	"goki.dev/vci/v2"
)

func (ft *Tree) OnInit() {
	ft.Node.OnInit()
	ft.FRoot = ft
	ft.NodeType = NodeType
	ft.OpenDepth = 4
}

func (fn *Node) OnInit() {
	fn.TreeView.OnInit()
	fn.HandleEvents()
	fn.SetStyles()
}

func (fn *Node) SetStyles() {
	fn.Style(func(s *styles.Style) {
		vcs := fn.Info.Vcs
		s.Font.Weight = styles.WeightNormal
		s.Font.Style = styles.FontNormal
		if fn.IsExec() && !fn.IsDir() {
			s.Font.Weight = styles.WeightBold // todo: somehow not working
		}
		if fn.Buf != nil {
			s.Font.Style = styles.FontItalic
		}
		switch {
		case vcs == vci.Untracked:
			s.Color = grr.Must1(colors.FromHex("#808080"))
		case vcs == vci.Modified:
			s.Color = grr.Must1(colors.FromHex("#4b7fd1"))
		case vcs == vci.Added:
			s.Color = grr.Must1(colors.FromHex("#008800"))
		case vcs == vci.Deleted:
			s.Color = grr.Must1(colors.FromHex("#ff4252"))
		case vcs == vci.Conflicted:
			s.Color = grr.Must1(colors.FromHex("#ce8020"))
		case vcs == vci.Updated:
			s.Color = grr.Must1(colors.FromHex("#008060"))
		case vcs == vci.Stored:
			s.Color = colors.Scheme.OnSurface
		}
	})
	fn.OnWidgetAdded(func(w gi.Widget) {
		switch w.PathFrom(fn) {
		case "parts":
			parts := w.(*gi.Layout)
			w.OnClick(func(e events.Event) {
				fn.OpenEmptyDir()
			})
			parts.OnDoubleClick(func(e events.Event) {
				if fn.OpenEmptyDir() {
					e.SetHandled()
				}
			})
		case "parts/branch":
			sw := w.(*gi.Switch)
			sw.Type = gi.SwitchCheckbox
			sw.SetIcons(icons.FolderOpen, icons.Folder, icons.Blank)
		}
	})
}
