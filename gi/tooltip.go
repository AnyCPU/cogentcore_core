// Copyright (c) 2023, The GoKi Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gi

import (
	"image"

	"goki.dev/colors"
	"goki.dev/girl/states"
	"goki.dev/girl/styles"
	"goki.dev/girl/units"
)

// TODO: rich tooltips

// NewTooltipFromScene returns a new Tooltip stage with given scene contents,
// in connection with given widget (which provides key context).
// Make further configuration choices using Set* methods, which
// can be chained directly after the New call.
// Use an appropriate Run call at the end to start the Stage running.
func NewTooltipFromScene(sc *Scene, ctx Widget) *Stage {
	return NewPopupStage(TooltipStage, sc, ctx)
}

// NewTooltip returns a new tooltip stage displaying the tooltip text
// for the given widget based on the widget's position and size.
func NewTooltip(w Widget) *Stage {
	return NewTooltipText(w, w.AsWidget().Tooltip)
}

// NewTooltipText returns a new tooltip stage displaying the given tooltip text
// for the given widget based on the widget's position and size.
func NewTooltipText(w Widget, tooltip string) *Stage {
	wb := w.AsWidget()
	bb := wb.WinBBox()
	return NewTooltipTextAt(w, tooltip, bb.Min, bb.Size())
}

// NewTooltipTextAt returns a new tooltip stage displaying the given tooltip text
// for the given widget at the given position with the given size.
func NewTooltipTextAt(w Widget, tooltip string, pos, sz image.Point) *Stage {
	return NewTooltipFromScene(NewTooltipScene(w, tooltip, pos, sz), w)
}

// NewTooltipScene returns a new tooltip scene for the given widget with the
// given tooltip based on the given context position and context size.
func NewTooltipScene(w Widget, tooltip string, pos, sz image.Point) *Scene {
	sc := NewScene(w.Name() + "-tooltip")
	// tooltip positioning uses the original scene geom as the context values
	sc.SceneGeom.Pos = pos
	sc.SceneGeom.Size = sz
	sc.Style(func(s *styles.Style) {
		s.Border.Radius = styles.BorderRadiusExtraSmall
		s.Grow.Set(1, 1)
		s.Max.X.Em(20)
		s.Overflow.Set(styles.OverflowVisible) // key for avoiding sizing errors when re-rendering with small pref size
		s.Padding.Set(units.Dp(8))
		s.Background = colors.C(colors.Scheme.InverseSurface)
		s.Color = colors.Scheme.InverseOnSurface
		s.BoxShadow = styles.BoxShadow1()
	})
	NewLabel(sc, "text").SetType(LabelBodyMedium).SetText(tooltip).
		Style(func(s *styles.Style) {
			s.Grow.Set(1, 0)
			s.Text.WhiteSpace = styles.WhiteSpaceNormal
			if s.Is(states.Selected) {
				s.Color = colors.Scheme.Select.OnContainer
			}
		})
	return sc
}
