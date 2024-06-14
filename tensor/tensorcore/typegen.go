// Code generated by "core generate"; DO NOT EDIT.

package tensorcore

import (
	"cogentcore.org/core/colors/colormap"
	"cogentcore.org/core/tensor"
	"cogentcore.org/core/tensor/stats/simat"
	"cogentcore.org/core/tensor/table"
	"cogentcore.org/core/tree"
	"cogentcore.org/core/types"
)

// SimMatGridType is the [types.Type] for [SimMatGrid]
var SimMatGridType = types.AddType(&types.Type{Name: "cogentcore.org/core/tensor/tensorcore.SimMatGrid", IDName: "sim-mat-grid", Doc: "SimMatGrid is a widget that displays a similarity / distance matrix\nwith tensor values as a grid of colored squares, and labels for rows and columns.", Directives: []types.Directive{{Tool: "types", Directive: "add"}}, Embeds: []types.Field{{Name: "TensorGrid"}}, Fields: []types.Field{{Name: "SimMat", Doc: "the similarity / distance matrix"}, {Name: "rowMaxSz"}, {Name: "rowMinBlank"}, {Name: "rowNGps"}, {Name: "colMaxSz"}, {Name: "colMinBlank"}, {Name: "colNGps"}}, Instance: &SimMatGrid{}})

// NewSimMatGrid returns a new [SimMatGrid] with the given optional parent:
// SimMatGrid is a widget that displays a similarity / distance matrix
// with tensor values as a grid of colored squares, and labels for rows and columns.
func NewSimMatGrid(parent ...tree.Node) *SimMatGrid { return tree.New[*SimMatGrid](parent...) }

// NodeType returns the [*types.Type] of [SimMatGrid]
func (t *SimMatGrid) NodeType() *types.Type { return SimMatGridType }

// New returns a new [*SimMatGrid] value
func (t *SimMatGrid) New() tree.Node { return &SimMatGrid{} }

// TableType is the [types.Type] for [Table]
var TableType = types.AddType(&types.Type{Name: "cogentcore.org/core/tensor/tensorcore.Table", IDName: "table", Doc: "Table provides a GUI widget for representing [table.Table] values.", Embeds: []types.Field{{Name: "ListBase"}}, Fields: []types.Field{{Name: "Table", Doc: "the idx view of the table that we're a view of"}, {Name: "TensorDisplay", Doc: "overall display options for tensor display"}, {Name: "ColumnTensorDisplay", Doc: "per column tensor display params"}, {Name: "ColumnTensorBlank", Doc: "per column blank tensor values"}, {Name: "NCols", Doc: "number of columns in table (as of last update)"}, {Name: "SortIndex", Doc: "current sort index"}, {Name: "SortDescending", Doc: "whether current sort order is descending"}, {Name: "headerWidths", Doc: "headerWidths has number of characters in each header, per visfields"}, {Name: "colMaxWidths", Doc: "colMaxWidths records maximum width in chars of string type fields"}, {Name: "BlankString", Doc: "\tblank values for out-of-range rows"}, {Name: "BlankFloat"}}, Instance: &Table{}})

// NewTable returns a new [Table] with the given optional parent:
// Table provides a GUI widget for representing [table.Table] values.
func NewTable(parent ...tree.Node) *Table { return tree.New[*Table](parent...) }

// NodeType returns the [*types.Type] of [Table]
func (t *Table) NodeType() *types.Type { return TableType }

// New returns a new [*Table] value
func (t *Table) New() tree.Node { return &Table{} }

// SetNCols sets the [Table.NCols]:
// number of columns in table (as of last update)
func (t *Table) SetNCols(v int) *Table { t.NCols = v; return t }

// SetSortIndex sets the [Table.SortIndex]:
// current sort index
func (t *Table) SetSortIndex(v int) *Table { t.SortIndex = v; return t }

// SetSortDescending sets the [Table.SortDescending]:
// whether current sort order is descending
func (t *Table) SetSortDescending(v bool) *Table { t.SortDescending = v; return t }

// SetBlankString sets the [Table.BlankString]:
//
//	blank values for out-of-range rows
func (t *Table) SetBlankString(v string) *Table { t.BlankString = v; return t }

// SetBlankFloat sets the [Table.BlankFloat]
func (t *Table) SetBlankFloat(v float64) *Table { t.BlankFloat = v; return t }

var _ = types.AddType(&types.Type{Name: "cogentcore.org/core/tensor/tensorcore.TensorLayout", IDName: "tensor-layout", Doc: "TensorLayout are layout options for displaying tensors", Directives: []types.Directive{{Tool: "types", Directive: "add"}}, Fields: []types.Field{{Name: "OddRow", Doc: "even-numbered dimensions are displayed as Y*X rectangles -- this determines along which dimension to display any remaining odd dimension: OddRow = true = organize vertically along row dimension, false = organize horizontally across column dimension"}, {Name: "TopZero", Doc: "if true, then the Y=0 coordinate is displayed from the top-down; otherwise the Y=0 coordinate is displayed from the bottom up, which is typical for emergent network patterns."}, {Name: "Image", Doc: "display the data as a bitmap image.  if a 2D tensor, then it will be a greyscale image.  if a 3D tensor with size of either the first or last dim = either 3 or 4, then it is a RGB(A) color image"}}})

var _ = types.AddType(&types.Type{Name: "cogentcore.org/core/tensor/tensorcore.TensorDisplay", IDName: "tensor-display", Doc: "TensorDisplay are options for displaying tensors", Directives: []types.Directive{{Tool: "types", Directive: "add"}}, Embeds: []types.Field{{Name: "TensorLayout"}}, Fields: []types.Field{{Name: "Range", Doc: "range to plot"}, {Name: "MinMax", Doc: "if not using fixed range, this is the actual range of data"}, {Name: "ColorMap", Doc: "the name of the color map to use in translating values to colors"}, {Name: "GridFill", Doc: "what proportion of grid square should be filled by color block -- 1 = all, .5 = half, etc"}, {Name: "DimExtra", Doc: "amount of extra space to add at dimension boundaries, as a proportion of total grid size"}, {Name: "GridMinSize", Doc: "minimum size for grid squares -- they will never be smaller than this"}, {Name: "GridMaxSize", Doc: "maximum size for grid squares -- they will never be larger than this"}, {Name: "TotPrefSize", Doc: "total preferred display size along largest dimension.\ngrid squares will be sized to fit within this size,\nsubject to harder GridMin / Max size constraints"}, {Name: "FontSize", Doc: "font size in standard point units for labels (e.g., SimMat)"}, {Name: "GridView", Doc: "our gridview, for update method"}}})

// TensorGridType is the [types.Type] for [TensorGrid]
var TensorGridType = types.AddType(&types.Type{Name: "cogentcore.org/core/tensor/tensorcore.TensorGrid", IDName: "tensor-grid", Doc: "TensorGrid is a widget that displays tensor values as a grid of colored squares.", Methods: []types.Method{{Name: "EditSettings", Directives: []types.Directive{{Tool: "types", Directive: "add"}}}}, Embeds: []types.Field{{Name: "WidgetBase"}}, Fields: []types.Field{{Name: "Tensor", Doc: "the tensor that we view"}, {Name: "Display", Doc: "display options"}, {Name: "ColorMap", Doc: "the actual colormap"}}, Instance: &TensorGrid{}})

// NewTensorGrid returns a new [TensorGrid] with the given optional parent:
// TensorGrid is a widget that displays tensor values as a grid of colored squares.
func NewTensorGrid(parent ...tree.Node) *TensorGrid { return tree.New[*TensorGrid](parent...) }

// NodeType returns the [*types.Type] of [TensorGrid]
func (t *TensorGrid) NodeType() *types.Type { return TensorGridType }

// New returns a new [*TensorGrid] value
func (t *TensorGrid) New() tree.Node { return &TensorGrid{} }

// SetDisplay sets the [TensorGrid.Display]:
// display options
func (t *TensorGrid) SetDisplay(v TensorDisplay) *TensorGrid { t.Display = v; return t }

// SetColorMap sets the [TensorGrid.ColorMap]:
// the actual colormap
func (t *TensorGrid) SetColorMap(v *colormap.Map) *TensorGrid { t.ColorMap = v; return t }

// TensorButtonType is the [types.Type] for [TensorButton]
var TensorButtonType = types.AddType(&types.Type{Name: "cogentcore.org/core/tensor/tensorcore.TensorButton", IDName: "tensor-button", Doc: "TensorButton represents a Tensor with a button for making a [TensorView]\nviewer for an [tensor.Tensor].", Embeds: []types.Field{{Name: "Button"}}, Fields: []types.Field{{Name: "Tensor"}}, Instance: &TensorButton{}})

// NewTensorButton returns a new [TensorButton] with the given optional parent:
// TensorButton represents a Tensor with a button for making a [TensorView]
// viewer for an [tensor.Tensor].
func NewTensorButton(parent ...tree.Node) *TensorButton { return tree.New[*TensorButton](parent...) }

// NodeType returns the [*types.Type] of [TensorButton]
func (t *TensorButton) NodeType() *types.Type { return TensorButtonType }

// New returns a new [*TensorButton] value
func (t *TensorButton) New() tree.Node { return &TensorButton{} }

// SetTensor sets the [TensorButton.Tensor]
func (t *TensorButton) SetTensor(v tensor.Tensor) *TensorButton { t.Tensor = v; return t }

// TableButtonType is the [types.Type] for [TableButton]
var TableButtonType = types.AddType(&types.Type{Name: "cogentcore.org/core/tensor/tensorcore.TableButton", IDName: "table-button", Doc: "TableButton presents a button that pulls up the [Table] viewer for a [table.Table].", Embeds: []types.Field{{Name: "Button"}}, Fields: []types.Field{{Name: "Table"}}, Instance: &TableButton{}})

// NewTableButton returns a new [TableButton] with the given optional parent:
// TableButton presents a button that pulls up the [Table] viewer for a [table.Table].
func NewTableButton(parent ...tree.Node) *TableButton { return tree.New[*TableButton](parent...) }

// NodeType returns the [*types.Type] of [TableButton]
func (t *TableButton) NodeType() *types.Type { return TableButtonType }

// New returns a new [*TableButton] value
func (t *TableButton) New() tree.Node { return &TableButton{} }

// SetTable sets the [TableButton.Table]
func (t *TableButton) SetTable(v *table.Table) *TableButton { t.Table = v; return t }

// SimMatButtonType is the [types.Type] for [SimMatButton]
var SimMatButtonType = types.AddType(&types.Type{Name: "cogentcore.org/core/tensor/tensorcore.SimMatButton", IDName: "sim-mat-button", Doc: "SimMatValue presents a button that pulls up the [SimMatGridView] viewer for a [table.Table].", Embeds: []types.Field{{Name: "Button"}}, Fields: []types.Field{{Name: "SimMat"}}, Instance: &SimMatButton{}})

// NewSimMatButton returns a new [SimMatButton] with the given optional parent:
// SimMatValue presents a button that pulls up the [SimMatGridView] viewer for a [table.Table].
func NewSimMatButton(parent ...tree.Node) *SimMatButton { return tree.New[*SimMatButton](parent...) }

// NodeType returns the [*types.Type] of [SimMatButton]
func (t *SimMatButton) NodeType() *types.Type { return SimMatButtonType }

// New returns a new [*SimMatButton] value
func (t *SimMatButton) New() tree.Node { return &SimMatButton{} }

// SetSimMat sets the [SimMatButton.SimMat]
func (t *SimMatButton) SetSimMat(v *simat.SimMat) *SimMatButton { t.SimMat = v; return t }
