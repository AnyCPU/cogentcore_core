// Copyright (c) 2023, The GoKi Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package colors

import (
	"image"
	"image/color"

	"goki.dev/mat32/v2"
)

// Full represents a fully specified color that can either be a solid color or
// a gradient. If Gradient is nil, it is a solid color; otherwise, it is a gradient.
// Solid should typically be set using the [Full.SetSolid] method to
// ensure that Gradient is nil and thus Solid will be taken into account.
type Full struct {
	Gradient *Gradient
	Solid    color.RGBA
}

// SolidFull returns a new [Full] from the given solid color.
func SolidFull(solid color.Color) Full {
	return Full{Solid: AsRGBA(solid)}
}

// GradientFull returns a new [Full] from the given gradient color.
func GradientFull(gradient *Gradient) Full {
	return Full{Gradient: gradient}
}

// IsNil returns whether the color is nil, checking both the gradient
// and the solid color.
func (f *Full) IsNil() bool {
	return f.Gradient == nil && IsNil(f.Solid)
}

// SolidOrNil returns the solid color if it is not non-nil, or nil otherwise.
// It is should be used by consumers that explicitly handle nil colors.
func (f *Full) SolidOrNil() color.Color {
	if IsNil(f.Solid) {
		return nil
	}
	return f.Solid
}

// SetSolid sets the color to the given solid [color.Color],
// also setting the gradient to nil.
func (f *Full) SetSolid(solid color.Color) {
	f.Solid = AsRGBA(solid)
	f.Gradient = nil
}

// SetSolid sets the color to the solid color with the given name,
// also setting the gradient to nil.
func (f *Full) SetName(name string) error {
	s, err := FromName(name)
	if err != nil {
		return err
	}
	f.Solid = s
	f.Gradient = nil
	return nil
}

// CopyFrom copies from the given full color, making new copies
// of the gradient stops instead of re-using pointers
func (f *Full) CopyFrom(cp Full) {
	f.Solid = cp.Solid
	if f.Gradient == nil && cp.Gradient == nil {
		return
	}
	if cp.Gradient == nil {
		f.Gradient = nil
		return
	}
	if f.Gradient == nil {
		f.Gradient = &Gradient{}
	}
	f.Gradient.CopyFrom(cp.Gradient)
}

// RenderColor returns the [Render] color for rendering, applying the given opacity and bounds.
func (f *Full) RenderColor(opacity float32, bounds image.Rectangle, transform mat32.Mat2) Render {
	if f.Gradient == nil {
		return SolidRender(ApplyOpacity(f.Solid, opacity))
	}
	return f.Gradient.RenderColor(opacity, bounds, transform)
}

// SetAny sets the color from the given value of any type in the given Context.
// It handles values of types [color.Color], [*Gradient], and string. If no Context
// is provided, it uses [BaseContext] with [Transparent].
func (f *Full) SetAny(val any, ctx ...Context) error {
	switch valv := val.(type) {
	case *Full:
		*f = *valv
	case Full:
		*f = valv
	case color.Color:
		f.Solid = AsRGBA(valv)
	case *Gradient:
		*f.Gradient = *valv
	case Gradient:
		*f.Gradient = valv
	case string:
		f.SetString(valv, ctx...)
	}
	return nil
}
