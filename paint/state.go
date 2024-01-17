// Copyright (c) 2018, The Goki Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package paint

import (
	"image"
	"log/slog"
	"sync"

	"cogentcore.org/core/mat32"
	"cogentcore.org/core/raster"
	"cogentcore.org/core/scan"
)

// The State holds all the current rendering state information used
// while painting -- a viewport just has one of these
type State struct {

	// current transform
	CurTransform mat32.Mat2

	// current path
	Path raster.Path

	// rasterizer -- stroke / fill rendering engine from raster
	Raster *raster.Dasher

	// scanner for scanx
	Scanner *scan.Scanner

	// spanner for scanx
	ImgSpanner *scan.ImgSpanner

	// starting point, for close path
	Start mat32.Vec2

	// current point
	Current mat32.Vec2

	// is current point current?
	HasCurrent bool

	// pointer to image to render into
	Image *image.RGBA

	// current mask
	Mask *image.Alpha

	// boundaries to restrict drawing to -- much faster than clip mask for basic square region exclusion -- used for restricting drawing
	Bounds image.Rectangle

	// bounding box of last object rendered -- computed by renderer during Fill or Stroke, grabbed by SVG objects
	LastRenderBBox image.Rectangle

	// stack of transforms
	TransformStack []mat32.Mat2

	// stack of bounds -- every render starts with a push onto this stack, and finishes with a pop
	BoundsStack []image.Rectangle

	// stack of clips, if needed
	ClipStack []*image.Alpha

	// mutex for overall rendering
	RenderMu sync.Mutex

	// mutex for final raster rendering -- only one at a time
	RasterMu sync.Mutex
}

// Init initializes State -- must be called whenever image size changes
func (rs *State) Init(width, height int, img *image.RGBA) {
	rs.CurTransform = mat32.Identity2()
	rs.Image = img
	// to use the golang.org/x/image/vector scanner, do this:
	// rs.Scanner = raster.NewScannerGV(width, height, img, img.Bounds())
	// and cut out painter:
	/*
		painter := scanFT.NewRGBAPainter(img)
		rs.Scanner = scanFT.NewScannerFT(width, height, painter)
	*/
	/*
		rs.CompSpanner = &scanx.CompressSpanner{}
		rs.CompSpanner.SetBounds(img.Bounds())
	*/
	rs.ImgSpanner = scan.NewImgSpanner(img)
	rs.Scanner = scan.NewScanner(rs.ImgSpanner, width, height)
	// rs.Scanner = scanx.NewScanner(rs.CompSpanner, width, height)
	rs.Raster = raster.NewDasher(width, height, rs.Scanner)
}

// PushTransform pushes current transform onto stack and apply new transform on top of it
// must protect within render mutex lock (see Lock version)
func (rs *State) PushTransform(xf mat32.Mat2) {
	if rs.TransformStack == nil {
		rs.TransformStack = make([]mat32.Mat2, 0)
	}
	rs.TransformStack = append(rs.TransformStack, rs.CurTransform)
	rs.CurTransform = xf.Mul(rs.CurTransform)
}

// PushTransformLock pushes current transform onto stack and apply new transform on top of it
// protects within render mutex lock
func (rs *State) PushTransformLock(xf mat32.Mat2) {
	rs.RenderMu.Lock()
	rs.PushTransform(xf)
	rs.RenderMu.Unlock()
}

// PopTransform pops transform off the stack and set to current transform
// must protect within render mutex lock (see Lock version)
func (rs *State) PopTransform() {
	sz := len(rs.TransformStack)
	if sz == 0 {
		slog.Error("programmer error: paint.State.PopTransform: stack is empty")
		rs.CurTransform = mat32.Identity2()
		return
	}
	rs.CurTransform = rs.TransformStack[sz-1]
	rs.TransformStack = rs.TransformStack[:sz-1]
}

// PopTransformLock pops transform off the stack and set to current transform
// protects within render mutex lock (see Lock version)
func (rs *State) PopTransformLock() {
	rs.RenderMu.Lock()
	rs.PopTransform()
	rs.RenderMu.Unlock()
}

// PushBounds pushes current bounds onto stack and set new bounds.
// this is the essential first step in rendering!
// any further actual rendering should always be surrounded
// by Lock() / Unlock() calls
func (rs *State) PushBounds(b image.Rectangle) {
	rs.RenderMu.Lock()
	defer rs.RenderMu.Unlock()

	if rs.BoundsStack == nil {
		rs.BoundsStack = make([]image.Rectangle, 0, 100)
	}
	if rs.Bounds.Empty() { // note: method name should be IsEmpty!
		rs.Bounds = rs.Image.Bounds()
	}
	rs.BoundsStack = append(rs.BoundsStack, rs.Bounds)
	// note: this does not fix the ghost trace from rendering..
	// bp1 := image.Rectangle{Min: image.Point{X: b.Min.X - 1, Y: b.Min.Y - 1}, Max: image.Point{X: b.Max.X + 1, Y: b.Max.Y + 1}}
	rs.Bounds = b
}

// Lock locks the render mutex -- must lock prior to rendering!
func (rs *State) Lock() {
	rs.RenderMu.Lock()
}

// Unlock unlocks the render mutex, locked with PushBounds --
// call this prior to children rendering etc.
func (rs *State) Unlock() {
	rs.RenderMu.Unlock()
}

// PopBounds pops bounds off the stack and set to current bounds
// must be equally balanced with corresponding PushBounds
func (rs *State) PopBounds() {
	rs.RenderMu.Lock()
	defer rs.RenderMu.Unlock()

	sz := len(rs.BoundsStack)
	if sz == 0 {
		slog.Error("programmer error: paint.State.PopBounds: stack is empty")
		rs.Bounds = rs.Image.Bounds()
		return
	}
	rs.Bounds = rs.BoundsStack[sz-1]
	rs.BoundsStack = rs.BoundsStack[:sz-1]
}

// PushClip pushes current Mask onto the clip stack
func (rs *State) PushClip() {
	if rs.Mask == nil {
		return
	}
	if rs.ClipStack == nil {
		rs.ClipStack = make([]*image.Alpha, 0, 10)
	}
	rs.ClipStack = append(rs.ClipStack, rs.Mask)
}

// PopClip pops Mask off the clip stack and set to current mask
func (rs *State) PopClip() {
	sz := len(rs.ClipStack)
	if sz == 0 {
		slog.Error("programmer error: paint.State.PopClip: stack is empty")
		rs.Mask = nil // implied
		return
	}
	rs.Mask = rs.ClipStack[sz-1]
	rs.ClipStack[sz-1] = nil
	rs.ClipStack = rs.ClipStack[:sz-1]
}
