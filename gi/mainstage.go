// Copyright (c) 2023, The GoKi Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gi

import (
	"fmt"
	"image"
	"log/slog"

	"goki.dev/goosi"
	"goki.dev/goosi/events"
	"goki.dev/ki/v2"
	"goki.dev/mat32/v2"
)

func (st *Stage) RenderCtx() *RenderContext {
	if st.StageMgr == nil {
		slog.Error("Stage has nil StageMgr", "stage", st.Name)
		return nil
	}
	return st.StageMgr.RenderCtx
}

// NewMainStage returns a new MainStage with given type and scene contents.
// Make further configuration choices using Set* methods, which
// can be chained directly after the NewMainStage call.
// Use an appropriate Run call at the end to start the Stage running.
func NewMainStage(typ StageTypes, sc *Scene) *Stage {
	st := &Stage{}
	st.SetType(typ)
	st.SetScene(sc)
	st.PopupMgr.Main = st
	return st
}

// NewWindow returns a new Window stage with given scene contents.
// Make further configuration choices using Set* methods, which
// can be chained directly after the New call.
// Use an appropriate Run call at the end to start the Stage running.
func (sc *Scene) NewWindow() *Stage {
	ms := NewMainStage(WindowStage, sc)
	ms.SetNewWindow(true)
	return ms
}

// NewWindow returns a new Window stage with given scene contents.
// Make further configuration choices using Set* methods, which
// can be chained directly after the New call.
// Use an appropriate Run call at the end to start the Stage running.
func (bd *Body) NewWindow() *Stage {
	return bd.Sc.NewWindow()
}

// NewDialog in dialogs.go

// NewSheet returns a new Sheet stage with given scene contents,
// in connection with given widget (which provides key context).
// for given side (e.g., Bottom or LeftSide).
// Make further configuration choices using Set* methods, which
// can be chained directly after the New call.
// Use an appropriate Run call at the end to start the Stage running.
func NewSheet(sc *Scene, side StageSides) *Stage {
	return NewMainStage(SheetStage, sc).SetSide(side).AsMain()
}

/////////////////////////////////////////////////////
//		Decorate

// SetWindowInsets updates the padding on the Scene
// to the inset values provided by the RenderWin window.
func (st *Stage) SetWindowInsets() {
	if st.StageMgr == nil {
		return
	}
	if st.StageMgr.RenderWin == nil {
		return
	}
	// insets := st.StageMgr.RenderWin.GoosiWin.Insets()
	// // fmt.Println(insets)
	// uv := func(val float32) units.Value {
	// 	return units.Custom(func(uc *units.Context) float32 {
	// 		return max(val, uc.Dp(12))
	// 	})
	// }
	// st.Scene.Style(func(s *styles.Style) {
	// 	s.Padding.Set(
	// 		uv(insets.Top),
	// 		uv(insets.Right),
	// 		uv(insets.Bottom),
	// 		uv(insets.Left),
	// 	)
	// })
}

// only called when !NewWindow
func (st *Stage) AddWindowDecor() *Stage {
	return st
}

func (st *Stage) AddDialogDecor() *Stage {
	return st
}

func (st *Stage) AddSheetDecor() *Stage {
	// todo: handle based on side
	return st
}

func (st *Stage) InheritBars() {
	st.Scene.InheritBarsWidget(st.Context)
}

// FirstWinManager creates a MainStageMgr for the first window
// to be able to get sizing information prior to having a RenderWin,
// based on the goosi App Screen Size. Only adds a RenderCtx.
func (st *Stage) FirstWinManager() *StageMgr {
	ms := &MainStageMgr{}
	ms.This = ms
	rc := &RenderContext{}
	ms.RenderCtx = rc
	scr := goosi.TheApp.Screen(0)
	rc.Size = scr.Geometry.Size()
	// fmt.Println("Screen Size:", rc.Size)
	rc.SetFlag(true, RenderVisible)
	rc.LogicalDPI = scr.LogicalDPI
	// fmt.Println("first win:", rc.LogicalDPI)
	return ms
}

// RunWindow runs a Window with current settings.
func (st *Stage) RunWindow() *Stage {
	st.AddWindowDecor() // sensitive to cases
	sc := st.Scene
	sc.ConfigSceneBars()
	sc.ConfigSceneWidgets()

	// note: need a StageMgr to get initial pref size
	if CurRenderWin == nil {
		st.StageMgr = st.FirstWinManager()
	} else {
		st.StageMgr = &CurRenderWin.StageMgr
	}
	sz := st.RenderCtx().Size
	// non-new full windows must take up the whole window
	// and thus don't consider pref size
	if st.NewWindow || !st.FullWindow {
		sz = sc.PrefSize(sz)
	}
	if WinRenderTrace {
		fmt.Println("MainStage.RunWindow: Window Size:", sz)
	}

	if st.NewWindow {
		sc.Resize(sz)
		win := st.NewRenderWin()
		if CurRenderWin == nil {
			CurRenderWin = win
		}
		st.SetWindowInsets()
		win.GoStartEventLoop()
		return st
	}
	if CurRenderWin == nil {
		sc.Resize(sz)
		CurRenderWin = st.NewRenderWin()
		st.SetWindowInsets()
		CurRenderWin.GoStartEventLoop()
		return st
	}
	if st.Context != nil {
		ms := st.Context.AsWidget().Sc.MainStageMgr()
		msc := ms.Top().AsMain().Scene
		sc.SceneGeom.Size = sz
		sc.FitInWindow(msc.SceneGeom) // does resize
		ms.Push(st)
	} else {
		msc := st.StageMgr.Top().AsMain().Scene
		sc.SceneGeom.Size = sz
		sc.FitInWindow(msc.SceneGeom) // does resize
		CurRenderWin.StageMgr.Push(st)
	}
	return st
}

// RunDialog runs a Dialog with current settings.
// RenderWin field will be set to the parent RenderWin window.
func (st *Stage) RunDialog() *Stage {
	ctx := st.Context.AsWidget()
	ms := ctx.Sc.MainStageMgr()

	// if our main stage manager is nil, we wait until our context is shown and then try again
	if ms == nil {
		slog.Error("RunDialog: CurRenderWin is nil")
		ctx.OnShow(func(e events.Event) {
			st.RunDialog()
		})
		return st
	}

	sc := st.Scene

	st.AddDialogDecor()
	sc.ConfigSceneBars()
	sc.ConfigSceneWidgets()
	sc.SceneGeom.Pos = ctx.ContextMenuPos(nil)

	st.StageMgr = ms // temporary
	winsz := ms.RenderCtx.Size

	sz := winsz
	// history-based stages always take up the whole window
	if !st.FullWindow {
		sz = sc.PrefSize(winsz)
		sz = sz.Add(image.Point{50, 50})
		sc.EventMgr.StartFocusFirst = true // fallback
	}
	if WinRenderTrace {
		slog.Info("MainStage.RunDialog", "size", sz)
	}

	if st.NewWindow && !goosi.TheApp.Platform().IsMobile() {
		sc.Resize(sz)
		st.Type = WindowStage            // critical: now is its own window!
		sc.SceneGeom.Pos = image.Point{} // ignore pos
		win := st.NewRenderWin()
		DialogRenderWins.Add(win)
		win.GoStartEventLoop()
		return st
	}
	winGeom := mat32.Geom2DInt{Size: winsz}
	sc.SceneGeom.Size = sz
	// fmt.Println("dlg:", sc.SceneGeom, "win:", winGeom)
	sc.FitInWindow(winGeom) // does resize

	ms.Push(st)
	return st
}

// RunSheet runs a Sheet with current settings.
// RenderWin field will be set to the parent RenderWin window.
func (st *Stage) RunSheet() *Stage {
	st.AddSheetDecor()
	st.Scene.ConfigSceneBars()
	st.Scene.ConfigSceneWidgets()

	if CurRenderWin == nil {
		// todo: error here -- must have main window!
		return nil
	}
	// todo: need some kind of linkage here for dialog relative to existing window
	// probably just CurRenderWin but it needs to be a stack or updated properly etc.
	CurRenderWin.StageMgr.Push(st)
	return st
}

func (st *Stage) NewRenderWin() *RenderWin {
	if st.Scene == nil {
		slog.Error("MainStage.NewRenderWin: Scene is nil")
	}
	name := st.Name
	title := st.Title
	opts := &goosi.NewWindowOptions{
		Title: title, Size: st.Scene.SceneGeom.Size, StdPixels: false,
	}
	wgp := WinGeomMgr.Pref(title, nil)
	if goosi.TheApp.Platform() != goosi.Offscreen && wgp != nil {
		WinGeomMgr.SettingStart()
		opts.Size = wgp.Size()
		opts.Pos = wgp.Pos()
		opts.StdPixels = false
		// fmt.Printf("got prefs for %v: size: %v pos: %v\n", name, opts.Size, opts.Pos)
		if _, found := AllRenderWins.FindName(name); found { // offset from existing
			opts.Pos.X += 20
			opts.Pos.Y += 20
		}
		if wgp.Fullscreen {
			opts.SetFullscreen()
		}
	}
	win := NewRenderWin(name, title, opts)
	WinGeomMgr.SettingEnd()
	if win == nil {
		return nil
	}
	if wgp != nil {
		win.SetFlag(true, WinHasGeomPrefs)
	}
	AllRenderWins.Add(win)
	MainRenderWins.Add(win)
	WinNewCloseStamp()
	win.StageMgr.Push(st)
	return win
}

func (st *Stage) Delete() {
	st.PopupMgr.CloseAll()
	if st.Scene != nil {
		st.Scene.Delete(ki.DestroyKids)
	}
	st.Scene = nil
	st.StageMgr = nil
}

func (st *Stage) Resize(sz image.Point) {
	if st.Scene == nil {
		return
	}
	switch st.Type {
	case WindowStage:
		st.SetWindowInsets()
		st.Scene.Resize(sz)
	case DialogStage:
		if st.FullWindow {
			st.Scene.Resize(sz)
		}
		// todo: other types fit in constraints
	}
}

// MainDoUpdate calls DoUpdate on our Scene and UpdateAll on our Popups
// returns stageMods = true if any Popup Stages have been modified
// and sceneMods = true if any Scenes have been modified.
func (st *Stage) MainDoUpdate() (stageMods, sceneMods bool) {
	if st.Scene == nil {
		return
	}
	stageMods, sceneMods = st.PopupMgr.UpdateAll()
	scMod := st.Scene.DoUpdate()
	sceneMods = sceneMods || scMod
	// if scMod {
	// 	fmt.Println("main scene mod:", st.Scene.Name)
	// }
	// if stageMods {
	// 	fmt.Println("pop stage mod:", st.Name)
	// }
	return
}

func (st *Stage) MainStageAdded(smi StageMgr) {
	st.StageMgr = smi.AsMainMgr()
}

// MainHandleEvent handles main stage events
func (st *Stage) MainHandleEvent(evi events.Event) {
	if st.Scene == nil {
		return
	}
	st.PopupMgr.HandleEvent(evi)
	if evi.IsHandled() || st.PopupMgr.TopIsModal() {
		if EventTrace && evi.Type() != events.MouseMove {
			fmt.Println("Event handled by popup:", evi)
		}
		return
	}
	evi.SetLocalOff(st.Scene.SceneGeom.Pos)
	st.Scene.EventMgr.HandleEvent(evi)
}

//////////////////////////////////////////////////////////////////////////////
//		Main StageMgr

func (sm *StageMgr) SetRenderWin(win *RenderWin) {
	sm.RenderWin = win
	sm.RenderCtx = &RenderContext{LogicalDPI: 96}
}

// resize resizes all main stages
func (sm *StageMgr) Resize(sz image.Point) {
	for _, kv := range sm.Stack.Order {
		st := kv.Val.AsMain()
		st.Resize(sz)
	}
}

func (sm *StageMgr) MainHandleEvent(evi events.Event) {
	n := sm.Stack.Len()
	for i := n - 1; i >= 0; i-- {
		st := sm.Stack.ValByIdx(i).AsMain()
		st.HandleEvent(evi)
		if evi.IsHandled() || st.Modal || st.Type == WindowStage {
			return
		}
	}
}

/*
todo: main menu on full win

// ConfigVLay creates and configures the vertical layout as first child of
// Scene, and installs MainMenu as first element of layout.
func (w *RenderWin) ConfigVLay() {
	sc := w.Scene
	updt := sc.UpdateStart()
	defer sc.UpdateEnd(updt)
	if !sc.HasChildren() {
		sc.NewChild(LayoutType, "main-vlay")
	}
	w.Scene.Frame = sc.Child(0).Embed(LayoutType).(*Layout)
	if !w.Scene.Frame.HasChildren() {
		w.Scene.Frame.NewChild(TypeMenuBar, "main-menu")
	}
	w.MainMenu = w.Scene.Frame.Child(0).(*MenuBar)
	w.MainMenu.MainMenu = true
	w.MainMenu.SetStretchMaxWidth()
}

// AddMainMenu installs MainMenu as first element of main layout
// used for dialogs that don't always have a main menu -- returns
// menubar -- safe to call even if there is a menubar
func (w *RenderWin) AddMainMenu() *MenuBar {
	sc := w.Scene
	updt := sc.UpdateStart()
	defer sc.UpdateEnd(updt)
	if !sc.HasChildren() {
		sc.NewChild(LayoutType, "main-vlay")
	}
	w.Scene.Frame = sc.Child(0).Embed(LayoutType).(*Layout)
	if !w.Scene.Frame.HasChildren() {
		w.MainMenu = w.Scene.Frame.NewChild(TypeMenuBar, "main-menu").(*MenuBar)
	} else {
		mmi := w.Scene.Frame.ChildByName("main-menu", 0)
		if mmi != nil {
			mm := mmi.(*MenuBar)
			w.MainMenu = mm
			return mm
		}
	}
	w.MainMenu = w.Scene.Frame.InsertNewChild(TypeMenuBar, 0, "main-menu").(*MenuBar)
	w.MainMenu.MainMenu = true
	w.MainMenu.SetStretchMaxWidth()
	return w.MainMenu
}

*/
