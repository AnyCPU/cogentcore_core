// Copyright (c) 2023, Cogent Core. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package laser

import (
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"reflect"
	"strings"
)

// FlatFieldsTypeFunc calls a function on all the primary fields of a given
// struct type, including those on anonymous embedded structs that this struct
// has, passing the current (embedded) type and StructField -- effectively
// flattens the reflect field list -- if fun returns false then iteration
// stops -- overall rval is false if iteration was stopped or there was an
// error (logged), true otherwise
func FlatFieldsTypeFunc(typ reflect.Type, fun func(typ reflect.Type, field reflect.StructField) bool) bool {
	return FlatFieldsTypeFuncIf(typ, nil, fun)
}

// FlatFieldsTypeFunc calls a function on all the primary fields of a given
// struct type, including those on anonymous embedded structs that this struct
// has, passing the current (embedded) type and StructField -- effectively
// flattens the reflect field list -- if fun returns false then iteration
// stops -- overall rval is false if iteration was stopped or there was an
// error (logged), true otherwise. If the given ifFun is non-nil, it is called
// on every embedded struct field to determine whether the fields of that embedded
// field should be handled (a return value of true indicates to continue down and
// a value of false indicates to not).
func FlatFieldsTypeFuncIf(typ reflect.Type, ifFun, fun func(typ reflect.Type, field reflect.StructField) bool) bool {
	typ = NonPtrType(typ)
	if typ.Kind() != reflect.Struct {
		log.Printf("laser.FlatFieldsTypeFunc: Must call on a struct type, not: %v\n", typ)
		return false
	}
	rval := true
	for i := 0; i < typ.NumField(); i++ {
		f := typ.Field(i)
		if f.Type.Kind() == reflect.Struct && f.Anonymous {
			if ifFun != nil {
				if !ifFun(typ, f) {
					continue
				}
			}
			rval = FlatFieldsTypeFunc(f.Type, fun) // no err here
			if !rval {
				break
			}
		} else {
			rval = fun(typ, f)
			if !rval {
				break
			}
		}
	}
	return rval
}

// AllFieldsTypeFunc calls a function on all the fields of a given struct type,
// including those on *any* fields of struct fields that this struct has -- if fun
// returns false then iteration stops -- overall rval is false if iteration
// was stopped or there was an error (logged), true otherwise.
func AllFieldsTypeFunc(typ reflect.Type, fun func(typ reflect.Type, field reflect.StructField) bool) bool {
	typ = NonPtrType(typ)
	if typ.Kind() != reflect.Struct {
		log.Printf("laser.AllFieldsTypeFunc: Must call on a struct type, not: %v\n", typ)
		return false
	}
	rval := true
	for i := 0; i < typ.NumField(); i++ {
		f := typ.Field(i)
		if f.Type.Kind() == reflect.Struct {
			rval = AllFieldsTypeFunc(f.Type, fun) // no err here
			if !rval {
				break
			}
		} else {
			rval = fun(typ, f)
			if !rval {
				break
			}
		}
	}
	return rval
}

// FlatFieldsValueFunc calls a function on all the primary fields of a
// given struct value (must pass a pointer to the struct) including those on
// anonymous embedded structs that this struct has, passing the current
// (embedded) type and StructField, which effectively flattens the reflect field list.
func FlatFieldsValueFunc(stru any, fun func(stru any, typ reflect.Type, field reflect.StructField, fieldVal reflect.Value) bool) bool {
	return FlatFieldsValueFuncIf(stru, nil, fun)
}

// FlatFieldsValueFunc calls a function on all the primary fields of a
// given struct value (must pass a pointer to the struct) including those on
// anonymous embedded structs that this struct has, passing the current
// (embedded) type and StructField, which effectively flattens the reflect field
// list. If the given ifFun is non-nil, it is called on every embedded struct field to
// determine whether the fields of that embedded field should be handled (a return value
// of true indicates to continue down and a value of false indicates to not).
func FlatFieldsValueFuncIf(stru any, ifFun, fun func(stru any, typ reflect.Type, field reflect.StructField, fieldVal reflect.Value) bool) bool {
	vv := reflect.ValueOf(stru)
	if stru == nil || vv.Kind() != reflect.Ptr {
		log.Printf("laser.FlatFieldsValueFunc: must pass a non-nil pointer to the struct: %v\n", stru)
		return false
	}
	v := NonPtrValue(vv)
	if !v.IsValid() {
		return true
	}
	typ := v.Type()
	if typ.Kind() != reflect.Struct {
		// log.Printf("laser.FlatFieldsValueFunc: non-pointer type is not a struct: %v\n", typ.String())
		return false
	}
	rval := true
	for i := 0; i < typ.NumField(); i++ {
		f := typ.Field(i)
		vf := v.Field(i)
		if !vf.CanInterface() {
			continue
		}
		vfi := vf.Interface()
		if vfi == stru {
			continue
		}
		if f.Type.Kind() == reflect.Struct && f.Anonymous {
			if ifFun != nil {
				if !ifFun(vfi, typ, f, vf) {
					continue
				}
			}
			// key to take addr here so next level is addressable
			rval = FlatFieldsValueFunc(PtrValue(vf).Interface(), fun)
			if !rval {
				break
			}
		} else {
			rval = fun(vfi, typ, f, vf)
			if !rval {
				break
			}
		}
	}
	return rval
}

// FlatFields returns a slice list of all the StructField type information for
// fields of given type and any embedded types -- returns nil on error
// (logged)
func FlatFields(typ reflect.Type) []reflect.StructField {
	ff := make([]reflect.StructField, 0)
	falseErr := FlatFieldsTypeFunc(typ, func(typ reflect.Type, field reflect.StructField) bool {
		ff = append(ff, field)
		return true
	})
	if falseErr == false {
		return nil
	}
	return ff
}

// AllFields returns a slice list of all the StructField type information for
// all elemental fields of given type and all embedded types -- returns nil on
// error (logged)
func AllFields(typ reflect.Type) []reflect.StructField {
	ff := make([]reflect.StructField, 0)
	falseErr := AllFieldsTypeFunc(typ, func(typ reflect.Type, field reflect.StructField) bool {
		ff = append(ff, field)
		return true
	})
	if falseErr == false {
		return nil
	}
	return ff
}

// AllFieldsN returns number of elemental fields in given type
func AllFieldsN(typ reflect.Type) int {
	n := 0
	falseErr := AllFieldsTypeFunc(typ, func(typ reflect.Type, field reflect.StructField) bool {
		n++
		return true
	})
	if falseErr == false {
		return 0
	}
	return n
}

// FlatFieldsVals returns a slice list of all the field reflect.Value's for
// fields of given struct (must pass a pointer to the struct) and any of its
// embedded structs -- returns nil on error (logged)
func FlatFieldVals(stru any) []reflect.Value {
	ff := make([]reflect.Value, 0)
	falseErr := FlatFieldsValueFunc(stru, func(stru any, typ reflect.Type, field reflect.StructField, fieldVal reflect.Value) bool {
		ff = append(ff, fieldVal)
		return true
	})
	if falseErr == false {
		return nil
	}
	return ff
}

// FlatFieldInterfaces returns a slice list of all the field interface{}
// values *as pointers to the field value* (i.e., calling Addr() on the Field
// Value) for fields of given struct (must pass a pointer to the struct) and
// any of its embedded structs -- returns nil on error (logged)
func FlatFieldInterfaces(stru any) []any {
	ff := make([]any, 0)
	falseErr := FlatFieldsValueFunc(stru, func(stru any, typ reflect.Type, field reflect.StructField, fieldVal reflect.Value) bool {
		ff = append(ff, PtrValue(fieldVal).Interface())
		return true
	})
	if falseErr == false {
		return nil
	}
	return ff
}

// FlatFieldByName returns field in type or embedded structs within type, by
// name -- native function already does flat version, so this is just for
// reference and consistency
func FlatFieldByName(typ reflect.Type, nm string) (reflect.StructField, bool) {
	return typ.FieldByName(nm)
}

// FieldByPath returns field in type or embedded structs within type, by a
// dot-separated path -- finds field by name for each level of the path, and
// recurses.
func FieldByPath(typ reflect.Type, path string) (reflect.StructField, bool) {
	pels := strings.Split(path, ".")
	ctyp := typ
	plen := len(pels)
	for i, pe := range pels {
		fld, ok := ctyp.FieldByName(pe)
		if !ok {
			log.Printf("laser.FieldByPath: field: %v not found in type: %v, starting from path: %v, in type: %v\n", pe, ctyp.String(), path, typ.String())
			return fld, false
		}
		if i == plen-1 {
			return fld, true
		}
		ctyp = fld.Type
	}
	return reflect.StructField{}, false
}

// FieldValueByPath returns field interface in type or embedded structs within
// type, by a dot-separated path -- finds field by name for each level of the
// path, and recurses.
func FieldValueByPath(stru any, path string) (reflect.Value, bool) {
	pels := strings.Split(path, ".")
	sval := reflect.ValueOf(stru)
	cval := sval
	typ := sval.Type()
	ctyp := typ
	plen := len(pels)
	for i, pe := range pels {
		_, ok := ctyp.FieldByName(pe)
		if !ok {
			log.Printf("laser.FieldValueByPath: field: %v not found in type: %v, starting from path: %v, in type: %v\n", pe, cval.Type().String(), path, typ.String())
			return cval, false
		}
		fval := cval.FieldByName(pe)
		if i == plen-1 {
			return fval, true
		}
		cval = fval
		ctyp = fval.Type()
	}
	return reflect.Value{}, false
}

// FlatFieldTag returns given tag value in field in type or embedded structs
// within type, by name -- empty string if not set or field not found
func FlatFieldTag(typ reflect.Type, nm, tag string) string {
	fld, ok := typ.FieldByName(nm)
	if !ok {
		return ""
	}
	return fld.Tag.Get(tag)
}

// FlatFieldValueByName finds field in object and embedded objects, by name,
// returning reflect.Value of field -- native version of Value function
// already does flat find, so this just provides a convenient wrapper
func FlatFieldValueByName(stru any, nm string) reflect.Value {
	vv := reflect.ValueOf(stru)
	if stru == nil || vv.Kind() != reflect.Ptr {
		log.Printf("laser.FlatFieldsValueFunc: must pass a non-nil pointer to the struct: %v\n", stru)
		return reflect.Value{}
	}
	v := NonPtrValue(vv)
	return v.FieldByName(nm)
}

// FlatFieldInterfaceByName finds field in object and embedded objects, by
// name, returning interface{} to pointer of field, or nil if not found
func FlatFieldInterfaceByName(stru any, nm string) any {
	ff := FlatFieldValueByName(stru, nm)
	if !ff.IsValid() {
		return nil
	}
	return PtrValue(ff).Interface()
}

// TypeEmbeds checks if given type embeds another type, at any level of
// recursive embedding (including being the type itself)
func TypeEmbeds(typ, embed reflect.Type) bool {
	typ = NonPtrType(typ)
	embed = NonPtrType(embed)
	if typ == embed {
		return true
	}
	for i := 0; i < typ.NumField(); i++ {
		f := typ.Field(i)
		if f.Type.Kind() == reflect.Struct && f.Anonymous {
			// fmt.Printf("typ %v anon struct %v\n", typ.Name(), f.Name)
			if f.Type == embed {
				return true
			}
			return TypeEmbeds(f.Type, embed)
		}
	}
	return false
}

// Embed returns the embedded struct of given type within given struct
func Embed(stru any, embed reflect.Type) any {
	if AnyIsNil(stru) {
		return nil
	}
	v := NonPtrValue(reflect.ValueOf(stru))
	typ := v.Type()
	if typ == embed {
		return PtrValue(v).Interface()
	}
	for i := 0; i < typ.NumField(); i++ {
		f := typ.Field(i)
		if f.Type.Kind() == reflect.Struct && f.Anonymous { // anon only avail on StructField fm typ
			vf := v.Field(i)
			vfpi := PtrValue(vf).Interface()
			if f.Type == embed {
				return vfpi
			}
			rv := Embed(vfpi, embed)
			if rv != nil {
				return rv
			}
		}
	}
	return nil
}

// EmbedImplements checks if given type implements given interface, or
// it embeds a type that does so -- must pass a type constructed like this:
// reflect.TypeOf((*gi.Node2D)(nil)).Elem() or just reflect.TypeOf(laser.BaseIface())
func EmbedImplements(typ, iface reflect.Type) bool {
	if iface.Kind() != reflect.Interface {
		log.Printf("laser.EmbedImplements -- type is not an interface: %v\n", iface)
		return false
	}
	if typ.Implements(iface) {
		return true
	}
	if reflect.PtrTo(typ).Implements(iface) { // typically need the pointer type to impl
		return true
	}
	typ = NonPtrType(typ)
	if typ.Implements(iface) { // try it all possible ways..
		return true
	}
	if typ.Kind() != reflect.Struct {
		return false
	}
	for i := 0; i < typ.NumField(); i++ {
		f := typ.Field(i)
		if f.Type.Kind() == reflect.Struct && f.Anonymous {
			rv := EmbedImplements(f.Type, iface)
			if rv {
				return true
			}
		}
	}
	return false
}

// SetFromDefaultTags sets values of fields in given struct based on
// `default:` default value field tags.
func SetFromDefaultTags(obj any) error {
	if AnyIsNil(obj) {
		return nil
	}
	ov := reflect.ValueOf(obj)
	if ov.Kind() == reflect.Pointer && ov.IsNil() {
		return nil
	}
	val := NonPtrValue(ov)
	typ := val.Type()
	for i := 0; i < typ.NumField(); i++ {
		f := typ.Field(i)
		fv := val.Field(i)
		def := f.Tag.Get("default")
		if NonPtrType(f.Type).Kind() == reflect.Struct && def == "" {
			SetFromDefaultTags(PtrValue(fv).Interface())
			continue
		}
		err := SetFromDefaultTag(fv, def)
		if err != nil {
			return fmt.Errorf("laser.SetFromDefaultTags: error setting field %q in object of type %q from val %q: %w", f.Name, typ.Name(), def, err)
		}
	}
	return nil
}

// SetFromDefaultTag sets the given value from the given default tag.
func SetFromDefaultTag(v reflect.Value, def string) error {
	def = FormatDefault(def)
	if def == "" {
		return nil
	}
	return SetRobust(PtrValue(v).Interface(), def)
}

// NonDefaultFields returns a map representing all of the fields of the given
// struct (or pointer to a struct) that have values different than their default
// values as specified by the `default:` struct tag. The resulting map is then typically
// saved using something like JSON or TOML. If a value has no default value, it
// checks whether its value is non-zero. If a field has a `save:"-"` tag, it wil
// not be included in the resulting map.
func NonDefaultFields(v any) map[string]any {
	res := map[string]any{}

	rv := NonPtrValue(reflect.ValueOf(v))
	if !rv.IsValid() {
		return nil
	}
	rt := rv.Type()
	nf := rt.NumField()
	for i := 0; i < nf; i++ {
		fv := rv.Field(i)
		ft := rt.Field(i)
		if ft.Tag.Get("save") == "-" {
			continue
		}
		def := ft.Tag.Get("default")
		if NonPtrType(ft.Type).Kind() == reflect.Struct && def == "" {
			sfm := NonDefaultFields(fv.Interface())
			if len(sfm) > 0 {
				res[ft.Name] = sfm
			}
			continue
		}
		def = FormatDefault(def)
		if def == "" {
			if !fv.IsZero() {
				res[ft.Name] = fv.Interface()
			}
			continue
		}
		dv := reflect.New(ft.Type)
		err := SetRobust(dv.Interface(), def)
		if err != nil {
			slog.Error("laser.NonDefaultFields: error getting value from default struct tag", "field", ft.Name, "type", rt, "defaultStructTag", def, "err", err)
			res[ft.Name] = fv.Interface()
			continue
		}
		if !reflect.DeepEqual(fv.Interface(), dv.Elem().Interface()) {
			res[ft.Name] = fv.Interface()
		}
	}
	return res
}

// FormatDefault converts the given `default:` struct tag string into a format suitable
// for being used as a value in [SetRobust]. If it returns "", the default value
// should not be used.
func FormatDefault(def string) string {
	if def == "" {
		return ""
	}
	if strings.ContainsAny(def, "{[") { // complex type, so don't split on commas and colons
		return strings.ReplaceAll(def, `'`, `"`) // allow single quote to work as double quote for JSON format
	}
	// we split on commas and colons so we get the first item of lists and ranges
	def = strings.Split(def, ",")[0]
	def = strings.Split(def, ":")[0]
	return def
}

// StructTags returns a map[string]string of the tag string from a reflect.StructTag value
// e.g., from StructField.Tag
func StructTags(tags reflect.StructTag) map[string]string {
	if len(tags) == 0 {
		return nil
	}
	flds := strings.Fields(string(tags))
	smap := make(map[string]string, len(flds))
	for _, fld := range flds {
		cli := strings.Index(fld, ":")
		if cli < 0 || len(fld) < cli+3 {
			continue
		}
		vl := strings.TrimSuffix(fld[cli+2:], `"`)
		smap[fld[:cli]] = vl
	}
	return smap
}

// StringJSON returns a JSON representation of item, as a string
// e.g., for printing / debugging etc.
func StringJSON(it any) string {
	b, _ := json.MarshalIndent(it, "", "  ")
	return string(b)
}

// SetField sets given field name on given struct object to given value,
// using very robust conversion routines to e.g., convert from strings to numbers, and
// vice-versa, automatically.  Returns error if not successfully set.
func SetField(obj any, field string, val any) error {
	fv := FlatFieldValueByName(obj, field)
	if !fv.IsValid() {
		return fmt.Errorf("laser.SetField: could not find field %q", field)
	}
	err := SetRobust(PtrValue(fv).Interface(), val)
	if err != nil {
		return fmt.Errorf("laser.SetField: SetRobust failed to set field %q to value: %v: %w", field, val, err)
	}
	return nil
}
