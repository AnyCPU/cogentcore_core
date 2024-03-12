// Code generated by "core generate"; DO NOT EDIT.

package svg

import (
	"cogentcore.org/core/enums"
	"cogentcore.org/core/ki"
)

var _NodeFlagsValues = []NodeFlags{1}

// NodeFlagsN is the highest valid value for type NodeFlags, plus one.
const NodeFlagsN NodeFlags = 2

var _NodeFlagsValueMap = map[string]NodeFlags{`IsDef`: 1}

var _NodeFlagsDescMap = map[NodeFlags]string{1: `Rendering means that the SVG is currently redrawing Can be useful to check for animations etc to decide whether to drive another update`}

var _NodeFlagsMap = map[NodeFlags]string{1: `IsDef`}

// String returns the string representation of this NodeFlags value.
func (i NodeFlags) String() string {
	return enums.BitFlagStringExtended(i, _NodeFlagsValues, ki.FlagsValues())
}

// BitIndexString returns the string representation of this NodeFlags value
// if it is a bit index value (typically an enum constant), and
// not an actual bit flag value.
func (i NodeFlags) BitIndexString() string {
	return enums.BitIndexStringExtended[NodeFlags, ki.Flags](i, _NodeFlagsMap)
}

// SetString sets the NodeFlags value from its string representation,
// and returns an error if the string is invalid.
func (i *NodeFlags) SetString(s string) error { *i = 0; return i.SetStringOr(s) }

// SetStringOr sets the NodeFlags value from its string representation
// while preserving any bit flags already set, and returns an
// error if the string is invalid.
func (i *NodeFlags) SetStringOr(s string) error {
	return enums.SetStringOrExtended(i, (*ki.Flags)(i), s, _NodeFlagsValueMap)
}

// Int64 returns the NodeFlags value as an int64.
func (i NodeFlags) Int64() int64 { return int64(i) }

// SetInt64 sets the NodeFlags value from an int64.
func (i *NodeFlags) SetInt64(in int64) { *i = NodeFlags(in) }

// Desc returns the description of the NodeFlags value.
func (i NodeFlags) Desc() string {
	return enums.DescExtended[NodeFlags, ki.Flags](i, _NodeFlagsDescMap)
}

// NodeFlagsValues returns all possible values for the type NodeFlags.
func NodeFlagsValues() []NodeFlags {
	return enums.ValuesGlobalExtended(_NodeFlagsValues, ki.FlagsValues())
}

// Values returns all possible values for the type NodeFlags.
func (i NodeFlags) Values() []enums.Enum {
	return enums.ValuesExtended(_NodeFlagsValues, ki.FlagsValues())
}

// HasFlag returns whether these bit flags have the given bit flag set.
func (i NodeFlags) HasFlag(f enums.BitFlag) bool { return enums.HasFlag((*int64)(&i), f) }

// SetFlag sets the value of the given flags in these flags to the given value.
func (i *NodeFlags) SetFlag(on bool, f ...enums.BitFlag) { enums.SetFlag((*int64)(i), on, f...) }

// MarshalText implements the [encoding.TextMarshaler] interface.
func (i NodeFlags) MarshalText() ([]byte, error) { return []byte(i.String()), nil }

// UnmarshalText implements the [encoding.TextUnmarshaler] interface.
func (i *NodeFlags) UnmarshalText(text []byte) error {
	return enums.UnmarshalText(i, text, "NodeFlags")
}

var _ViewBoxAlignsValues = []ViewBoxAligns{0, 1, 2, 3}

// ViewBoxAlignsN is the highest valid value for type ViewBoxAligns, plus one.
const ViewBoxAlignsN ViewBoxAligns = 4

var _ViewBoxAlignsValueMap = map[string]ViewBoxAligns{`mid`: 0, `none`: 1, `min`: 2, `max`: 3}

var _ViewBoxAlignsDescMap = map[ViewBoxAligns]string{0: `align ViewBox.Min with midpoint of Viewport (default)`, 1: `do not preserve uniform scaling (if either X or Y is None, both are treated as such). In this case, the Meet / Slice value is ignored. This is the same as FitFill from styles.ObjectFits`, 2: `align ViewBox.Min with top / left of Viewport`, 3: `align ViewBox.Min+Size with bottom / right of Viewport`}

var _ViewBoxAlignsMap = map[ViewBoxAligns]string{0: `mid`, 1: `none`, 2: `min`, 3: `max`}

// String returns the string representation of this ViewBoxAligns value.
func (i ViewBoxAligns) String() string { return enums.String(i, _ViewBoxAlignsMap) }

// SetString sets the ViewBoxAligns value from its string representation,
// and returns an error if the string is invalid.
func (i *ViewBoxAligns) SetString(s string) error {
	return enums.SetString(i, s, _ViewBoxAlignsValueMap, "ViewBoxAligns")
}

// Int64 returns the ViewBoxAligns value as an int64.
func (i ViewBoxAligns) Int64() int64 { return int64(i) }

// SetInt64 sets the ViewBoxAligns value from an int64.
func (i *ViewBoxAligns) SetInt64(in int64) { *i = ViewBoxAligns(in) }

// Desc returns the description of the ViewBoxAligns value.
func (i ViewBoxAligns) Desc() string { return enums.Desc(i, _ViewBoxAlignsDescMap) }

// ViewBoxAlignsValues returns all possible values for the type ViewBoxAligns.
func ViewBoxAlignsValues() []ViewBoxAligns { return _ViewBoxAlignsValues }

// Values returns all possible values for the type ViewBoxAligns.
func (i ViewBoxAligns) Values() []enums.Enum { return enums.Values(_ViewBoxAlignsValues) }

// MarshalText implements the [encoding.TextMarshaler] interface.
func (i ViewBoxAligns) MarshalText() ([]byte, error) { return []byte(i.String()), nil }

// UnmarshalText implements the [encoding.TextUnmarshaler] interface.
func (i *ViewBoxAligns) UnmarshalText(text []byte) error {
	return enums.UnmarshalText(i, text, "ViewBoxAligns")
}

var _ViewBoxMeetOrSliceValues = []ViewBoxMeetOrSlice{0, 1}

// ViewBoxMeetOrSliceN is the highest valid value for type ViewBoxMeetOrSlice, plus one.
const ViewBoxMeetOrSliceN ViewBoxMeetOrSlice = 2

var _ViewBoxMeetOrSliceValueMap = map[string]ViewBoxMeetOrSlice{`meet`: 0, `slice`: 1}

var _ViewBoxMeetOrSliceDescMap = map[ViewBoxMeetOrSlice]string{0: `Meet only applies if Align != None (i.e., only for uniform scaling), and means the entire ViewBox is visible within Viewport, and it is scaled up as much as possible to meet the align constraints. This is the same as FitContain from styles.ObjectFits`, 1: `Slice only applies if Align != None (i.e., only for uniform scaling), and means the entire ViewBox is covered by the ViewBox, and the ViewBox is scaled down as much as possible, while still meeting the align constraints. This is the same as FitCover from styles.ObjectFits`}

var _ViewBoxMeetOrSliceMap = map[ViewBoxMeetOrSlice]string{0: `meet`, 1: `slice`}

// String returns the string representation of this ViewBoxMeetOrSlice value.
func (i ViewBoxMeetOrSlice) String() string { return enums.String(i, _ViewBoxMeetOrSliceMap) }

// SetString sets the ViewBoxMeetOrSlice value from its string representation,
// and returns an error if the string is invalid.
func (i *ViewBoxMeetOrSlice) SetString(s string) error {
	return enums.SetString(i, s, _ViewBoxMeetOrSliceValueMap, "ViewBoxMeetOrSlice")
}

// Int64 returns the ViewBoxMeetOrSlice value as an int64.
func (i ViewBoxMeetOrSlice) Int64() int64 { return int64(i) }

// SetInt64 sets the ViewBoxMeetOrSlice value from an int64.
func (i *ViewBoxMeetOrSlice) SetInt64(in int64) { *i = ViewBoxMeetOrSlice(in) }

// Desc returns the description of the ViewBoxMeetOrSlice value.
func (i ViewBoxMeetOrSlice) Desc() string { return enums.Desc(i, _ViewBoxMeetOrSliceDescMap) }

// ViewBoxMeetOrSliceValues returns all possible values for the type ViewBoxMeetOrSlice.
func ViewBoxMeetOrSliceValues() []ViewBoxMeetOrSlice { return _ViewBoxMeetOrSliceValues }

// Values returns all possible values for the type ViewBoxMeetOrSlice.
func (i ViewBoxMeetOrSlice) Values() []enums.Enum { return enums.Values(_ViewBoxMeetOrSliceValues) }

// MarshalText implements the [encoding.TextMarshaler] interface.
func (i ViewBoxMeetOrSlice) MarshalText() ([]byte, error) { return []byte(i.String()), nil }

// UnmarshalText implements the [encoding.TextUnmarshaler] interface.
func (i *ViewBoxMeetOrSlice) UnmarshalText(text []byte) error {
	return enums.UnmarshalText(i, text, "ViewBoxMeetOrSlice")
}
