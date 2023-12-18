// Copyright (c) 2018, The GoKi Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package styles

import (
	"image"
	"image/color"

	"goki.dev/colors"
	"goki.dev/girl/units"
	"goki.dev/mat32/v2"
)

// Paint provides the styling parameters for SVG-style rendering
type Paint struct { //gti:add

	// prop: display:none -- node and everything below it are off, non-rendering
	Off bool

	// todo big enum of how to display item -- controls layout etc
	Display bool `xml:"display"`

	// stroke (line drawing) parameters
	StrokeStyle Stroke

	// fill (region filling) parameters
	FillStyle Fill

	// font also has global opacity setting, along with generic color, background-color settings, which can be copied into stroke / fill as needed
	FontStyle FontRender

	// font also has global opacity setting, along with generic color, background-color settings, which can be copied into stroke / fill as needed
	TextStyle Text

	// prop: vector-effect = various rendering special effects settings
	VecEff VectorEffects `xml:"vector-effect"`

	// prop: transform = our additions to transform -- pushed to render state
	Transform mat32.Mat2 `xml:"transform"`

	// units context -- parameters necessary for anchoring relative units
	UnContext units.Context `xml:"-"`

	// have the styles already been set?
	StyleSet bool

	PropsNil   bool
	dotsSet    bool
	lastUnCtxt units.Context
}

func (pc *Paint) Defaults() {
	pc.Off = false
	pc.Display = true
	pc.StyleSet = false
	pc.StrokeStyle.Defaults()
	pc.FillStyle.Defaults()
	pc.FontStyle.Defaults()
	pc.TextStyle.Defaults()
	pc.Transform = mat32.Identity2D()
}

// CopyStyleFrom copies styles from another paint
func (pc *Paint) CopyStyleFrom(cp *Paint) {
	pc.Off = cp.Off
	pc.Display = cp.Display
	pc.UnContext = cp.UnContext
	pc.StrokeStyle = cp.StrokeStyle
	pc.FillStyle = cp.FillStyle
	pc.FontStyle = cp.FontStyle
	pc.TextStyle = cp.TextStyle
	pc.VecEff = cp.VecEff
}

// InheritFields from parent: Manual inheriting of values is much faster than
// automatic version!
func (pc *Paint) InheritFields(par *Paint) {
	pc.FontStyle.InheritFields(&par.FontStyle)
	pc.TextStyle.InheritFields(&par.TextStyle)
}

// SetStyleProps sets paint values based on given property map (name: value
// pairs), inheriting elements as appropriate from parent, and also having a
// default style for the "initial" setting
func (pc *Paint) SetStyleProps(par *Paint, props map[string]any, ctxt colors.Context) {
	if !pc.StyleSet && par != nil { // first time
		pc.InheritFields(par)
	}
	pc.StyleFromProps(par, props, ctxt)

	pc.PropsNil = (len(props) == 0)
	pc.StyleSet = true
}

func (pc *Paint) FromStyle(st *Style) {
	pc.UnContext = st.UnContext
	pc.FontStyle = *st.FontRender()
	pc.TextStyle = st.Text
}

// ToDotsImpl runs ToDots on unit values, to compile down to raw pixels
func (pc *Paint) ToDotsImpl(uc *units.Context) {
	pc.StrokeStyle.ToDots(uc)
	pc.FillStyle.ToDots(uc)
	pc.FontStyle.ToDots(uc)
	pc.TextStyle.ToDots(uc)
}

// SetUnitContextExt sets the unit context for external usage of paint
// outside of a Viewport, based on overall size of painting canvas.
// caches everything out in terms of raw pixel dots for rendering
// call at start of render.
func (pc *Paint) SetUnitContextExt(size image.Point) {
	pc.UnContext.Defaults()
	// paint (SVG) context is always 96 = 1to1
	pc.UnContext.DPI = 96
	// TODO: maybe should have different values for these sizes?
	pc.UnContext.SetSizes(float32(size.X), float32(size.Y), float32(size.X), float32(size.Y), float32(size.X), float32(size.Y))
	pc.FontStyle.SetUnitContext(&pc.UnContext)
	pc.ToDotsImpl(&pc.UnContext)
	pc.dotsSet = true
}

// ToDots runs ToDots on unit values, to compile down to raw pixels
func (pc *Paint) ToDots() {
	if !(pc.dotsSet && pc.UnContext == pc.lastUnCtxt && pc.PropsNil) {
		pc.ToDotsImpl(&pc.UnContext)
		pc.dotsSet = true
		pc.lastUnCtxt = pc.UnContext
	}
}

/////////////////////////////////////////////////////////////////
//  enums

type FillRules int32 //enums:enum -trim-prefix FillRule

const (
	FillRuleNonZero FillRules = iota
	FillRuleEvenOdd
)

// VectorEffects contains special effects for rendering
type VectorEffects int32 //enums:enum -trim-prefix VecEff

const (
	VecEffNone VectorEffects = iota

	// VecEffNonScalingStroke means that the stroke width is not affected by
	// transform properties
	VecEffNonScalingStroke
)

// IMPORTANT: any changes here must be updated below in StyleFillFuncs

// Fill contains all the properties for filling a region
type Fill struct {

	// prop: fill = fill color image specification; filling is off if nil
	Color image.Image `xml:"fill"`

	// prop: fill-opacity = global alpha opacity / transparency factor between 0 and 1
	Opacity float32 `xml:"fill-opacity"`

	// prop: fill-rule = rule for how to fill more complex shapes with crossing lines
	Rule FillRules `xml:"fill-rule"`
}

// Defaults initializes default values for paint fill
func (pf *Fill) Defaults() {
	pf.Color = colors.C(color.Black)
	pf.Rule = FillRuleNonZero
	pf.Opacity = 1.0
}

// ToDots runs ToDots on unit values, to compile down to raw pixels
func (fs *Fill) ToDots(uc *units.Context) {
}

////////////////////////////////////////////////////////////////////////////////////
// Stroke

// end-cap of a line: stroke-linecap property in SVG
type LineCaps int32 //enums:enum -trim-prefix LineCap

const (
	LineCapButt LineCaps = iota
	LineCapRound
	LineCapSquare
	// rasterx extension
	LineCapCubic
	// rasterx extension
	LineCapQuadratic
)

// the way in which lines are joined together: stroke-linejoin property in SVG
type LineJoins int32 //enums:enum -trim-prefix LineJoin

const (
	LineJoinMiter LineJoins = iota
	LineJoinMiterClip
	LineJoinRound
	LineJoinBevel
	LineJoinArcs
	// rasterx extension
	LineJoinArcsClip
)

// IMPORTANT: any changes here must be updated below in StyleStrokeFuncs

// Stroke contains all the properties for painting a line
type Stroke struct {

	// prop: stroke = stroke color image specification; stroking is off if nil
	Color image.Image `xml:"stroke"`

	// prop: stroke-opacity = global alpha opacity / transparency factor between 0 and 1
	Opacity float32 `xml:"stroke-opacity"`

	// prop: stroke-width = line width
	Width units.Value `xml:"stroke-width"`

	// prop: stroke-min-width = minimum line width used for rendering -- if width is > 0, then this is the smallest line width -- this value is NOT subject to transforms so is in absolute dot values, and is ignored if vector-effects non-scaling-stroke is used -- this is an extension of the SVG / CSS standard
	MinWidth units.Value `xml:"stroke-min-width"`

	// prop: stroke-dasharray = dash pattern, in terms of alternating on and off distances -- e.g., = 4 pixels on, 4 pixels off.  Currently only supporting raw pixel numbers, but in principle should support units.
	Dashes []float32 `xml:"stroke-dasharray"`

	// prop: stroke-linecap = how to draw the end cap of lines
	Cap LineCaps `xml:"stroke-linecap"`

	// prop: stroke-linejoin = how to join line segments
	Join LineJoins `xml:"stroke-linejoin"`

	// prop: stroke-miterlimit = limit of how far to miter -- must be 1 or larger
	MiterLimit float32 `xml:"stroke-miterlimit" min:"1"`
}

// Defaults initializes default values for paint stroke
func (ps *Stroke) Defaults() {
	// stroking is off by default in svg
	ps.Color = nil
	ps.Width.Px(1)
	ps.MinWidth.Dot(.5)
	ps.Cap = LineCapButt
	ps.Join = LineJoinMiter // Miter not yet supported, but that is the default -- falls back on bevel
	ps.MiterLimit = 10.0
	ps.Opacity = 1.0
}

// ToDots runs ToDots on unit values, to compile down to raw pixels
func (ss *Stroke) ToDots(uc *units.Context) {
	ss.Width.ToDots(uc)
	ss.MinWidth.ToDots(uc)
}
