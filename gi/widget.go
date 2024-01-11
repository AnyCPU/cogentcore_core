// Copyright (c) 2023, The Goki Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gi

//go:generate goki generate

import (
	"image"
	"log"
	"sync"

	"goki.dev/abilities"
	"goki.dev/enums"
	"goki.dev/events"
	"goki.dev/gti"
	"goki.dev/ki"
	"goki.dev/states"
	"goki.dev/styles"
)

// Widget is the interface for all GoGi Widget Nodes
type Widget interface {
	ki.Ki

	// OnWidgetAdded adds a function to call when a widget is added
	// as a child to the widget or any of its children.
	OnWidgetAdded(f func(w Widget)) *WidgetBase

	// Style sets the styling of the widget by adding a Styler function
	Style(s func(s *styles.Style)) *WidgetBase

	// AsWidget returns the WidgetBase embedded field for any Widget node.
	// The Widget interface defines only methods that can be overridden
	// or need to be called on other nodes.  Everything else that is common
	// to all Widgets is in the WidgetBase.
	AsWidget() *WidgetBase

	// Config configures the widget, primarily configuring its Parts.
	// it does _not_ call Config on children, just self.
	// ApplyStyle must generally be called after Config - it is called
	// automatically when Scene is first shown, but must be called
	// manually thereafter as needed after configuration changes.
	// See Update for a convenience function that does both.
	// ConfigScene on Scene handles full tree configuration.
	// This config calls UpdateStart / End, and SetNeedsLayout,
	// and calls ConfigWidget to do the actual configuration,
	// so it does not need to manage this housekeeping.
	// Thus, this Config call is typically never changed, and
	// all custom configuration should happen in ConfigWidget.
	Config()

	// ConfigWidget does the actual configuration of the widget,
	// primarily configuring its Parts.
	// All configuration should be robust to multiple calls
	// (i.e., use Parts.ConfigChildren with Config).
	// Outer Config call handles all the other infrastructure,
	// so this call just does the core configuration.
	ConfigWidget()

	// Update calls Config and ApplyStyle on this widget.
	// This should be called if any config options are changed
	// while the Scene is being viewed.
	Update()

	// StateIs returns whether the widget has the given [states.States] flag set
	StateIs(flag states.States) bool

	// AbilityIs returns whether the widget has the given [abilities.Abilities] flag set
	AbilityIs(flag abilities.Abilities) bool

	// SetState sets the given [states.State] flags to the given value
	SetState(on bool, state ...states.States) *WidgetBase

	// SetAbilities sets the given [abilities.Abilities] flags to the given value
	SetAbilities(on bool, able ...abilities.Abilities) *WidgetBase

	// ApplyStyle applies style functions to the widget based on current state.
	// It is typically not overridden; instead, call Style to apply custom styling.
	ApplyStyle()

	// SizeUp (bottom-up) gathers Actual sizes from our Children & Parts,
	// based on Styles.Min / Max sizes and actual content sizing
	// (e.g., text size).  Flexible elements (e.g., Label, Flex Wrap,
	// TopAppBar) should reserve the _minimum_ size possible at this stage,
	// and then Grow based on SizeDown allocation.
	SizeUp()

	// SizeDown (top-down, multiple iterations possible) provides top-down
	// size allocations based initially on Scene available size and
	// the SizeUp Actual sizes.  If there is extra space available, it is
	// allocated according to the Grow factors.
	// Flexible elements (e.g., Flex Wrap layouts and Label with word wrap)
	// update their Actual size based on available Alloc size (re-wrap),
	// to fit the allocated shape vs. the initial bottom-up guess.
	// However, do NOT grow the Actual size to match Alloc at this stage,
	// as Actual sizes must always represent the minimums (see Position).
	// Returns true if any change in Actual size occurred.
	SizeDown(iter int) bool

	// SizeFinal: (bottom-up) similar to SizeUp but done at the end of the
	// Sizing phase: first grows widget Actual sizes based on their Grow
	// factors, up to their Alloc sizes.  Then gathers this updated final
	// actual Size information for layouts to register their actual sizes
	// prior to positioning, which requires accurate Actual vs. Alloc
	// sizes to perform correct alignment calculations.
	SizeFinal()

	// Position uses the final sizes to set relative positions within layouts
	// according to alignment settings, and Grow elements to their actual
	// Alloc size per Styles settings and widget-specific behavior.
	Position()

	// ScenePos computes scene-based absolute positions and final BBox
	// bounding boxes for rendering, based on relative positions from
	// Position step and parents accumulated position and scroll offset.
	// This is the only step needed when scrolling (very fast).
	ScenePos()

	// Render performs actual rendering pass.  Bracket the render calls in
	// PushBounds / PopBounds and a false from PushBounds indicates that
	// the node is invisible and should not be rendered.
	// If Parts are present, RenderParts is called by default.
	// Layouts or other widgets that manage children should call RenderChildren.
	Render()

	// On adds an event listener function for the given event type
	On(etype events.Types, fun func(e events.Event)) *WidgetBase

	// Helper functions for common event types
	// TODO(kai/menu): should we have these in the Widget interface?
	// we need them for OnWidgetAdded functions

	// OnClick adds an event listener function for [events.Click] events
	OnClick(fun func(e events.Event)) *WidgetBase

	// HandleEvent sends the given event to all Listeners for that event type.
	// It also checks if the State has changed and calls ApplyStyle if so.
	// If more significant Config level changes are needed due to an event,
	// the event handler must do this itself.
	HandleEvent(e events.Event)

	// Send sends an NEW event of given type to this widget,
	// optionally starting from values in the given original event
	// (recommended to include where possible).
	// Do NOT send an existing event using this method if you
	// want the Handled state to persist throughout the call chain;
	// call HandleEvent directly for any existing events.
	Send(e events.Types, orig ...events.Event)

	// WidgetTooltip returns the tooltip text that should be used for this
	// widget. By default, it just returns [WidgetBase.Tooltip], but widgets
	// can override it to do different things. For example, buttons add their
	// shortcut to the tooltip here.
	WidgetTooltip() string

	// AddContextMenu adds the given context menu to [WidgetBase.ContextMenus].
	// It is the main way that code should modify a widget's context menus.
	// Context menu functions are run in reverse order.
	AddContextMenu(menu func(m *Scene)) *WidgetBase

	// ApplyContextMenus adds the [Widget.ContextMenus] to the given menu scene
	// in reverse order.
	ApplyContextMenus(m *Scene)

	// ContextMenuPos returns the default position for popup menus;
	// by default in the middle its Bounding Box, but can be adapted as
	// appropriate for different widgets.
	ContextMenuPos(e events.Event) image.Point

	// ShowContextMenu displays the context menu of various actions
	// to perform on a Widget, activated by default on the ShowContextMenu
	// event, triggered by a Right mouse click.
	// Returns immediately, and actions are all executed directly
	// (later) via the action signals. Calls ContextMenu and
	// ContextMenuPos.
	ShowContextMenu(e events.Event)

	// IsVisible returns true if a node is visible for rendering according
	// to the [states.Invisible] flag on it or any of its parents.
	// This flag is also set by [styles.DisplayNone] during [ApplyStyle].
	// This does *not* check for an empty TotalBBox, indicating that the widget
	// is out of render range -- that is done by [PushBounds] prior to rendering.
	// Non-visible nodes are automatically not rendered and do not get
	// window events.
	// This call recursively calls the parent, which is typically a short path.
	IsVisible() bool

	// todo: revisit this -- in general anything with a largish image (including svg,
	// SubScene, but not Icon) should get put on a list so the RenderWin Drawer just
	// directly uploads its image.

	// IsDirectWinUpload returns true if this is a node that does a direct window upload
	// e.g., for gi3d.Scene which renders directly to the window texture for maximum efficiency
	IsDirectWinUpload() bool

	// DirectWinUpload does a direct upload of contents to a window
	// Drawer compositing image, which will then be used for drawing
	// the window during a Publish() event (triggered by the window Update
	// event).  This is called by the scene in its Update signal processing
	// routine on nodes that respond true to IsDirectWinUpload().
	// The node is also free to update itself of its own accord at any point.
	DirectWinUpload()
}

// WidgetBase is the base type for all Widget Widget elements, which are
// managed by a containing Layout, and use all 5 rendering passes.  All
// elemental widgets must support the ReadOnly and Selected states in a
// reasonable way (Selected only essential when also ReadOnly), so they can
// function appropriately in a chooser (e.g., SliceView or TableView) -- this
// includes toggling selection on left mouse press.
type WidgetBase struct {
	ki.Node

	// Tooltip is the text for the tooltip for this widget,
	// which can use HTML formatting.
	Tooltip string

	// Parts are a separate tree of sub-widgets that implement discrete parts
	// of a widget.  Positions are relative to the parent widget.
	// These are fully managed by the parent widget
	Parts *Layout `copy:"-" json:"-" xml:"-" set:"-"`

	// Geom has the full layout geometry for size and position of this Widget
	Geom GeomState `edit:"-" copy:"-" json:"-" xml:"-" set:"-"`

	// If true, Override the computed styles and allow directly editing Style
	OverrideStyle bool `copy:"-" json:"-" xml:"-" set:"-"`

	// Styles are styling settings for this widget.
	// These are set in SetApplyStyle which should be called after any Config
	// change (e.g., as done by the Update method).  See Stylers for functions
	// that set all of the styles, ordered from initial base defaults to later
	// added overrides.
	Styles styles.Style `copy:"-" json:"-" xml:"-" set:"-"`

	// Stylers are a slice of functions that are called in sequential
	// ascending order (so the last added styler is called last and
	// thus overrides all other functions) to style the element.
	// These should be set using Style function, which can be called
	// by end-user and internal code. FirstStylers and FinalStylers
	// are called before and after these stylers, respectively.
	Stylers []func(s *styles.Style) `copy:"-" json:"-" xml:"-" set:"-"`

	// FirstStylers are a slice of functions that are called in sequential
	// ascending order (so the last added styler is called last and
	// thus overrides all other functions) to style the element.
	// These should be set using StyleFirst function, which can be called
	// by end-user and internal code. These stylers are called before
	// Stylers and FinalStylers.
	FirstStylers []func(s *styles.Style) `copy:"-" json:"-" xml:"-" set:"-"`

	// FinalStylers are a slice of functions that are called in sequential
	// ascending order (so the last added styler is called last and
	// thus overrides all other functions) to style the element.
	// These should be set using StyleFinal function, which can be called
	// by end-user and internal code. These stylers are called after
	// FirstStylers and Stylers.
	FinalStylers []func(s *styles.Style) `copy:"-" json:"-" xml:"-" set:"-"`

	// A slice of functions to call on all widgets that are added as children
	// to this widget or its children.  These functions are called in sequential
	// ascending order, so the last added one is called last and thus can
	// override anything set by the other ones. These should be set using
	// OnWidgetAdded, which can be called by both end-user and internal code.
	OnWidgetAdders []func(w Widget) `copy:"-" json:"-" xml:"-" set:"-"`

	// Listeners are event listener functions for processing events on this widget.
	// type specific Listeners are added in OnInit when the widget is initialized.
	Listeners events.Listeners `copy:"-" json:"-" xml:"-" set:"-"`

	// PriorityEvents has event type(s) that this widget gets sent first.
	// Events are sent in depth-first order, so this enables outer container
	// widgets to get first access to these events.
	PriorityEvents []events.Types `set:"-"`

	// ContextMenus is a slice of menu functions to call to construct
	// the widget's context menu on an [events.ContextMenu]. The
	// functions are called in reverse order such that the elements
	// added in the last function are the first in the menu.
	// Context menus should be added through [Widget.AddContextMenu].
	ContextMenus []func(m *Scene) `json:"-" xml:"-" set:"-"`

	// Sc is the overall Scene to which we belong. It is automatically
	// by widgets whenever they are added to another widget parent.
	// It is passed to most Config, Layout, and Render functions as
	// a convenience.
	Sc *Scene `copy:"-" json:"-" xml:"-" set:"-"`

	// mutex protecting the Style field
	StyMu sync.RWMutex `copy:"-" view:"-" json:"-" xml:"-" set:"-"`

	// mutex protecting the BBox fields
	BBoxMu sync.RWMutex `copy:"-" view:"-" json:"-" xml:"-" set:"-"`
}

func (wb *WidgetBase) FlagType() enums.BitFlagSetter {
	return (*WidgetFlags)(&wb.Flags)
}

// OnInit for WidgetBase is not called by the usual ki mechanism,
// but should be called by every Widget type in its own OnInit
// to establish all the default styling and event handling
// that applies to all widgets.
func (wb *WidgetBase) OnInit() {
	wb.SetStyles()
	wb.HandleEvents()
}

// OnAdd is called when widgets are added to a parent.
// It sets the scene of the widget to its widget parent.
// It should be called by all other OnAdd functions defined
// by widget types.
func (wb *WidgetBase) OnAdd() {
	if pwb := wb.ParentWidget(); pwb != nil {
		wb.Sc = pwb.Sc
	}
}

// SetScene sets the Scene pointer for this widget and all of its children.
// This can be necessary when creating widgets outside the usual "NewWidget" paradigm,
// e.g., when reading from a JSON file.
func (wb *WidgetBase) SetScene(sc *Scene) {
	wb.WidgetWalkPre(func(kwi Widget, kwb *WidgetBase) bool {
		kwb.Sc = sc
		return ki.Continue
	})
}

func (wb *WidgetBase) OnChildAdded(child ki.Ki) {
	w, _ := AsWidget(child)
	if w == nil {
		return
	}
	for _, f := range wb.OnWidgetAdders {
		f(w)
	}
}

// OnWidgetAdded adds a function to call when a widget is added
// as a child to the widget or any of its children.
func (wb *WidgetBase) OnWidgetAdded(fun func(w Widget)) *WidgetBase {
	wb.OnWidgetAdders = append(wb.OnWidgetAdders, fun)
	return wb
}

// AsWidget returns the given Ki object
// as a Widget interface and a WidgetBase.
func AsWidget(k ki.Ki) (Widget, *WidgetBase) {
	if k == nil || k.This() == nil {
		return nil, nil
	}
	if w, ok := k.This().(Widget); ok {
		return w, w.AsWidget()
	}
	return nil, nil
}

func (wb *WidgetBase) AsWidget() *WidgetBase {
	return wb
}

// AsWidgetBase returns the given Ki object as a WidgetBase, or nil.
// for direct use of the return value in cases where that is needed.
func AsWidgetBase(k ki.Ki) *WidgetBase {
	_, wb := AsWidget(k)
	return wb
}

func (wb *WidgetBase) CopyFieldsFrom(frm any) {
	fr, ok := frm.(*WidgetBase)
	if !ok {
		log.Printf("GoGi node of type: %v needs a CopyFieldsFrom method defined\n", wb.KiType().Name)
		return
	}
	wb.Tooltip = fr.Tooltip
	wb.Styles.CopyFrom(&fr.Styles)
	wb.Stylers = fr.Stylers
	wb.Listeners = fr.Listeners // direct copy -- functions..
	wb.ContextMenus = fr.ContextMenus
}

func (wb *WidgetBase) Destroy() {
	wb.DeleteParts()
	wb.Node.Destroy()
}

func (wb *WidgetBase) DeleteParts() {
	if wb.Parts != nil {
		wb.Parts.DeleteChildren(true) // first delete all my children
	}
	wb.Parts = nil
}

func (wb *WidgetBase) BaseType() *gti.Type {
	return WidgetBaseType
}

// NewParts makes the Parts layout if not already there,
// with given layout orientation
func (wb *WidgetBase) NewParts() *Layout {
	if wb.Parts != nil {
		return wb.Parts
	}
	parts := ki.NewRoot[*Layout]("parts")
	ki.SetParent(parts, wb.This())
	parts.SetFlag(true, ki.Field)
	wb.Parts = parts
	return parts
}

// ParentWidget returns the parent as a [WidgetBase] or nil
// if this is the root and has no parent.
func (wb *WidgetBase) ParentWidget() *WidgetBase {
	if wb.Par == nil {
		return nil
	}
	return wb.Par.(Widget).AsWidget()
}

// ParentWidgetIf returns the nearest widget parent
// of the widget for which the given function returns true.
// It returns nil if no such parent is found.
func (wb *WidgetBase) ParentWidgetIf(fun func(p *WidgetBase) bool) *WidgetBase {
	cur := wb
	for {
		par := cur.Par
		if par == nil {
			return nil
		}
		pwi, ok := par.(Widget)
		if !ok {
			return nil
		}
		pwb := pwi.AsWidget()
		if fun(pwb) {
			return pwb
		}
		cur = pwb
	}
}

// IsVisible returns true if a node is visible for rendering according
// to the [states.Invisible] flag on it or any of its parents.
// This flag is also set by [styles.DisplayNone] during [ApplyStyle].
// This does *not* check for an empty TotalBBox, indicating that the widget
// is out of render range -- that is done by [PushBounds] prior to rendering.
// Non-visible nodes are automatically not rendered and do not get
// window events.
// This call recursively calls the parent, which is typically a short path.
func (wb *WidgetBase) IsVisible() bool {
	if wb == nil || wb.This() == nil || wb.Is(ki.Deleted) || wb.StateIs(states.Invisible) || wb.Sc == nil {
		return false
	}
	if wb.Par == nil || wb.Par.This() == nil {
		return true
	}
	return wb.Par.This().(Widget).IsVisible()
}

func (wb *WidgetBase) IsDirectWinUpload() bool {
	return false
}

func (wb *WidgetBase) DirectWinUpload() {
}

// WalkPreNode extends WalkPre to Parts -- key for getting full Update protection!
func (wb *WidgetBase) WalkPreNode(fun func(ki.Ki) bool) {
	if wb.Parts == nil {
		return
	}
	wb.Parts.WalkPre(fun)
}

// WidgetKidsIter iterates through the Kids, as widgets, calling the given function.
// Return ki.Continue (true) to continue, and ki.Break (false) to terminate.
func (wb *WidgetBase) WidgetKidsIter(fun func(i int, kwi Widget, kwb *WidgetBase) bool) {
	for i, k := range wb.Kids {
		i := i
		kwi, kwb := AsWidget(k)
		if kwi == nil || kwi.This() == nil || kwi.Is(ki.Deleted) {
			break
		}
		cont := fun(i, kwi, kwb)
		if !cont {
			break
		}
	}
}

// VisibleKidsIter iterates through the Kids, as widgets, calling the given function,
// excluding any with the *local* states.Invisible flag set (does not check parents).
// This is used e.g., for layout functions to exclude non-visible direct children.
// Return ki.Continue (true) to continue, and ki.Break (false) to terminate.
func (wb *WidgetBase) VisibleKidsIter(fun func(i int, kwi Widget, kwb *WidgetBase) bool) {
	for i, k := range wb.Kids {
		i := i
		kwi, kwb := AsWidget(k)
		if kwi == nil || kwi.This() == nil || kwi.Is(ki.Deleted) {
			break
		}
		if kwb.StateIs(states.Invisible) {
			continue
		}
		cont := fun(i, kwi, kwb)
		if !cont {
			break
		}
	}
}

// WidgetWalkPre is a version of the ki WalkPre iterator that automatically filters
// nil or deleted items and operates on Widget types.
// Return ki.Continue (true) to continue, and ki.Break (false) to terminate.
func (wb *WidgetBase) WidgetWalkPre(fun func(kwi Widget, kwb *WidgetBase) bool) {
	wb.WalkPre(func(k ki.Ki) bool {
		kwi, kwb := AsWidget(k)
		if kwi == nil || kwi.This() == nil || kwi.Is(ki.Deleted) {
			return ki.Break
		}
		return fun(kwi, kwb)
	})
}

// WidgetNextVisible returns the next visible node in the tree as a Widget,
// nil if no more.
func (wb *WidgetBase) WidgetNextVisible() (Widget, *WidgetBase) {
	nwi, nwb := AsWidget(ki.Next(wb.This()))
	for nwb != nil && !nwb.IsVisible() {
		nwi, nwb = AsWidget(ki.Next(nwb.This()))
	}
	return nwi, nwb
}

// WidgetPrevVisible returns the previous visible node in the tree as a Widget,
// nil if no more.
func (wb *WidgetBase) WidgetPrevVisible() (Widget, *WidgetBase) {
	nwi, nwb := AsWidget(ki.Prev(wb.This()))
	for nwb != nil && !nwb.IsVisible() {
		nwi, nwb = AsWidget(ki.Prev(nwb.This()))
	}
	return nwi, nwb
}

// WidgetTooltip is the base implementation of [Widget.WidgetTooltip],
// which just returns [WidgetBase.Tooltip].
func (wb *WidgetBase) WidgetTooltip() string {
	return wb.Tooltip
}
