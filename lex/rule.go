// Copyright (c) 2018, The GoKi Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package lex

import (
	"fmt"
	"io"
	"reflect"
	"strings"
	"text/tabwriter"

	"github.com/goki/ki"
	"github.com/goki/ki/indent"
	"github.com/goki/ki/kit"
	"github.com/goki/pi/token"
)

// Lexer is the interface type for lexers -- likely not necessary except is essential
// for defining the BaseIface for gui in making new nodes
type Lexer interface {
	ki.Ki

	// Validate checks for any errors in the rules and issues warnings,
	// returns true if valid (no err) and false if invalid (errs)
	Validate(ls *State) bool

	// Lex tries to apply rule to given input state, returns true if matched, false if not
	Lex(ls *State) *Rule

	// AsLexRule returns object as a lex.Rule
	AsLexRule() *Rule
}

// lex.Rule operates on the text input to produce the lexical tokens
// it is assembled into a lexical grammar structure to perform lexing
//
// Lexing is done line-by-line -- you must push and pop states to
// coordinate across multiple lines, e.g., for multi-line comments
//
// In general it is best to keep lexing as simple as possible and
// leave the more complex things for the parsing step.
type Rule struct {
	ki.Node
	Desc      string       `desc:"description / comments about this rule"`
	Token     token.Tokens `desc:"the token value that this rule generates -- use None for non-terminals"`
	Match     Matches      `desc:"the lexical match that we look for to engage this rule"`
	String    string       `desc:"if action is LexMatch, this is the string we match"`
	Off       int          `desc:"offset into the input to look for a match: 0 = current char, 1 = next one, etc"`
	Acts      []Actions    `desc:"the action(s) to perform, in order, if there is a match -- these are performed prior to iterating over child nodes"`
	PushState string       `desc:"the state to push if our action is PushState -- note that State matching is on String, not this value"`
	TokEff    token.Tokens `view:"-" json:"-" desc:"effective token based on input -- e.g., for number is the type of number"`
	MatchLen  int          `view:"-" json:"-" desc:"length of source that matched -- if Next is called, this is what will be skipped to"`
}

var KiT_Rule = kit.Types.AddType(&Rule{}, RuleProps)

func (lr *Rule) BaseIface() reflect.Type {
	return reflect.TypeOf((*Lexer)(nil)).Elem()
}

func (lr *Rule) AsLexRule() *Rule {
	return lr.This().Embed(KiT_Rule).(*Rule)
}

// Validate checks for any errors in the rules and issues warnings,
// returns true if valid (no err) and false if invalid (errs)
func (lr *Rule) Validate(ls *State) bool {
	valid := true
	if !lr.IsRoot() {
		switch lr.Match {
		case StrName:
			fallthrough
		case String:
			if len(lr.String) == 0 {
				valid = false
				ls.Error(0, fmt.Sprintf("lex.Rule: match = String or StrName but String is empty, in: %v\n", lr.PathUnique()))
			}
		case CurState:
			for _, act := range lr.Acts {
				if act == Next {
					valid = false
					ls.Error(0, fmt.Sprintf("lex.Rule: match = CurState cannot have Action = Next -- no src match, in: %v\n", lr.PathUnique()))
				}
			}
			if len(lr.String) == 0 {
				ls.Error(0, fmt.Sprintf("lex.Rule: match = CurState must have state to match in String -- is empty, in: %v\n", lr.PathUnique()))
			}
			if len(lr.PushState) > 0 {
				ls.Error(0, fmt.Sprintf("lex.Rule: match = CurState has non-empty PushState -- must have state to match in String instead, in: %v\n", lr.PathUnique()))
			}
		}
	}

	if !lr.HasChildren() && len(lr.Acts) == 0 {
		valid = false
		ls.Error(0, fmt.Sprintf("lex.Rule: has no children and no action -- does nothing, in: %v\n", lr.PathUnique()))
	}

	hasPos := false
	for _, act := range lr.Acts {
		if act >= Name && act <= EOL {
			hasPos = true
		}
		if act == Next && hasPos {
			valid = false
			ls.Error(0, fmt.Sprintf("lex.Rule: action = Next incompatible with action that reads item such as Name, Number, Quoted, in: %v\n", lr.PathUnique()))
		}
	}

	if lr.Token.Cat() == token.Keyword && lr.Match != StrName {
		valid = false
		ls.Error(0, fmt.Sprintf("lex.Rule: Keyword token must use StrName to match entire name, in: %v\n", lr.PathUnique()))
	}

	// now we iterate over our kids
	for _, klri := range lr.Kids {
		klr := klri.Embed(KiT_Rule).(*Rule)
		if !klr.Validate(ls) {
			valid = false
		}
	}
	return valid
}

// LexStart is called on the top-level lex node to start lexing process for one step
func (lr *Rule) LexStart(ls *State) *Rule {
	cpos := ls.Pos
	rval := lr.Lex(ls)
	if !ls.AtEol() && cpos == ls.Pos {
		msg := fmt.Sprintf("did not advance position -- need more rules to match current input: %v", string(ls.Src[cpos:]))
		ls.Error(cpos, msg)
		return nil
	}
	return rval
}

// Lex tries to apply rule to given input state, returns lowest-level rule that matched, nil if none
func (lr *Rule) Lex(ls *State) *Rule {
	if !lr.IsMatch(ls) {
		return nil
	}
	st := ls.Pos // starting pos that we're consuming
	lr.TokEff = lr.Token
	for _, act := range lr.Acts {
		lr.DoAct(ls, act)
	}
	ed := ls.Pos // our ending state
	if ed > st {
		ls.Add(lr.TokEff, st, ed)
	}
	if !lr.HasChildren() {
		return lr
	}

	// now we iterate over our kids
	for _, klri := range lr.Kids {
		klr := klri.Embed(KiT_Rule).(*Rule)
		if mrule := klr.Lex(ls); mrule != nil { // first to match takes it -- order matters!
			return mrule
		}
	}

	// if kids don't match and we don't have any actions, we are just a grouper
	// and thus we depend entirely on kids matching
	if len(lr.Acts) == 0 {
		return nil
	}

	return lr
}

// IsMatch tests if the rule matches for current input state, returns true if so, false if not
func (lr *Rule) IsMatch(ls *State) bool {
	if lr.IsRoot() { // root always matches
		return true
	}
	switch lr.Match {
	case String:
		sz := len(lr.String)
		str, ok := ls.String(lr.Off, sz)
		if !ok {
			return false
		}
		if str != lr.String {
			return false
		}
		lr.MatchLen = lr.Off + sz
		return true
	case StrName:
		cp := ls.Pos
		ls.Pos += lr.Off
		st := ls.Pos
		ls.ReadName()
		ed := ls.Pos
		ls.Pos = cp
		nsz := ed - st
		sz := len(lr.String)
		if nsz != sz {
			return false
		}
		str := string(ls.Src[st:ed])
		if str != lr.String {
			return false
		}
		lr.MatchLen = lr.Off + sz
		return true
	case Letter:
		rn, ok := ls.Rune(lr.Off)
		if !ok {
			return false
		}
		if IsLetter(rn) {
			lr.MatchLen = lr.Off + 1
			return true
		}
		return false
	case Digit:
		rn, ok := ls.Rune(lr.Off)
		if !ok {
			return false
		}
		if IsDigit(rn) {
			lr.MatchLen = lr.Off + 1
			return true
		}
		return false
	case WhiteSpace:
		rn, ok := ls.Rune(lr.Off)
		if !ok {
			return false
		}
		if IsWhiteSpace(rn) {
			lr.MatchLen = lr.Off + 1
			return true
		}
		return false
	case CurState:
		if ls.CurState() == lr.String {
			lr.MatchLen = 0
			return true
		}
		return false
	case AnyRune:
		_, ok := ls.Rune(lr.Off)
		if !ok {
			return false
		}
		lr.MatchLen = lr.Off + 1
		return true
	}
	return false
}

// DoAct performs given action
func (lr *Rule) DoAct(ls *State, act Actions) {
	switch act {
	case Next:
		ls.Next(lr.MatchLen)
	case Name:
		ls.ReadName()
	case Number:
		lr.TokEff = ls.ReadNumber()
	case Quoted:
		ls.ReadQuoted()
	case QuotedRaw:
		ls.ReadQuoted() // todo: raw!
	case EOL:
		ls.Pos = len(ls.Src)
	case PushState:
		ls.PushState(lr.PushState)
	case PopState:
		ls.PopState()
	}
}

///////////////////////////////////////////////////////////////////////
//  Non-lexing functions

// Find looks for rules in the tree that contain given string in String or Name fields
func (lr *Rule) Find(find string) []*Rule {
	var res []*Rule
	lr.FuncDownMeFirst(0, lr.This(), func(k ki.Ki, level int, d interface{}) bool {
		lri := k.Embed(KiT_Rule).(*Rule)
		if strings.Contains(lri.String, find) || strings.Contains(lri.Nm, find) {
			res = append(res, lri)
		}
		return true
	})
	return res
}

// WriteGrammar outputs the lexer rules as a formatted grammar in a BNF-like format
// it is called recursively
func (lr *Rule) WriteGrammar(writer io.Writer, depth int) {
	if lr.IsRoot() {
		for _, k := range lr.Kids {
			lri := k.Embed(KiT_Rule).(*Rule)
			lri.WriteGrammar(writer, depth)
		}
	} else {
		ind := indent.Tabs(depth)
		gpstr := ""
		if lr.HasChildren() {
			gpstr = " {"
		}
		offstr := ""
		if lr.Off > 0 {
			offstr = fmt.Sprintf("+%d:", lr.Off)
		}
		actstr := ""
		if len(lr.Acts) > 0 {
			actstr = "\t do: "
			for _, ac := range lr.Acts {
				if ac == PushState {
					actstr += ac.String() + ": " + lr.PushState + "; "
				} else {
					actstr += ac.String() + "; "
				}
			}
		}
		if lr.Desc != "" {
			fmt.Fprintf(writer, "%v// %v %v \n", ind, lr.Nm, lr.Desc)
		}
		if (lr.Match >= Letter && lr.Match <= WhiteSpace) || lr.Match == AnyRune {
			fmt.Fprintf(writer, "%v%v:\t\t %v\t\t if %v%v%v%v\n", ind, lr.Nm, lr.Token, offstr, lr.Match, actstr, gpstr)
		} else {
			fmt.Fprintf(writer, "%v%v:\t\t %v\t\t if %v%v == \"%v\"%v%v\n", ind, lr.Nm, lr.Token, offstr, lr.Match, lr.String, actstr, gpstr)
		}
		if lr.HasChildren() {
			w := tabwriter.NewWriter(writer, 4, 4, 2, ' ', 0)
			for _, k := range lr.Kids {
				lri := k.Embed(KiT_Rule).(*Rule)
				lri.WriteGrammar(w, depth+1)
			}
			w.Flush()
			fmt.Fprintf(writer, "%v}\n", ind)
		}
	}
}

var RuleProps = ki.Props{
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
