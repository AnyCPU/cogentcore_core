// Copyright (c) 2018, The Goki Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package cursor defines the oswin cursor interface and standard system
// cursors that are supported across platforms
package cursor

//go:generate goki generate

import (
	"goki.dev/goki/enums"
)

// Cursor manages the mouse cursor / pointer appearance.  Currently only a
// fixed set of standard cursors are supported, but in the future it will be
// possible to set the cursor from an image / svg.
type Cursor interface {

	// Current returns the current cursor as an enum, which is a
	// [goki.dev/goki/cursors.Cursor]
	// by default, but could be something else if you are extending
	// the default cursor set.
	Current() enums.Enum

	// Set sets the active cursor to the given cursor as an enum, which is typically
	// a [cursors.Cursor], unless you are extending the default cursor set, in
	// which case it should be a type you defined. The string version of the
	// enum value must correspond to a filename of the form "name.svg" in
	// [goki.dev/goki/cursors.Cursors]; this will be satisfied automatically by all
	// [cursor.Cursor] values.
	Set(cursor enums.Enum) error

	// IsVisible returns whether cursor is currently visible (according to [Cursor.Hide] and [Cursor.Show] actions)
	IsVisible() bool

	// Hide hides the cursor if it is not already hidden.
	Hide()

	// Show shows the cursor after a hide if it is hidden.
	Show()

	// SetSize sets the size that cursors are rendered at.
	SetSize(size int)
}

// CursorBase provides the common infrastructure for the [Cursor] interface,
// to be extended on desktop platforms. It can also be used as an empty
// implementation of the [Cursor] interface on mobile platforms, as they
// do not have cursors.
type CursorBase struct {
	// Cur is the current cursor, which is maintained by the standard methods.
	Cur enums.Enum

	// Vis is whether the cursor is visible; be sure to initialize to true!
	Vis bool

	// Size is the size that cursors are rendered at
	Size int
}

// CursorBase should be a valid cursor so that it can be used directly in mobile
var _ Cursor = (*CursorBase)(nil)

func (c *CursorBase) Current() enums.Enum {
	return c.Cur
}

func (c *CursorBase) Set(cursor enums.Enum) error {
	c.Cur = cursor
	return nil
}

func (c *CursorBase) IsVisible() bool {
	return c.Vis
}

func (c *CursorBase) Hide() {
	c.Vis = false
}

func (c *CursorBase) Show() {
	c.Vis = true
}

func (c *CursorBase) SetSize(size int) {
	c.Size = size
}
