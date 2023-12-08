// Copyright (c) 2018, The GoKi Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package svg

//go:generate goki generate

import (
	"image"
	"image/color"
	"strings"
	"sync"

	"goki.dev/colors"
	"goki.dev/girl/paint"
	"goki.dev/girl/styles"
	"goki.dev/girl/units"
	"goki.dev/ki/v2"
	"goki.dev/mat32/v2"
)

// SVG is an SVG object.
type SVG struct {
	// Name is the name of the SVG -- e.g., the filename if loaded
	Name string

	// the title of the svg
	Title string `xml:"title"`

	// the description of the svg
	Desc string `xml:"desc"`

	// fill the viewport with background-color
	Fill bool

	// color to fill background if Fill set
	BackgroundColor colors.Full

	// Color can be set to provide a default Fill and Stroke Color value
	Color colors.Full

	// Size is size of image, Pos is offset within any parent viewport.  Node bounding boxes are based on 0 Pos offset within Pixels image
	Geom mat32.Geom2DInt

	// physical width of the drawing, e.g., when printed -- does not affect rendering -- metadata
	PhysWidth units.Value

	// physical height of the drawing, e.g., when printed -- does not affect rendering -- metadata
	PhysHeight units.Value

	// prop: norm = install a transform that renormalizes so that the specified ViewBox exactly fits within the allocated SVG size
	Norm bool

	// prop: invert-y = when doing Norm transform, also flip the Y axis so that the smallest Y value is at the bottom of the SVG box, instead of being at the top as it is by default
	InvertY bool

	// render state for rendering
	RenderState paint.State `copy:"-" json:"-" xml:"-" view:"-"`

	// live pixels that we render into
	Pixels *image.RGBA `copy:"-" json:"-" xml:"-" view:"-"`

	// all defs defined elements go here (gradients, symbols, etc)
	Defs Group

	// root of the svg tree -- top-level viewbox and paint style here
	Root SVGNode

	// map of def names to index -- uses starting index to find element -- always updated after each search
	DefIdxs map[string]int `view:"-" json:"-" xml:"-"`

	// map of unique numeric ids for all elements -- used for allocating new unique id numbers, appended to end of elements -- see NewUniqueId, GatherIds
	UniqueIds map[int]struct{} `view:"-" json:"-" xml:"-"`

	// mutex for protecting rendering
	RenderMu sync.Mutex `view:"-" json:"-" xml:"-"`
}

// NewSVG creates a SVG with Pixels Image of the specified width and height
func NewSVG(width, height int) *SVG {
	sv := &SVG{}
	sv.Config(width, height)
	return sv
}

// Config configures the SVG, setting image to given size
// and initializing all relevant fields.
func (sv *SVG) Config(width, height int) {
	sz := image.Point{width, height}
	sv.Geom.Size = sz
	sv.BackgroundColor.SetSolid(colors.White)
	sv.Pixels = image.NewRGBA(image.Rectangle{Max: sz})
	sv.RenderState.Init(width, height, sv.Pixels)
	sv.Root.InitName(&sv.Root, "svg")
	sv.Defs.InitName(&sv.Defs, "defs")
}

// Resize resizes the viewport, creating a new image -- updates Geom Size
func (sv *SVG) Resize(nwsz image.Point) {
	if nwsz.X == 0 || nwsz.Y == 0 {
		return
	}
	if sv.Root.Ths == nil {
		sv.Config(nwsz.X, nwsz.Y)
		return
	}
	if sv.Pixels != nil {
		ib := sv.Pixels.Bounds().Size()
		if ib == nwsz {
			sv.Geom.Size = nwsz // make sure
			return              // already good
		}
	}
	if sv.Pixels != nil {
		sv.Pixels = nil
	}
	sv.Pixels = image.NewRGBA(image.Rectangle{Max: nwsz})
	sv.RenderState.Init(nwsz.X, nwsz.Y, sv.Pixels)
	sv.Geom.Size = nwsz // make sure
}

func (sv *SVG) CopyFrom(fr *SVG) {
	sv.Title = fr.Title
	sv.Desc = fr.Desc
	sv.Fill = fr.Fill
	sv.BackgroundColor = fr.BackgroundColor
	sv.Geom = fr.Geom
	sv.Norm = fr.Norm
	sv.InvertY = fr.InvertY
	sv.Defs.CopyFrom(&fr.Defs)
	sv.Root.CopyFrom(&fr.Root)
	sv.UniqueIds = nil
}

// DeleteAll deletes any existing elements in this svg
func (sv *SVG) DeleteAll() {
	if sv.Root.This() == nil {
		return
	}
	updt := sv.Root.UpdateStart() // don't really need update logic here
	sv.Root.Paint.Defaults()
	sv.Root.DeleteChildren(ki.DestroyKids)
	sv.Defs.DeleteChildren(ki.DestroyKids)
	sv.Root.UpdateEnd(updt)
}

// Base returns the current Color activated in the context.
// Color has support for special color names that are relative to
// this current color.
func (sv *SVG) Base() color.RGBA {
	return sv.BackgroundColor.Solid
}

// FullByURL finds a Node by an element name (URL-like path), and
// attempts to convert it to a Gradient -- if successful, returns ColorSpec on that.
// Used for colorspec styling based on url() value.
func (sv *SVG) FullByURL(url string) *colors.Full {
	if sv == nil {
		return nil
	}
	val := url[4:]
	val = strings.TrimPrefix(strings.TrimSuffix(val, ")"), "#")
	def := sv.FindDefByName(val)
	if def != nil {
		if grad, ok := def.(*Gradient); ok {
			return &grad.Grad
		}
	}
	ne := sv.FindNamedElement(val)
	if grad, ok := ne.(*Gradient); ok {
		return &grad.Grad
	}
	return nil
}

func (sv *SVG) Style() {
	// set the Defs flags
	sv.Defs.WalkPre(func(k ki.Ki) bool {
		ni := k.(Node)
		if ni == nil || ni.Is(ki.Deleted) || ni.Is(ki.Destroyed) {
			return ki.Break
		}
		ni.SetFlag(true, IsDef)
		return ki.Continue
	})

	sv.Root.Paint.Defaults()
	if !sv.Color.IsNil() {
		sv.Root.SetProp("stroke", sv.Color)
		sv.Root.SetProp("fill", sv.Color)
	}
	sv.SetUnitContext(&sv.Root.Paint, mat32.Vec2{}, mat32.Vec2{})

	sv.Root.WalkPre(func(k ki.Ki) bool {
		ni := k.(Node)
		if ni == nil || ni.Is(ki.Deleted) || ni.Is(ki.Destroyed) {
			return ki.Break
		}
		ni.Style(sv)
		return ki.Continue
	})
}

func (sv *SVG) Render() {
	sv.RenderMu.Lock()
	defer sv.RenderMu.Unlock()

	sv.Style()

	rs := &sv.RenderState
	rs.PushBounds(sv.Pixels.Bounds())
	if sv.Fill {
		sv.FillViewport()
	}
	if sv.Norm {
		sv.SetNormTransform()
	}
	sv.Root.Render(sv)
	// rs.PushTransform(sv.Root.Paint.Transform)
	// for _, kid := range sv.Root.Kids {
	// 	ni := kid.(Node)
	// 	ni.Render()
	// }
	// rs.PopTransform()
	rs.PopBounds()
}

func (sv *SVG) FillViewport() {
	pc := &paint.Context{&sv.RenderState, &sv.Root.Paint}
	pc.Lock()
	pc.FillBox(mat32.Vec2Zero, mat32.NewVec2FmPoint(sv.Geom.Size), &sv.BackgroundColor)
	pc.Unlock()
}

// SetNormTransform sets a scaling transform to make the entire viewbox to fit the viewport
func (sv *SVG) SetNormTransform() {
	pc := &sv.Root.Paint
	pc.Transform = mat32.Identity2D()
	vb := &sv.Root.ViewBox
	if vb.Size != mat32.Vec2Zero {
		// todo: deal with all the other options!
		vpsX := float32(sv.Geom.Size.X) / vb.Size.X
		vpsY := float32(sv.Geom.Size.Y) / vb.Size.Y
		if sv.InvertY {
			vpsY *= -1
		}
		pc.Transform = pc.Transform.Scale(vpsX, vpsY).Translate(-vb.Min.X, -vb.Min.Y)
		if sv.InvertY {
			pc.Transform.Y0 = -pc.Transform.Y0
		}
	}
}

// SetDPITransform sets a scaling transform to compensate for
// a given LogicalDPI factor.
// svg rendering is done within a 96 DPI context.
func (sv *SVG) SetDPITransform(logicalDPI float32) {
	pc := &sv.Root.Paint
	dpisc := logicalDPI / 96.0
	pc.Transform = mat32.Scale2D(dpisc, dpisc)
}

/*
// todo:  for gi wrapper node:
//
// func (sv *SVG) OnInit() {
// 	sv.AddStyler(func(s *styles.Style) {
// 		if par := sv.ParentWidget(); par != nil {
// 			sv.Paint.FillStyle.Color.SetColor(par.Style.Color)
// 			sv.Paint.StrokeStyle.Color.SetColor(par.Style.Color)
// 		}
// 	})
// }

// func (sv *SVG) Style2D() {
// 	if nv, err := sv.PropTry("norm"); err == nil {
// 		sv.Norm, _ = kit.ToBool(nv)
// 	}
// 	if iv, err := sv.PropTry("invert-y"); err == nil {
// 		sv.InvertY, _ = kit.ToBool(iv)
// 	}
// }

var SVGProps = ki.Props{
	ki.EnumTypeFlag: TypeSVGFlags,
	"Toolbar": ki.PropSlice{
		{"OpenXML", ki.Props{
			"label": "Open...",
			"desc":  "Open SVG XML-formatted file",
			"icon":  icons.Open,
			"Args": ki.PropSlice{
				{"File Name", ki.Props{
					"ext": ".svg",
				}},
			},
		}},
		{"SaveXML", ki.Props{
			"label": "SaveAs...",
			"desc":  "Save SVG content to an XML-formatted file.",
			"icon":  icons.SaveAs,
			"Args": ki.PropSlice{
				{"File Name", ki.Props{
					"ext": ".svg",
				}},
			},
		}},
	},
}
*/

//////////////////////////////////////////////////////////////
// 	SVGNode

// SVGNode represents the root of an SVG tree
type SVGNode struct {
	Group

	// viewbox defines the coordinate system for the drawing -- these units are mapped into the screen space allocated for the SVG during rendering
	ViewBox ViewBox
}

func (g *SVGNode) CopyFieldsFrom(frm any) {
	fr := frm.(*SVGNode)
	g.NodeBase.CopyFieldsFrom(&fr.NodeBase)
	g.ViewBox = fr.ViewBox
}

func (g *SVGNode) SVGName() string { return "svg" }

func (g *SVGNode) EnforceSVGName() bool { return false }

func (g *SVGNode) NodeBBox(sv *SVG) image.Rectangle {
	// todo: return viewbox
	return sv.Geom.SizeRect()
}

// SetUnitContext sets the unit context based on size of viewport, element,
// and parent element (from bbox) and then caches everything out in terms of raw pixel
// dots for rendering -- call at start of render
func (sv *SVG) SetUnitContext(pc *styles.Paint, el, par mat32.Vec2) {
	pc.UnContext.Defaults()
	pc.UnContext.DPI = 96 // paint (SVG) context is always 96 = 1to1
	if sv.RenderState.Image != nil {
		sz := sv.RenderState.Image.Bounds().Size()
		pc.UnContext.SetSizes(float32(sz.X), float32(sz.Y), el.X, el.Y, par.X, par.Y)
	} else {
		pc.UnContext.SetSizes(0, 0, el.X, el.Y, par.X, par.Y)
	}
	pc.FontStyle.SetUnitContext(&pc.UnContext)
	pc.ToDots()
}
