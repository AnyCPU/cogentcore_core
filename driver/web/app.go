// Copyright 2023 The Goki Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build js

// Package web implements goosi interfaces on the web through WASM
package web

import (
	"image"
	"strings"
	"syscall/js"

	"goki.dev/goosi"
	"goki.dev/goosi/clip"
	"goki.dev/goosi/cursor"
	"goki.dev/goosi/driver/base"
	"goki.dev/goosi/events"
	"goki.dev/goosi/events/key"
	"goki.dev/jsfs"
)

// TheApp is the single [goosi.App] for the web platform
var TheApp = &App{AppSingle: base.NewAppSingle[*Drawer, *Window]()}

// App is the [goosi.App] implementation for the web platform
type App struct { //gti:add
	base.AppSingle[*Drawer, *Window]

	// SystemPlatform is the underlying system SystemPlatform (Android, iOS, etc)
	SystemPlatform goosi.Platforms

	// KeyMods are the current key mods
	KeyMods key.Modifiers
}

func init() {
	jsfs.Config(js.Global().Get("fs"))
}

// Main is called from main thread when it is time to start running the
// main loop. When function f returns, the app ends automatically.
func Main(f func(goosi.App)) {
	TheApp.Drawer = &Drawer{}
	base.Main(f, TheApp, &TheApp.App)
}

// NewWindow creates a new window with the given options.
// It waits for the underlying system window to be created first.
// Also, it hides all other windows and shows the new one.
func (a *App) NewWindow(opts *goosi.NewWindowOptions) (goosi.Window, error) {
	defer func() { goosi.HandleRecover(recover()) }()

	if goosi.InitScreenLogicalDPIFunc != nil {
		goosi.InitScreenLogicalDPIFunc()
	}
	a.Win = &Window{base.NewWindowSingle(a, opts)}
	a.Win.This = a.Win
	a.SetSystemWindow()

	go a.Win.WinLoop()

	return a.Win, nil
}

// SetSystemWindow sets the underlying system window information.
func (a *App) SetSystemWindow() {
	defer func() { goosi.HandleRecover(recover()) }()

	a.AddEventListeners()

	ua := js.Global().Get("navigator").Get("userAgent").String()
	lua := strings.ToLower(ua)
	if strings.Contains(lua, "android") {
		a.SystemPlatform = goosi.Android
	} else if strings.Contains(lua, "ipad") || strings.Contains(lua, "iphone") || strings.Contains(lua, "ipod") {
		a.SystemPlatform = goosi.IOS
	} else {
		// TODO(kai/web): more specific desktop platform
		a.SystemPlatform = goosi.Windows
	}

	a.Resize()
	a.Win.EvMgr.Window(events.WinShow)
	a.Win.EvMgr.Window(events.ScreenUpdate)
	a.Win.EvMgr.Window(events.WinFocus)
}

// Resize updates the app sizing information and sends a Resize event.
func (a *App) Resize() {
	a.Scrn.DevicePixelRatio = float32(js.Global().Get("devicePixelRatio").Float())
	dpi := 160 * a.Scrn.DevicePixelRatio
	a.Scrn.PhysicalDPI = dpi
	a.Scrn.LogicalDPI = dpi

	w, h := js.Global().Get("innerWidth").Int(), js.Global().Get("innerHeight").Int()
	sz := image.Pt(w, h)
	a.Scrn.Geometry.Max = sz
	a.Scrn.PixSize = image.Pt(int(float32(sz.X)*a.Scrn.DevicePixelRatio), int(float32(sz.Y)*a.Scrn.DevicePixelRatio))
	physX := 25.4 * float32(w) / dpi
	physY := 25.4 * float32(h) / dpi
	a.Scrn.PhysicalSize = image.Pt(int(physX), int(physY))

	canvas := js.Global().Get("document").Call("getElementById", "app")
	canvas.Set("width", a.Scrn.PixSize.X)
	canvas.Set("height", a.Scrn.PixSize.Y)

	a.Drawer.Image = image.NewRGBA(image.Rectangle{Max: a.Scrn.PixSize})

	a.Win.EvMgr.WindowResize()
}

func (a *App) DataDir() string {
	// TODO(kai): implement web filesystem
	return "/data"
}

func (a *App) Platform() goosi.Platforms {
	return goosi.Web
}

func (a *App) OpenURL(url string) {
	js.Global().Call("open", url)
}

func (a *App) ClipBoard(win goosi.Window) clip.Board {
	return TheClip
}

func (a *App) Cursor(win goosi.Window) cursor.Cursor {
	return TheCursor
}

func (a *App) IsDark() bool {
	return js.Global().Get("matchMedia").Truthy() &&
		js.Global().Call("matchMedia", "(prefers-color-scheme: dark)").Get("matches").Truthy()
}

func (a *App) ShowVirtualKeyboard(typ goosi.VirtualKeyboardTypes) {
	js.Global().Get("document").Call("getElementById", "text-field").Call("focus")
}

func (a *App) HideVirtualKeyboard() {
	js.Global().Get("document").Call("getElementById", "text-field").Call("blur")
}
