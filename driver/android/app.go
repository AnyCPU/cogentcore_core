// Copyright 2023 The Goki Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build android

package android

import (
	"log"

	vk "github.com/goki/vulkan"
	"goki.dev/goosi"
	"goki.dev/goosi/clip"
	"goki.dev/goosi/driver/base"
	"goki.dev/goosi/events"
	"goki.dev/vgpu/v2/vdraw"
	"goki.dev/vgpu/v2/vgpu"
)

func Init() {
	goosi.OnSystemWindowCreated = make(chan struct{})
	TheApp.InitVk()
	base.Init(TheApp, &TheApp.App)
}

// TheApp is the single [goosi.App] for the Android platform
var TheApp = &App{AppSingle: base.NewAppSingle[*vdraw.Drawer, *Window]()}

// App is the [goosi.App] implementation for the Android platform
type App struct { //gti:add
	base.AppSingle[*vdraw.Drawer, *Window]

	// GPU is the system GPU used for the app
	GPU *vgpu.GPU

	// Winptr is the pointer to the underlying system window
	Winptr uintptr
}

// InitVk initializes Vulkan things for the app
func (a *App) InitVk() {
	err := vk.SetDefaultGetInstanceProcAddr()
	if err != nil {
		// TODO(kai): maybe implement better error handling here
		log.Fatalln("goosi/driver/android.App.InitVk: failed to set Vulkan DefaultGetInstanceProcAddr")
	}
	err = vk.Init()
	if err != nil {
		log.Fatalln("goosi/driver/android.App.InitVk: failed to initialize vulkan")
	}

	winext := vk.GetRequiredInstanceExtensions()
	a.GPU = vgpu.NewGPU()
	a.GPU.AddInstanceExt(winext...)
	a.GPU.Config(a.Name())
}

// DestroyVk destroys vulkan things (the drawer and surface of the window) for when the app becomes invisible
func (a *App) DestroyVk() {
	a.Mu.Lock()
	defer a.Mu.Unlock()
	vk.DeviceWaitIdle(a.Draw.Surf.Device.Device)
	a.Draw.Destroy()
	a.Draw.Surf.Destroy()
	a.Draw = nil
}

// FullDestroyVk destroys all vulkan things for when the app is fully quit
func (a *App) FullDestroyVk() {
	a.Mu.Lock()
	defer a.Mu.Unlock()
	a.GPU.Destroy()
}

// NewWindow creates a new window with the given options.
// It waits for the underlying system window to be created first.
// Also, it hides all other windows and shows the new one.
func (a *App) NewWindow(opts *goosi.NewWindowOptions) (goosi.Window, error) {
	defer func() { goosi.HandleRecover(recover()) }()

	a.Mu.Lock()
	defer a.Mu.Unlock()
	a.Win = &Window{base.NewWindowSingle(a, opts)}
	a.Win.This = a.Win
	a.EvMgr.Window(events.WinShow)
	a.EvMgr.Window(events.WinFocus)

	go a.Win.WinLoop()

	return a.Win, nil
}

// SetSystemWindow sets the underlying system window pointer, surface, system, and drawer.
// It should only be called when [App.Mu] is already locked.
func (a *App) SetSystemWindow(winptr uintptr) error {
	defer func() { goosi.HandleRecover(recover()) }()
	var vsf vk.Surface
	// we have to remake the surface, system, and drawer every time someone reopens the window
	// because the operating system changes the underlying window
	ret := vk.CreateWindowSurface(a.GPU.Instance, winptr, nil, &vsf)
	if err := vk.Error(ret); err != nil {
		return err
	}
	sf := vgpu.NewSurface(a.GPU, vsf)

	sys := a.GPU.NewGraphicsSystem(a.Name(), &sf.Device)
	sys.ConfigRender(&sf.Format, vgpu.UndefType)
	sf.SetRender(&sys.Render)
	// sys.Mem.Vars.NDescs = vgpu.MaxTexturesPerSet
	sys.Config()
	a.Draw = &vdraw.Drawer{
		Sys:     *sys,
		YIsDown: true,
	}
	// a.Draw.ConfigSys()
	a.Draw.ConfigSurface(sf, vgpu.MaxTexturesPerSet)

	a.Winptr = winptr

	// if the window already exists, we are coming back to it, so we need to show it
	// again and send a screen update
	if a.Win != nil {
		a.EvMgr.Window(events.WinShow)
		a.EvMgr.Window(events.ScreenUpdate)
	}
	return nil
}

func (a *App) DataDir() string {
	return "/data/data"
}

func (a *App) Platform() goosi.Platforms {
	return goosi.Android
}

func (a *App) OpenURL(url string) {
	// TODO(kai): implement OpenURL on Android
}

func (a *App) ClipBoard(win goosi.Window) clip.Board {
	// TODO(kai): implement clipboard on Android
	return &clip.BoardBase{}
}
