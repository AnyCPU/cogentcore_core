// Copyright (c) 2018, The Goki Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package giv

import (
	"goki.dev/gi/v2/gi"
	"goki.dev/gi/v2/keyfun"
	"goki.dev/girl/styles"
	"goki.dev/goosi/events"
	"goki.dev/icons"
)

// TODO: make base simplified preferences view, improve organization of information, and maybe add titles

// PrefsView opens a view of user preferences
func PrefsView(pf *gi.GeneralSettingsData) {
	if gi.ActivateExistingMainWindow(pf) {
		return
	}
	d := gi.NewBody("gogi-prefs")
	d.SetTitle("GoGi Preferences")
	d.Sc.Data = pf
	d.AddAppBar(func(tb *gi.Toolbar) {
		NewFuncButton(tb, pf.UpdateAll).SetIcon(icons.Refresh)
		gi.NewSeparator(tb)
		NewFuncButton(tb, pf.LightMode)
		NewFuncButton(tb, pf.DarkMode)
		gi.NewSeparator(tb)
		sz := NewFuncButton(tb, pf.SaveZoom).SetIcon(icons.ZoomIn)
		sz.Args[0].SetValue(true)
		NewFuncButton(tb, pf.ScreenInfo).SetShowReturn(true).SetIcon(icons.Info)
		NewFuncButton(tb, pf.VersionInfo).SetShowReturn(true).SetIcon(icons.Info)
		gi.NewSeparator(tb)
		NewFuncButton(tb, pf.EditKeyMaps).SetIcon(icons.Keyboard)
		NewFuncButton(tb, pf.EditHiStyles).SetIcon(icons.InkHighlighter)
		NewFuncButton(tb, pf.EditDetailed).SetIcon(icons.Description)
		NewFuncButton(tb, pf.EditDebug).SetIcon(icons.BugReport)
		tb.AddOverflowMenu(func(m *gi.Scene) {
			NewFuncButton(m, pf.Open).SetKey(keyfun.Open)
			NewFuncButton(m, pf.Delete).SetConfirm(true)
			NewFuncButton(m, pf.DeleteSavedWindowGeoms).SetConfirm(true).SetIcon(icons.Delete)
			gi.NewSeparator(tb)
		})
	})
	sv := NewStructView(d)
	sv.SetStruct(pf)
	sv.OnChange(func(e events.Event) {
		pf.Changed = true
		tab := d.GetTopAppBar()
		if tab != nil {
			tab.UpdateBar()
		}
		pf.Apply()
		pf.Save()
		pf.UpdateAll()
	})
	d.NewWindow().Run()
}

// PrefsDetView opens a view of user detailed preferences
func PrefsDetView(pf *gi.PrefsDetailed) {
	if gi.ActivateExistingMainWindow(pf) {
		return
	}

	d := gi.NewBody("gogi-prefs-det").SetTitle("GoGi Detailed Preferences")

	sv := NewStructView(d, "sv")
	sv.SetStruct(pf)

	d.Sc.Data = pf

	d.AddAppBar(func(tb *gi.Toolbar) {
		NewFuncButton(tb, pf.Apply).SetIcon(icons.Refresh)
		gi.NewSeparator(tb)
		NewFuncButton(tb, pf.Save).SetKey(keyfun.Save).
			StyleFirst(func(s *styles.Style) { s.SetEnabled(pf.Changed) })
		tb.AddOverflowMenu(func(m *gi.Scene) {
			NewFuncButton(m, pf.Open).SetKey(keyfun.Open)
			gi.NewSeparator(tb)
		})
	})

	d.NewWindow().Run()
}

// PrefsDbgView opens a view of user debugging preferences
func PrefsDbgView(pf *gi.PrefsDebug) {
	if gi.ActivateExistingMainWindow(pf) {
		return
	}
	d := gi.NewBody("gogi-prefs-dbg")
	d.Title = "GoGi Debugging Preferences"

	sv := NewStructView(d, "sv")
	sv.SetStruct(pf)

	d.Sc.Data = pf

	d.AddAppBar(func(tb *gi.Toolbar) {
		NewFuncButton(tb, pf.Profile).SetIcon(icons.LabProfile)
	})

	d.NewWindow().Run()
}
