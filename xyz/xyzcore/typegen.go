// Code generated by "core generate"; DO NOT EDIT.

package xyzcore

import (
	"cogentcore.org/core/tree"
	"cogentcore.org/core/types"
	"cogentcore.org/core/xyz"
)

// ManipPointType is the [types.Type] for [ManipPoint]
var ManipPointType = types.AddType(&types.Type{Name: "cogentcore.org/core/xyz/xyzcore.ManipPoint", IDName: "manip-point", Doc: "ManipPoint is a manipulation control point.", Directives: []types.Directive{{Tool: "core", Directive: "no-new"}}, Embeds: []types.Field{{Name: "Solid"}}, Instance: &ManipPoint{}})

// SceneType is the [types.Type] for [Scene]
var SceneType = types.AddType(&types.Type{Name: "cogentcore.org/core/xyz/xyzcore.Scene", IDName: "scene", Doc: "Scene is a core.Widget that manages a xyz.Scene,\nproviding the basic rendering logic for the 3D scene\nin the 2D core GUI context.", Embeds: []types.Field{{Name: "WidgetBase"}}, Fields: []types.Field{{Name: "XYZ", Doc: "XYZ is the 3D xyz.Scene"}, {Name: "SelectionMode", Doc: "how to deal with selection / manipulation events"}, {Name: "CurrentSelected", Doc: "currently selected node"}, {Name: "CurrentManipPoint", Doc: "currently selected manipulation control point"}, {Name: "SelectionParams", Doc: "parameters for selection / manipulation box"}}, Instance: &Scene{}})

// NewScene returns a new [Scene] with the given optional parent:
// Scene is a core.Widget that manages a xyz.Scene,
// providing the basic rendering logic for the 3D scene
// in the 2D core GUI context.
func NewScene(parent ...tree.Node) *Scene { return tree.New[Scene](parent...) }

// SetSelectionMode sets the [Scene.SelectionMode]:
// how to deal with selection / manipulation events
func (t *Scene) SetSelectionMode(v SelectionModes) *Scene { t.SelectionMode = v; return t }

// SetCurrentSelected sets the [Scene.CurrentSelected]:
// currently selected node
func (t *Scene) SetCurrentSelected(v xyz.Node) *Scene { t.CurrentSelected = v; return t }

// SetCurrentManipPoint sets the [Scene.CurrentManipPoint]:
// currently selected manipulation control point
func (t *Scene) SetCurrentManipPoint(v *ManipPoint) *Scene { t.CurrentManipPoint = v; return t }

// SetSelectionParams sets the [Scene.SelectionParams]:
// parameters for selection / manipulation box
func (t *Scene) SetSelectionParams(v SelectionParams) *Scene { t.SelectionParams = v; return t }

// SceneEditorType is the [types.Type] for [SceneEditor]
var SceneEditorType = types.AddType(&types.Type{Name: "cogentcore.org/core/xyz/xyzcore.SceneEditor", IDName: "scene-editor", Doc: "SceneEditor provides a toolbar controller and manipulation abilities\nfor a [Scene].", Embeds: []types.Field{{Name: "Frame"}}, Instance: &SceneEditor{}})

// NewSceneEditor returns a new [SceneEditor] with the given optional parent:
// SceneEditor provides a toolbar controller and manipulation abilities
// for a [Scene].
func NewSceneEditor(parent ...tree.Node) *SceneEditor { return tree.New[SceneEditor](parent...) }

// MeshButtonType is the [types.Type] for [MeshButton]
var MeshButtonType = types.AddType(&types.Type{Name: "cogentcore.org/core/xyz/xyzcore.MeshButton", IDName: "mesh-button", Doc: "MeshButton represents an [xyz.MeshName] value with a button.", Embeds: []types.Field{{Name: "Button"}}, Fields: []types.Field{{Name: "MeshName"}}, Instance: &MeshButton{}})

// NewMeshButton returns a new [MeshButton] with the given optional parent:
// MeshButton represents an [xyz.MeshName] value with a button.
func NewMeshButton(parent ...tree.Node) *MeshButton { return tree.New[MeshButton](parent...) }

// SetMeshName sets the [MeshButton.MeshName]
func (t *MeshButton) SetMeshName(v string) *MeshButton { t.MeshName = v; return t }
