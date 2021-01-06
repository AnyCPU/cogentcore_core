// Copyright (c) 2018, The GoKi Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package oswin

import (
	"fmt"
	"image"
	"sync/atomic"
	"time"

	"github.com/goki/ki/kit"
	"github.com/goki/ki/nptime"
)

// GoGi event structure is derived from go.wde and golang/x/mobile/event
//
// GoGi requires event type enum for widgets to request what events to
// receive, and we add an overall interface with base support for time and
// marking events as processed, which is critical for simplifying logic and
// preventing unintended multiple effects
//
// OSWin deals exclusively in raw "dot" pixel integer coordinates (as in
// go.wde) -- abstraction to different DPI etc takes place higher up in the
// system

// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
   Copyright 2012 the go.wde authors

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
*/

// EventType determines which type of GUI event is being sent -- need this for
// indexing into different event signalers based on event type, and sending
// event type in signals -- critical to break up different event types into
// the right categories needed for different types of widgets -- e.g., most do
// not need move or scroll events, so those are separated.
type EventType int64

const (
	// MouseEvent includes all mouse button actions, but not move or drag
	MouseEvent EventType = iota

	// MouseMoveEvent is when the mouse is moving but no button is down
	MouseMoveEvent

	// MouseDragEvent is when the mouse is moving and there is a button down
	MouseDragEvent

	// MouseScrollEvent is for mouse scroll wheel events
	MouseScrollEvent

	// MouseFocusEvent is for mouse focus (enter / exit of widget area) --
	// generated by gi.Window based on mouse move events
	MouseFocusEvent

	// MouseHoverEvent is for mouse hover -- generated by gi.Window based on
	// mouse events
	MouseHoverEvent

	// KeyEvent for key pressed or released -- fine-grained data about each
	// key as it happens
	KeyEvent

	// KeyChordEvent is only generated when a non-modifier key is released,
	// and it also contains a string representation of the full chord,
	// suitable for translation into keyboard commands, emacs-style etc
	KeyChordEvent

	// TouchEvent is a generic touch-based event
	TouchEvent

	// MagnifyEvent is a touch-based magnify event (e.g., pinch)
	MagnifyEvent

	// RotateEvent is a touch-based rotate event
	RotateEvent

	// WindowEvent reports any changes in the window size, orientation,
	// iconify, close, open, paint -- these are all "internal" events
	// from OS to GUI system, and not sent to widgets
	WindowEvent

	// WindowResizeEvent is specifically for window resize events which need
	// special treatment -- this is an internal event not sent to widgets
	WindowResizeEvent

	// WindowPaintEvent is specifically for window paint events which need
	// special treatment -- this is an internal event not sent to widgets
	WindowPaintEvent

	// WindowShowEvent is a synthetic event sent to widget consumers,
	// sent *only once* when window is shown for the very first time
	WindowShowEvent

	// WindowFocusEvent is a synthetic event sent to widget consumers,
	// sent when window focus changes (action is Focus / DeFocus)
	WindowFocusEvent

	// DNDEvent is for the Drag-n-Drop (DND) drop event
	DNDEvent
	// DNDMoveEvent is when the DND position has changed
	DNDMoveEvent
	// DNDFocusEvent is for Enter / Exit events of the DND into / out of a given widget
	DNDFocusEvent

	// OSEvent is an operating system generated event (app level typically)
	OSEvent
	// OSOpenFilesEvent is an event telling app to open given files
	OSOpenFilesEvent

	// CustomEventType is a user-defined event with a data interface{} field
	CustomEventType

	// number of event types
	EventTypeN
)

//go:generate stringer -type=EventType

var KiT_EventType = kit.Enums.AddEnum(EventTypeN, kit.NotBitFlag, nil)

// Event is the interface for oswin GUI events.  also includes Stringer
// to get a string description of the event
type Event interface {
	fmt.Stringer

	// Type returns the type of event associated with given event
	Type() EventType

	// HasPos returns true if the event has a window position where it takes place
	HasPos() bool

	// Pos returns the position in raw display dots (pixels) where event took place -- needed for sending events to the right place
	Pos() image.Point

	// OnFocus returns true if the event operates only on focus item (e.g., keyboard events)
	OnFocus() bool

	// OnWinFocus returns true if the event operates only when the window has focus
	OnWinFocus() bool

	// Time returns the time at which the event was generated, in UnixNano nanosecond units
	Time() time.Time

	// IsProcessed returns whether this event has already been processed
	IsProcessed() bool

	// SetProcessed marks the event as having been processed
	SetProcessed()

	// Init sets the time to now, and any other init -- done just prior to event delivery
	Init()

	// SetTime sets the event time to Now
	SetTime()
}

//////////////////////////////////////////////////////////////////////
// EventBase

// EventBase is the base type for events -- records time and whether event has
// been processed by a receiver of the event -- in which case it is skipped
type EventBase struct {
	// GenTime records the time when the event was first generated, using more
	// efficient nptime struct
	GenTime nptime.Time

	// Processed indicates if the event has been processed by an end receiver,
	// and thus should no longer be processed by other possible receivers.
	// Atomic operations are used to encode a 0 or 1, so it is an int32.
	Processed int32
}

// SetTime sets the event time to Now
func (ev *EventBase) SetTime() {
	ev.GenTime.Now()
}

func (ev *EventBase) Init() {
	ev.SetTime()
}

func (ev EventBase) Time() time.Time {
	return ev.GenTime.Time()
}

func (ev EventBase) IsProcessed() bool {
	return atomic.LoadInt32(&ev.Processed) != 0
}

func (ev *EventBase) SetProcessed() {
	atomic.StoreInt32(&ev.Processed, int32(1))
}

func (ev *EventBase) ClearProcessed() {
	atomic.StoreInt32(&ev.Processed, int32(0))
}

func (ev EventBase) String() string {
	return fmt.Sprintf("Event at Time: %v", ev.Time())
}

func (ev EventBase) OnWinFocus() bool {
	return true
}

//////////////////////////////////////////////////////////////////////
// CustomEvent

// CustomEvent is a user-specified event that can be sent and received
// as needed, and contains a Data field for arbitrary data, and
// optional position and focus parameters
type CustomEvent struct {
	EventBase
	Data     interface{}
	PosAvail bool        `desc:"set to true if position is available"`
	Where    image.Point `desc:"position info if relevant -- set PosAvail"`
	Focus    bool        `desc:"set to true if this event should be sent to widget in focus"`
}

func (ce CustomEvent) Type() EventType {
	return CustomEventType
}

func (ce CustomEvent) String() string {
	return fmt.Sprintf("Type: %v Data: %v  Time: %v", ce.Type(), ce.Data, ce.Time())
}

func (ce CustomEvent) HasPos() bool {
	return ce.PosAvail
}

func (ce CustomEvent) Pos() image.Point {
	return ce.Where
}

func (ce CustomEvent) OnFocus() bool {
	return ce.Focus
}

func (ce CustomEvent) OnWinFocus() bool {
	return false
}

// SendCustomEvent sends a new custom event to given window, with
// given data -- constructs the event and sends it. For other params
// you can follow these steps yourself..
func SendCustomEvent(win Window, data interface{}) {
	ce := &CustomEvent{Data: data}
	ce.Init()
	win.Send(ce)
}

//////////////////////////////////////////////////////////////////////
// EventDeque

// EventDeque is an infinitely buffered double-ended queue of events.
type EventDeque interface {
	// Send adds an event to the end of the deque. They are returned by
	// NextEvent in FIFO order.
	Send(event Event)

	// SendFirst adds an event to the start of the deque. They are returned by
	// NextEvent in LIFO order, and have priority over events sent via Send.
	SendFirst(event Event)

	// NextEvent returns the next event in the deque. It blocks until such an
	// event has been sent.
	NextEvent() Event

	// PollEvent returns the next event in the deque if available, returns true
	// returns false and does not wait if no events currently available
	PollEvent() (Event, bool)
}
