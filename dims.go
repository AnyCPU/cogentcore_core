// Copyright 2019 The Goki Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mat32

//go:generate enumgen

// Dims is a list of vector dimension (component) names
type Dims int32 //enums:enum

const (
	X Dims = iota
	Y
	Z
	W
)

// OtherDim returns the other dimension for 2D X,Y
func OtherDim(d Dims) Dims {
	switch d {
	case X:
		return Y
	default:
		return X
	}
}

func (d Dims) Other() Dims {
	return OtherDim(d)
}
