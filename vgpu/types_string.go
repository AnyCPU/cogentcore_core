// Code generated by "stringer -type=Types"; DO NOT EDIT.

package vgpu

import (
	"errors"
	"strconv"
)

var _ = errors.New("dummy error")

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[UndefType-0]
	_ = x[Bool32-1]
	_ = x[Int32-2]
	_ = x[Int32Vec2-3]
	_ = x[Int32Vec4-4]
	_ = x[Uint32-5]
	_ = x[Uint32Vec2-6]
	_ = x[Uint32Vec4-7]
	_ = x[Float32-8]
	_ = x[Float32Vec2-9]
	_ = x[Float32Vec3-10]
	_ = x[Float32Vec4-11]
	_ = x[Float64-12]
	_ = x[Float64Vec2-13]
	_ = x[Float64Vec3-14]
	_ = x[Float64Vec4-15]
	_ = x[Float32Mat4-16]
	_ = x[ImageRGBA32-17]
	_ = x[Struct-18]
	_ = x[TypesN-19]
}

const _Types_name = "UndefTypeBool32Int32Int32Vec2Int32Vec4Uint32Uint32Vec2Uint32Vec4Float32Float32Vec2Float32Vec3Float32Vec4Float64Float64Vec2Float64Vec3Float64Vec4Float32Mat4ImageRGBA32StructTypesN"

var _Types_index = [...]uint8{0, 9, 15, 20, 29, 38, 44, 54, 64, 71, 82, 93, 104, 111, 122, 133, 144, 155, 166, 172, 178}

func (i Types) String() string {
	if i < 0 || i >= Types(len(_Types_index)-1) {
		return "Types(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _Types_name[_Types_index[i]:_Types_index[i+1]]
}

func (i *Types) FromString(s string) error {
	for j := 0; j < len(_Types_index)-1; j++ {
		if s == _Types_name[_Types_index[j]:_Types_index[j+1]] {
			*i = Types(j)
			return nil
		}
	}
	return errors.New("String: " + s + " is not a valid option for type: Types")
}
