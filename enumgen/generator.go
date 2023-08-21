// Copyright (c) 2023, The GoKi Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Based on http://github.com/dmarkham/enumer and
// golang.org/x/tools/cmd/stringer:

// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package enumgen

import (
	"bytes"
	"errors"
	"fmt"
	"go/ast"
	"go/format"
	"os"
	"strings"

	"golang.org/x/tools/go/packages"
)

// Generator holds the state of the generator.
// It is primarily used to buffer
// the output for [format.Source].
type Generator struct {
	Config Config                 // The configuration information
	Buf    bytes.Buffer           // The accumulated output.
	Pkg    *Package               // The package we are scanning.
	Types  map[*ast.TypeSpec]bool // The enum types; the value is whether they are a bit flag or not
}

// NewGenerator returns a new generator with the
// given configuration information.
func NewGenerator(config Config) *Generator {
	return &Generator{Config: config}
}

// ParsePackage parses the single package located in the configuration directory.
func (g *Generator) ParsePackage() error {
	cfg := &packages.Config{
		Mode: packages.NeedName | packages.NeedFiles | packages.NeedCompiledGoFiles | packages.NeedImports | packages.NeedTypes | packages.NeedTypesSizes | packages.NeedSyntax | packages.NeedTypesInfo,
		// TODO: Need to think about constants in test files. Maybe write type_string_test.go
		// in a separate pass? For later.
		Tests: false,
	}
	pkgs, err := packages.Load(cfg, g.Config.Dir)
	if err != nil {
		return err
	}
	if len(pkgs) != 1 {
		return fmt.Errorf("expected 1 package, but found %d packages", len(pkgs))
	}
	g.AddPackage(pkgs[0])
	return nil
}

// AddPackage adds a type-checked Package and its syntax files to the generator.
func (g *Generator) AddPackage(pkg *packages.Package) {
	g.Pkg = &Package{
		Name:  pkg.Name,
		Defs:  pkg.TypesInfo.Defs,
		Files: make([]*File, len(pkg.Syntax)),
	}

	for i, file := range pkg.Syntax {
		g.Pkg.Files[i] = &File{
			File: file,
			Pkg:  g.Pkg,
		}
	}
}

// Printf prints the formatted string to the
// accumulated output in [Generator.Buf]
func (g *Generator) Printf(format string, args ...any) {
	fmt.Fprintf(&g.Buf, format, args...)
}

// PrintHeader prints the header and package clause
// to the accumulated output
func (g *Generator) PrintHeader() {
	g.Printf("// Code generated by \"enumgen %s\"; DO NOT EDIT.\n", strings.Join(os.Args[1:], " "))
	g.Printf("\n")
	if g.Config.Comment != "" {
		g.Printf("// %s\n", g.Config.Comment)
	}
	g.Printf("package %s", g.Pkg.Name)
	g.Printf("\n")
	g.Printf("import (\n")
	g.Printf("\t\"fmt\"\n")
	g.Printf("\t\"strings\"\n")
	g.Printf("\t\"strconv\"\n")
	g.Printf("\t\"errors\"\n")
	g.Printf("\t\"sync/atomic\"\n")
	if g.Config.SQL {
		g.Printf("\t\"database/sql/driver\"\n")
	}
	if g.Config.JSON {
		g.Printf("\t\"encoding/json\"\n")
	}
	if g.Config.GQLGEN {
		g.Printf("\t\"io\"\n")
	}
	g.Printf("\t\"github.com/goki/enums/enums\"\n")
	g.Printf(")\n")
}

// FindEnumTypes goes through all of the types in the package
// and finds all integer (signed or unsigned) types labeled with enums:enum
// or enums:bitflag. It stores the resulting types in [Generator.Types].
func (g *Generator) FindEnumTypes() error {
	g.Types = map[*ast.TypeSpec]bool{}
	for _, file := range g.Pkg.Files {
		var terr error
		ast.Inspect(file.File, func(n ast.Node) bool {
			if terr != nil {
				return false
			}
			cont, err := g.InspectForType(n)
			if err != nil {
				terr = err
			}
			return cont
		})
		if terr != nil {
			return fmt.Errorf("FindEnumTypes: error finding enum types: %w", terr)
		}
	}
	return nil
}

// InspectForType looks at the given AST node and adds it
// to [Generator.Types] if it is marked with an appropriate
// comment directive. It returns whether the AST inspector should
// continue, and an error if there is one. It should only
// be called in [ast.Inspect].
func (g *Generator) InspectForType(n ast.Node) (bool, error) {
	typ, ok := n.(*ast.TypeSpec)
	if !ok {
		return true, nil
	}
	if typ.Comment == nil {
		return true, nil
	}
	for _, c := range typ.Comment.List {
		if strings.HasPrefix(c.Text, "//enums:") {
			d := strings.TrimPrefix(c.Text, "//enums:")
			switch d {
			case "enum":
				g.Types[typ] = false
			case "bitflag":
				g.Types[typ] = true
			default:
				return false, errors.New("unrecognized enums directive: '" + c.Text + "'")
			}
		}
	}
	return true, nil
}

// Generate produces the enum methods for the types
// stored in [Generator.Types].
func (g *Generator) Generate() error {
	for typ, bitflag := range g.Types {
		values := make([]Value, 0, 100)
		typeName := typ.Name.String()
		for _, file := range g.Pkg.Files {
			file.Config = g.Config
			// Set the state for this run of the walker.
			file.TypeName = typeName
			file.BitFlag = bitflag
			file.Values = nil
			if file.File != nil {
				var terr error
				ast.Inspect(file.File, func(n ast.Node) bool {
					if terr != nil {
						return false
					}
					cont, err := file.GenDecl(n)
					if err != nil {
						terr = err
					}
					return cont
				})
				if terr != nil {
					return fmt.Errorf("Generate: error parsing declaration clauses: %w", terr)
				}
				values = append(values, file.Values...)
			}
		}

		if len(values) == 0 {
			return errors.New("no values defined for type " + typeName)
		}

		g.TrimValueNames(values)

		g.TransformValueNames(values)

		g.PrefixValueNames(values)

		runs := SplitIntoRuns(values)
		// The decision of which pattern to use depends on the number of
		// runs in the numbers. If there's only one, it's easy. For more than
		// one, there's a tradeoff between complexity and size of the data
		// and code vs. the simplicity of a map. A map takes more space,
		// but so does the code. The decision here (crossover at 10) is
		// arbitrary, but considers that for large numbers of runs the cost
		// of the linear scan in the switch might become important, and
		// rather than use yet another algorithm such as binary search,
		// we punt and use a map. In any case, the likelihood of a map
		// being necessary for any realistic example other than bitmasks
		// is very low. And bitmasks probably deserve their own analysis,
		// to be done some other day.
		const runsThreshold = 10
		switch {
		case len(runs) == 1:
			g.BuildOneRun(runs, typeName)
		case len(runs) <= runsThreshold:
			g.BuildMultipleRuns(runs, typeName)
		default:
			g.BuildMap(runs, typeName)
		}

		g.BuildNoOpOrderChangeDetect(runs, typeName)

		g.BuildBasicExtras(runs, typeName, runsThreshold)
		if bitflag {
			g.BuildBitFlagMethods(runs, typeName)
		}

		if g.Config.JSON {
			g.BuildJSONMethods(runs, typeName, runsThreshold)
		}
		if g.Config.Text {
			g.BuildTextMethods(runs, typeName, runsThreshold)
		}
		if g.Config.YAML {
			g.BuildYAMLMethods(runs, typeName, runsThreshold)
		}
		if g.Config.SQL {
			g.addValueAndScanMethod(typeName)
		}
		if g.Config.GQLGEN {
			g.buildGQLGenMethods(runs, typeName)
		}
	}
	return nil
}

// Format returns the contents of the Generator's buffer
// ([Generator.Buf]) with gofmt applied.
func (g *Generator) Format() ([]byte, error) {
	src, err := format.Source(g.Buf.Bytes())
	if err != nil {
		// Should never happen, but can arise when developing this code.
		// The user can compile the output to see the error.
		return g.Buf.Bytes(), errors.New("internal error: invalid Go generated: " + err.Error() + "; compile the package to analyze the error")
	}
	return src, nil
}

// Write formats the data in the the Generator's buffer
// ([Generator.Buf]) and writes it to the file specified by
// [Generator.Config.Output].
func (g *Generator) Write() error {
	b, err := g.Format()
	if err != nil {
		return fmt.Errorf("Generator.Write: error formatting code: %w", err)
	}
	err = os.WriteFile(g.Config.Output, b, 0666)
	if err != nil {
		return fmt.Errorf("Generator.Write: error writing file: %w", err)
	}
	return nil
}
