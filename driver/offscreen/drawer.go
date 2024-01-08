// Copyright 2023 The Goki Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package offscreen

import (
	"image"
	"image/color"
	"image/draw"

	"goki.dev/goosi"
	"goki.dev/mat32/v2"
)

// Drawer is the implementation of [goosi.Drawer] for the offscreen platform
type Drawer struct {
	// MaxTxts is the max number of textures
	MaxTxts int

	// Image is the target render image
	Image *image.RGBA

	// Images is a stack of images indexed by render scene index and then layer number
	Images [][]*image.RGBA
}

// SetMaxTextures updates the max number of textures for drawing
// Must call this prior to doing any allocation of images.
func (dw *Drawer) SetMaxTextures(maxTextures int) {
	dw.MaxTxts = maxTextures
}

// MaxTextures returns the max number of textures for drawing
func (dw *Drawer) MaxTextures() int {
	return dw.MaxTxts
}

// DestBounds returns the bounds of the render destination
func (dw *Drawer) DestBounds() image.Rectangle {
	return TheApp.Scrn.Geometry
}

// SetGoImage sets given Go image as a drawing source to given image index,
// and layer, used in subsequent Draw methods.
// A standard Go image is rendered upright on a standard surface.
// Set flipY to true to flip.
func (dw *Drawer) SetGoImage(idx, layer int, img image.Image, flipY bool) {
	for len(dw.Images) <= idx {
		dw.Images = append(dw.Images, nil)
	}
	imgs := &dw.Images[idx]
	for len(*imgs) <= layer {
		*imgs = append(*imgs, nil)
	}
	(*imgs)[layer] = img.(*image.RGBA)
}

// ConfigImageDefaultFormat configures the draw image at the given index
// to fit the default image format specified by the given width, height,
// and number of layers.
func (dw *Drawer) ConfigImageDefaultFormat(idx int, width int, height int, layers int) {
	dw.Image = image.NewRGBA(image.Rect(0, 0, width, height))
}

// SyncImages must be called after images have been updated, to sync
// memory up to the GPU.
func (dw *Drawer) SyncImages() {}

// Copy copies texture at given index and layer to render target.
// dp is the destination point,
// sr is the source region (set to image.ZR zero rect for all),
// op is the drawing operation: Src = copy source directly (blit),
// Over = alpha blend with existing
// flipY = flipY axis when drawing this image
func (dw *Drawer) Copy(idx, layer int, dp image.Point, sr image.Rectangle, op draw.Op, flipY bool) error {
	img := dw.Images[idx][layer]
	draw.Draw(dw.Image, image.Rectangle{dp, dp.Add(img.Rect.Size())}, img, sr.Min, op)
	return nil
}

// UseTextureSet selects the descriptor set to use --
// choose this based on the bank of 16
// texture values if number of textures > MaxTexturesPerSet.
func (dw *Drawer) UseTextureSet(descIdx int) {}

// StartDraw starts image drawing rendering process on render target
// No images can be added or set after this point.
// descIdx is the descriptor set to use -- choose this based on the bank of 16
// texture values if number of textures > MaxTexturesPerSet.
// This is a no-op on offscreen; if rendering logic is done here instead of
// EndDraw, everything is delayed by one render because Scale and Copy are
// called after StartDraw but before EndDraw, and we need them to be called
// before actually rendering the image to the capture channel.
func (dw *Drawer) StartDraw(descIdx int) {
	// no-op
}

// EndDraw ends image drawing rendering process on render target.
// This is the function that actually sends the image to the capture channel.
func (dw *Drawer) EndDraw() {
	if !goosi.NeedsCapture {
		return
	}
	goosi.CaptureImage <- dw.Image
}

// Fill fills given color to to render target.
// src2dst is the transform mapping source to destination
// coordinates (translation, scaling),
// reg is the region to fill
// op is the drawing operation: Src = copy source directly (blit),
// Over = alpha blend with existing
func (dw *Drawer) Fill(clr color.Color, src2dst mat32.Mat3, reg image.Rectangle, op draw.Op) error {
	draw.Draw(dw.Image, reg, image.NewUniform(clr), image.Point{}, op)
	return nil
}

func (dw *Drawer) Surface() any {
	// no-op
	return nil
}
