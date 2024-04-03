// Copyright (c) 2018, Cogent Core. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tex

import (
	_ "embed"
	"strings"
	"unicode"

	"cogentcore.org/core/fi"
	"cogentcore.org/core/glop/indent"
	"cogentcore.org/core/pi"
	"cogentcore.org/core/pi/langs"
	"cogentcore.org/core/pi/langs/bibtex"
	"cogentcore.org/core/pi/lex"
	"cogentcore.org/core/pi/syms"
)

//go:embed tex.pi
var parserBytes []byte

// TexLang implements the Lang interface for the Tex / LaTeX language
type TexLang struct {
	Pr *pi.Parser

	// bibliography files that have been loaded, keyed by file path from bibfile metadata stored in filestate
	Bibs bibtex.Files
}

// TheTexLang is the instance variable providing support for the Go language
var TheTexLang = TexLang{}

func init() {
	pi.StandardLangProps[fi.TeX].Lang = &TheTexLang
	langs.ParserBytes[fi.TeX] = parserBytes
}

func (tl *TexLang) Parser() *pi.Parser {
	if tl.Pr != nil {
		return tl.Pr
	}
	lp, _ := pi.LangSupport.Props(fi.TeX)
	if lp.Parser == nil {
		pi.LangSupport.OpenStandard()
	}
	tl.Pr = lp.Parser
	if tl.Pr == nil {
		return nil
	}
	return tl.Pr
}

func (tl *TexLang) ParseFile(fss *pi.FileStates, txt []byte) {
	pr := tl.Parser()
	if pr == nil {
		return
	}
	pfs := fss.StartProc(txt) // current processing one
	pr.LexAll(pfs)
	tl.OpenBibfile(fss, pfs)
	fss.EndProc() // now done
	// no parser
}

func (tl *TexLang) LexLine(fs *pi.FileState, line int, txt []rune) lex.Line {
	pr := tl.Parser()
	if pr == nil {
		return nil
	}
	return pr.LexLine(fs, line, txt)
}

func (tl *TexLang) ParseLine(fs *pi.FileState, line int) *pi.FileState {
	// n/a
	return nil
}

func (tl *TexLang) HiLine(fss *pi.FileStates, line int, txt []rune) lex.Line {
	fs := fss.Done()
	return tl.LexLine(fs, line, txt)
}

func (tl *TexLang) ParseDir(fs *pi.FileState, path string, opts pi.LangDirOpts) *syms.Symbol {
	// n/a
	return nil
}

// IndentLine returns the indentation level for given line based on
// previous line's indentation level, and any delta change based on
// e.g., brackets starting or ending the previous or current line, or
// other language-specific keywords.  See lex.BracketIndentLine for example.
// Indent level is in increments of tabSz for spaces, and tabs for tabs.
// Operates on rune source with markup lex tags per line.
func (tl *TexLang) IndentLine(fs *pi.FileStates, src [][]rune, tags []lex.Line, ln int, tabSz int) (pInd, delInd, pLn int, ichr indent.Char) {
	pInd, pLn, ichr = lex.PrevLineIndent(src, tags, ln, tabSz)

	curUnd, _ := lex.LineStartEndBracket(src[ln], tags[ln])
	_, prvInd := lex.LineStartEndBracket(src[pLn], tags[pLn])

	delInd = 0
	switch {
	case prvInd && curUnd:
		delInd = 0 // offset
	case prvInd:
		delInd = 1 // indent
	case curUnd:
		delInd = -1 // undent
	}

	pst := lex.FirstNonSpaceRune(src[pLn])
	cst := lex.FirstNonSpaceRune(src[ln])

	pbeg := false
	if pst >= 0 {
		sts := string(src[pLn][pst:])
		if strings.HasPrefix(sts, "\\begin{") {
			pbeg = true
		}
	}

	cend := false
	if cst >= 0 {
		sts := string(src[ln][cst:])
		if strings.HasPrefix(sts, "\\end{") {
			cend = true
		}
	}

	switch {
	case pbeg && cend:
		delInd = 0
	case pbeg:
		delInd = 1
	case cend:
		delInd = -1
	}

	if pInd == 0 && delInd < 0 { // error..
		delInd = 0
	}
	return
}

// AutoBracket returns what to do when a user types a starting bracket character
// (bracket, brace, paren) while typing.
// pos = position where bra will be inserted, and curLn is the current line
// match = insert the matching ket, and newLine = insert a new line.
func (tl *TexLang) AutoBracket(fs *pi.FileStates, bra rune, pos lex.Pos, curLn []rune) (match, newLine bool) {
	lnLen := len(curLn)
	match = pos.Ch == lnLen || unicode.IsSpace(curLn[pos.Ch]) // at end or if space after
	newLine = false
	return
}
