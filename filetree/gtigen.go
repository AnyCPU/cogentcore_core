// Code generated by "goki generate"; DO NOT EDIT.

package filetree

import (
	"sync"
	"time"

	"goki.dev/gi/v2/gi"
	"goki.dev/gi/v2/giv"
	"goki.dev/girl/units"
	"goki.dev/goosi/events"
	"goki.dev/gti"
	"goki.dev/icons"
	"goki.dev/ki/v2"
	"goki.dev/mat32/v2"
	"goki.dev/ordmap"
	"goki.dev/vci/v2"
	"gopkg.in/fsnotify.v1"
)

// NodeType is the [gti.Type] for [Node]
var NodeType = gti.AddType(&gti.Type{
	Name:      "goki.dev/gi/v2/filetree.Node",
	ShortName: "filetree.Node",
	IDName:    "node",
	Doc:       "Node represents a file in the file system, as a TreeView node.\nThe name of the node is the name of the file.\nFolders have children containing further nodes.",
	Directives: gti.Directives{
		&gti.Directive{Tool: "goki", Directive: "embedder", Args: []string{}},
	},
	Fields: ordmap.Make([]ordmap.KeyVal[string, *gti.Field]{
		{"FPath", &gti.Field{Name: "FPath", Type: "goki.dev/gi/v2/gi.FileName", LocalType: "gi.FileName", Doc: "full path to this file", Directives: gti.Directives{}, Tag: "edit:\"-\" set:\"-\" json:\"-\" xml:\"-\" copy:\"-\""}},
		{"Info", &gti.Field{Name: "Info", Type: "goki.dev/fi.FileInfo", LocalType: "fi.FileInfo", Doc: "full standard file info about this file", Directives: gti.Directives{}, Tag: "edit:\"-\" set:\"-\" json:\"-\" xml:\"-\" copy:\"-\""}},
		{"Buf", &gti.Field{Name: "Buf", Type: "*goki.dev/gi/v2/texteditor.Buf", LocalType: "*texteditor.Buf", Doc: "file buffer for editing this file", Directives: gti.Directives{}, Tag: "edit:\"-\" set:\"-\" json:\"-\" xml:\"-\" copy:\"-\""}},
		{"FRoot", &gti.Field{Name: "FRoot", Type: "*goki.dev/gi/v2/filetree.Tree", LocalType: "*Tree", Doc: "root of the tree -- has global state", Directives: gti.Directives{}, Tag: "edit:\"-\" set:\"-\" json:\"-\" xml:\"-\" copy:\"-\""}},
		{"DirRepo", &gti.Field{Name: "DirRepo", Type: "goki.dev/vci/v2.Repo", LocalType: "vci.Repo", Doc: "version control system repository for this directory,\nonly non-nil if this is the highest-level directory in the tree under vcs control", Directives: gti.Directives{}, Tag: "edit:\"-\" set:\"-\" json:\"-\" xml:\"-\" copy:\"-\""}},
		{"RepoFiles", &gti.Field{Name: "RepoFiles", Type: "goki.dev/vci/v2.Files", LocalType: "vci.Files", Doc: "version control system repository file status -- only valid during ReadDir", Directives: gti.Directives{}, Tag: "edit:\"-\" set:\"-\" json:\"-\" xml:\"-\" copy:\"-\""}},
	}),
	Embeds: ordmap.Make([]ordmap.KeyVal[string, *gti.Field]{
		{"TreeView", &gti.Field{Name: "TreeView", Type: "goki.dev/gi/v2/giv.TreeView", LocalType: "giv.TreeView", Doc: "", Directives: gti.Directives{}, Tag: ""}},
	}),
	Methods: ordmap.Make([]ordmap.KeyVal[string, *gti.Method]{
		{"OpenFilesDefault", &gti.Method{Name: "OpenFilesDefault", Doc: "OpenFilesDefault opens selected files with default app for that file type (os defined).\nruns open on Mac, xdg-open on Linux, and start on Windows", Directives: gti.Directives{
			&gti.Directive{Tool: "gti", Directive: "add", Args: []string{}},
		}, Args: ordmap.Make([]ordmap.KeyVal[string, *gti.Field]{}), Returns: ordmap.Make([]ordmap.KeyVal[string, *gti.Field]{})}},
		{"DuplicateFiles", &gti.Method{Name: "DuplicateFiles", Doc: "DuplicateFiles makes a copy of selected files", Directives: gti.Directives{
			&gti.Directive{Tool: "gti", Directive: "add", Args: []string{}},
		}, Args: ordmap.Make([]ordmap.KeyVal[string, *gti.Field]{}), Returns: ordmap.Make([]ordmap.KeyVal[string, *gti.Field]{})}},
		{"DeleteFiles", &gti.Method{Name: "DeleteFiles", Doc: "deletes any selected files or directories. If any directory is selected,\nall files and subdirectories in that directory are also deleted.", Directives: gti.Directives{
			&gti.Directive{Tool: "gti", Directive: "add", Args: []string{}},
		}, Args: ordmap.Make([]ordmap.KeyVal[string, *gti.Field]{}), Returns: ordmap.Make([]ordmap.KeyVal[string, *gti.Field]{})}},
		{"RenameFiles", &gti.Method{Name: "RenameFiles", Doc: "renames any selected files", Directives: gti.Directives{
			&gti.Directive{Tool: "gti", Directive: "add", Args: []string{}},
		}, Args: ordmap.Make([]ordmap.KeyVal[string, *gti.Field]{}), Returns: ordmap.Make([]ordmap.KeyVal[string, *gti.Field]{})}},
		{"RenameFile", &gti.Method{Name: "RenameFile", Doc: "RenameFile renames file to new name", Directives: gti.Directives{
			&gti.Directive{Tool: "gti", Directive: "add", Args: []string{}},
		}, Args: ordmap.Make([]ordmap.KeyVal[string, *gti.Field]{
			{"newpath", &gti.Field{Name: "newpath", Type: "string", LocalType: "string", Doc: "", Directives: gti.Directives{}, Tag: ""}},
		}), Returns: ordmap.Make([]ordmap.KeyVal[string, *gti.Field]{
			{"error", &gti.Field{Name: "error", Type: "error", LocalType: "error", Doc: "", Directives: gti.Directives{}, Tag: ""}},
		})}},
		{"NewFiles", &gti.Method{Name: "NewFiles", Doc: "NewFiles makes a new file in selected directory", Directives: gti.Directives{
			&gti.Directive{Tool: "gti", Directive: "add", Args: []string{}},
		}, Args: ordmap.Make([]ordmap.KeyVal[string, *gti.Field]{
			{"filename", &gti.Field{Name: "filename", Type: "string", LocalType: "string", Doc: "", Directives: gti.Directives{}, Tag: ""}},
			{"addToVCS", &gti.Field{Name: "addToVCS", Type: "bool", LocalType: "bool", Doc: "", Directives: gti.Directives{}, Tag: ""}},
		}), Returns: ordmap.Make([]ordmap.KeyVal[string, *gti.Field]{})}},
		{"NewFile", &gti.Method{Name: "NewFile", Doc: "NewFile makes a new file in this directory node", Directives: gti.Directives{
			&gti.Directive{Tool: "gti", Directive: "add", Args: []string{}},
		}, Args: ordmap.Make([]ordmap.KeyVal[string, *gti.Field]{
			{"filename", &gti.Field{Name: "filename", Type: "string", LocalType: "string", Doc: "", Directives: gti.Directives{}, Tag: ""}},
			{"addToVCS", &gti.Field{Name: "addToVCS", Type: "bool", LocalType: "bool", Doc: "", Directives: gti.Directives{}, Tag: ""}},
		}), Returns: ordmap.Make([]ordmap.KeyVal[string, *gti.Field]{})}},
		{"NewFolders", &gti.Method{Name: "NewFolders", Doc: "makes a new folder in the given selected directory", Directives: gti.Directives{
			&gti.Directive{Tool: "gti", Directive: "add", Args: []string{}},
		}, Args: ordmap.Make([]ordmap.KeyVal[string, *gti.Field]{
			{"foldername", &gti.Field{Name: "foldername", Type: "string", LocalType: "string", Doc: "", Directives: gti.Directives{}, Tag: ""}},
		}), Returns: ordmap.Make([]ordmap.KeyVal[string, *gti.Field]{})}},
		{"NewFolder", &gti.Method{Name: "NewFolder", Doc: "NewFolder makes a new folder (directory) in this directory node", Directives: gti.Directives{
			&gti.Directive{Tool: "gti", Directive: "add", Args: []string{}},
		}, Args: ordmap.Make([]ordmap.KeyVal[string, *gti.Field]{
			{"foldername", &gti.Field{Name: "foldername", Type: "string", LocalType: "string", Doc: "", Directives: gti.Directives{}, Tag: ""}},
		}), Returns: ordmap.Make([]ordmap.KeyVal[string, *gti.Field]{})}},
		{"ShowFileInfo", &gti.Method{Name: "ShowFileInfo", Doc: "Shows file information about selected file(s)", Directives: gti.Directives{
			&gti.Directive{Tool: "gti", Directive: "add", Args: []string{}},
		}, Args: ordmap.Make([]ordmap.KeyVal[string, *gti.Field]{}), Returns: ordmap.Make([]ordmap.KeyVal[string, *gti.Field]{})}},
		{"SortBys", &gti.Method{Name: "SortBys", Doc: "SortBys determines how to sort the selected files in the directory.\nDefault is alpha by name, optionally can be sorted by modification time.", Directives: gti.Directives{
			&gti.Directive{Tool: "gti", Directive: "add", Args: []string{}},
		}, Args: ordmap.Make([]ordmap.KeyVal[string, *gti.Field]{
			{"modTime", &gti.Field{Name: "modTime", Type: "bool", LocalType: "bool", Doc: "", Directives: gti.Directives{}, Tag: ""}},
		}), Returns: ordmap.Make([]ordmap.KeyVal[string, *gti.Field]{})}},
		{"OpenAll", &gti.Method{Name: "OpenAll", Doc: "OpenAll opens all directories under this one", Directives: gti.Directives{
			&gti.Directive{Tool: "gti", Directive: "add", Args: []string{}},
		}, Args: ordmap.Make([]ordmap.KeyVal[string, *gti.Field]{}), Returns: ordmap.Make([]ordmap.KeyVal[string, *gti.Field]{})}},
		{"CloseAll", &gti.Method{Name: "CloseAll", Doc: "CloseAll closes all directories under this one, this included", Directives: gti.Directives{
			&gti.Directive{Tool: "gti", Directive: "add", Args: []string{}},
		}, Args: ordmap.Make([]ordmap.KeyVal[string, *gti.Field]{}), Returns: ordmap.Make([]ordmap.KeyVal[string, *gti.Field]{})}},
		{"RemoveFromExterns", &gti.Method{Name: "RemoveFromExterns", Doc: "RemoveFromExterns removes file from list of external files", Directives: gti.Directives{
			&gti.Directive{Tool: "gti", Directive: "add", Args: []string{}},
		}, Args: ordmap.Make([]ordmap.KeyVal[string, *gti.Field]{}), Returns: ordmap.Make([]ordmap.KeyVal[string, *gti.Field]{})}},
		{"AddToVCSSel", &gti.Method{Name: "AddToVCSSel", Doc: "AddToVCSSel adds selected files to version control system", Directives: gti.Directives{
			&gti.Directive{Tool: "gti", Directive: "add", Args: []string{}},
		}, Args: ordmap.Make([]ordmap.KeyVal[string, *gti.Field]{}), Returns: ordmap.Make([]ordmap.KeyVal[string, *gti.Field]{})}},
		{"DeleteFromVCSSel", &gti.Method{Name: "DeleteFromVCSSel", Doc: "DeleteFromVCSSel removes selected files from version control system", Directives: gti.Directives{
			&gti.Directive{Tool: "gti", Directive: "add", Args: []string{}},
		}, Args: ordmap.Make([]ordmap.KeyVal[string, *gti.Field]{}), Returns: ordmap.Make([]ordmap.KeyVal[string, *gti.Field]{})}},
		{"CommitToVCSSel", &gti.Method{Name: "CommitToVCSSel", Doc: "CommitToVCSSel commits to version control system based on last selected file", Directives: gti.Directives{
			&gti.Directive{Tool: "gti", Directive: "add", Args: []string{}},
		}, Args: ordmap.Make([]ordmap.KeyVal[string, *gti.Field]{}), Returns: ordmap.Make([]ordmap.KeyVal[string, *gti.Field]{})}},
		{"RevertVCSSel", &gti.Method{Name: "RevertVCSSel", Doc: "RevertVCSSel removes selected files from version control system", Directives: gti.Directives{
			&gti.Directive{Tool: "gti", Directive: "add", Args: []string{}},
		}, Args: ordmap.Make([]ordmap.KeyVal[string, *gti.Field]{}), Returns: ordmap.Make([]ordmap.KeyVal[string, *gti.Field]{})}},
		{"DiffVCSSel", &gti.Method{Name: "DiffVCSSel", Doc: "DiffVCSSel shows the diffs between two versions of selected files, given by the\nrevision specifiers -- if empty, defaults to A = current HEAD, B = current WC file.\n-1, -2 etc also work as universal ways of specifying prior revisions.\nDiffs are shown in a DiffViewDialog.", Directives: gti.Directives{
			&gti.Directive{Tool: "gti", Directive: "add", Args: []string{}},
		}, Args: ordmap.Make([]ordmap.KeyVal[string, *gti.Field]{
			{"rev_a", &gti.Field{Name: "rev_a", Type: "string", LocalType: "string", Doc: "", Directives: gti.Directives{}, Tag: ""}},
			{"rev_b", &gti.Field{Name: "rev_b", Type: "string", LocalType: "string", Doc: "", Directives: gti.Directives{}, Tag: ""}},
		}), Returns: ordmap.Make([]ordmap.KeyVal[string, *gti.Field]{})}},
		{"LogVCSSel", &gti.Method{Name: "LogVCSSel", Doc: "LogVCSSel shows the VCS log of commits for selected files, optionally with a\nsince date qualifier: If since is non-empty, it should be\na date-like expression that the VCS will understand, such as\n1/1/2020, yesterday, last year, etc.  SVN only understands a\nnumber as a maximum number of items to return.\nIf allFiles is true, then the log will show revisions for all files, not just\nthis one.\nReturns the Log and also shows it in a VCSLogView which supports further actions.", Directives: gti.Directives{
			&gti.Directive{Tool: "gti", Directive: "add", Args: []string{}},
		}, Args: ordmap.Make([]ordmap.KeyVal[string, *gti.Field]{
			{"allFiles", &gti.Field{Name: "allFiles", Type: "bool", LocalType: "bool", Doc: "", Directives: gti.Directives{}, Tag: ""}},
			{"since", &gti.Field{Name: "since", Type: "string", LocalType: "string", Doc: "", Directives: gti.Directives{}, Tag: ""}},
		}), Returns: ordmap.Make([]ordmap.KeyVal[string, *gti.Field]{})}},
		{"BlameVCSSel", &gti.Method{Name: "BlameVCSSel", Doc: "BlameVCSSel shows the VCS blame report for this file, reporting for each line\nthe revision and author of the last change.", Directives: gti.Directives{
			&gti.Directive{Tool: "gti", Directive: "add", Args: []string{}},
		}, Args: ordmap.Make([]ordmap.KeyVal[string, *gti.Field]{}), Returns: ordmap.Make([]ordmap.KeyVal[string, *gti.Field]{})}},
	}),
	Instance: &Node{},
})

// NewNode adds a new [Node] with the given name
// to the given parent. If the name is unspecified, it defaults
// to the ID (kebab-case) name of the type, plus the
// [ki.Ki.NumLifetimeChildren] of the given parent.
func NewNode(par ki.Ki, name ...string) *Node {
	return par.NewChild(NodeType, name...).(*Node)
}

// KiType returns the [*gti.Type] of [Node]
func (t *Node) KiType() *gti.Type {
	return NodeType
}

// New returns a new [*Node] value
func (t *Node) New() ki.Ki {
	return &Node{}
}

// NodeEmbedder is an interface that all types that embed Node satisfy
type NodeEmbedder interface {
	AsNode() *Node
}

// AsNode returns the given value as a value of type Node if the type
// of the given value embeds Node, or nil otherwise
func AsNode(k ki.Ki) *Node {
	if k == nil || k.This() == nil {
		return nil
	}
	if t, ok := k.(NodeEmbedder); ok {
		return t.AsNode()
	}
	return nil
}

// AsNode satisfies the [NodeEmbedder] interface
func (t *Node) AsNode() *Node {
	return t
}

// SetTooltip sets the [Node.Tooltip]
func (t *Node) SetTooltip(v string) *Node {
	t.Tooltip = v
	return t
}

// SetClass sets the [Node.Class]
func (t *Node) SetClass(v string) *Node {
	t.Class = v
	return t
}

// SetPriorityEvents sets the [Node.PriorityEvents]
func (t *Node) SetPriorityEvents(v []events.Types) *Node {
	t.PriorityEvents = v
	return t
}

// SetCustomContextMenu sets the [Node.CustomContextMenu]
func (t *Node) SetCustomContextMenu(v func(m *gi.Scene)) *Node {
	t.CustomContextMenu = v
	return t
}

// SetText sets the [Node.Text]
func (t *Node) SetText(v string) *Node {
	t.Text = v
	return t
}

// SetIcon sets the [Node.Icon]
func (t *Node) SetIcon(v icons.Icon) *Node {
	t.Icon = v
	return t
}

// SetIndent sets the [Node.Indent]
func (t *Node) SetIndent(v units.Value) *Node {
	t.Indent = v
	return t
}

// SetOpenDepth sets the [Node.OpenDepth]
func (t *Node) SetOpenDepth(v int) *Node {
	t.OpenDepth = v
	return t
}

// SetViewIdx sets the [Node.ViewIdx]
func (t *Node) SetViewIdx(v int) *Node {
	t.ViewIdx = v
	return t
}

// SetWidgetSize sets the [Node.WidgetSize]
func (t *Node) SetWidgetSize(v mat32.Vec2) *Node {
	t.WidgetSize = v
	return t
}

// SetRootView sets the [Node.RootView]
func (t *Node) SetRootView(v *giv.TreeView) *Node {
	t.RootView = v
	return t
}

// SetSelectedNodes sets the [Node.SelectedNodes]
func (t *Node) SetSelectedNodes(v []giv.TreeViewer) *Node {
	t.SelectedNodes = v
	return t
}

// TreeType is the [gti.Type] for [Tree]
var TreeType = gti.AddType(&gti.Type{
	Name:       "goki.dev/gi/v2/filetree.Tree",
	ShortName:  "filetree.Tree",
	IDName:     "tree",
	Doc:        "Tree is the root of a tree representing files in a given directory\n(and subdirectories thereof), and has some overall management state for how to\nview things.  The Tree can be viewed by a TreeView to provide a GUI\ninterface into it.",
	Directives: gti.Directives{},
	Fields: ordmap.Make([]ordmap.KeyVal[string, *gti.Field]{
		{"ExtFiles", &gti.Field{Name: "ExtFiles", Type: "[]string", LocalType: "[]string", Doc: "external files outside the root path of the tree -- abs paths are stored -- these are shown in the first sub-node if present -- use AddExtFile to add and update", Directives: gti.Directives{}, Tag: ""}},
		{"Dirs", &gti.Field{Name: "Dirs", Type: "goki.dev/gi/v2/filetree.DirFlagMap", LocalType: "DirFlagMap", Doc: "records state of directories within the tree (encoded using paths relative to root),\ne.g., open (have been opened by the user) -- can persist this to restore prior view of a tree", Directives: gti.Directives{}, Tag: ""}},
		{"DirsOnTop", &gti.Field{Name: "DirsOnTop", Type: "bool", LocalType: "bool", Doc: "if true, then all directories are placed at the top of the tree view\notherwise everything is mixed", Directives: gti.Directives{}, Tag: ""}},
		{"NodeType", &gti.Field{Name: "NodeType", Type: "*goki.dev/gti.Type", LocalType: "*gti.Type", Doc: "type of node to create -- defaults to giv.Node but can use custom node types", Directives: gti.Directives{}, Tag: "view:\"-\" json:\"-\" xml:\"-\""}},
		{"InOpenAll", &gti.Field{Name: "InOpenAll", Type: "bool", LocalType: "bool", Doc: "if true, we are in midst of an OpenAll call -- nodes should open all dirs", Directives: gti.Directives{}, Tag: ""}},
		{"Watcher", &gti.Field{Name: "Watcher", Type: "*gopkg.in/fsnotify.v1.Watcher", LocalType: "*fsnotify.Watcher", Doc: "change notify for all dirs", Directives: gti.Directives{}, Tag: "view:\"-\""}},
		{"DoneWatcher", &gti.Field{Name: "DoneWatcher", Type: "chan bool", LocalType: "chan bool", Doc: "channel to close watcher watcher", Directives: gti.Directives{}, Tag: "view:\"-\""}},
		{"WatchedPaths", &gti.Field{Name: "WatchedPaths", Type: "map[string]bool", LocalType: "map[string]bool", Doc: "map of paths that have been added to watcher -- only active if bool = true", Directives: gti.Directives{}, Tag: "view:\"-\""}},
		{"LastWatchUpdt", &gti.Field{Name: "LastWatchUpdt", Type: "string", LocalType: "string", Doc: "last path updated by watcher", Directives: gti.Directives{}, Tag: "view:\"-\""}},
		{"LastWatchTime", &gti.Field{Name: "LastWatchTime", Type: "time.Time", LocalType: "time.Time", Doc: "timestamp of last update", Directives: gti.Directives{}, Tag: "view:\"-\""}},
		{"UpdtMu", &gti.Field{Name: "UpdtMu", Type: "sync.Mutex", LocalType: "sync.Mutex", Doc: "Update mutex", Directives: gti.Directives{}, Tag: "view:\"-\""}},
	}),
	Embeds: ordmap.Make([]ordmap.KeyVal[string, *gti.Field]{
		{"Node", &gti.Field{Name: "Node", Type: "goki.dev/gi/v2/filetree.Node", LocalType: "Node", Doc: "", Directives: gti.Directives{}, Tag: ""}},
	}),
	Methods:  ordmap.Make([]ordmap.KeyVal[string, *gti.Method]{}),
	Instance: &Tree{},
})

// NewTree adds a new [Tree] with the given name
// to the given parent. If the name is unspecified, it defaults
// to the ID (kebab-case) name of the type, plus the
// [ki.Ki.NumLifetimeChildren] of the given parent.
func NewTree(par ki.Ki, name ...string) *Tree {
	return par.NewChild(TreeType, name...).(*Tree)
}

// KiType returns the [*gti.Type] of [Tree]
func (t *Tree) KiType() *gti.Type {
	return TreeType
}

// New returns a new [*Tree] value
func (t *Tree) New() ki.Ki {
	return &Tree{}
}

// SetExtFiles sets the [Tree.ExtFiles]:
// external files outside the root path of the tree -- abs paths are stored -- these are shown in the first sub-node if present -- use AddExtFile to add and update
func (t *Tree) SetExtFiles(v []string) *Tree {
	t.ExtFiles = v
	return t
}

// SetDirs sets the [Tree.Dirs]:
// records state of directories within the tree (encoded using paths relative to root),
// e.g., open (have been opened by the user) -- can persist this to restore prior view of a tree
func (t *Tree) SetDirs(v DirFlagMap) *Tree {
	t.Dirs = v
	return t
}

// SetDirsOnTop sets the [Tree.DirsOnTop]:
// if true, then all directories are placed at the top of the tree view
// otherwise everything is mixed
func (t *Tree) SetDirsOnTop(v bool) *Tree {
	t.DirsOnTop = v
	return t
}

// SetNodeType sets the [Tree.NodeType]:
// type of node to create -- defaults to giv.Node but can use custom node types
func (t *Tree) SetNodeType(v *gti.Type) *Tree {
	t.NodeType = v
	return t
}

// SetInOpenAll sets the [Tree.InOpenAll]:
// if true, we are in midst of an OpenAll call -- nodes should open all dirs
func (t *Tree) SetInOpenAll(v bool) *Tree {
	t.InOpenAll = v
	return t
}

// SetWatcher sets the [Tree.Watcher]:
// change notify for all dirs
func (t *Tree) SetWatcher(v *fsnotify.Watcher) *Tree {
	t.Watcher = v
	return t
}

// SetDoneWatcher sets the [Tree.DoneWatcher]:
// channel to close watcher watcher
func (t *Tree) SetDoneWatcher(v chan bool) *Tree {
	t.DoneWatcher = v
	return t
}

// SetWatchedPaths sets the [Tree.WatchedPaths]:
// map of paths that have been added to watcher -- only active if bool = true
func (t *Tree) SetWatchedPaths(v map[string]bool) *Tree {
	t.WatchedPaths = v
	return t
}

// SetLastWatchUpdt sets the [Tree.LastWatchUpdt]:
// last path updated by watcher
func (t *Tree) SetLastWatchUpdt(v string) *Tree {
	t.LastWatchUpdt = v
	return t
}

// SetLastWatchTime sets the [Tree.LastWatchTime]:
// timestamp of last update
func (t *Tree) SetLastWatchTime(v time.Time) *Tree {
	t.LastWatchTime = v
	return t
}

// SetUpdtMu sets the [Tree.UpdtMu]:
// Update mutex
func (t *Tree) SetUpdtMu(v sync.Mutex) *Tree {
	t.UpdtMu = v
	return t
}

// SetTooltip sets the [Tree.Tooltip]
func (t *Tree) SetTooltip(v string) *Tree {
	t.Tooltip = v
	return t
}

// SetClass sets the [Tree.Class]
func (t *Tree) SetClass(v string) *Tree {
	t.Class = v
	return t
}

// SetPriorityEvents sets the [Tree.PriorityEvents]
func (t *Tree) SetPriorityEvents(v []events.Types) *Tree {
	t.PriorityEvents = v
	return t
}

// SetCustomContextMenu sets the [Tree.CustomContextMenu]
func (t *Tree) SetCustomContextMenu(v func(m *gi.Scene)) *Tree {
	t.CustomContextMenu = v
	return t
}

// SetText sets the [Tree.Text]
func (t *Tree) SetText(v string) *Tree {
	t.Text = v
	return t
}

// SetIcon sets the [Tree.Icon]
func (t *Tree) SetIcon(v icons.Icon) *Tree {
	t.Icon = v
	return t
}

// SetIndent sets the [Tree.Indent]
func (t *Tree) SetIndent(v units.Value) *Tree {
	t.Indent = v
	return t
}

// SetOpenDepth sets the [Tree.OpenDepth]
func (t *Tree) SetOpenDepth(v int) *Tree {
	t.OpenDepth = v
	return t
}

// SetViewIdx sets the [Tree.ViewIdx]
func (t *Tree) SetViewIdx(v int) *Tree {
	t.ViewIdx = v
	return t
}

// SetWidgetSize sets the [Tree.WidgetSize]
func (t *Tree) SetWidgetSize(v mat32.Vec2) *Tree {
	t.WidgetSize = v
	return t
}

// SetRootView sets the [Tree.RootView]
func (t *Tree) SetRootView(v *giv.TreeView) *Tree {
	t.RootView = v
	return t
}

// SetSelectedNodes sets the [Tree.SelectedNodes]
func (t *Tree) SetSelectedNodes(v []giv.TreeViewer) *Tree {
	t.SelectedNodes = v
	return t
}

// VCSLogViewType is the [gti.Type] for [VCSLogView]
var VCSLogViewType = gti.AddType(&gti.Type{
	Name:       "goki.dev/gi/v2/filetree.VCSLogView",
	ShortName:  "filetree.VCSLogView",
	IDName:     "vcs-log-view",
	Doc:        "VCSLogView is a view of the VCS log data",
	Directives: gti.Directives{},
	Fields: ordmap.Make([]ordmap.KeyVal[string, *gti.Field]{
		{"Log", &gti.Field{Name: "Log", Type: "goki.dev/vci/v2.Log", LocalType: "vci.Log", Doc: "current log", Directives: gti.Directives{}, Tag: ""}},
		{"File", &gti.Field{Name: "File", Type: "string", LocalType: "string", Doc: "file that this is a log of -- if blank then it is entire repository", Directives: gti.Directives{}, Tag: ""}},
		{"Since", &gti.Field{Name: "Since", Type: "string", LocalType: "string", Doc: "date expression for how long ago to include log entries from", Directives: gti.Directives{}, Tag: ""}},
		{"Repo", &gti.Field{Name: "Repo", Type: "goki.dev/vci/v2.Repo", LocalType: "vci.Repo", Doc: "version control system repository", Directives: gti.Directives{}, Tag: "json:\"-\" xml:\"-\" copy:\"-\""}},
		{"RevA", &gti.Field{Name: "RevA", Type: "string", LocalType: "string", Doc: "revision A -- defaults to HEAD", Directives: gti.Directives{}, Tag: "set:\"-\""}},
		{"RevB", &gti.Field{Name: "RevB", Type: "string", LocalType: "string", Doc: "revision B -- blank means current working copy", Directives: gti.Directives{}, Tag: "set:\"-\""}},
		{"SetA", &gti.Field{Name: "SetA", Type: "bool", LocalType: "bool", Doc: "double-click will set the A revision -- else B", Directives: gti.Directives{}, Tag: ""}},
	}),
	Embeds: ordmap.Make([]ordmap.KeyVal[string, *gti.Field]{
		{"Layout", &gti.Field{Name: "Layout", Type: "goki.dev/gi/v2/gi.Layout", LocalType: "gi.Layout", Doc: "", Directives: gti.Directives{}, Tag: ""}},
	}),
	Methods:  ordmap.Make([]ordmap.KeyVal[string, *gti.Method]{}),
	Instance: &VCSLogView{},
})

// NewVCSLogView adds a new [VCSLogView] with the given name
// to the given parent. If the name is unspecified, it defaults
// to the ID (kebab-case) name of the type, plus the
// [ki.Ki.NumLifetimeChildren] of the given parent.
func NewVCSLogView(par ki.Ki, name ...string) *VCSLogView {
	return par.NewChild(VCSLogViewType, name...).(*VCSLogView)
}

// KiType returns the [*gti.Type] of [VCSLogView]
func (t *VCSLogView) KiType() *gti.Type {
	return VCSLogViewType
}

// New returns a new [*VCSLogView] value
func (t *VCSLogView) New() ki.Ki {
	return &VCSLogView{}
}

// SetLog sets the [VCSLogView.Log]:
// current log
func (t *VCSLogView) SetLog(v vci.Log) *VCSLogView {
	t.Log = v
	return t
}

// SetFile sets the [VCSLogView.File]:
// file that this is a log of -- if blank then it is entire repository
func (t *VCSLogView) SetFile(v string) *VCSLogView {
	t.File = v
	return t
}

// SetSince sets the [VCSLogView.Since]:
// date expression for how long ago to include log entries from
func (t *VCSLogView) SetSince(v string) *VCSLogView {
	t.Since = v
	return t
}

// SetRepo sets the [VCSLogView.Repo]:
// version control system repository
func (t *VCSLogView) SetRepo(v vci.Repo) *VCSLogView {
	t.Repo = v
	return t
}

// SetSetA sets the [VCSLogView.SetA]:
// double-click will set the A revision -- else B
func (t *VCSLogView) SetSetA(v bool) *VCSLogView {
	t.SetA = v
	return t
}

// SetTooltip sets the [VCSLogView.Tooltip]
func (t *VCSLogView) SetTooltip(v string) *VCSLogView {
	t.Tooltip = v
	return t
}

// SetClass sets the [VCSLogView.Class]
func (t *VCSLogView) SetClass(v string) *VCSLogView {
	t.Class = v
	return t
}

// SetPriorityEvents sets the [VCSLogView.PriorityEvents]
func (t *VCSLogView) SetPriorityEvents(v []events.Types) *VCSLogView {
	t.PriorityEvents = v
	return t
}

// SetCustomContextMenu sets the [VCSLogView.CustomContextMenu]
func (t *VCSLogView) SetCustomContextMenu(v func(m *gi.Scene)) *VCSLogView {
	t.CustomContextMenu = v
	return t
}

// SetStackTop sets the [VCSLogView.StackTop]
func (t *VCSLogView) SetStackTop(v int) *VCSLogView {
	t.StackTop = v
	return t
}
