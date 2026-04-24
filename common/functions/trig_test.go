package functions

import (
	"mattwach/rpngo/common/rpn"
	"testing"
)

func TestASin(t *testing.T) {
	data := []rpn.UnitTestExecData{
		{
			Args:    []string{"asin"},
			WantErr: rpn.ErrStackEmpty,
		},
		{
			Args:    []string{"true", "asin"},
			WantErr: rpn.ErrExpectedANumber,
		},
		{
			Args: []string{"1", "asin", "3", "round"},
			Want: []string{"1.571"},
		},
		{
			Args: []string{"deg", "1", "asin", "3", "round"},
			Want: []string{"90"},
		},
		{
			Args: []string{"grad", "1", "asin", "3", "round"},
			Want: []string{"100"},
		},
		{
			Args: []string{"2+i", "asin", "3", "round"},
			Want: []string{"1.063+1.469i"},
		},
	}
	rpn.UnitTestExecAll(t, data, func(r *rpn.RPN) { RegisterAll(r) })
}

func TestSin(t *testing.T) {
	data := []rpn.UnitTestExecData{
		{
			Args:    []string{"sin"},
			WantErr: rpn.ErrStackEmpty,
		},
		{
			Args:    []string{"true", "sin"},
			WantErr: rpn.ErrExpectedANumber,
		},
		{
			Args: []string{"1", "sin", "3", "round"},
			Want: []string{"0.841"},
		},
		{
			Args: []string{"deg", "90", "sin", "3", "round"},
			Want: []string{"1"},
		},
		{
			Args: []string{"grad", "100", "sin", "3", "round"},
			Want: []string{"1"},
		},
		{
			Args: []string{"1+i", "sin", "3", "round"},
			Want: []string{"1.298+0.635i"},
		},
	}
	rpn.UnitTestExecAll(t, data, func(r *rpn.RPN) { RegisterAll(r) })
}

func TestACos(t *testing.T) {
	data := []rpn.UnitTestExecData{
		{
			Args:    []string{"acos"},
			WantErr: rpn.ErrStackEmpty,
		},
		{
			Args:    []string{"true", "acos"},
			WantErr: rpn.ErrExpectedANumber,
		},
		{
			Args: []string{"1", "acos", "3", "round"},
			Want: []string{"0"},
		},
		{
			Args: []string{"deg", "0", "acos", "3", "round"},
			Want: []string{"90"},
		},
		{
			Args: []string{"grad", "0", "acos", "3", "round"},
			Want: []string{"100"},
		},
		{
			Args: []string{"2+i", "acos", "3", "round"},
			Want: []string{"0.507-1.469i"},
		},
	}
	rpn.UnitTestExecAll(t, data, func(r *rpn.RPN) { RegisterAll(r) })
}

func TestCos(t *testing.T) {
	data := []rpn.UnitTestExecData{
		{
			Args:    []string{"cos"},
			WantErr: rpn.ErrStackEmpty,
		},
		{
			Args:    []string{"true", "cos"},
			WantErr: rpn.ErrExpectedANumber,
		},
		{
			Args: []string{"1", "cos", "3", "round"},
			Want: []string{"0.54"},
		},
		{
			Args: []string{"deg", "89.999", "cos", "3", "round"},
			Want: []string{"0"},
		},
		{
			Args: []string{"grad", "99.999", "cos", "3", "round"},
			Want: []string{"0"},
		},
		{
			Args: []string{"21+i", "cos", "3", "round"},
			Want: []string{"-0.845-0.983i"},
		},
	}
	rpn.UnitTestExecAll(t, data, func(r *rpn.RPN) { RegisterAll(r) })
}

func TestTan(t *testing.T) {
	data := []rpn.UnitTestExecData{
		{
			Args:    []string{"tan"},
			WantErr: rpn.ErrStackEmpty,
		},
		{
			Args:    []string{"true", "tan"},
			WantErr: rpn.ErrExpectedANumber,
		},
		{
			Args: []string{"1", "tan", "3", "round"},
			Want: []string{"1.557"},
		},
		{
			Args: []string{"deg", "45", "tan", "3", "round"},
			Want: []string{"1"},
		},
		{
			Args: []string{"grad", "50", "tan", "3", "round"},
			Want: []string{"1"},
		},
		{
			Args: []string{"1.5+2.1i", "tan", "3", "round"},
			Want: []string{"0.004+1.03i"},
		},
	}
	rpn.UnitTestExecAll(t, data, func(r *rpn.RPN) { RegisterAll(r) })
}

func TestATan(t *testing.T) {
	data := []rpn.UnitTestExecData{
		{
			Args:    []string{"atan"},
			WantErr: rpn.ErrStackEmpty,
		},
		{
			Args:    []string{"true", "atan"},
			WantErr: rpn.ErrExpectedANumber,
		},
		{
			Args: []string{"2", "atan", "3", "round"},
			Want: []string{"1.107"},
		},
		{
			Args: []string{"deg", "1", "atan", "3", "round"},
			Want: []string{"45"},
		},
		{
			Args: []string{"grad", "1", "atan", "3", "round"},
			Want: []string{"50"},
		},
		{
			Args: []string{"1+2i", "atan", "3", "round"},
			Want: []string{"1.339+0.402i"},
		},
	}
	rpn.UnitTestExecAll(t, data, func(r *rpn.RPN) { RegisterAll(r) })
}
