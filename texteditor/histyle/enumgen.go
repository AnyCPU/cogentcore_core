// Code generated by "core generate -add-types"; DO NOT EDIT.

package histyle

import (
	"errors"
	"log"
	"strconv"

	"cogentcore.org/core/enums"
)

var _TrileanValues = []Trilean{0, 1, 2}

// TrileanN is the highest valid value
// for type Trilean, plus one.
const TrileanN Trilean = 3

// An "invalid array index" compiler error signifies that the constant values have changed.
// Re-run the enumgen command to generate them again.
func _TrileanNoOp() {
	var x [1]struct{}
	_ = x[Pass-(0)]
	_ = x[Yes-(1)]
	_ = x[No-(2)]
}

var _TrileanNameToValueMap = map[string]Trilean{
	`Pass`: 0,
	`Yes`:  1,
	`No`:   2,
}

var _TrileanDescMap = map[Trilean]string{
	0: ``,
	1: ``,
	2: ``,
}

var _TrileanMap = map[Trilean]string{
	0: `Pass`,
	1: `Yes`,
	2: `No`,
}

// String returns the string representation
// of this Trilean value.
func (i Trilean) String() string {
	if str, ok := _TrileanMap[i]; ok {
		return str
	}
	return strconv.FormatInt(int64(i), 10)
}

// SetString sets the Trilean value from its
// string representation, and returns an
// error if the string is invalid.
func (i *Trilean) SetString(s string) error {
	if val, ok := _TrileanNameToValueMap[s]; ok {
		*i = val
		return nil
	}
	return errors.New(s + " is not a valid value for type Trilean")
}

// Int64 returns the Trilean value as an int64.
func (i Trilean) Int64() int64 {
	return int64(i)
}

// SetInt64 sets the Trilean value from an int64.
func (i *Trilean) SetInt64(in int64) {
	*i = Trilean(in)
}

// Desc returns the description of the Trilean value.
func (i Trilean) Desc() string {
	if str, ok := _TrileanDescMap[i]; ok {
		return str
	}
	return i.String()
}

// TrileanValues returns all possible values
// for the type Trilean.
func TrileanValues() []Trilean {
	return _TrileanValues
}

// Values returns all possible values
// for the type Trilean.
func (i Trilean) Values() []enums.Enum {
	res := make([]enums.Enum, len(_TrileanValues))
	for i, d := range _TrileanValues {
		res[i] = d
	}
	return res
}

// IsValid returns whether the value is a
// valid option for type Trilean.
func (i Trilean) IsValid() bool {
	_, ok := _TrileanMap[i]
	return ok
}

// MarshalText implements the [encoding.TextMarshaler] interface.
func (i Trilean) MarshalText() ([]byte, error) {
	return []byte(i.String()), nil
}

// UnmarshalText implements the [encoding.TextUnmarshaler] interface.
func (i *Trilean) UnmarshalText(text []byte) error {
	if err := i.SetString(string(text)); err != nil {
		log.Println("Trilean.UnmarshalText:", err)
	}
	return nil
}
