package rpn

import "testing"

func TestAngleUnits(t *testing.T) {
	data := []UnitTestExecData{
		{
			Args: []string{"getangle"},
			Want: []string{"'rad'"},
		},
		{
			Args: []string{"'rad'", "setangle", "getangle"},
			Want: []string{"'rad'"},
		},
		{
			Args: []string{"'deg'", "setangle", "getangle"},
			Want: []string{"'deg'"},
		},
		{
			Args: []string{"'grad'", "setangle", "getangle"},
			Want: []string{"'grad'"},
		},
		{
			Args:    []string{"'foo'", "setangle"},
			WantErr: ErrChooseDegRadOGrad,
		},
		{
			Args:    []string{"5", "setangle"},
			WantErr: ErrExpectedAString,
		},
	}
	UnitTestExecAll(t, data, func(r *RPN) {})
}
