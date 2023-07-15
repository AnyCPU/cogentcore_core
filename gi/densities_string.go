// Code generated by "stringer -type=Densities"; DO NOT EDIT.

package gi

import (
	"errors"
	"strconv"
)

var _ = errors.New("dummy error")

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[DensityCompact-0]
	_ = x[DensityMedium-1]
	_ = x[DensitySpread-2]
	_ = x[DensitiesN-3]
}

const _Densities_name = "DensityCompactDensityMediumDensitySpreadDensitiesN"

var _Densities_index = [...]uint8{0, 14, 27, 40, 50}

func (i Densities) String() string {
	if i < 0 || i >= Densities(len(_Densities_index)-1) {
		return "Densities(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _Densities_name[_Densities_index[i]:_Densities_index[i+1]]
}

func (i *Densities) FromString(s string) error {
	for j := 0; j < len(_Densities_index)-1; j++ {
		if s == _Densities_name[_Densities_index[j]:_Densities_index[j+1]] {
			*i = Densities(j)
			return nil
		}
	}
	return errors.New("String: " + s + " is not a valid option for type: Densities")
}
