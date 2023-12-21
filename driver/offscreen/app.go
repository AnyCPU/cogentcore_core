// Copyright 2023 The Goki Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package offscreen provides placeholder implementations of goosi interfaces
// to allow for offscreen testing and capturing of apps.
package offscreen

import (
	"image"
	"path/filepath"

	"goki.dev/goosi"
	"goki.dev/goosi/driver/base"
	"goki.dev/goosi/events"
)

// TheApp is the single [goosi.App] for the offscreen platform
var TheApp = &App{base.NewAppSingle[*Drawer, *Window]()}

// App is the [goosi.App] implementation for the offscreen platform
type App struct { //gti:add
	base.AppSingle[*Drawer, *Window]
}

// Main is called from main thread when it is time to start running the
// main loop. When function f returns, the app ends automatically.
func Main(f func(goosi.App)) {
	TheApp.Drawer = &Drawer{}
	TheApp.GetScreens()
	base.Main(f, TheApp, &TheApp.App)
}

// NewWindow creates a new window with the given options.
// It waits for the underlying system window to be created first.
// Also, it hides all other windows and shows the new one.
func (a *App) NewWindow(opts *goosi.NewWindowOptions) (goosi.Window, error) {
	defer func() { base.HandleRecover(recover()) }()

	if goosi.InitScreenLogicalDPIFunc != nil {
		goosi.InitScreenLogicalDPIFunc()
	}
	a.Win = &Window{base.NewWindowSingle(a, opts)}
	a.Win.This = a.Win
	a.Scrn.PixSize = opts.Size
	a.GetScreens()

	a.Win.EvMgr.WindowResize()
	a.Win.EvMgr.Window(events.WinShow)
	a.Win.EvMgr.Window(events.ScreenUpdate)
	a.Win.EvMgr.Window(events.WinFocus)

	go a.Win.WinLoop()

	return a.Win, nil
}

func (a *App) GetScreens() {
	if a.Scrn.PixSize.X == 0 {
		a.Scrn.PixSize.X = 800
	}
	if a.Scrn.PixSize.Y == 0 {
		a.Scrn.PixSize.Y = 600
	}

	a.Scrn.DevicePixelRatio = 1
	a.Scrn.Geometry.Max = a.Scrn.PixSize
	dpi := float32(160)
	a.Scrn.PhysicalDPI = dpi
	a.Scrn.LogicalDPI = dpi

	physX := 25.4 * float32(a.Scrn.PixSize.X) / dpi
	physY := 25.4 * float32(a.Scrn.PixSize.Y) / dpi
	a.Scrn.PhysicalSize = image.Pt(int(physX), int(physY))

	a.Drawer.Image = image.NewRGBA(image.Rectangle{Max: a.Scrn.PixSize})
}

func (a *App) DataDir() string {
	// TODO(kai): figure out a better solution to offscreen prefs dir
	return filepath.Join(".", "tmpDataDir")
}

func (a *App) Platform() goosi.Platforms {
	return goosi.Offscreen
}
