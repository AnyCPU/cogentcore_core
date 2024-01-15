// Copyright (c) 2023, The Goki Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package matcolor

//go:generate goki generate

import "image/color"

// Scheme contains the colors for one color scheme (ex: light or dark).
// To generate a scheme, use [NewScheme].
type Scheme struct {

	// Primary is the primary color applied to important elements
	Primary Accent

	// Secondary is the secondary color applied to less important elements
	Secondary Accent

	// Tertiary is the tertiary color applied as an accent to highlight elements and create contrast between other colors
	Tertiary Accent

	// Select is the selection color applied to selected or highlighted elements and text
	Select Accent

	// Error is the error color applied to elements that indicate an error or danger
	Error Accent

	// Success is the color applied to elements that indicate success
	Success Accent

	// Warn is the color applied to elements that indicate a warning
	Warn Accent

	// an optional map of custom accent colors
	Custom map[string]Accent

	// SurfaceDim is the color applied to elements that will always have the dimmest surface color (see Surface for more information)
	SurfaceDim color.RGBA

	// Surface is the color applied to contained areas, like the background of an app
	Surface color.RGBA

	// SurfaceBright is the color applied to elements that will always have the brightest surface color (see Surface for more information)
	SurfaceBright color.RGBA

	// SurfaceContainerLowest is the color applied to surface container elements that have the lowest emphasis (see SurfaceContainer for more information)
	SurfaceContainerLowest color.RGBA

	// SurfaceContainerLow is the color applied to surface container elements that have lower emphasis (see SurfaceContainer for more information)
	SurfaceContainerLow color.RGBA

	// SurfaceContainer is the color applied to container elements that contrast elements with the surface color
	SurfaceContainer color.RGBA

	// SurfaceContainerHigh is the color applied to surface container elements that have higher emphasis (see SurfaceContainer for more information)
	SurfaceContainerHigh color.RGBA

	// SurfaceContainerHighest is the color applied to surface container elements that have the highest emphasis (see SurfaceContainer for more information)
	SurfaceContainerHighest color.RGBA

	// SurfaceVariant is the color applied to contained areas that contrast standard Surface elements
	SurfaceVariant color.RGBA

	// OnSurface is the color applied to content on top of Surface elements
	OnSurface color.RGBA

	// OnSurfaceVariant is the color applied to content on top of SurfaceVariant elements
	OnSurfaceVariant color.RGBA

	// InverseSurface is the color applied to elements to make them the reverse color of the surrounding elements and create a contrasting effect
	InverseSurface color.RGBA

	// InverseOnSurface is the color applied to content on top of InverseSurface
	InverseOnSurface color.RGBA

	// InversePrimary is the color applied to interactive elements on top of InverseSurface
	InversePrimary color.RGBA

	// Background is the color applied to the background of the app and other low-emphasis areas
	Background color.RGBA

	// OnBackground is the color applied to content on top of Background
	OnBackground color.RGBA

	// Outline is the color applied to borders to create emphasized boundaries that need to have sufficient contrast
	Outline color.RGBA

	// OutlineVariant is the color applied to create decorative boundaries
	OutlineVariant color.RGBA

	// Shadow is the color applied to shadows
	Shadow color.RGBA

	// SurfaceTint is the color applied to tint surfaces
	SurfaceTint color.RGBA

	// Scrim is the color applied to scrims (semi-transparent overlays)
	Scrim color.RGBA

	// // PrimaryFixed is a primary fill color that stays the same regardless of color scheme type (light/dark)
	// PrimaryFixed color.RGBA `desc:"PrimaryFixed is a primary fill color that stays the same regardless of color scheme type (light/dark)"`

	// // PrimaryFixedDim is a higher-emphasis, dimmer primary fill color that stays the same regardless of color scheme type (light/dark)
	// PrimaryFixedDim color.RGBA `desc:"PrimaryFixedDim is a higher-emphasis, dimmer primary fill color that stays the same regardless of color scheme type (light/dark)"`

	// // OnPrimaryFixed is the color applied to high-emphasis content on top of PrimaryFixed
	// OnPrimaryFixed color.RGBA `desc:"OnPrimaryFixed is the color applied to high-emphasis content on top of PrimaryFixed"`

	// // OnPrimaryFixedVariant is the color applied to low-emphasis content on top of PrimaryFixed
	// OnPrimaryFixedVariant color.RGBA `desc:"OnPrimaryFixedVariant is the color applied to low-emphasis content on top of PrimaryFixed"`

	// // SecondaryFixed is a secondary fill color that stays the same regardless of color scheme type (light/dark)
	// SecondaryFixed color.RGBA `desc:"SecondaryFixed is a secondary fill color that stays the same regardless of color scheme type (light/dark)"`

	// // SecondaryFixedDim is a higher-emphasis, dimmer secondary fill color that stays the same regardless of color scheme type (light/dark)
	// SecondaryFixedDim color.RGBA `desc:"SecondaryFixedDim is a higher-emphasis, dimmer secondary fill color that stays the same regardless of color scheme type (light/dark)"`

	// // OnSecondaryFixed is the color applied to high-emphasis content on top of SecondaryFixed
	// OnSecondaryFixed color.RGBA `desc:"OnSecondaryFixed is the color applied to high-emphasis content on top of SecondaryFixed"`

	// // OnSecondaryFixedVariant is the color applied to low-emphasis content on top of SecondaryFixed
	// OnSecondaryFixedVariant color.RGBA `desc:"OnSecondaryFixedVariant is the color applied to low-emphasis content on top of SecondaryFixed"`

	// // TertiaryFixed is a tertiary fill color that stays the same regardless of color scheme type (light/dark)
	// TertiaryFixed color.RGBA `desc:"TertiaryFixed is a tertiary fill color that stays the same regardless of color scheme type (light/dark)"`

	// // TertiaryFixedDim is a higher-emphasis, dimmer tertiary fill color that stays the same regardless of color scheme type (light/dark)
	// TertiaryFixedDim color.RGBA `desc:"TertiaryFixedDim is a higher-emphasis, dimmer tertiary fill color that stays the same regardless of color scheme type (light/dark)"`

	// // OnTertiaryFixed is the color applied to high-emphasis content on top of TertiaryFixed
	// OnTertiaryFixed color.RGBA `desc:"OnTertiaryFixed is the color applied to high-emphasis content on top of TertiaryFixed"`

	// // OnTertiaryFixedVariant is the color applied to low-emphasis content on top of TertiaryFixed
	// OnTertiaryFixedVariant color.RGBA `desc:"OnTertiaryFixedVariant is the color applied to low-emphasis content on top of TertiaryFixed"`
}

// NewLightScheme returns a new light-themed [Scheme]
// based on the given [Palette].
func NewLightScheme(p *Palette) Scheme {
	s := Scheme{
		Primary:   NewAccentLight(p.Primary),
		Secondary: NewAccentLight(p.Secondary),
		Tertiary:  NewAccentLight(p.Tertiary),
		Select:    NewAccentLight(p.Select),
		Error:     NewAccentLight(p.Error),
		Success:   NewAccentLight(p.Success),
		Warn:      NewAccentLight(p.Warn),
		Custom:    map[string]Accent{},

		SurfaceDim:    p.Neutral.AbsTone(87),
		Surface:       p.Neutral.AbsTone(98),
		SurfaceBright: p.Neutral.AbsTone(98),

		SurfaceContainerLowest:  p.Neutral.AbsTone(100),
		SurfaceContainerLow:     p.Neutral.AbsTone(96),
		SurfaceContainer:        p.Neutral.AbsTone(94),
		SurfaceContainerHigh:    p.Neutral.AbsTone(92),
		SurfaceContainerHighest: p.Neutral.AbsTone(90),

		SurfaceVariant:   p.NeutralVariant.AbsTone(90),
		OnSurface:        p.NeutralVariant.AbsTone(10),
		OnSurfaceVariant: p.NeutralVariant.AbsTone(30),

		InverseSurface:   p.Neutral.AbsTone(20),
		InverseOnSurface: p.Neutral.AbsTone(95),
		InversePrimary:   p.Primary.AbsTone(80),

		Background:   p.Neutral.AbsTone(98),
		OnBackground: p.Neutral.AbsTone(10),

		Outline:        p.NeutralVariant.AbsTone(50),
		OutlineVariant: p.NeutralVariant.AbsTone(80),

		Shadow:      p.Neutral.AbsTone(0),
		SurfaceTint: p.Primary.AbsTone(40),
		Scrim:       p.Neutral.AbsTone(0),
	}
	for nm, c := range p.Custom {
		s.Custom[nm] = NewAccentLight(c)
	}
	return s
	// TODO: maybe fixed colors
}

// NewDarkScheme returns a new dark-themed [Scheme]
// based on the given [Palette].
func NewDarkScheme(p *Palette) Scheme {
	s := Scheme{
		Primary:   NewAccentDark(p.Primary),
		Secondary: NewAccentDark(p.Secondary),
		Tertiary:  NewAccentDark(p.Tertiary),
		Select:    NewAccentDark(p.Select),
		Error:     NewAccentDark(p.Error),
		Success:   NewAccentDark(p.Success),
		Warn:      NewAccentDark(p.Warn),
		Custom:    map[string]Accent{},

		SurfaceDim:    p.Neutral.AbsTone(6),
		Surface:       p.Neutral.AbsTone(6),
		SurfaceBright: p.Neutral.AbsTone(24),

		SurfaceContainerLowest:  p.Neutral.AbsTone(4),
		SurfaceContainerLow:     p.Neutral.AbsTone(10),
		SurfaceContainer:        p.Neutral.AbsTone(12),
		SurfaceContainerHigh:    p.Neutral.AbsTone(17),
		SurfaceContainerHighest: p.Neutral.AbsTone(22),

		SurfaceVariant:   p.NeutralVariant.AbsTone(30),
		OnSurface:        p.NeutralVariant.AbsTone(90),
		OnSurfaceVariant: p.NeutralVariant.AbsTone(80),

		InverseSurface:   p.Neutral.AbsTone(90),
		InverseOnSurface: p.Neutral.AbsTone(20),
		InversePrimary:   p.Primary.AbsTone(40),

		Background:   p.Neutral.AbsTone(6),
		OnBackground: p.Neutral.AbsTone(90),

		Outline:        p.NeutralVariant.AbsTone(60),
		OutlineVariant: p.NeutralVariant.AbsTone(30),

		Shadow:      p.Neutral.AbsTone(0),
		SurfaceTint: p.Primary.AbsTone(80),
		Scrim:       p.Neutral.AbsTone(0),
	}
	for nm, c := range p.Custom {
		s.Custom[nm] = NewAccentDark(c)
	}
	return s
	// TODO: custom and fixed colors
}

// SchemeIsDark is whether the currently active color scheme
// is a dark-themed or light-themed color scheme. In almost
// all cases, it should be set via [goki.dev/colors.SetScheme],
// not directly.
var SchemeIsDark = false
