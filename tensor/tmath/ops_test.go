// Copyright (c) 2024, Cogent Core. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tmath

import (
	"testing"

	"cogentcore.org/core/tensor"
	"github.com/stretchr/testify/assert"
)

func TestOps(t *testing.T) {
	scalar := tensor.NewFloat64Scalar(-5.5)
	scb := scalar.Clone()
	scb.SetFloat1D(-4.0, 0)
	scout := scalar.Clone()

	vals := []float64{-1.507556722888818, -1.2060453783110545, -0.9045340337332908, -0.6030226891555273, -0.3015113445777635, 0.1, 0.3015113445777635, 0.603022689155527, 0.904534033733291, 1.2060453783110545, 1.507556722888818, .3}

	oned := tensor.NewNumberFromValues(vals...)
	oneout := oned.Clone()

	cell2d := tensor.NewFloat32(5, 12)
	_, cells := cell2d.Shape().RowCellSize()
	assert.Equal(t, cells, 12)
	tensor.VectorizeThreaded(1, tensor.NFirstLen, func(idx int, tsr ...tensor.Tensor) {
		ci := idx % cells
		cell2d.SetFloat1D(oned.Float1D(ci), idx)
	}, cell2d)
	// cell2d.DeleteRows(3, 1)
	cellout := cell2d.Clone()
	_ = cellout

	AddOut(scalar, scb, scout)
	assert.Equal(t, -5.5+-4, scout.Float1D(0))

	AddOut(scalar, oned, oneout)
	for i, v := range vals {
		assert.Equal(t, v+-5.5, oneout.Float1D(i))
	}

	AddOut(oned, oned, oneout)
	for i, v := range vals {
		assert.Equal(t, v+v, oneout.Float1D(i))
	}

	AddOut(cell2d, oned, cellout)
	for ri := range 5 {
		for i, v := range vals {
			assert.InDelta(t, v+v, cellout.FloatRowCell(ri, i), 1.0e-6)
		}
	}

	SubOut(scalar, scb, scout)
	assert.Equal(t, -5.5 - -4, scout.Float1D(0))

	SubOut(scb, scalar, scout)
	assert.Equal(t, -4 - -5.5, scout.Float1D(0))

	SubOut(scalar, oned, oneout)
	for i, v := range vals {
		assert.Equal(t, -5.5-v, oneout.Float1D(i))
	}

	SubOut(oned, scalar, oneout)
	for i, v := range vals {
		assert.Equal(t, v - -5.5, oneout.Float1D(i))
	}

	SubOut(oned, oned, oneout)
	for i, v := range vals {
		assert.Equal(t, v-v, oneout.Float1D(i))
	}

	SubOut(cell2d, oned, cellout)
	for ri := range 5 {
		for i, v := range vals {
			assert.InDelta(t, v-v, cellout.FloatRowCell(ri, i), 1.0e-6)
		}
	}

	MulOut(scalar, scb, scout)
	assert.Equal(t, -5.5*-4, scout.Float1D(0))

	MulOut(scalar, oned, oneout)
	for i, v := range vals {
		assert.Equal(t, v*-5.5, oneout.Float1D(i))
	}

	MulOut(oned, oned, oneout)
	for i, v := range vals {
		assert.Equal(t, v*v, oneout.Float1D(i))
	}

	MulOut(cell2d, oned, cellout)
	for ri := range 5 {
		for i, v := range vals {
			assert.InDelta(t, v*v, cellout.FloatRowCell(ri, i), 1.0e-6)
		}
	}

	DivOut(scalar, scb, scout)
	assert.Equal(t, -5.5/-4, scout.Float1D(0))

	DivOut(scb, scalar, scout)
	assert.Equal(t, -4/-5.5, scout.Float1D(0))

	DivOut(scalar, oned, oneout)
	for i, v := range vals {
		assert.Equal(t, -5.5/v, oneout.Float1D(i))
	}

	DivOut(oned, scalar, oneout)
	for i, v := range vals {
		assert.Equal(t, v/-5.5, oneout.Float1D(i))
	}

	DivOut(oned, oned, oneout)
	for i, v := range vals {
		assert.Equal(t, v/v, oneout.Float1D(i))
	}

	DivOut(cell2d, oned, cellout)
	for ri := range 5 {
		for i, v := range vals {
			assert.InDelta(t, v/v, cellout.FloatRowCell(ri, i), 1.0e-6)
		}
	}

	onedc := tensor.Clone(oned)
	AddAssign(onedc, scalar)
	for i, v := range vals {
		assert.Equal(t, v+-5.5, onedc.Float1D(i))
	}

	SubAssign(onedc, scalar)
	for i, v := range vals {
		assert.InDelta(t, v, onedc.Float1D(i), 1.0e-8)
	}

	MulAssign(onedc, scalar)
	for i, v := range vals {
		assert.InDelta(t, v*-5.5, onedc.Float1D(i), 1.0e-7)
	}

	DivAssign(onedc, scalar)
	for i, v := range vals {
		assert.InDelta(t, v, onedc.Float1D(i), 1.0e-7)
	}

	Inc(onedc)
	for i, v := range vals {
		assert.InDelta(t, v+1, onedc.Float1D(i), 1.0e-7)
	}

	Dec(onedc)
	for i, v := range vals {
		assert.InDelta(t, v, onedc.Float1D(i), 1.0e-7)
	}
}

func TestBoolOps(t *testing.T) {
	ar := tensor.NewSliceInts(12)
	// fmt.Println(v)
	bo := tensor.NewBool()
	sc := tensor.NewIntScalar(6)

	GreaterOut(ar, sc, bo)
	// fmt.Println(bo)
	for i, v := range ar.Values {
		assert.Equal(t, v > 6, bo.Bool1D(i))
	}

	LessOut(ar, sc, bo)
	// fmt.Println(bo)
	for i, v := range ar.Values {
		assert.Equal(t, v < 6, bo.Bool1D(i))
	}
}
