// Copyright (c) 2023, The GoKi Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gi

import (
	"fmt"
	"log"

	"goki.dev/icons"
	"goki.dev/ki/v2"
)

// HasSc checks that the Sc Scene has been set.
// Called prior to using -- logs an error if not.
// todo: need slog Debug mode for this kind of thing.
func (wb *WidgetBase) HasSc() bool {
	if wb.This() == nil || wb.Sc == nil {
		log.Printf("gi.WidgetBase.ReConfig: object or scene is nil\n") // todo: slog.Debug
		return false
	}
	return true
}

// ReConfig is a convenience method for reconfiguring a widget after changes
// have been made.  In general it is more efficient to call Set* methods that
// automatically determine if Config is needed.
// The plain Config method is used during initial configuration,
// called by the Scene and caches the Sc pointer.
func (wb *WidgetBase) ReConfig() {
	if !wb.HasSc() {
		return
	}
	wi := wb.This().(Widget)
	wi.Config(wb.Sc)
	wi.ApplyStyle(wb.Sc)
}

func (wb *WidgetBase) Config(sc *Scene) {
	if wb.This() == nil {
		fmt.Println("ERROR: nil this in config")
		return
	}
	wi := wb.This().(Widget)
	updt := wi.UpdateStart()
	wb.Sc = sc
	wi.ConfigWidget(sc) // where everything actually happens
	wb.UpdateEnd(updt)
	wb.SetNeedsLayout(sc, updt)
}

func (wb *WidgetBase) ConfigWidget(sc *Scene) {
	// this must be defined for each widget type
}

// ConfigPartsIconLabel adds to config to create parts, of icon
// and label left-to right in a row, based on whether items are nil or empty
func (wb *WidgetBase) ConfigPartsIconLabel(config *ki.Config, icnm icons.Icon, txt string) (icIdx, lbIdx int) {
	icIdx = -1
	lbIdx = -1
	if icnm.IsValid() {
		icIdx = len(*config)
		config.Add(IconType, "icon")
		if txt != "" {
			config.Add(SpaceType, "space")
		}
	}
	if txt != "" {
		lbIdx = len(*config)
		config.Add(LabelType, "label")
	}
	return
}

// ConfigPartsSetIconLabel sets the icon and text values in parts, and get
// part style props, using given props if not set in object props
func (wb *WidgetBase) ConfigPartsSetIconLabel(icnm icons.Icon, txt string, icIdx, lbIdx int) {
	if icIdx >= 0 {
		ic := wb.Parts.Child(icIdx).(*Icon)
		ic.SetIcon(icnm)
	}
	if lbIdx >= 0 {
		lbl := wb.Parts.Child(lbIdx).(*Label)
		if lbl.Text != txt {
			lbl.SetText(txt)
			lbl.Config(wb.Sc) // this is essential
		}
	}
}
