// Copyright (c) 2023, Cogent Core. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package events

import (
	"fmt"
	"image"

	"cogentcore.org/core/events/key"
	"cogentcore.org/core/gox/nptime"
	"cogentcore.org/core/mat32"
	"cogentcore.org/core/mimedata"
)

// TraceWindowPaint prints out a . for each WindowPaint event
// - for other window events, * for mouse move events.
// Makes it easier to see what is going on in the overall flow.
var TraceWindowPaint = false

// Mgr manages the event construction and sending process,
// for its parent window.  Caches state as needed
// to generate derived events such as MouseDrag.
type Mgr struct {
	// Deque is the event queue
	Deque Deque

	// flag for ignoring mouse events when disabling mouse movement
	ResettingPos bool

	// Last has the prior state for key variables
	Last MgrState

	// PaintCount is used for printing paint events as .
	PaintCount int
}

// MgrState tracks basic event state over time
// to enable recognition and full data for generating events.
type MgrState struct {
	// last mouse button event type (down or up)
	MouseButtonType Types

	// last mouse button
	MouseButton Buttons

	// time of MouseDown
	MouseDownTime nptime.Time

	// position at MouseDown
	MouseDownPos image.Point

	// position of mouse from move events
	MousePos image.Point

	// time of last move
	MouseMoveTime nptime.Time

	// keyboard modifiers (Shift, Alt, etc)
	Mods key.Modifiers

	// Key event code
	Key key.Codes
}

///////////////////////////////////////////////////////////////
//  New Events

// SendKey processes a basic key event and sends it
func (em *Mgr) Key(typ Types, rn rune, code key.Codes, mods key.Modifiers) {
	ev := NewKey(typ, rn, code, mods)
	em.Last.Mods = mods
	em.Last.Key = code
	ev.Init()
	em.Deque.Send(ev)

	_, mapped := key.CodeRuneMap[code]

	if typ == KeyDown && ev.Code < key.CodeLeftControl &&
		(ev.HasAnyModifier(key.Control, key.Meta) || !mapped || ev.Code == key.CodeTab) {
		che := NewKey(KeyChord, rn, code, mods)
		che.Init()
		em.Deque.Send(che)
	}
}

// KeyChord processes a basic KeyChord event and sends it
func (em *Mgr) KeyChord(rn rune, code key.Codes, mods key.Modifiers) {
	ev := NewKey(KeyChord, rn, code, mods)
	// no further processing of these
	ev.Init()
	em.Deque.Send(ev)
}

// MouseButton creates and sends a mouse button event with given values
func (em *Mgr) MouseButton(typ Types, but Buttons, where image.Point, mods key.Modifiers) {
	ev := NewMouse(typ, but, where, mods)
	em.Last.Mods = mods
	em.Last.MouseButtonType = typ
	em.Last.MouseButton = but
	em.Last.MousePos = where
	ev.Init()
	if typ == MouseDown {
		em.Last.MouseDownPos = where
		em.Last.MouseDownTime = ev.GenTime
		em.Last.MouseMoveTime = ev.GenTime
	}
	em.Deque.Send(ev)
}

// MouseMove creates and sends a mouse move or drag event with given values
func (em *Mgr) MouseMove(where image.Point) {
	lastPos := em.Last.MousePos
	var ev *Mouse
	if em.Last.MouseButtonType == MouseDown {
		ev = NewMouseDrag(em.Last.MouseButton, where, lastPos, em.Last.MouseDownPos, em.Last.Mods)
		ev.StTime = em.Last.MouseDownTime
		ev.PrvTime = em.Last.MouseMoveTime
	} else {
		ev = NewMouseMove(em.Last.MouseButton, where, lastPos, em.Last.Mods)
		ev.PrvTime = em.Last.MouseMoveTime
	}
	ev.Init()
	em.Last.MouseMoveTime = ev.GenTime
	// if em.Win.IsCursorEnabled() {
	em.Last.MousePos = where
	// }
	if TraceWindowPaint {
		fmt.Printf("*")
	}
	em.Deque.Send(ev)
}

// Scroll creates and sends a scroll event with given values
func (em *Mgr) Scroll(where image.Point, delta mat32.Vec2) {
	ev := NewScroll(where, delta, em.Last.Mods)
	ev.Init()
	em.Deque.Send(ev)
}

// DropExternal creates and sends a Drop event with given values
func (em *Mgr) DropExternal(where image.Point, md mimedata.Mimes) {
	ev := NewExternalDrop(Drop, em.Last.MouseButton, where, em.Last.Mods, md)
	em.Last.MousePos = where
	ev.Init()
	em.Deque.Send(ev)
}

// Touch creates and sends a touch event with the given values.
// It also creates and sends a corresponding mouse event.
func (em *Mgr) Touch(typ Types, seq Sequence, where image.Point) {
	ev := NewTouch(typ, seq, where)
	ev.Init()
	em.Deque.Send(ev)

	if typ == TouchStart {
		em.MouseButton(MouseDown, Left, where, 0) // TODO: modifiers
	} else if typ == TouchEnd {
		em.MouseButton(MouseUp, Left, where, 0) // TODO: modifiers
	} else {
		em.MouseMove(where)
	}
}

// Magnify creates and sends a [TouchMagnify] event with the given values.
func (em *Mgr) Magnify(scaleFactor float32, where image.Point) {
	ev := NewMagnify(scaleFactor, where)
	ev.Init()
	em.Deque.Send(ev)
}

//	func (em *Mgr) DND(act dnd.Actions, where image.Point, data mimedata.Mimes) {
//		ev := dnd.NewEvent(act, where, em.Last.Mods)
//		ev.Data = data
//		ev.Init()
//		em.Deque.Send(ev)
//	}

func (em *Mgr) Window(act WinActions) {
	ev := NewWindow(act)
	ev.Init()
	if TraceWindowPaint {
		fmt.Printf("-")
	}
	em.Deque.SendFirst(ev)
}

func (em *Mgr) WindowPaint() {
	ev := NewWindowPaint()
	ev.Init()
	if TraceWindowPaint {
		fmt.Printf(".")
		em.PaintCount++
		if em.PaintCount > 60 {
			fmt.Println("")
			em.PaintCount = 0
		}
	}
	em.Deque.SendFirst(ev) // separate channel for window!
}

func (em *Mgr) WindowResize() {
	ev := NewWindowResize()
	ev.Init()
	if TraceWindowPaint {
		fmt.Printf("r")
	}
	em.Deque.SendFirst(ev)
}

func (em *Mgr) Custom(data any) {
	ce := &CustomEvent{}
	ce.Typ = Custom
	ce.Data = data
	ce.Init()
	em.Deque.Send(ce)
}
