// Copyright (c) 2023, Cogent Core. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package filetree

//go:generate core generate

import (
	"fmt"
	"log"
	"log/slog"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"cogentcore.org/core/base/errors"
	"cogentcore.org/core/base/fileinfo"
	"cogentcore.org/core/base/fsx"
	"cogentcore.org/core/base/vcs"
	"cogentcore.org/core/colors"
	"cogentcore.org/core/colors/gradient"
	"cogentcore.org/core/core"
	"cogentcore.org/core/events"
	"cogentcore.org/core/events/key"
	"cogentcore.org/core/icons"
	"cogentcore.org/core/keymap"
	"cogentcore.org/core/styles"
	"cogentcore.org/core/styles/units"
	"cogentcore.org/core/texteditor"
	"cogentcore.org/core/texteditor/histyle"
	"cogentcore.org/core/tree"
	"cogentcore.org/core/types"
)

// NodeHiStyle is the default style for syntax highlighting to use for
// file node buffers
var NodeHiStyle = histyle.StyleDefault

// Node represents a file in the file system, as a [core.Tree] node.
// The name of the node is the name of the file.
// Folders have children containing further nodes.
type Node struct { //core:embedder
	core.Tree

	// Filepath is the full path to this file.
	Filepath core.Filename `edit:"-" set:"-" json:"-" xml:"-" copier:"-"`

	// Info is the full standard file info about this file.
	Info fileinfo.FileInfo `edit:"-" set:"-" json:"-" xml:"-" copier:"-"`

	// Buffer is the file buffer for editing this file.
	Buffer *texteditor.Buffer `edit:"-" set:"-" json:"-" xml:"-" copier:"-"`

	// FileRoot is the root [Tree] of the tree, which has global state.
	FileRoot *Tree `edit:"-" set:"-" json:"-" xml:"-" copier:"-"`

	// DirRepo is the version control system repository for this directory,
	// only non-nil if this is the highest-level directory in the tree under vcs control.
	DirRepo vcs.Repo `edit:"-" set:"-" json:"-" xml:"-" copier:"-"`

	// RepoFiles has the version control system repository file status; only valid during SetPath.
	RepoFiles vcs.Files `edit:"-" set:"-" json:"-" xml:"-" copier:"-"`
}

func (fn *Node) Init() {
	fn.Tree.Init()
	fn.ContextMenus = nil // do not include tree
	fn.AddContextMenu(fn.ContextMenu)
	fn.Styler(func(s *styles.Style) {
		status := fn.Info.VCS
		s.Font.Weight = styles.WeightNormal
		s.Font.Style = styles.FontNormal
		if fn.IsExec() && !fn.IsDir() {
			s.Font.Weight = styles.WeightBold // todo: somehow not working
		}
		if fn.Buffer != nil {
			s.Font.Style = styles.Italic
		}
		switch {
		case status == vcs.Untracked:
			s.Color = errors.Must1(gradient.FromString("#808080"))
		case status == vcs.Modified:
			s.Color = errors.Must1(gradient.FromString("#4b7fd1"))
		case status == vcs.Added:
			s.Color = errors.Must1(gradient.FromString("#008800"))
		case status == vcs.Deleted:
			s.Color = errors.Must1(gradient.FromString("#ff4252"))
		case status == vcs.Conflicted:
			s.Color = errors.Must1(gradient.FromString("#ce8020"))
		case status == vcs.Updated:
			s.Color = errors.Must1(gradient.FromString("#008060"))
		case status == vcs.Stored:
			s.Color = colors.C(colors.Scheme.OnSurface)
		}
	})
	fn.On(events.KeyChord, func(e events.Event) {
		if core.DebugSettings.KeyEventTrace {
			fmt.Printf("Tree KeyInput: %v\n", fn.Path())
		}
		kf := keymap.Of(e.KeyChord())
		selMode := events.SelectModeBits(e.Modifiers())

		if selMode == events.SelectOne {
			if fn.SelectMode {
				selMode = events.ExtendContinuous
			}
		}

		// first all the keys that work for ReadOnly and active
		if !fn.IsReadOnly() && !e.IsHandled() {
			switch kf {
			case keymap.Delete:
				fn.DeleteFiles()
				e.SetHandled()
			case keymap.Backspace:
				fn.DeleteFiles()
				e.SetHandled()
			case keymap.Duplicate:
				fn.DuplicateFiles()
				e.SetHandled()
			case keymap.Insert: // New File
				core.CallFunc(fn, fn.NewFile)
				e.SetHandled()
			case keymap.InsertAfter: // New Folder
				core.CallFunc(fn, fn.NewFolder)
				e.SetHandled()
			}
		}
	})

	fn.Parts.Styler(func(s *styles.Style) {
		s.Gap.X.Em(0.4)
	})
	fn.Parts.OnClick(func(e events.Event) {
		if !e.HasAnyModifier(key.Control, key.Meta, key.Alt, key.Shift) {
			fn.OpenEmptyDir()
		}
	})
	fn.Parts.OnDoubleClick(func(e events.Event) {
		if fn.FileRoot != nil && fn.FileRoot.DoubleClickFunc != nil {
			fn.FileRoot.DoubleClickFunc(e)
		} else {
			if fn.IsDir() && fn.OpenEmptyDir() {
				e.SetHandled()
			}
		}
	})
	core.AddChildInit(fn.Parts, "branch", func(w *core.Switch) {
		w.SetType(core.SwitchCheckbox)
		w.SetIcons(icons.FolderOpen, icons.Folder, icons.Blank)
		core.AddChildInit(w, "stack", func(w *core.Frame) {
			f := func(name string) {
				core.AddChildInit(w, name, func(w *core.Icon) {
					w.Styler(func(s *styles.Style) {
						s.Min.Set(units.Em(1))
					})
				})
			}
			f("icon-on")
			f("icon-off")
			f("icon-indeterminate")
		})
	})
}

func (fn *Node) BaseType() *types.Type {
	return fn.NodeType()
}

// IsDir returns true if file is a directory (folder)
func (fn *Node) IsDir() bool {
	return fn.Info.IsDir()
}

// IsIrregular  returns true if file is a special "Irregular" node
func (fn *Node) IsIrregular() bool {
	return (fn.Info.Mode & os.ModeIrregular) != 0
}

// IsExternal returns true if file is external to main file tree
func (fn *Node) IsExternal() bool {
	isExt := false
	fn.WalkUp(func(k tree.Node) bool {
		sfn := AsNode(k)
		if sfn == nil {
			return tree.Break
		}
		if sfn.IsIrregular() {
			isExt = true
			return tree.Break
		}
		return tree.Continue
	})
	return isExt
}

// HasClosedParent returns true if node has a parent node with !IsOpen flag set
func (fn *Node) HasClosedParent() bool {
	hasClosed := false
	fn.WalkUpParent(func(k tree.Node) bool {
		sfn := AsNode(k)
		if sfn == nil {
			return tree.Break
		}
		if !sfn.IsOpen() {
			hasClosed = true
			return tree.Break
		}
		return tree.Continue
	})
	return hasClosed
}

// IsExec returns true if file is an executable file
func (fn *Node) IsExec() bool {
	return fn.Info.IsExec()
}

// IsOpen returns true if file is flagged as open
func (fn *Node) IsOpen() bool {
	return !fn.Closed
}

// IsChanged returns true if the file is open and has been changed (edited) since last EditDone
func (fn *Node) IsChanged() bool {
	return fn.Buffer != nil && fn.Buffer.Changed
}

// IsNotSaved returns true if the file is open and has been changed (edited) since last Save
func (fn *Node) IsNotSaved() bool {
	return fn.Buffer != nil && fn.Buffer.NotSaved
}

// IsAutoSave returns true if file is an auto-save file (starts and ends with #)
func (fn *Node) IsAutoSave() bool {
	return strings.HasPrefix(fn.Info.Name, "#") && strings.HasSuffix(fn.Info.Name, "#")
}

// MyRelPath returns the relative path from root for this node
func (fn *Node) MyRelPath() string {
	if fn.IsIrregular() || fn.FileRoot == nil {
		return fn.Name
	}
	return fsx.RelativeFilePath(string(fn.Filepath), string(fn.FileRoot.Filepath))
}

// SetPath sets the current node to represent the given path.
// This then calls [SyncDir] to synchronize the tree with the file
// system tree at this path.
func (fn *Node) SetPath(path string) error {
	_, fnm := filepath.Split(path)
	fn.SetText(fnm)
	pth, err := filepath.Abs(path)
	if err != nil {
		return err
	}
	fn.Filepath = core.Filename(pth)
	err = fn.Info.InitFile(string(fn.Filepath))
	if err != nil {
		log.Printf("core.Tree: could not read directory: %v err: %v\n", fn.Filepath, err)
		return err
	}

	fn.SyncDir()
	return nil
}

// SyncDir synchronizes the current directory node with all the files
// contained within the directory on the filesystem, using the efficient
// Plan-based diff-driven updating of only what is different.
func (fn *Node) SyncDir() {
	fn.DetectVCSRepo(true) // update files
	path := string(fn.Filepath)
	repo, rnode := fn.Repo()
	fn.Open() // ensure
	plan := fn.PlanOfFiles(path)
	hasExtFiles := false
	if fn.This == fn.FileRoot.This {
		if len(fn.FileRoot.ExternalFiles) > 0 {
			plan = append(tree.TypePlan{{Type: fn.FileRoot.FileNodeType, Name: ExternalFilesName}}, plan...)
			hasExtFiles = true
		}
	}
	mods := tree.Update(fn, plan)
	// always go through kids, regardless of mods
	for _, sfk := range fn.Children {
		sf := AsNode(sfk)
		sf.FileRoot = fn.FileRoot
		if hasExtFiles && sf.Name == ExternalFilesName {
			fn.FileRoot.SyncExternalFiles(sf)
			continue
		}
		fp := filepath.Join(path, sf.Name)
		// if sf.Buf != nil {
		// 	fmt.Printf("fp: %v  nm: %v\n", fp, sf.Nm)
		// }
		sf.SetNodePath(fp)
		if sf.IsDir() {
			sf.Info.VCS = vcs.Stored // always
		} else if repo != nil {
			rstat := rnode.RepoFiles.Status(repo, string(sf.Filepath))
			sf.Info.VCS = rstat
		} else {
			sf.Info.VCS = vcs.Stored
		}
	}
	if mods {
		root := fn.FileRoot
		fn.Update()
		if root != nil {
			root.TreeChanged(nil)
		}
	}
}

// PlanOfFiles returns a tree.TypePlan for building nodes based on
// files immediately within given path.
func (fn *Node) PlanOfFiles(path string) tree.TypePlan {
	plan1 := tree.TypePlan{}
	plan2 := tree.TypePlan{}
	typ := fn.FileRoot.FileNodeType
	filepath.Walk(path, func(pth string, info os.FileInfo, err error) error {
		if err != nil {
			emsg := fmt.Sprintf("filetree.Node PlanFilesIn Path %q: Error: %v", path, err)
			log.Println(emsg)
			return nil // ignore
		}
		if pth == path { // proceed..
			return nil
		}
		_, fnm := filepath.Split(pth)
		if fn.FileRoot.DirsOnTop {
			if info.IsDir() {
				plan1.Add(typ, fnm)
			} else {
				plan2.Add(typ, fnm)
			}
		} else {
			plan1.Add(typ, fnm)
		}
		if info.IsDir() {
			return filepath.SkipDir
		}
		return nil
	})
	modSort := fn.FileRoot.DirSortByModTime(core.Filename(path))
	if fn.FileRoot.DirsOnTop {
		if modSort {
			fn.SortPlanByModTime(plan2) // just sort files, not dirs
		}
		plan1 = append(plan1, plan2...)
	} else {
		if modSort {
			fn.SortPlanByModTime(plan1) // all
		}
	}
	return plan1
}

// SortPlanByModTime sorts given plan list by mod time
func (fn *Node) SortPlanByModTime(confg tree.TypePlan) {
	sort.Slice(confg, func(i, j int) bool {
		ifn, _ := os.Stat(filepath.Join(string(fn.Filepath), confg[i].Name))
		jfn, _ := os.Stat(filepath.Join(string(fn.Filepath), confg[j].Name))
		return ifn.ModTime().After(jfn.ModTime()) // descending
	})
}

func (fn *Node) SetFileIcon() {
	ic, hasic := fn.Info.FindIcon()
	if !hasic {
		ic = icons.Blank
	}
	if _, ok := fn.Branch(); !ok {
		fn.Update()
	}
	if bp, ok := fn.Branch(); ok {
		if bp.IconIndeterminate != ic {
			bp.SetIcons(icons.FolderOpen, icons.Folder, ic)
			fn.Update()
		}
	}
}

// SetNodePath sets the path for given node and updates it based on associated file
func (fn *Node) SetNodePath(path string) error {
	pth, err := filepath.Abs(path)
	if err != nil {
		return err
	}
	fn.Filepath = core.Filename(pth)
	err = fn.InitFileInfo()
	if err != nil {
		return err
	}
	if fn.IsDir() && !fn.IsIrregular() {
		openAll := fn.FileRoot.inOpenAll && !fn.Info.IsHidden()
		if openAll || fn.FileRoot.IsDirOpen(fn.Filepath) {
			fn.SetPath(string(fn.Filepath)) // keep going down..
		}
	}
	fn.SetFileIcon()
	return nil
}

// InitFileInfo initializes file info
func (fn *Node) InitFileInfo() error {
	effpath, err := filepath.EvalSymlinks(string(fn.Filepath))
	if err != nil {
		// this happens too often for links -- skip
		// log.Printf("filetree.Node Path: %v could not be opened -- error: %v\n", fn.FPath, err)
		return err
	}
	fn.Filepath = core.Filename(effpath)
	err = fn.Info.InitFile(string(fn.Filepath))
	if err != nil {
		emsg := fmt.Errorf("filetree.Node InitFileInfo Path %q: Error: %v", fn.Filepath, err)
		log.Println(emsg)
		return emsg
	}
	return nil
}

// UpdateNode updates information in node based on its associated file in FPath.
func (fn *Node) UpdateNode() error {
	err := fn.InitFileInfo()
	if err != nil {
		return err
	}
	if fn.IsIrregular() {
		return nil
	}
	// fmt.Println(fn, "update node start")
	if fn.IsDir() {
		openAll := fn.FileRoot.inOpenAll && !fn.Info.IsHidden()
		if openAll || fn.FileRoot.IsDirOpen(fn.Filepath) {
			// fmt.Printf("set open: %s\n", fn.FPath)
			fn.Open()
			repo, rnode := fn.Repo()
			if repo != nil {
				rnode.UpdateRepoFiles()
			}
			fn.SyncDir()
		}
	} else {
		repo, _ := fn.Repo()
		if repo != nil {
			fn.Info.VCS, _ = repo.Status(string(fn.Filepath))
		}
		fn.Update()
		fn.SetFileIcon()
	}
	// fmt.Println(fn, "update node end")
	return nil
}

// SelectedFunc runsthe given function on all selected nodes in reverse order.
func (fn *Node) SelectedFunc(fun func(n *Node)) {
	sels := fn.SelectedViews()
	for i := len(sels) - 1; i >= 0; i-- {
		sn := AsNode(sels[i])
		if sn == nil {
			continue
		}
		fun(sn)
	}
}

// OpenDirs opens directories for selected views
func (fn *Node) OpenDirs() {
	fn.SelectedFunc(func(sn *Node) {
		sn.OpenDir()
	})
}

func (fn *Node) OnOpen() {
	fn.OpenDir()
}

func (fn *Node) CanOpen() bool {
	return fn.HasChildren() || fn.IsDir()
}

// OpenDir opens given directory node
func (fn *Node) OpenDir() {
	// fmt.Printf("fn: %s opened\n", fn.FPath)
	fn.FileRoot.SetDirOpen(fn.Filepath)
	fn.UpdateNode()
}

func (fn *Node) OnClose() {
	fn.CloseDir()
}

// CloseDir closes given directory node -- updates memory state
func (fn *Node) CloseDir() {
	// fmt.Printf("fn: %s closed\n", fn.FPath)
	fn.FileRoot.SetDirClosed(fn.Filepath)
	// note: not doing anything with open files within directory..
}

// OpenEmptyDir will attempt to open a directory that has no children
// which presumably was not processed originally
func (fn *Node) OpenEmptyDir() bool {
	if fn.IsDir() && !fn.HasChildren() {
		fn.OpenDir()
		fn.NeedsLayout()
		return true
	}
	return false
}

// SortBys determines how to sort the selected files in the directory.
// Default is alpha by name, optionally can be sorted by modification time.
func (fn *Node) SortBys(modTime bool) { //types:add
	fn.SelectedFunc(func(sn *Node) {
		sn.SortBy(modTime)
	})
}

// SortBy determines how to sort the files in the directory -- default is alpha by name,
// optionally can be sorted by modification time.
func (fn *Node) SortBy(modTime bool) {
	fn.FileRoot.SetDirSortBy(fn.Filepath, modTime)
	fn.NeedsLayout()
}

// OpenAll opens all directories under this one
func (fn *Node) OpenAll() { //types:add
	fn.FileRoot.inOpenAll = true // causes chaining of opening
	fn.Tree.OpenAll()
	fn.FileRoot.inOpenAll = false
}

// CloseAll closes all directories under this one, this included
func (fn *Node) CloseAll() { //types:add
	fn.WidgetWalkDown(func(wi core.Widget, wb *core.WidgetBase) bool {
		sfn := AsNode(wi)
		if sfn == nil {
			return tree.Continue
		}
		if sfn.IsDir() {
			sfn.Close()
		}
		return tree.Continue
	})
}

// OpenBuf opens the file in its buffer if it is not already open.
// returns true if file is newly opened
func (fn *Node) OpenBuf() (bool, error) {
	if fn.IsDir() {
		err := fmt.Errorf("filetree.Node cannot open directory in editor: %v", fn.Filepath)
		log.Println(err)
		return false, err
	}
	if fn.Buffer != nil {
		if fn.Buffer.Filename == fn.Filepath { // close resets filename
			return false, nil
		}
	} else {
		fn.Buffer = texteditor.NewBuffer()
		fn.Buffer.OnChange(func(e events.Event) {
			if fn.Info.VCS == vcs.Stored {
				fn.Info.VCS = vcs.Modified
			}
		})
	}
	fn.Buffer.Hi.Style = NodeHiStyle
	return true, fn.Buffer.Open(fn.Filepath)
}

// RemoveFromExterns removes file from list of external files
func (fn *Node) RemoveFromExterns() { //types:add
	fn.SelectedFunc(func(sn *Node) {
		if !sn.IsExternal() {
			return
		}
		sn.FileRoot.RemoveExternalFile(string(sn.Filepath))
		sn.CloseBuf()
		sn.Delete()
	})
}

// CloseBuf closes the file in its buffer if it is open.
// returns true if closed.
func (fn *Node) CloseBuf() bool {
	if fn.Buffer == nil {
		return false
	}
	fn.Buffer.Close(nil)
	fn.Buffer = nil
	return true
}

// RelPath returns the relative path from node for given full path
func (fn *Node) RelPath(fpath core.Filename) string {
	return fsx.RelativeFilePath(string(fpath), string(fn.Filepath))
}

// DirsTo opens all the directories above the given filename, and returns the node
// for element at given path (can be a file or directory itself -- not opened -- just returned)
func (fn *Node) DirsTo(path string) (*Node, error) {
	pth, err := filepath.Abs(path)
	if err != nil {
		log.Printf("filetree.Node DirsTo path %v could not be turned into an absolute path: %v\n", path, err)
		return nil, err
	}
	rpath := fn.RelPath(core.Filename(pth))
	if rpath == "." {
		return fn, nil
	}
	dirs := strings.Split(rpath, string(filepath.Separator))
	cfn := fn
	sz := len(dirs)
	for i := 0; i < sz; i++ {
		dr := dirs[i]
		sfni := cfn.ChildByName(dr, 0)
		if sfni == nil {
			if i == sz-1 { // ok for terminal -- might not exist yet
				return cfn, nil
			} else {
				err = fmt.Errorf("filetree.Node could not find node %v in: %v, orig: %v, rel: %v", dr, cfn.Filepath, pth, rpath)
				// slog.Error(err.Error()) // note: this is expected sometimes
				return nil, err
			}
		}
		sfn := AsNode(sfni)
		if sfn.IsDir() || i == sz-1 {
			if i < sz-1 && !sfn.IsOpen() {
				sfn.OpenDir()
				sfn.UpdateNode()
			} else {
				cfn = sfn
			}
		} else {
			err := fmt.Errorf("filetree.Node non-terminal node %v is not a directory in: %v", dr, cfn.Filepath)
			slog.Error(err.Error())
			return nil, err
		}
		cfn = sfn
	}
	return cfn, nil
}
