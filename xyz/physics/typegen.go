// Code generated by "core generate -add-types"; DO NOT EDIT.

package physics

import (
	"cogentcore.org/core/math32"
	"cogentcore.org/core/tree"
	"cogentcore.org/core/types"
)

var _ = types.AddType(&types.Type{Name: "cogentcore.org/core/xyz/physics.BBox", IDName: "b-box", Doc: "BBox contains bounding box and other gross object properties", Fields: []types.Field{{Name: "BBox", Doc: "bounding box in world coords (Axis-Aligned Bounding Box = AABB)"}, {Name: "VelBBox", Doc: "velocity-projected bounding box in world coords: extend BBox to include future position of moving bodies -- collision must be made on this basis"}, {Name: "BSphere", Doc: "bounding sphere in local coords"}, {Name: "Area", Doc: "area"}, {Name: "Volume", Doc: "volume"}}})

var _ = types.AddType(&types.Type{Name: "cogentcore.org/core/xyz/physics.Body", IDName: "body", Doc: "Body is the common interface for all body types", Methods: []types.Method{{Name: "AsBodyBase", Doc: "AsBodyBase returns the body as a BodyBase", Returns: []string{"BodyBase"}}}})

// BodyBaseType is the [types.Type] for [BodyBase]
var BodyBaseType = types.AddType(&types.Type{Name: "cogentcore.org/core/xyz/physics.BodyBase", IDName: "body-base", Doc: "BodyBase is the base type for all specific Body types", Embeds: []types.Field{{Name: "NodeBase"}}, Fields: []types.Field{{Name: "Rigid", Doc: "rigid body properties, including mass, bounce, friction etc"}, {Name: "Vis", Doc: "visualization name -- looks up an entry in the scene library that provides the visual representation of this body"}, {Name: "Color", Doc: "default color of body for basic InitLibrary configuration"}}, Instance: &BodyBase{}})

// NewBodyBase returns a new [BodyBase] with the given optional parent:
// BodyBase is the base type for all specific Body types
func NewBodyBase(parent ...tree.Node) *BodyBase { return tree.New[BodyBase](parent...) }

// SetRigid sets the [BodyBase.Rigid]:
// rigid body properties, including mass, bounce, friction etc
func (t *BodyBase) SetRigid(v Rigid) *BodyBase { t.Rigid = v; return t }

// SetVis sets the [BodyBase.Vis]:
// visualization name -- looks up an entry in the scene library that provides the visual representation of this body
func (t *BodyBase) SetVis(v string) *BodyBase { t.Vis = v; return t }

// SetColor sets the [BodyBase.Color]:
// default color of body for basic InitLibrary configuration
func (t *BodyBase) SetColor(v string) *BodyBase { t.Color = v; return t }

// BoxType is the [types.Type] for [Box]
var BoxType = types.AddType(&types.Type{Name: "cogentcore.org/core/xyz/physics.Box", IDName: "box", Doc: "Box is a box body shape", Embeds: []types.Field{{Name: "BodyBase"}}, Fields: []types.Field{{Name: "Size", Doc: "size of box in each dimension (units arbitrary, as long as they are all consistent -- meters is typical)"}}, Instance: &Box{}})

// NewBox returns a new [Box] with the given optional parent:
// Box is a box body shape
func NewBox(parent ...tree.Node) *Box { return tree.New[Box](parent...) }

// SetSize sets the [Box.Size]:
// size of box in each dimension (units arbitrary, as long as they are all consistent -- meters is typical)
func (t *Box) SetSize(v math32.Vector3) *Box { t.Size = v; return t }

// CapsuleType is the [types.Type] for [Capsule]
var CapsuleType = types.AddType(&types.Type{Name: "cogentcore.org/core/xyz/physics.Capsule", IDName: "capsule", Doc: "Capsule is a generalized cylinder body shape, with hemispheres at each end,\nwith separate radii for top and bottom.", Embeds: []types.Field{{Name: "BodyBase"}}, Fields: []types.Field{{Name: "Height", Doc: "height of the cylinder portion of the capsule"}, {Name: "TopRad", Doc: "radius of the top hemisphere"}, {Name: "BotRad", Doc: "radius of the bottom hemisphere"}}, Instance: &Capsule{}})

// NewCapsule returns a new [Capsule] with the given optional parent:
// Capsule is a generalized cylinder body shape, with hemispheres at each end,
// with separate radii for top and bottom.
func NewCapsule(parent ...tree.Node) *Capsule { return tree.New[Capsule](parent...) }

// SetHeight sets the [Capsule.Height]:
// height of the cylinder portion of the capsule
func (t *Capsule) SetHeight(v float32) *Capsule { t.Height = v; return t }

// SetTopRad sets the [Capsule.TopRad]:
// radius of the top hemisphere
func (t *Capsule) SetTopRad(v float32) *Capsule { t.TopRad = v; return t }

// SetBotRad sets the [Capsule.BotRad]:
// radius of the bottom hemisphere
func (t *Capsule) SetBotRad(v float32) *Capsule { t.BotRad = v; return t }

var _ = types.AddType(&types.Type{Name: "cogentcore.org/core/xyz/physics.Contact", IDName: "contact", Doc: "Contact is one pairwise point of contact between two bodies.\nContacts are represented in spherical terms relative to the\nspherical BBox of A and B.", Fields: []types.Field{{Name: "A", Doc: "one body"}, {Name: "B", Doc: "the other body"}, {Name: "NormB", Doc: "normal pointing from center of B to center of A"}, {Name: "PtB", Doc: "point on spherical shell of B where A is contacting"}, {Name: "Dist", Doc: "distance from PtB along NormB to contact point on spherical shell of A"}}})

var _ = types.AddType(&types.Type{Name: "cogentcore.org/core/xyz/physics.Contacts", IDName: "contacts", Doc: "Contacts is a slice list of contacts"})

// CylinderType is the [types.Type] for [Cylinder]
var CylinderType = types.AddType(&types.Type{Name: "cogentcore.org/core/xyz/physics.Cylinder", IDName: "cylinder", Doc: "Cylinder is a generalized cylinder body shape, with separate radii for top and bottom.\nA cone has a zero radius at one end.", Embeds: []types.Field{{Name: "BodyBase"}}, Fields: []types.Field{{Name: "Height", Doc: "height of the cylinder"}, {Name: "TopRad", Doc: "radius of the top -- set to 0 for a cone"}, {Name: "BotRad", Doc: "radius of the bottom"}}, Instance: &Cylinder{}})

// NewCylinder returns a new [Cylinder] with the given optional parent:
// Cylinder is a generalized cylinder body shape, with separate radii for top and bottom.
// A cone has a zero radius at one end.
func NewCylinder(parent ...tree.Node) *Cylinder { return tree.New[Cylinder](parent...) }

// SetHeight sets the [Cylinder.Height]:
// height of the cylinder
func (t *Cylinder) SetHeight(v float32) *Cylinder { t.Height = v; return t }

// SetTopRad sets the [Cylinder.TopRad]:
// radius of the top -- set to 0 for a cone
func (t *Cylinder) SetTopRad(v float32) *Cylinder { t.TopRad = v; return t }

// SetBotRad sets the [Cylinder.BotRad]:
// radius of the bottom
func (t *Cylinder) SetBotRad(v float32) *Cylinder { t.BotRad = v; return t }

// GroupType is the [types.Type] for [Group]
var GroupType = types.AddType(&types.Type{Name: "cogentcore.org/core/xyz/physics.Group", IDName: "group", Doc: "Group is a container of bodies, joints, or other groups\nit should be used strategically to partition the space\nand its BBox is used to optimize tree-based collision detection.\nUse a group for the top-level World node as well.", Embeds: []types.Field{{Name: "NodeBase"}}, Instance: &Group{}})

// NewGroup returns a new [Group] with the given optional parent:
// Group is a container of bodies, joints, or other groups
// it should be used strategically to partition the space
// and its BBox is used to optimize tree-based collision detection.
// Use a group for the top-level World node as well.
func NewGroup(parent ...tree.Node) *Group { return tree.New[Group](parent...) }

var _ = types.AddType(&types.Type{Name: "cogentcore.org/core/xyz/physics.BodyPoint", IDName: "body-point", Doc: "BodyPoint contains a Body and a Point on that body", Fields: []types.Field{{Name: "Body"}, {Name: "Point"}}})

var _ = types.AddType(&types.Type{Name: "cogentcore.org/core/xyz/physics.Node", IDName: "node", Doc: "Node is the common interface for all nodes.", Methods: []types.Method{{Name: "AsNodeBase", Doc: "AsNodeBase returns a generic NodeBase for our node -- gives generic\naccess to all the base-level data structures without needing interface methods.", Returns: []string{"NodeBase"}}, {Name: "AsBody", Doc: "AsBody returns a generic Body interface for our node -- nil if not a Body", Returns: []string{"Body"}}, {Name: "GroupBBox", Doc: "GroupBBox sets bounding boxes for groups based on groups or bodies.\ncalled in a FuncDownMeLast traversal."}, {Name: "InitAbs", Doc: "InitAbs sets current Abs physical state parameters from Initial values\nwhich are local, relative to parent -- is passed the parent (nil = top).\nBody nodes should also set their bounding boxes.\nCalled in a FuncDownMeFirst traversal.", Args: []string{"par"}}, {Name: "RelToAbs", Doc: "RelToAbs updates current world Abs physical state parameters\nbased on Rel values added to updated Abs values at higher levels.\nAbs.LinVel is updated from the resulting change from prior position.\nThis is useful for manual updating of relative positions (scripted movement).\nIt is passed the parent (nil = top).\nBody nodes should also update their bounding boxes.\nCalled in a FuncDownMeFirst traversal.", Args: []string{"par"}}, {Name: "Step", Doc: "Step computes one update of the world Abs physical state parameters,\nusing *current* velocities -- add forces prior to calling.\nUse this for physics-based state updates.\nBody nodes should also update their bounding boxes.", Args: []string{"step"}}}})

// NodeBaseType is the [types.Type] for [NodeBase]
var NodeBaseType = types.AddType(&types.Type{Name: "cogentcore.org/core/xyz/physics.NodeBase", IDName: "node-base", Doc: "NodeBase is the basic node, which has position, rotation, velocity\nand computed bounding boxes, etc.\nThere are only three different kinds of Nodes: Group, Body, and Joint", Embeds: []types.Field{{Name: "NodeBase"}}, Fields: []types.Field{{Name: "Dynamic", Doc: "Dynamic is whether this node can move. If it is false, then this is a Static node.\nAny top-level group that is not Dynamic is immediately pruned from further consideration,\nso top-level groups should be separated into Dynamic and Static nodes at the start."}, {Name: "Initial", Doc: "initial position, orientation, velocity in *local* coordinates (relative to parent)"}, {Name: "Rel", Doc: "current relative (local) position, orientation, velocity -- only change these values, as abs values are computed therefrom"}, {Name: "Abs", Doc: "current absolute (world) position, orientation, velocity"}, {Name: "BBox", Doc: "bounding box in world coordinates (aggregated for groups)"}}, Instance: &NodeBase{}})

// NewNodeBase returns a new [NodeBase] with the given optional parent:
// NodeBase is the basic node, which has position, rotation, velocity
// and computed bounding boxes, etc.
// There are only three different kinds of Nodes: Group, Body, and Joint
func NewNodeBase(parent ...tree.Node) *NodeBase { return tree.New[NodeBase](parent...) }

// SetDynamic sets the [NodeBase.Dynamic]:
// Dynamic is whether this node can move. If it is false, then this is a Static node.
// Any top-level group that is not Dynamic is immediately pruned from further consideration,
// so top-level groups should be separated into Dynamic and Static nodes at the start.
func (t *NodeBase) SetDynamic(v bool) *NodeBase { t.Dynamic = v; return t }

// SetInitial sets the [NodeBase.Initial]:
// initial position, orientation, velocity in *local* coordinates (relative to parent)
func (t *NodeBase) SetInitial(v State) *NodeBase { t.Initial = v; return t }

// SetRel sets the [NodeBase.Rel]:
// current relative (local) position, orientation, velocity -- only change these values, as abs values are computed therefrom
func (t *NodeBase) SetRel(v State) *NodeBase { t.Rel = v; return t }

var _ = types.AddType(&types.Type{Name: "cogentcore.org/core/xyz/physics.Rigid", IDName: "rigid", Doc: "Rigid contains the full specification of a given object's basic physics\nproperties including position, orientation, velocity.  These", Fields: []types.Field{{Name: "InvMass", Doc: "1/mass -- 0 for no mass"}, {Name: "Bounce", Doc: "COR or coefficient of restitution -- how elastic is the collision i.e., final velocity / initial velocity"}, {Name: "Friction", Doc: "friction coefficient -- how much friction is generated by transverse motion"}, {Name: "Force", Doc: "record of computed force vector from last iteration"}, {Name: "RotInertia", Doc: "Last calculated rotational inertia matrix in local coords"}}})

// SphereType is the [types.Type] for [Sphere]
var SphereType = types.AddType(&types.Type{Name: "cogentcore.org/core/xyz/physics.Sphere", IDName: "sphere", Doc: "Sphere is a spherical body shape.", Embeds: []types.Field{{Name: "BodyBase"}}, Fields: []types.Field{{Name: "Radius", Doc: "radius"}}, Instance: &Sphere{}})

// NewSphere returns a new [Sphere] with the given optional parent:
// Sphere is a spherical body shape.
func NewSphere(parent ...tree.Node) *Sphere { return tree.New[Sphere](parent...) }

// SetRadius sets the [Sphere.Radius]:
// radius
func (t *Sphere) SetRadius(v float32) *Sphere { t.Radius = v; return t }

var _ = types.AddType(&types.Type{Name: "cogentcore.org/core/xyz/physics.State", IDName: "state", Doc: "State contains the basic physical state including position, orientation, velocity.\nThese are only the values that can be either relative or absolute -- other physical\nstate values such as Mass should go in Rigid.", Fields: []types.Field{{Name: "Pos", Doc: "position of center of mass of object"}, {Name: "Quat", Doc: "rotation specified as a Quat"}, {Name: "LinVel", Doc: "linear velocity"}, {Name: "AngVel", Doc: "angular velocity"}}})
