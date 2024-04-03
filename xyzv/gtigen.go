// Code generated by "core generate"; DO NOT EDIT.

package xyzv

import (
	"cogentcore.org/core/gti"
	"cogentcore.org/core/ki"
	"cogentcore.org/core/xyz"
)

// ManipPtType is the [gti.Type] for [ManipPt]
var ManipPtType = gti.AddType(&gti.Type{Name: "cogentcore.org/core/xyzv.ManipPt", IDName: "manip-pt", Doc: "ManipPt is a manipulation control point", Directives: []gti.Directive{{Tool: "core", Directive: "no-new"}}, Embeds: []gti.Field{{Name: "Solid"}}, Instance: &ManipPt{}})

// KiType returns the [*gti.Type] of [ManipPt]
func (t *ManipPt) KiType() *gti.Type { return ManipPtType }

// New returns a new [*ManipPt] value
func (t *ManipPt) New() ki.Ki { return &ManipPt{} }

// SetMat sets the [ManipPt.Mat]
func (t *ManipPt) SetMat(v xyz.Material) *ManipPt { t.Mat = v; return t }

// SceneType is the [gti.Type] for [Scene]
var SceneType = gti.AddType(&gti.Type{Name: "cogentcore.org/core/xyzv.Scene", IDName: "scene", Doc: "Scene is a gi.Widget that manages a xyz.Scene,\nproviding the basic rendering logic for the 3D scene\nin the 2D gi gui context.", Embeds: []gti.Field{{Name: "WidgetBase"}}, Fields: []gti.Field{{Name: "XYZ", Doc: "XYZ is the 3D xyz.Scene"}, {Name: "SelMode", Doc: "how to deal with selection / manipulation events"}, {Name: "CurSel", Doc: "currently selected node"}, {Name: "CurManipPt", Doc: "currently selected manipulation control point"}, {Name: "SelParams", Doc: "parameters for selection / manipulation box"}}, Instance: &Scene{}})

// NewScene adds a new [Scene] with the given name to the given parent:
// Scene is a gi.Widget that manages a xyz.Scene,
// providing the basic rendering logic for the 3D scene
// in the 2D gi gui context.
func NewScene(par ki.Ki, name ...string) *Scene {
	return par.NewChild(SceneType, name...).(*Scene)
}

// KiType returns the [*gti.Type] of [Scene]
func (t *Scene) KiType() *gti.Type { return SceneType }

// New returns a new [*Scene] value
func (t *Scene) New() ki.Ki { return &Scene{} }

// SetSelMode sets the [Scene.SelMode]:
// how to deal with selection / manipulation events
func (t *Scene) SetSelMode(v SelectionModes) *Scene { t.SelectionMode = v; return t }

// SetCurSel sets the [Scene.CurSel]:
// currently selected node
func (t *Scene) SetCurSel(v xyz.Node) *Scene { t.CurrentSelected = v; return t }

// SetCurManipPt sets the [Scene.CurManipPt]:
// currently selected manipulation control point
func (t *Scene) SetCurManipPt(v *ManipPt) *Scene { t.CurrentManipPoint = v; return t }

// SetSelParams sets the [Scene.SelParams]:
// parameters for selection / manipulation box
func (t *Scene) SetSelParams(v SelectionParams) *Scene { t.SelectionParams = v; return t }

// SetTooltip sets the [Scene.Tooltip]
func (t *Scene) SetTooltip(v string) *Scene { t.Tooltip = v; return t }

// SceneViewType is the [gti.Type] for [SceneView]
var SceneViewType = gti.AddType(&gti.Type{Name: "cogentcore.org/core/xyzv.SceneView", IDName: "scene-view", Doc: "SceneView provides a toolbar controller for an xyz.Scene,\nand manipulation abilities.", Embeds: []gti.Field{{Name: "Layout"}}, Instance: &SceneView{}})

// NewSceneView adds a new [SceneView] with the given name to the given parent:
// SceneView provides a toolbar controller for an xyz.Scene,
// and manipulation abilities.
func NewSceneView(par ki.Ki, name ...string) *SceneView {
	return par.NewChild(SceneViewType, name...).(*SceneView)
}

// KiType returns the [*gti.Type] of [SceneView]
func (t *SceneView) KiType() *gti.Type { return SceneViewType }

// New returns a new [*SceneView] value
func (t *SceneView) New() ki.Ki { return &SceneView{} }

// SetTooltip sets the [SceneView.Tooltip]
func (t *SceneView) SetTooltip(v string) *SceneView { t.Tooltip = v; return t }

// SetStackTop sets the [SceneView.StackTop]
func (t *SceneView) SetStackTop(v int) *SceneView { t.StackTop = v; return t }
