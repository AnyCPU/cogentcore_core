// Copyright (c) 2024, Cogent Core. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// This tool prints out information about your available WebGPU devices.
package main

import (
	"fmt"

	"cogentcore.org/core/base/reflectx"
	"cogentcore.org/core/gpu"
	"github.com/cogentcore/webgpu/wgpu"
)

func main() {
	instance := wgpu.CreateInstance(nil)

	gpus := instance.EnumerateAdapters(nil)
	gp := gpu.NewGPU(nil)
	gpIndex := gp.SelectGPU(gpus)
	props := gpus[gpIndex].GetInfo()
	fmt.Println("Default WebGPU Adapter number:", gpIndex, "  Type:", props.AdapterType.String(), "  Backend:", props.BackendType.String())
	fmt.Println("Set GPU_DEVICE environment variable to select a different device")

	for i, a := range gpus {
		props := a.GetInfo()
		fmt.Println("\n#####################################################################")
		fmt.Println("WebGPU Adapter number:", i, "  Type:", props.AdapterType.String(), "  Backend:", props.BackendType.String())
		fmt.Println(reflectx.StringJSON(props))
		limits := a.GetLimits()
		fmt.Println(reflectx.StringJSON(limits))
	}
}
