// Copyright (c) 2018, The GoKi Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package svg

import (
	"goki.dev/girl/styles"
	"goki.dev/girl/units"
	"goki.dev/ki/v2"
	"goki.dev/mat32/v2"
)

// Rect is a SVG rectangle, optionally with rounded corners
type Rect struct {
	NodeBase

	// position of the top-left of the rectangle
	Pos mat32.Vec2 `xml:"{x,y}"`

	// size of the rectangle
	Size mat32.Vec2 `xml:"{width,height}"`

	// radii for curved corners, as a proportion of width, height
	Radius mat32.Vec2 `xml:"{rx,ry}"`
}

// AddNewRect adds a new rectangle to given parent node, with given name, pos, and size.
func AddNewRect(parent ki.Ki, name string, x, y, sx, sy float32) *Rect {
	g := parent.NewChild(RectType, name).(*Rect)
	g.Pos.Set(x, y)
	g.Size.Set(sx, sy)
	return g
}

func (g *Rect) SVGName() string { return "rect" }

func (g *Rect) CopyFieldsFrom(frm any) {
	fr := frm.(*Rect)
	g.NodeBase.CopyFieldsFrom(&fr.NodeBase)
	g.Pos = fr.Pos
	g.Size = fr.Size
	g.Radius = fr.Radius
}

func (g *Rect) SetPos(pos mat32.Vec2) {
	g.Pos = pos
}

func (g *Rect) SetSize(sz mat32.Vec2) {
	g.Size = sz
}

func (g *Rect) LocalBBox() mat32.Box2 {
	bb := mat32.Box2{}
	hlw := 0.5 * g.LocalLineWidth()
	bb.Min = g.Pos.SubScalar(hlw)
	bb.Max = g.Pos.Add(g.Size).AddScalar(hlw)
	return bb
}

func (g *Rect) Render(sv *SVG) {
	vis, rs := g.PushXForm(sv)
	if !vis {
		return
	}
	pc := &g.Paint
	rs.Lock()
	// TODO: figure out a better way to do this
	bs := styles.Border{}
	bs.Style.Set(styles.BorderSolid)
	bs.Width.Set(pc.StrokeStyle.Width)
	bs.Color.Set(pc.StrokeStyle.Color.Solid)
	bs.Radius.Set(units.Px(g.Radius.X))
	if g.Radius.X == 0 && g.Radius.Y == 0 {
		pc.DrawRectangle(rs, g.Pos.X, g.Pos.Y, g.Size.X, g.Size.Y)
	} else {
		// todo: only supports 1 radius right now -- easy to add another
		// SidesTODO: also support different radii for each corner
		pc.DrawRoundedRectangle(rs, g.Pos.X, g.Pos.Y, g.Size.X, g.Size.Y, styles.NewSideFloats(g.Radius.X))
	}
	pc.FillStrokeClear(rs)
	rs.Unlock()
	g.BBoxes(sv)
	g.RenderChildren(sv)
	rs.PopXFormLock()
}

// ApplyXForm applies the given 2D transform to the geometry of this node
// each node must define this for itself
func (g *Rect) ApplyXForm(sv *SVG, xf mat32.Mat2) {
	rot := xf.ExtractRot()
	if rot != 0 || !g.Paint.XForm.IsIdentity() {
		g.Paint.XForm = g.Paint.XForm.Mul(xf)
		g.SetProp("transform", g.Paint.XForm.String())
	} else {
		g.Pos = xf.MulVec2AsPt(g.Pos)
		g.Size = xf.MulVec2AsVec(g.Size)
		g.GradientApplyXForm(sv, xf)
	}
}

// ApplyDeltaXForm applies the given 2D delta transforms to the geometry of this node
// relative to given point.  Trans translation and point are in top-level coordinates,
// so must be transformed into local coords first.
// Point is upper left corner of selection box that anchors the translation and scaling,
// and for rotation it is the center point around which to rotate
func (g *Rect) ApplyDeltaXForm(sv *SVG, trans mat32.Vec2, scale mat32.Vec2, rot float32, pt mat32.Vec2) {
	crot := g.Paint.XForm.ExtractRot()
	if rot != 0 || crot != 0 {
		xf, lpt := g.DeltaXForm(trans, scale, rot, pt, false) // exclude self
		g.Paint.XForm = g.Paint.XForm.MulCtr(xf, lpt)
		g.SetProp("transform", g.Paint.XForm.String())
	} else {
		xf, lpt := g.DeltaXForm(trans, scale, rot, pt, true) // include self
		g.Pos = xf.MulVec2AsPtCtr(g.Pos, lpt)
		g.Size = xf.MulVec2AsVec(g.Size)
		g.GradientApplyXFormPt(sv, xf, lpt)
	}
}

// WriteGeom writes the geometry of the node to a slice of floating point numbers
// the length and ordering of which is specific to each node type.
// Slice must be passed and will be resized if not the correct length.
func (g *Rect) WriteGeom(sv *SVG, dat *[]float32) {
	SetFloat32SliceLen(dat, 4+6)
	(*dat)[0] = g.Pos.X
	(*dat)[1] = g.Pos.Y
	(*dat)[2] = g.Size.X
	(*dat)[3] = g.Size.Y
	g.WriteXForm(*dat, 4)
	g.GradientWritePts(sv, dat)
}

// ReadGeom reads the geometry of the node from a slice of floating point numbers
// the length and ordering of which is specific to each node type.
func (g *Rect) ReadGeom(sv *SVG, dat []float32) {
	g.Pos.X = dat[0]
	g.Pos.Y = dat[1]
	g.Size.X = dat[2]
	g.Size.Y = dat[3]
	g.ReadXForm(dat, 4)
	g.GradientReadPts(sv, dat)
}
