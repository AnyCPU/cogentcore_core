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

var TextMethodsTmpl = template.Must(template.New("TextMethods").Parse(
	`
// MarshalText implements the [encoding.TextMarshaler] interface.
func (i {{.Name}}) MarshalText() ([]byte, error) { return []byte(i.String()), nil }

// UnmarshalText implements the [encoding.TextUnmarshaler] interface.
func (i *{{.Name}}) UnmarshalText(text []byte) error { return enums.UnmarshalText(i, text, "{{.Name}}") }
`))

func (g *Generator) BuildTextMethods(runs []Value, typ *Type) {
	g.ExecTmpl(TextMethodsTmpl, typ)
}

var JSONMethodsTmpl = template.Must(template.New("JSONMethods").Parse(
	`
// MarshalJSON implements the [json.Marshaler] interface.
func (i {{.Name}}) MarshalJSON() ([]byte, error) { return json.Marshal(i.String()) }

// UnmarshalJSON implements the [json.Unmarshaler] interface.
func (i *{{.Name}}) UnmarshalJSON(data []byte) error { return enums.UnmarshalJSON(i, data, "{{.Name}}") }
`))

func (g *Generator) BuildJSONMethods(runs []Value, typ *Type) {
	g.ExecTmpl(JSONMethodsTmpl, typ)
}

var YAMLMethodsTmpl = template.Must(template.New("YAMLMethods").Parse(
	`
// MarshalYAML implements the [yaml.Marshaler] interface.
func (i {{.Name}}) MarshalYAML() (any, error) { return i.String(), nil }

// UnmarshalYAML implements the [yaml.Unmarshaler] interface.
func (i *{{.Name}}) UnmarshalYAML(value *yaml.Node) error {
	var s string
	if err := n.Decode(&s); err != nil {
		return err
	}
	if err := i.SetString(s); err != nil {
		log.Println("{{.Name}}.UnmarshalYAML:", err)
	}
	return nil
}
`))

func (g *Generator) BuildYAMLMethods(runs []Value, typ *Type) {
	g.ExecTmpl(YAMLMethodsTmpl, typ)
}
