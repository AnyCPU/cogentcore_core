// Copyright (c) 2018, The GoKi Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package svg

import (
	"goki.dev/colors"
	"goki.dev/girl/styles"
	"goki.dev/girl/units"
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

func (g *Rect) SVGName() string { return "rect" }

func (g *Rect) OnInit() {
	g.Size.Set(1, 1)
}

func (g *Rect) CopyFieldsFrom(frm any) {
	fr := frm.(*Rect)
	g.NodeBase.CopyFieldsFrom(&fr.NodeBase)
	g.Pos = fr.Pos
	g.Size = fr.Size
	g.Radius = fr.Radius
}

func (g *Rect) LocalBBox() mat32.Box2 {
	bb := mat32.Box2{}
	hlw := 0.5 * g.LocalLineWidth()
	bb.Min = g.Pos.SubScalar(hlw)
	bb.Max = g.Pos.Add(g.Size).AddScalar(hlw)
	return bb
}

func (g *Rect) Render(sv *SVG) {
	vis, pc := g.PushTransform(sv)
	if !vis {
		return
	}
	pc.Lock()
	// TODO: figure out a better way to do this
	bs := styles.Border{}
	bs.Style.Set(styles.BorderSolid)
	bs.Width.Set(pc.StrokeStyle.Width)
	bs.Color.Set(colors.ToUniform(pc.StrokeStyle.Color))
	bs.Radius.Set(units.Px(g.Radius.X))
	if g.Radius.X == 0 && g.Radius.Y == 0 {
		pc.DrawRectangle(g.Pos.X, g.Pos.Y, g.Size.X, g.Size.Y)
	} else {
		// todo: only supports 1 radius right now -- easy to add another
		// SidesTODO: also support different radii for each corner
		pc.DrawRoundedRectangle(g.Pos.X, g.Pos.Y, g.Size.X, g.Size.Y, styles.NewSideFloats(g.Radius.X))
	}
	pc.FillStrokeClear()
	pc.Unlock()
	g.BBoxes(sv)
	g.RenderChildren(sv)
	pc.PopTransformLock()
}

// ApplyTransform applies the given 2D transform to the geometry of this node
// each node must define this for itself
func (g *Rect) ApplyTransform(sv *SVG, xf mat32.Mat2) {
	rot := xf.ExtractRot()
	if rot != 0 || !g.Paint.Transform.IsIdentity() {
		g.Paint.Transform = g.Paint.Transform.Mul(xf)
		g.SetProp("transform", g.Paint.Transform.String())
	} else {
		g.Pos = xf.MulVec2AsPt(g.Pos)
		g.Size = xf.MulVec2AsVec(g.Size)
		g.GradientApplyTransform(sv, xf)
	}
}

// ApplyDeltaTransform applies the given 2D delta transforms to the geometry of this node
// relative to given point.  Trans translation and point are in top-level coordinates,
// so must be transformed into local coords first.
// Point is upper left corner of selection box that anchors the translation and scaling,
// and for rotation it is the center point around which to rotate
func (g *Rect) ApplyDeltaTransform(sv *SVG, trans mat32.Vec2, scale mat32.Vec2, rot float32, pt mat32.Vec2) {
	crot := g.Paint.Transform.ExtractRot()
	if rot != 0 || crot != 0 {
		xf, lpt := g.DeltaTransform(trans, scale, rot, pt, false) // exclude self
		g.Paint.Transform = g.Paint.Transform.MulCtr(xf, lpt)
		g.SetProp("transform", g.Paint.Transform.String())
	} else {
		xf, lpt := g.DeltaTransform(trans, scale, rot, pt, true) // include self
		g.Pos = xf.MulVec2AsPtCtr(g.Pos, lpt)
		g.Size = xf.MulVec2AsVec(g.Size)
		g.GradientApplyTransformPt(sv, xf, lpt)
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
	g.WriteTransform(*dat, 4)
	g.GradientWritePts(sv, dat)
}

// ReadGeom reads the geometry of the node from a slice of floating point numbers
// the length and ordering of which is specific to each node type.
func (g *Rect) ReadGeom(sv *SVG, dat []float32) {
	g.Pos.X = dat[0]
	g.Pos.Y = dat[1]
	g.Size.X = dat[2]
	g.Size.Y = dat[3]
	g.ReadTransform(dat, 4)
	g.GradientReadPts(sv, dat)
}
