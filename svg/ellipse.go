// Copyright (c) 2018, The Goki Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package svg

import (
	"goki.dev/goki/mat32"
)

// Ellipse is a SVG ellipse
type Ellipse struct {
	NodeBase

	// position of the center of the ellipse
	Pos mat32.Vec2 `xml:"{cx,cy}" set:"-"`

	// radii of the ellipse in the horizontal, vertical axes
	Radii mat32.Vec2 `xml:"{rx,ry}"`
}

func (g *Ellipse) SVGName() string { return "ellipse" }

func (g *Ellipse) OnInit() {
	g.Radii.Set(1, 1)
}

func (g *Ellipse) CopyFieldsFrom(frm any) {
	fr := frm.(*Ellipse)
	g.NodeBase.CopyFieldsFrom(&fr.NodeBase)
	g.Pos = fr.Pos
	g.Radii = fr.Radii
}

func (g *Ellipse) SetPos(pos mat32.Vec2) *Ellipse {
	g.Pos = pos.Sub(g.Radii)
	return g
}

func (g *Ellipse) SetSize(sz mat32.Vec2) *Ellipse {
	g.Radii = sz.MulScalar(0.5)
	return g
}

func (g *Ellipse) LocalBBox() mat32.Box2 {
	bb := mat32.Box2{}
	hlw := 0.5 * g.LocalLineWidth()
	bb.Min = g.Pos.Sub(g.Radii.AddScalar(hlw))
	bb.Max = g.Pos.Add(g.Radii.AddScalar(hlw))
	return bb
}

func (g *Ellipse) Render(sv *SVG) {
	vis, pc := g.PushTransform(sv)
	if !vis {
		return
	}
	pc.Lock()
	pc.DrawEllipse(g.Pos.X, g.Pos.Y, g.Radii.X, g.Radii.Y)
	pc.FillStrokeClear()
	pc.Unlock()

	g.BBoxes(sv)
	g.RenderChildren(sv)

	pc.PopTransformLock()
}

// ApplyTransform applies the given 2D transform to the geometry of this node
// each node must define this for itself
func (g *Ellipse) ApplyTransform(sv *SVG, xf mat32.Mat2) {
	rot := xf.ExtractRot()
	if rot != 0 || !g.Paint.Transform.IsIdentity() {
		g.Paint.Transform = g.Paint.Transform.Mul(xf)
		g.SetProp("transform", g.Paint.Transform.String())
	} else {
		g.Pos = xf.MulVec2AsPt(g.Pos)
		g.Radii = xf.MulVec2AsVec(g.Radii)
		g.GradientApplyTransform(sv, xf)
	}
}

// ApplyDeltaTransform applies the given 2D delta transforms to the geometry of this node
// relative to given point.  Trans translation and point are in top-level coordinates,
// so must be transformed into local coords first.
// Point is upper left corner of selection box that anchors the translation and scaling,
// and for rotation it is the center point around which to rotate
func (g *Ellipse) ApplyDeltaTransform(sv *SVG, trans mat32.Vec2, scale mat32.Vec2, rot float32, pt mat32.Vec2) {
	crot := g.Paint.Transform.ExtractRot()
	if rot != 0 || crot != 0 {
		xf, lpt := g.DeltaTransform(trans, scale, rot, pt, false) // exclude self
		mat := g.Paint.Transform.MulCtr(xf, lpt)
		g.Paint.Transform = mat
		g.SetProp("transform", g.Paint.Transform.String())
	} else {
		xf, lpt := g.DeltaTransform(trans, scale, rot, pt, true) // include self
		g.Pos = xf.MulVec2AsPtCtr(g.Pos, lpt)
		g.Radii = xf.MulVec2AsVec(g.Radii)
		g.GradientApplyTransformPt(sv, xf, lpt)
	}
}

// WriteGeom writes the geometry of the node to a slice of floating point numbers
// the length and ordering of which is specific to each node type.
// Slice must be passed and will be resized if not the correct length.
func (g *Ellipse) WriteGeom(sv *SVG, dat *[]float32) {
	SetFloat32SliceLen(dat, 4+6)
	(*dat)[0] = g.Pos.X
	(*dat)[1] = g.Pos.Y
	(*dat)[2] = g.Radii.X
	(*dat)[3] = g.Radii.Y
	g.WriteTransform(*dat, 4)
	g.GradientWritePts(sv, dat)
}

// ReadGeom reads the geometry of the node from a slice of floating point numbers
// the length and ordering of which is specific to each node type.
func (g *Ellipse) ReadGeom(sv *SVG, dat []float32) {
	g.Pos.X = dat[0]
	g.Pos.Y = dat[1]
	g.Radii.X = dat[2]
	g.Radii.Y = dat[3]
	g.ReadTransform(dat, 4)
	g.GradientReadPts(sv, dat)
}
