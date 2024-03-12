// Code generated by "core generate"; DO NOT EDIT.

package keyfun

import (
	"cogentcore.org/core/enums"
)

var _FunsValues = []Funs{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32, 33, 34, 35, 36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 47, 48, 49, 50, 51, 52, 53, 54, 55, 56, 57, 58, 59, 60, 61, 62, 63, 64, 65}

// FunsN is the highest valid value for type Funs, plus one.
const FunsN Funs = 66

var _FunsValueMap = map[string]Funs{`Nil`: 0, `MoveUp`: 1, `MoveDown`: 2, `MoveRight`: 3, `MoveLeft`: 4, `PageUp`: 5, `PageDown`: 6, `Home`: 7, `End`: 8, `DocHome`: 9, `DocEnd`: 10, `WordRight`: 11, `WordLeft`: 12, `FocusNext`: 13, `FocusPrev`: 14, `Enter`: 15, `Accept`: 16, `CancelSelect`: 17, `SelectMode`: 18, `SelectAll`: 19, `Abort`: 20, `Copy`: 21, `Cut`: 22, `Paste`: 23, `PasteHist`: 24, `Backspace`: 25, `BackspaceWord`: 26, `Delete`: 27, `DeleteWord`: 28, `Kill`: 29, `Duplicate`: 30, `Transpose`: 31, `TransposeWord`: 32, `Undo`: 33, `Redo`: 34, `Insert`: 35, `InsertAfter`: 36, `ZoomOut`: 37, `ZoomIn`: 38, `Prefs`: 39, `Refresh`: 40, `Recenter`: 41, `Complete`: 42, `Lookup`: 43, `Search`: 44, `Find`: 45, `Replace`: 46, `Jump`: 47, `HistPrev`: 48, `HistNext`: 49, `Menu`: 50, `WinFocusNext`: 51, `WinClose`: 52, `WinSnapshot`: 53, `Inspector`: 54, `New`: 55, `NewAlt1`: 56, `NewAlt2`: 57, `Open`: 58, `OpenAlt1`: 59, `OpenAlt2`: 60, `Save`: 61, `SaveAs`: 62, `SaveAlt`: 63, `CloseAlt1`: 64, `CloseAlt2`: 65}

var _FunsDescMap = map[Funs]string{0: ``, 1: ``, 2: ``, 3: ``, 4: ``, 5: ``, 6: ``, 7: `PageRight PageLeft`, 8: ``, 9: ``, 10: ``, 11: ``, 12: ``, 13: ``, 14: ``, 15: ``, 16: ``, 17: ``, 18: ``, 19: ``, 20: ``, 21: `EditItem`, 22: ``, 23: ``, 24: ``, 25: ``, 26: ``, 27: ``, 28: ``, 29: ``, 30: ``, 31: ``, 32: ``, 33: ``, 34: ``, 35: ``, 36: ``, 37: ``, 38: ``, 39: ``, 40: ``, 41: ``, 42: ``, 43: ``, 44: ``, 45: ``, 46: ``, 47: ``, 48: ``, 49: ``, 50: ``, 51: ``, 52: ``, 53: ``, 54: ``, 55: `Below are menu specific functions -- use these as shortcuts for menu buttons allows uniqueness of mapping and easy customization of all key buttons`, 56: ``, 57: ``, 58: ``, 59: ``, 60: ``, 61: ``, 62: ``, 63: ``, 64: ``, 65: ``}

var _FunsMap = map[Funs]string{0: `Nil`, 1: `MoveUp`, 2: `MoveDown`, 3: `MoveRight`, 4: `MoveLeft`, 5: `PageUp`, 6: `PageDown`, 7: `Home`, 8: `End`, 9: `DocHome`, 10: `DocEnd`, 11: `WordRight`, 12: `WordLeft`, 13: `FocusNext`, 14: `FocusPrev`, 15: `Enter`, 16: `Accept`, 17: `CancelSelect`, 18: `SelectMode`, 19: `SelectAll`, 20: `Abort`, 21: `Copy`, 22: `Cut`, 23: `Paste`, 24: `PasteHist`, 25: `Backspace`, 26: `BackspaceWord`, 27: `Delete`, 28: `DeleteWord`, 29: `Kill`, 30: `Duplicate`, 31: `Transpose`, 32: `TransposeWord`, 33: `Undo`, 34: `Redo`, 35: `Insert`, 36: `InsertAfter`, 37: `ZoomOut`, 38: `ZoomIn`, 39: `Prefs`, 40: `Refresh`, 41: `Recenter`, 42: `Complete`, 43: `Lookup`, 44: `Search`, 45: `Find`, 46: `Replace`, 47: `Jump`, 48: `HistPrev`, 49: `HistNext`, 50: `Menu`, 51: `WinFocusNext`, 52: `WinClose`, 53: `WinSnapshot`, 54: `Inspector`, 55: `New`, 56: `NewAlt1`, 57: `NewAlt2`, 58: `Open`, 59: `OpenAlt1`, 60: `OpenAlt2`, 61: `Save`, 62: `SaveAs`, 63: `SaveAlt`, 64: `CloseAlt1`, 65: `CloseAlt2`}

// String returns the string representation of this Funs value.
func (i Funs) String() string { return enums.String(i, _FunsMap) }

// SetString sets the Funs value from its string representation,
// and returns an error if the string is invalid.
func (i *Funs) SetString(s string) error { return enums.SetString(i, s, _FunsValueMap, "Funs") }

// Int64 returns the Funs value as an int64.
func (i Funs) Int64() int64 { return int64(i) }

// SetInt64 sets the Funs value from an int64.
func (i *Funs) SetInt64(in int64) { *i = Funs(in) }

// Desc returns the description of the Funs value.
func (i Funs) Desc() string { return enums.Desc(i, _FunsDescMap) }

// FunsValues returns all possible values for the type Funs.
func FunsValues() []Funs { return _FunsValues }

// Values returns all possible values for the type Funs.
func (i Funs) Values() []enums.Enum { return enums.Values(_FunsValues) }

// MarshalText implements the [encoding.TextMarshaler] interface.
func (i Funs) MarshalText() ([]byte, error) { return []byte(i.String()), nil }

// UnmarshalText implements the [encoding.TextUnmarshaler] interface.
func (i *Funs) UnmarshalText(text []byte) error { return enums.UnmarshalText(i, text, "Funs") }
