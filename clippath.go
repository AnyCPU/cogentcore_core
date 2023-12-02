// Copyright (c) 2018, The GoKi Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package svg

// todo: needs to be impl

// ClipPath is used for holding a path that renders as a clip path
type ClipPath struct {
	NodeBase
}

func (g *ClipPath) SVGName() string { return "clippath" }

func (g *ClipPath) CopyFieldsFrom(frm any) {
	fr := frm.(*ClipPath)
	g.NodeBase.CopyFieldsFrom(&fr.NodeBase)
}
