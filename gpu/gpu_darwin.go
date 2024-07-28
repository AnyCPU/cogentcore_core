// Copyright (c) 2022, Cogent Core. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build darwin

package gpu

func PlatformDefaults(gp *GPU) {
	// gp.DeviceExts = append(gp.DeviceExts, []string{"VK_KHR_portability_subset"}...)
	// gp.DeviceExts = append(gp.DeviceExts)
	// gp.InstanceExts = append(gp.InstanceExts, vk.KhrGetPhysicalDeviceProperties2ExtensionName)
	// gp.InstanceExts = append(gp.InstanceExts, vk.KhrPortabilityEnumerationExtensionName)
	//
	// portFeatures := unsafe.Pointer(&vk.PhysicalDevicePortabilitySubsetFeatures{
	// 	SType:                                  vk.StructureTypePhysicalDevicePortabilitySubsetFeatures,
	// 	ConstantAlphaColorBlendFactors:         vk.True,
	// 	Events:                                 vk.True,
	// 	TextureViewFormatReinterpretation:        vk.True,
	// 	TextureViewFormatSwizzle:                 vk.True,
	// 	TextureView2DOn3DTexture:                   vk.False,
	// 	MultisampleArrayTexture:                  vk.True,
	// 	MutableComparisonSamplers:              vk.True,
	// 	PointPolygons:                          vk.False,
	// 	SamplerMipLodBias:                      vk.False,
	// 	SeparateStencilMaskRef:                 vk.True,
	// 	ShaderSampleRateInterpolationFunctions: vk.True,
	// 	TessellationIsolines:                   vk.False,
	// 	TessellationPointMode:                  vk.False,
	// 	TriangleFans:                           vk.False,
	// 	VertexAttributeAccessBeyondStride:      vk.True,
	// })
	//
	// gp.DeviceFeaturesNeeded = &vk.PhysicalDeviceVulkan12Features{
	// 	SType:                                        vk.StructureTypePhysicalDeviceVulkan12Features,
	// 	DescriptorBindingVariableDescriptorCount:     vk.True,
	// 	DescriptorBindingPartiallyBound:              vk.True,
	// 	RuntimeDescriptorArray:                       vk.True,
	// 	DescriptorIndexing:                           vk.True, // might not be needed?  not for phong or vdraw
	// 	DescriptorBindingSampledTextureUpdateAfterBind: vk.True, // might not be needed?  not for phong or vdraw
	// 	PNext: portFeatures,
	// }
}
