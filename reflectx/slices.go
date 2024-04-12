// Copyright (c) 2023, Cogent Core. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package reflectx

import (
	"fmt"
	"log"
	"reflect"
	"sort"
	"strings"
	"time"
)

// This file contains helpful functions for dealing with slices, in the reflect
// system

// MakeSlice makes a map that is actually addressable, getting around the hidden
// interface{} that reflect.MakeSlice makes, by calling UnhideAnyValue (from ptrs.go)
func MakeSlice(typ reflect.Type, len, cap int) reflect.Value {
	return UnhideAnyValue(reflect.MakeSlice(typ, len, cap))
}

// SliceElType returns the type of the elements for the given slice (which can be
// a pointer to a slice or a direct slice) -- just Elem() of slice type, but using
// this function makes it more explicit what is going on.  And it uses
// OnePtrUnderlyingValue to get past any interface wrapping.
func SliceElType(sl any) reflect.Type {
	return NonPointerValue(OnePointerUnderlyingValue(reflect.ValueOf(sl))).Type().Elem()
}

// SliceElValue returns a reflect.Value of the Slice element type.
// isPtr indicates that the value is a Pointer type, in which case
// a concrete element has been made and the value is a pointer to it.
func SliceElValue(sl any) reflect.Value {
	typ := SliceElType(sl)
	isPtr := typ.Kind() == reflect.Pointer
	val := reflect.New(NonPointerType(typ)) // make the concrete el
	if !isPtr {
		val = val.Elem() // use concrete value
	}
	return val
}

// SliceNewAt inserts a new blank element at given index in the slice.
// -1 means the end.
func SliceNewAt(sl any, idx int) {
	svl := OnePointerValue(reflect.ValueOf(sl))
	svnp := NonPointerValue(svl)
	val := SliceElValue(sl)
	sz := svnp.Len()
	svnp = reflect.Append(svnp, val)
	if idx >= 0 && idx < sz {
		reflect.Copy(svnp.Slice(idx+1, sz+1), svnp.Slice(idx, sz))
		svnp.Index(idx).Set(val)
	}
	svl.Elem().Set(svnp)
}

// SliceDeleteAt deletes element at given index from slice
func SliceDeleteAt(sl any, idx int) {
	svl := OnePointerValue(reflect.ValueOf(sl))
	svnp := NonPointerValue(svl)
	svtyp := svnp.Type()
	nval := reflect.New(svtyp.Elem())
	sz := svnp.Len()
	reflect.Copy(svnp.Slice(idx, sz-1), svnp.Slice(idx+1, sz))
	svnp.Index(sz - 1).Set(nval.Elem())
	svl.Elem().Set(svnp.Slice(0, sz-1))
}

// todo: this is kind of a hack but needed for pi/filecat/FileTime

// SortInter looks for Int() method for sorting
type SortInter interface {
	Int() int64
}

// SliceSort sorts a slice of basic values (see StructSliceSort for sorting a
// slice-of-struct using a specific field), using float, int, string conversions
// (first fmt.Stringer String()) and supporting time.Time directly as well.
func SliceSort(sl any, ascending bool) error {
	sv := reflect.ValueOf(sl)
	svnp := NonPointerValue(sv)
	if svnp.Len() == 0 {
		return nil
	}
	eltyp := SliceElType(sl)
	elnptyp := NonPointerType(eltyp)
	vk := elnptyp.Kind()
	elval := OnePointerValue(svnp.Index(0))
	elif := elval.Interface()

	if _, ok := elif.(SortInter); ok {
		sort.Slice(svnp.Interface(), func(i, j int) bool {
			iv := NonPointerValue(svnp.Index(i)).Interface().(SortInter).Int()
			jv := NonPointerValue(svnp.Index(j)).Interface().(SortInter).Int()
			if ascending {
				return iv < jv
			}
			return jv < iv
		})
		return nil
	}

	// try all the numeric types first!
	switch {
	case vk >= reflect.Int && vk <= reflect.Int64:
		sort.Slice(svnp.Interface(), func(i, j int) bool {
			iv := NonPointerValue(svnp.Index(i)).Int()
			jv := NonPointerValue(svnp.Index(j)).Int()
			if ascending {
				return iv < jv
			}
			return iv > jv
		})
		return nil
	case vk >= reflect.Uint && vk <= reflect.Uint64:
		sort.Slice(svnp.Interface(), func(i, j int) bool {
			iv := NonPointerValue(svnp.Index(i)).Uint()
			jv := NonPointerValue(svnp.Index(j)).Uint()
			if ascending {
				return iv < jv
			}
			return iv > jv
		})
		return nil
	case vk >= reflect.Float32 && vk <= reflect.Float64:
		sort.Slice(svnp.Interface(), func(i, j int) bool {
			iv := NonPointerValue(svnp.Index(i)).Float()
			jv := NonPointerValue(svnp.Index(j)).Float()
			if ascending {
				return iv < jv
			}
			return iv > jv
		})
		return nil
	case vk == reflect.Struct && ShortTypeName(elnptyp) == "time.Time":
		sort.Slice(svnp.Interface(), func(i, j int) bool {
			iv := NonPointerValue(svnp.Index(i)).Interface().(time.Time)
			jv := NonPointerValue(svnp.Index(j)).Interface().(time.Time)
			if ascending {
				return iv.Before(jv)
			}
			return jv.Before(iv)
		})
	}

	// this stringer case will likely pick up most of the rest
	switch elif.(type) {
	case fmt.Stringer:
		sort.Slice(svnp.Interface(), func(i, j int) bool {
			iv := NonPointerValue(svnp.Index(i)).Interface().(fmt.Stringer).String()
			jv := NonPointerValue(svnp.Index(j)).Interface().(fmt.Stringer).String()
			if ascending {
				return iv < jv
			}
			return iv > jv
		})
		return nil
	}

	// last resort!
	switch {
	case vk == reflect.String:
		sort.Slice(svnp.Interface(), func(i, j int) bool {
			iv := NonPointerValue(svnp.Index(i)).String()
			jv := NonPointerValue(svnp.Index(j)).String()
			if ascending {
				return strings.ToLower(iv) < strings.ToLower(jv)
			}
			return strings.ToLower(iv) > strings.ToLower(jv)
		})
		return nil
	}

	err := fmt.Errorf("SortSlice: unable to sort elements of type: %v", eltyp.String())
	log.Println(err)
	return err
}

// StructSliceSort sorts a slice of a struct according to the given field
// indexes and sort direction, using float, int, string
// conversions (first fmt.Stringer String()) and supporting time.Time directly
// as well.  There is no direct method for checking the field indexes so those
// are assumed to be accurate -- will panic if not!
func StructSliceSort(struSlice any, fldIndex []int, ascending bool) error {
	sv := reflect.ValueOf(struSlice)
	svnp := NonPointerValue(sv)
	if svnp.Len() == 0 {
		return nil
	}
	struTyp := SliceElType(struSlice)
	struNpTyp := NonPointerType(struTyp)
	fld := struNpTyp.FieldByIndex(fldIndex) // not easy to check.
	vk := fld.Type.Kind()
	struVal := OnePointerValue(svnp.Index(0))
	fldVal := struVal.Elem().FieldByIndex(fldIndex)
	fldIf := fldVal.Interface()

	if _, ok := fldIf.(SortInter); ok {
		sort.Slice(svnp.Interface(), func(i, j int) bool {
			ival := OnePointerValue(svnp.Index(i))
			iv := ival.Elem().FieldByIndex(fldIndex).Interface().(SortInter).Int()
			jval := OnePointerValue(svnp.Index(j))
			jv := jval.Elem().FieldByIndex(fldIndex).Interface().(SortInter).Int()
			if ascending {
				return iv < jv
			}
			return jv < iv
		})
		return nil
	}

	// try all the numeric types first!
	switch {
	case vk >= reflect.Int && vk <= reflect.Int64:
		sort.Slice(svnp.Interface(), func(i, j int) bool {
			ival := OnePointerValue(svnp.Index(i))
			iv := ival.Elem().FieldByIndex(fldIndex).Int()
			jval := OnePointerValue(svnp.Index(j))
			jv := jval.Elem().FieldByIndex(fldIndex).Int()
			if ascending {
				return iv < jv
			}
			return iv > jv
		})
		return nil
	case vk >= reflect.Uint && vk <= reflect.Uint64:
		sort.Slice(svnp.Interface(), func(i, j int) bool {
			ival := OnePointerValue(svnp.Index(i))
			iv := ival.Elem().FieldByIndex(fldIndex).Uint()
			jval := OnePointerValue(svnp.Index(j))
			jv := jval.Elem().FieldByIndex(fldIndex).Uint()
			if ascending {
				return iv < jv
			}
			return iv > jv
		})
		return nil
	case vk >= reflect.Float32 && vk <= reflect.Float64:
		sort.Slice(svnp.Interface(), func(i, j int) bool {
			ival := OnePointerValue(svnp.Index(i))
			iv := ival.Elem().FieldByIndex(fldIndex).Float()
			jval := OnePointerValue(svnp.Index(j))
			jv := jval.Elem().FieldByIndex(fldIndex).Float()
			if ascending {
				return iv < jv
			}
			return iv > jv
		})
		return nil
	case vk == reflect.Struct && ShortTypeName(fld.Type) == "time.Time":
		sort.Slice(svnp.Interface(), func(i, j int) bool {
			ival := OnePointerValue(svnp.Index(i))
			iv := ival.Elem().FieldByIndex(fldIndex).Interface().(time.Time)
			jval := OnePointerValue(svnp.Index(j))
			jv := jval.Elem().FieldByIndex(fldIndex).Interface().(time.Time)
			if ascending {
				return iv.Before(jv)
			}
			return jv.Before(iv)
		})
	}

	// this stringer case will likely pick up most of the rest
	switch fldIf.(type) {
	case fmt.Stringer:
		sort.Slice(svnp.Interface(), func(i, j int) bool {
			ival := OnePointerValue(svnp.Index(i))
			iv := ival.Elem().FieldByIndex(fldIndex).Interface().(fmt.Stringer).String()
			jval := OnePointerValue(svnp.Index(j))
			jv := jval.Elem().FieldByIndex(fldIndex).Interface().(fmt.Stringer).String()
			if ascending {
				return iv < jv
			}
			return iv > jv
		})
		return nil
	}

	// last resort!
	switch {
	case vk == reflect.String:
		sort.Slice(svnp.Interface(), func(i, j int) bool {
			ival := OnePointerValue(svnp.Index(i))
			iv := ival.Elem().FieldByIndex(fldIndex).String()
			jval := OnePointerValue(svnp.Index(j))
			jv := jval.Elem().FieldByIndex(fldIndex).String()
			if ascending {
				return strings.ToLower(iv) < strings.ToLower(jv)
			}
			return strings.ToLower(iv) > strings.ToLower(jv)
		})
		return nil
	}

	err := fmt.Errorf("SortStructSlice: unable to sort on field of type: %v\n", fld.Type.String())
	log.Println(err)
	return err
}

// ValueSliceSort sorts a slice of reflect.Values using basic types where possible
func ValueSliceSort(sl []reflect.Value, ascending bool) error {
	if len(sl) == 0 {
		return nil
	}
	felval := sl[0] // reflect.Value
	eltyp := felval.Type()
	elnptyp := NonPointerType(eltyp)
	vk := elnptyp.Kind()
	elval := OnePointerValue(felval)
	elif := elval.Interface()

	// try all the numeric types first!
	switch {
	case vk >= reflect.Int && vk <= reflect.Int64:
		sort.Slice(sl, func(i, j int) bool {
			iv := NonPointerValue(sl[i]).Int()
			jv := NonPointerValue(sl[j]).Int()
			if ascending {
				return iv < jv
			}
			return iv > jv
		})
		return nil
	case vk >= reflect.Uint && vk <= reflect.Uint64:
		sort.Slice(sl, func(i, j int) bool {
			iv := NonPointerValue(sl[i]).Uint()
			jv := NonPointerValue(sl[j]).Uint()
			if ascending {
				return iv < jv
			}
			return iv > jv
		})
		return nil
	case vk >= reflect.Float32 && vk <= reflect.Float64:
		sort.Slice(sl, func(i, j int) bool {
			iv := NonPointerValue(sl[i]).Float()
			jv := NonPointerValue(sl[j]).Float()
			if ascending {
				return iv < jv
			}
			return iv > jv
		})
		return nil
	case vk == reflect.Struct && ShortTypeName(elnptyp) == "time.Time":
		sort.Slice(sl, func(i, j int) bool {
			iv := NonPointerValue(sl[i]).Interface().(time.Time)
			jv := NonPointerValue(sl[j]).Interface().(time.Time)
			if ascending {
				return iv.Before(jv)
			}
			return jv.Before(iv)
		})
	}

	// this stringer case will likely pick up most of the rest
	switch elif.(type) {
	case fmt.Stringer:
		sort.Slice(sl, func(i, j int) bool {
			iv := NonPointerValue(sl[i]).Interface().(fmt.Stringer).String()
			jv := NonPointerValue(sl[j]).Interface().(fmt.Stringer).String()
			if ascending {
				return iv < jv
			}
			return iv > jv
		})
		return nil
	}

	// last resort!
	switch {
	case vk == reflect.String:
		sort.Slice(sl, func(i, j int) bool {
			iv := NonPointerValue(sl[i]).String()
			jv := NonPointerValue(sl[j]).String()
			if ascending {
				return strings.ToLower(iv) < strings.ToLower(jv)
			}
			return strings.ToLower(iv) > strings.ToLower(jv)
		})
		return nil
	}

	err := fmt.Errorf("ValueSliceSort: unable to sort elements of type: %v", eltyp.String())
	log.Println(err)
	return err
}

// CopySliceRobust robustly copies slices using SetRobust method for the elements.
func CopySliceRobust(to, fm any) error {
	tov := reflect.ValueOf(to)
	fmv := reflect.ValueOf(fm)
	tonp := NonPointerValue(tov)
	fmnp := NonPointerValue(fmv)
	totyp := tonp.Type()
	// eltyp := SliceElType(tonp)
	if totyp.Kind() != reflect.Slice {
		err := fmt.Errorf("reflectx.CopySliceRobust: 'to' is not slice, is: %v", totyp.String())
		log.Println(err)
		return err
	}
	fmtyp := fmnp.Type()
	if fmtyp.Kind() != reflect.Slice {
		err := fmt.Errorf("reflectx.CopySliceRobust: 'from' is not slice, is: %v", fmtyp.String())
		log.Println(err)
		return err
	}
	fmlen := fmnp.Len()
	if tonp.IsNil() {
		OnePointerValue(tonp).Elem().Set(MakeSlice(totyp, fmlen, fmlen).Elem())
	} else {
		if tonp.Len() > fmlen {
			tonp.SetLen(fmlen)
		}
	}
	for i := 0; i < fmlen; i++ {
		tolen := tonp.Len()
		if i >= tolen {
			SliceNewAt(to, i)
		}
		SetRobust(PointerValue(tonp.Index(i)).Interface(), fmnp.Index(i).Interface())
	}
	return nil
}
