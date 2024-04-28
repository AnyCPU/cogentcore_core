// Copyright (c) 2023, Cogent Core. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package reflectx

import (
	"path"
	"reflect"
	"strings"

	"cogentcore.org/core/base/strcase"
)

// LongTypeName returns the long, full package-path qualified type name.
// This is guaranteed to be unique and used for internal storage of
// several maps to avoid any conflicts.  It is also very quick to compute.
func LongTypeName(typ reflect.Type) string {
	nptyp := NonPointerType(typ)
	nm := nptyp.Name()
	if nm != "" {
		return nptyp.PkgPath() + "." + nm
	}
	return typ.String()
}

// ShortTypeName returns the short version of a package-qualified type name
// which just has the last element of the path.  This is what is used in
// standard Go programming, and is is used for the key to lookup reflect.Type
// names -- i.e., this is what you should save in a JSON file.
// The potential naming conflict is worth the brevity, and typically a given
// file will only contain mutually compatible, non-conflicting types.
// This is cached in ShortNames because the path.Base computation is apparently
// a bit slow.
func ShortTypeName(typ reflect.Type) string {
	nptyp := NonPointerType(typ)
	nm := nptyp.Name()
	if nm != "" {
		return path.Base(nptyp.PkgPath()) + "." + nm
	}
	return typ.String()
}

// FriendlyTypeName returns a user-friendly version of the name of the given type.
// It transforms it into sentence case, excludes the package, and converts various
// builtin types into more friendly forms (eg: "int" to "Number").
func FriendlyTypeName(typ reflect.Type) string {
	nptyp := NonPointerType(typ)
	nm := nptyp.Name()

	// if it is named, we use that
	if nm != "" {
		switch nm {
		case "string":
			return "Text"
		case "float32", "float64", "int", "int8", "int16", "int32", "int64", "uint", "uint8", "uint16", "uint32", "uint64", "uintptr":
			return "Number"
		}
		return strcase.ToSentence(nm)
	}

	// otherwise, we fall back on Kind
	switch nptyp.Kind() {
	case reflect.Slice, reflect.Array, reflect.Map:
		bnm := FriendlyTypeName(nptyp.Elem())
		if strings.HasSuffix(bnm, "s") {
			return "List of " + bnm
		} else if strings.Contains(bnm, "Function of") {
			return strings.ReplaceAll(bnm, "Function of", "Functions of") + "s"
		} else {
			return bnm + "s"
		}
	case reflect.Func:
		str := "Function"
		ni := nptyp.NumIn()
		if ni > 0 {
			str += " of"
		}
		for i := 0; i < ni; i++ {
			str += " " + FriendlyTypeName(nptyp.In(i))
			if ni == 2 && i == 0 {
				str += " and"
			} else if i == ni-2 {
				str += ", and"
			} else if i < ni-1 {
				str += ","
			}
		}
		return str
	}
	if nptyp.String() == "interface {}" {
		return "Value"
	}
	return nptyp.String()
}

// CloneToType creates a new pointer to the given type
// and uses [SetRobust] to copy an existing value
// (of potentially another type) to it.
func CloneToType(typ reflect.Type, val any) reflect.Value {
	vn := reflect.New(typ)
	evi := vn.Interface()
	SetRobust(evi, val)
	return vn
}
