// Copyright (c) 2023, Cogent Core. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Based on http://github.com/dmarkham/enumer and
// golang.org/x/tools/cmd/stringer:

// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package enumgen

import "text/template"

// BuildBitFlagMethods builds methods specific to bit flag types.
func (g *Generator) BuildBitFlagMethods(runs []Value, typ *Type) {
	g.Printf("\n")

	g.ExecTmpl(HasFlagMethodTmpl, typ)
	g.ExecTmpl(SetFlagMethodTmpl, typ)
}

var HasFlagMethodTmpl = template.Must(template.New("HasFlagMethod").Parse(
	`// HasFlag returns whether these bit flags have the given bit flag set.
func (i {{.Name}}) HasFlag(f enums.BitFlag) bool {
	return atomic.LoadInt64((*int64)(&i))&(1<<uint32(f.Int64())) != 0
}
`))

var SetFlagMethodTmpl = template.Must(template.New("SetFlagMethod").Parse(
	`// SetFlag sets the value of the given flags in these flags to the given value.
func (i *{{.Name}}) SetFlag(on bool, f ...enums.BitFlag) { enums.SetFlag((*int64)(i), on, f...) }
`))

var StringMethodBitFlagTmpl = template.Must(template.New("StringMethodBitFlag").Parse(
	`// String returns the string representation of this {{.Name}} value.
func (i {{.Name}}) String() string {
	{{- if eq .Extends ""}} return enums.BitFlagString(i, _{{.Name}}Values)
	{{- else}} return enums.BitFlagStringExtended(i, _{{.Name}}Values, {{.Extends}}Values()) {{end}} }
`))

var SetStringMethodBitFlagTmpl = template.Must(template.New("SetStringMethodBitFlag").Parse(
	`// SetString sets the {{.Name}} value from its string representation,
// and returns an error if the string is invalid.
func (i *{{.Name}}) SetString(s string) error { *i = 0; return i.SetStringOr(s) }
`))

var SetStringOrMethodBitFlagTmpl = template.Must(template.New("SetStringOrMethodBitFlag").Parse(
	`// SetStringOr sets the {{.Name}} value from its string representation
// while preserving any bit flags already set, and returns an
// error if the string is invalid.
func (i *{{.Name}}) SetStringOr(s string) error {
	flgs := strings.Split(s, "|")
	for _, flg := range flgs {
		if val, ok := _{{.Name}}ValueMap[flg]; ok {
			i.SetFlag(true, &val)
		{{if .Config.AcceptLower}} } else if val, ok := _{{.Name}}ValueMap[strings.ToLower(flg)]; ok {
			i.SetFlag(true, &val)
		{{end}} } else if flg == "" {
			continue
		} else { {{if eq .Extends ""}}
			return fmt.Errorf("%q is not a valid value for type {{.Name}}", flg){{else}}
			err := (*{{.Extends}})(i).SetStringOr(flg)
			if err != nil {
				return err
			}{{end}}
		}
	}
	return nil
}
`))
