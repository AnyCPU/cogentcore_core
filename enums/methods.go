// Copyright (c) 2024, Cogent Core. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package enums

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"sync/atomic"

	"cogentcore.org/core/glop/num"
)

// This file contains implementations of enumgen methods.

// EnumConstraint is the generic type constraint that all enums satisfy.
type EnumConstraint interface {
	Enum
	num.Integer
}

// BitFlagConstraint is the generic type constraint that all bit flags satisfy.
type BitFlagConstraint interface {
	BitFlag
	num.Integer
}

// String returns the string representation of the given
// enum value with the given map.
func String[T EnumConstraint](i T, m map[T]string) string {
	if str, ok := m[i]; ok {
		return str
	}
	return strconv.FormatInt(int64(i), 10)
}

// StringExtended returns the string representation of the given enum value
// with the given map, with the enum type extending the given other enum type.
func StringExtended[T, E EnumConstraint](i T, m map[T]string) string {
	if str, ok := m[i]; ok {
		return str
	}
	return E(i).String()
}

// BitIndexStringExtended returns the string representation of the given bit flag enum
// bit index value with the given map, with the bit flag type extending the given other
// bit flag type.
func BitIndexStringExtended[T, E BitFlagConstraint](i T, m map[T]string) string {
	if str, ok := m[i]; ok {
		return str
	}
	return E(i).BitIndexString()
}

// BitFlagString returns the string representation of the given bit flag value
// with the given values available.
func BitFlagString[T BitFlagConstraint](i T, values []T) string {
	str := ""
	for _, ie := range values {
		if i.HasFlag(ie) {
			ies := ie.BitIndexString()
			if str == "" {
				str = ies
			} else {
				str += "|" + ies
			}
		}
	}
	return str
}

// BitFlagStringExtended returns the string representation of the given bit flag value
// with the given values available, with the bit flag type extending the other given
// bit flag type that has the given values (extendedValues) available.
func BitFlagStringExtended[T, E BitFlagConstraint](i T, values []T, extendedValues []E) string {
	str := ""
	for _, ie := range extendedValues {
		if i.HasFlag(ie) {
			ies := ie.BitIndexString()
			if str == "" {
				str = ies
			} else {
				str += "|" + ies
			}
		}
	}
	for _, ie := range values {
		if i.HasFlag(ie) {
			ies := ie.BitIndexString()
			if str == "" {
				str = ies
			} else {
				str += "|" + ies
			}
		}
	}
	return str
}

// SetString sets the given enum value from its string representation, the map from
// enum names to values, and the name of the enum type, which is used for the error message.
func SetString[T EnumConstraint](i *T, s string, valueMap map[string]T, typeName string) error {
	if val, ok := valueMap[s]; ok {
		*i = val
		return nil
	}
	return errors.New(s + " is not a valid value for type " + typeName)
}

// SetStringLower sets the given enum value from its string representation, the map from
// enum names to values, and the name of the enum type, which is used for the error message.
// It also tries the lowercase version of the given string if the original version fails.
func SetStringLower[T EnumConstraint](i *T, s string, valueMap map[string]T, typeName string) error {
	if val, ok := valueMap[s]; ok {
		*i = val
		return nil
	}
	if val, ok := valueMap[strings.ToLower(s)]; ok {
		*i = val
		return nil
	}
	return errors.New(s + " is not a valid value for type " + typeName)
}

// SetStringExtended sets the given enum value from its string representation and the map from
// enum names to values, with the enum type extending the other given enum type. It also takes
// the enum value in terms of the extended enum type (ie).
func SetStringExtended[T EnumConstraint, E EnumSetter](i *T, ie E, s string, valueMap map[string]T) error {
	if val, ok := valueMap[s]; ok {
		*i = val
		return nil
	}
	return ie.SetString(s)
}

// SetStringLowerExtended sets the given enum value from its string representation and the map from
// enum names to values, with the enum type extending the other given enum type. It also takes
// the enum value in terms of the extended enum type (ie). It also tries the lowercase version
// of the given string if the original version fails.
func SetStringLowerExtended[T EnumConstraint, E EnumSetter](i *T, ie E, s string, valueMap map[string]T) error {
	if val, ok := valueMap[s]; ok {
		*i = val
		return nil
	}
	if val, ok := valueMap[strings.ToLower(s)]; ok {
		*i = val
		return nil
	}
	return ie.SetString(s)
}

// SetStringOr sets the given bit flag value from its string representation while
// preserving any bit flags already set.
func SetStringOr[T BitFlagConstraint, S BitFlagSetter](i S, s string, valueMap map[string]T) error {
	flags := strings.Split(s, "|")
	for _, flag := range flags {
		if val, ok := valueMap[flag]; ok {
			i.SetFlag(true, val)
		} else if flag == "" {
			continue
		} else {
			return fmt.Errorf("%q is not a valid value for type States", flag)
		}
	}
	return nil
}

// SetFlag sets the value of the given flags in these flags to the given value.
func SetFlag(i *int64, on bool, f ...BitFlag) {
	var mask int64
	for _, v := range f {
		mask |= 1 << v.Int64()
	}
	in := *i
	if on {
		in |= mask
		atomic.StoreInt64(i, in)
	} else {
		in &^= mask
		atomic.StoreInt64(i, in)
	}
}
