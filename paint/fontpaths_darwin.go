// Copyright (c) 2023, Cogent Core. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build darwin

package paint

// FontPaths contains the filepaths in which fonts are stored for the current platform.
var FontPaths = []string{"/System/Library/Fonts", "/Library/Fonts"}
