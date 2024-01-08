// Copyright 2023 The Goki Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Based on golang.org/x/exp/shiny:
// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package base

import (
	"image"
	"sync"
	"time"

	"goki.dev/girl/styles"
	"goki.dev/goosi"
	"goki.dev/goosi/events"
)

// Window contains the data and logic common to all implementations of [goosi.Window].
// A Window is associated with a corresponding [goosi.App] type.
type Window[A goosi.App] struct { //gti:add

	// This is the Window as a [goosi.Window] interface, which preserves the actual identity
	// of the window when calling interface methods in the base Window.
	This goosi.Window `view:"-"`

	// App is the [goosi.App] associated with the window.
	App A

	// Mu is the main mutex protecting access to window operations, including [Window.RunOnWin] functions.
	Mu sync.Mutex `view:"-"`

	// WinClose is a channel on which a single is sent to indicate that the
	// window should close.
	WinClose chan struct{} `view:"-"`

	// CloseReqFunc is the function to call on a close request
	CloseReqFunc func(win goosi.Window)

	// CloseCleanFunc is the function to call to close the window
	CloseCleanFunc func(win goosi.Window)

	// Nm is the name of the window
	Nm string `label:"Name"`

	// Titl is the title of the window
	Titl string `label:"Title"`

	// Flgs contains the flags associated with the window
	Flgs goosi.WindowFlags `label:"Flags"`

	// FPS is the FPS (frames per second) for rendering the window
	FPS int

	// EvMgr is the event manager for the window
	EvMgr events.Mgr `label:"Event manger"`

	// DestroyGPUFunc should be set to a function that will destroy GPU resources
	// in the main thread prior to destroying the drawer
	// and the surface; otherwise it is difficult to
	// ensure that the proper ordering of destruction applies.
	DestroyGPUFunc func()

	// CursorEnabled is whether the cursor is currently enabled
	CursorEnabled bool
}

// NewWindow makes a new [Window] for the given app with the given options.
func NewWindow[A goosi.App](a A, opts *goosi.NewWindowOptions) Window[A] {
	return Window[A]{
		WinClose:      make(chan struct{}),
		App:           a,
		Titl:          opts.GetTitle(),
		Flgs:          opts.Flags,
		FPS:           60,
		CursorEnabled: true,
	}
}

// WinLoop runs the window's own locked processing loop.
func (w *Window[A]) WinLoop() {
	defer func() { goosi.HandleRecover(recover()) }()

	var winPaint *time.Ticker
	if w.FPS > 0 {
		winPaint = time.NewTicker(time.Second / time.Duration(w.FPS))
	} else {
		winPaint = &time.Ticker{C: make(chan time.Time)} // nop
	}
outer:
	for {
		select {
		case <-w.WinClose:
			winPaint.Stop()
			break outer
		case <-winPaint.C:
			if w.This.IsClosed() {
				break outer
			}
			w.EvMgr.WindowPaint()
		}
	}
}

func (w *Window[A]) Lock() bool {
	if w.This.IsClosed() {
		return false
	}
	w.Mu.Lock()
	return true
}

func (w *Window[A]) Unlock() {
	w.Mu.Unlock()
}

func (w *Window[A]) Name() string {
	return w.Nm
}

func (w *Window[A]) SetName(name string) {
	w.Nm = name
}

func (w *Window[A]) Title() string {
	return w.Titl
}

func (w *Window[A]) SetTitle(title string) {
	if w.This.IsClosed() {
		return
	}
	w.Titl = title
}

func (w *Window[A]) SetIcon(images []image.Image) {
	// no-op by default
}

func (w *Window[A]) Flags() goosi.WindowFlags {
	return w.Flgs
}

func (w *Window[A]) Is(flag goosi.WindowFlags) bool {
	return w.Flgs.HasFlag(flag)
}

func (w *Window[A]) IsClosed() bool {
	return w == nil || w.This == nil || w.This.Drawer() == nil
}

func (w *Window[A]) IsVisible() bool {
	return !w.This.IsClosed() && !w.Is(goosi.Minimized)
}

func (w *Window[A]) SetFPS(fps int) {
	w.FPS = fps
}

func (w *Window[A]) EventMgr() *events.Mgr {
	return &w.EvMgr
}

func (w *Window[A]) SetDestroyGPUResourcesFunc(f func()) {
	w.DestroyGPUFunc = f
}

func (w *Window[A]) Insets() styles.SideFloats {
	// no-op by default
	return styles.NewSideFloats()
}

func (w *Window[A]) SetCloseReqFunc(fun func(win goosi.Window)) {
	w.Mu.Lock()
	defer w.Mu.Unlock()
	w.CloseReqFunc = fun
}

func (w *Window[A]) SetCloseCleanFunc(fun func(win goosi.Window)) {
	w.Mu.Lock()
	defer w.Mu.Unlock()
	w.CloseCleanFunc = fun
}

func (w *Window[A]) CloseReq() {
	if w.App.IsQuitting() {
		w.This.Close()
	}
	if w.CloseReqFunc != nil {
		w.CloseReqFunc(w.This)
	} else {
		w.This.Close()
	}
}

func (w *Window[A]) CloseClean() {
	if w.CloseCleanFunc != nil {
		w.CloseCleanFunc(w.This)
	}
}

func (w *Window[A]) Close() {
	// base implementation doesn't actually close any system windows,
	// but platform-specific implementations can
	w.EvMgr.Window(events.WinClose)

	w.Mu.Lock()
	defer w.Mu.Unlock()

	w.WinClose <- struct{}{}
	w.CloseClean()
	w.App.RemoveWindow(w.This)
}

func (w *Window[A]) MainMenu() goosi.MainMenu {
	// nil by default
	return nil
}

func (w *Window[A]) SetCursorEnabled(enabled, raw bool) {
	w.CursorEnabled = enabled
}

func (w *Window[A]) IsCursorEnabled() bool {
	return w.CursorEnabled
}

func (w *Window[A]) SetMousePos(x, y float64) {
	// no-op by default
}

func (w *Window[A]) SetTitleBarIsDark(isDark bool) {
	// no-op by default
}
