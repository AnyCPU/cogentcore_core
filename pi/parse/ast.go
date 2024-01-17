// Copyright (c) 2018, Cogent Core. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package parse does the parsing stage after lexing
package parse

//go:generate core generate

import (
	"fmt"
	"io"

	"cogentcore.org/core/glop/indent"
	"cogentcore.org/core/ki"
	"cogentcore.org/core/pi/lex"
	"cogentcore.org/core/pi/syms"
)

// Ast is a node in the abstract syntax tree generated by the parsing step
// the name of the node (from ki.Node) is the type of the element
// (e.g., expr, stmt, etc)
// These nodes are generated by the parse.Rule's by matching tokens
type Ast struct {
	ki.Node

	// region in source lexical tokens corresponding to this Ast node -- Ch = index in lex lines
	TokReg lex.Reg `set:"-"`

	// region in source file corresponding to this Ast node
	SrcReg lex.Reg `set:"-"`

	// source code corresponding to this Ast node
	Src string `set:"-"`

	// stack of symbols created for this node
	Syms syms.SymStack `set:"-"`
}

func (ast *Ast) Destroy() {
	ast.Syms.ClearAst()
	ast.Syms = nil
	ast.Node.Destroy()
}

// ChildAst returns the Child at given index as an Ast.
// Will panic if index is invalid -- use Try if unsure.
func (ast *Ast) ChildAst(idx int) *Ast {
	return ast.Child(idx).(*Ast)
}

// ParAst returns the Parent as an Ast.
func (ast *Ast) ParAst() *Ast {
	if ast.Par == nil {
		return nil
	}
	pki := ast.Par.This()
	if pki == nil {
		return nil
	}
	return pki.(*Ast)
}

// NextAst returns the next node in the Ast tree, or nil if none
func (ast *Ast) NextAst() *Ast {
	nxti := ki.Next(ast)
	if nxti == nil {
		return nil
	}
	return nxti.(*Ast)
}

// NextSiblingAst returns the next sibling node in the Ast tree, or nil if none
func (ast *Ast) NextSiblingAst() *Ast {
	nxti := ki.NextSibling(ast)
	if nxti == nil {
		return nil
	}
	return nxti.(*Ast)
}

// PrevAst returns the previous node in the Ast tree, or nil if none
func (ast *Ast) PrevAst() *Ast {
	nxti := ki.Prev(ast)
	if nxti == nil {
		return nil
	}
	return nxti.(*Ast)
}

// SetTokReg sets the token region for this rule to given region
func (ast *Ast) SetTokReg(reg lex.Reg, src *lex.File) {
	ast.TokReg = reg
	ast.SrcReg = src.TokenSrcReg(ast.TokReg)
	ast.Src = src.RegSrc(ast.SrcReg)
}

// SetTokRegEnd updates the ending token region to given position --
// token regions are typically over-extended and get narrowed as tokens actually match
func (ast *Ast) SetTokRegEnd(pos lex.Pos, src *lex.File) {
	ast.TokReg.Ed = pos
	ast.SrcReg = src.TokenSrcReg(ast.TokReg)
	ast.Src = src.RegSrc(ast.SrcReg)
}

// WriteTree writes the AST tree data to the writer -- not attempting to re-render
// source code -- just for debugging etc
func (ast *Ast) WriteTree(out io.Writer, depth int) {
	ind := indent.Tabs(depth)
	fmt.Fprintf(out, "%v%v: %v\n", ind, ast.Nm, ast.Src)
	for _, k := range ast.Kids {
		ai := k.(*Ast)
		ai.WriteTree(out, depth+1)
	}
}

var AstProps = ki.Props{
	"StructViewFields": ki.Props{ // hide in view
		"Flag":  `view:"-"`,
		"Props": `view:"-"`,
	},
	// "CallMethods": ki.PropSlice{
	// 	{"SaveAs", ki.Props{
	// 		"Args": ki.PropSlice{
	// 			{"File Name", ki.Props{
	// 				"default-field": "Filename",
	// 			}},
	// 		},
	// 	}},
	// },
}
