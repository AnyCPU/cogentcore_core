// Copyright (c) 2019, Cogent Core. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package vgpu

import (
	vk "github.com/goki/vulkan"
)

// See: https://www.khronos.org/opengl/wiki/Data_Type_(GLSL)

// Types is a list of supported GPU data types, which can be stored
// properly aligned in device memory, and used by the shader code.
// Note that a Vec3 or arrays of single scalar values such as Float32
// are not well supported outside of Vertex due to the std410 convention:
// http://www.opengl.org/registry/doc/glspec45.core.pdf#page=159
// The Struct type is particularly challenging as each member
// must be aligned in general on a 16 byte boundary (i.e., vec4)
// (unless all elements are exactly 4 bytes, which might work?).
// Go automatically aligns members to 8 bytes on 64 bit machines,
// but that doesn't quite cut it.
type Types int32 //enums:enum

const (
	UndefType Types = iota
	Bool32

	Int16
	Uint16

	Int32
	Int32Vec2
	Int32Vec4

	Uint32
	Uint32Vec2
	Uint32Vec4

	Float32
	Float32Vec2
	Float32Vec3 // note: only use for vertex data -- not properly aligned for uniforms
	Float32Vec4

	Float64
	Float64Vec2
	Float64Vec3
	Float64Vec4

	Float32Mat4 // std transform matrix: mat32.Mat4 works directly
	Float32Mat3 // std transform matrix: mat32.Mat3 works directly

	ImageRGBA32 // 32 bits with 8 bits per component of R,G,B,A -- std image format

	Depth32      // standard float32 depth buffer
	Depth24Sten8 // standard 24 bit float with 8 bit stencil

	Struct
)

// VkFormat returns the Vulkan VkFormat for given type
func (tp Types) VkFormat() vk.Format {
	return VulkanTypes[tp]
}

// VkIndexType returns the Vulkan vk.IndexType for var
// must be either Uint16 or Uint32
func (tp Types) VkIndexType() vk.IndexType {
	if tp == Uint16 {
		return vk.IndexTypeUint16
	}
	return vk.IndexTypeUint32
}

// Bytes returns number of bytes for this type
func (tp Types) Bytes() int {
	switch tp {
	case Float32Mat4:
		return 64
	case Float32Mat3:
		return 36
	}
	if vf, has := VulkanTypes[tp]; has {
		return FormatSizes[vf]
	}
	return 0
}

// FormatSizes gives size of known vulkan formats in bytes
var FormatSizes = map[vk.Format]int{
	vk.FormatUndefined:          0,
	vk.FormatR16Sint:            2,
	vk.FormatR16Uint:            2,
	vk.FormatR32Sint:            4,
	vk.FormatR32g32Sint:         8,
	vk.FormatR32g32b32a32Sint:   16,
	vk.FormatR32Uint:            4,
	vk.FormatR32g32Uint:         8,
	vk.FormatR32g32b32a32Uint:   16,
	vk.FormatR32Sfloat:          4,
	vk.FormatR32g32Sfloat:       8,
	vk.FormatR32g32b32Sfloat:    12,
	vk.FormatR32g32b32a32Sfloat: 16,
	vk.FormatR64Sfloat:          8,
	vk.FormatR64g64Sfloat:       16,
	vk.FormatR64g64b64Sfloat:    24,
	vk.FormatR64g64b64a64Sfloat: 32,
	vk.FormatR8g8b8a8Srgb:       4,
	vk.FormatR8g8b8a8Unorm:      4,
	vk.FormatD32Sfloat:          4,
	vk.FormatD24UnormS8Uint:     4,
}

// VulkanTypes maps vgpu.Types to vulkan types
var VulkanTypes = map[Types]vk.Format{
	UndefType:    vk.FormatUndefined,
	Bool32:       vk.FormatR32Uint,
	Int16:        vk.FormatR16Sint,
	Uint16:       vk.FormatR16Uint,
	Int32:        vk.FormatR32Sint,
	Int32Vec2:    vk.FormatR32g32Sint,
	Int32Vec4:    vk.FormatR32g32b32a32Sint,
	Uint32:       vk.FormatR32Uint,
	Uint32Vec2:   vk.FormatR32g32Uint,
	Uint32Vec4:   vk.FormatR32g32b32a32Uint,
	Float32:      vk.FormatR32Sfloat,
	Float32Vec2:  vk.FormatR32g32Sfloat,
	Float32Vec3:  vk.FormatR32g32b32Sfloat,
	Float32Vec4:  vk.FormatR32g32b32a32Sfloat,
	Float64:      vk.FormatR64Sfloat,
	Float64Vec2:  vk.FormatR64g64Sfloat,
	Float64Vec3:  vk.FormatR64g64b64Sfloat,
	Float64Vec4:  vk.FormatR64g64b64a64Sfloat,
	ImageRGBA32:  vk.FormatR8g8b8a8Srgb,
	Depth32:      vk.FormatD32Sfloat,
	Depth24Sten8: vk.FormatD24UnormS8Uint,
}

// most commonly available formats: https://vulkan.gpuinfo.org/listsurfaceformats.php

// ImageFormatNames translates image format into human-readable string
// for most commonly-available formats
var ImageFormatNames = map[vk.Format]string{
	vk.FormatR8g8b8a8Srgb:           "RGBA 8bit sRGB colorspace",
	vk.FormatR8g8b8a8Unorm:          "RGBA 8bit unsigned linear colorspace",
	vk.FormatR5g6b5UnormPack16:      "RGB 5bit (pack 16bit total) unsigned linear colorspace",
	vk.FormatA2b10g10r10UnormPack32: "ABGR 10bit, 2bit alpha (pack 32bit total), unsigned linear colorspace",
	vk.FormatB8g8r8a8Srgb:           "BGRA 8bit sRGB colorspace",
	vk.FormatB8g8r8a8Unorm:          "BGRA 8bit unsigned linear colorspace",
	vk.FormatR16g16b16a16Sfloat:     "RGBA 16bit signed floating point linear colorspace",
	vk.FormatA2r10g10b10UnormPack32: "ARGB 10bit, 2bit alpha (pack 32bit total), unsigned linear colorspace",
}
