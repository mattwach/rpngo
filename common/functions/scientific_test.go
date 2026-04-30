package functions

import (
	"mattwach/rpngo/common/rpn"
	"testing"
)

func TestPower(t *testing.T) {
	data := []rpn.UnitTestExecData{
		{
			Args:    []string{"**"},
			WantErr: rpn.ErrStackEmpty,
		},
		{
			Args:    []string{"2", "**"},
			WantErr: rpn.ErrStackEmpty,
		},
		{
			Args: []string{"2", "3", "**"},
			Want: []string{"8"},
		},
		{
			Args: []string{"-6", "2", "**"},
			Want: []string{"36"},
		},
		{
			Args: []string{"2d", "3", "**"},
			Want: []string{"8"},
		},
		{
			Args: []string{"2", "3d", "**"},
			Want: []string{"8"},
		},
		{
			Args: []string{"2d", "3d", "**"},
			Want: []string{"8d"},
		},
		{
			Args:    []string{"2", "true", "**"},
			WantErr: rpn.ErrExpectedANumber,
		},
		{
			Args:    []string{"true", "2", "**"},
			WantErr: rpn.ErrExpectedANumber,
		},
	}
	rpn.UnitTestExecAll(t, data, func(r *rpn.RPN) { RegisterAll(r) })
}

func TestSquareRoot(t *testing.T) {
	data := []rpn.UnitTestExecData{
		{
			Args:    []string{"sqrt"},
			WantErr: rpn.ErrStackEmpty,
		},
		{
			Args: []string{"4", "sqrt"},
			Want: []string{"2"},
		},
		{
			Args: []string{"4d", "sqrt"},
			Want: []string{"2d"},
		},
		{
			Args:    []string{"true", "sqrt"},
			WantErr: rpn.ErrExpectedANumber,
		},
	}
	rpn.UnitTestExecAll(t, data, func(r *rpn.RPN) { RegisterAll(r) })
}

func TestAbs(t *testing.T) {
	data := []rpn.UnitTestExecData{
		{
			Args:    []string{"abs"},
			WantErr: rpn.ErrStackEmpty,
		},
		{
			Args: []string{"4", "abs"},
			Want: []string{"4"},
		},
		{
			Args: []string{"-4", "abs"},
			Want: []string{"4"},
		},
		{
			Args: []string{"4d", "abs"},
			Want: []string{"4d"},
		},
		{
			Args: []string{"-4d", "abs"},
			Want: []string{"4d"},
		},
		{
			Args:    []string{"true", "abs"},
			WantErr: rpn.ErrExpectedANumber,
		},
	}
	rpn.UnitTestExecAll(t, data, func(r *rpn.RPN) { RegisterAll(r) })
}

func TestSquare(t *testing.T) {
	data := []rpn.UnitTestExecData{
		{
			Args:    []string{"sq"},
			WantErr: rpn.ErrStackEmpty,
		},
		{
			Args: []string{"4", "sq"},
			Want: []string{"16"},
		},
		{
			Args: []string{"4d", "sq"},
			Want: []string{"16d"},
		},
		{
			Args:    []string{"true", "sq"},
			WantErr: rpn.ErrExpectedANumber,
		},
	}
	rpn.UnitTestExecAll(t, data, func(r *rpn.RPN) { RegisterAll(r) })
}

func TestLog(t *testing.T) {
	data := []rpn.UnitTestExecData{
		{
			Args:    []string{"log"},
			WantErr: rpn.ErrStackEmpty,
		},
		{
			Args: []string{"4", "log", "3", "round"},
			Want: []string{"1.386"},
		},
		{
			Args: []string{"4d", "log", "3", "round"},
			Want: []string{"1.386"},
		},
		{
			Args:    []string{"true", "log"},
			WantErr: rpn.ErrExpectedANumber,
		},
	}
	rpn.UnitTestExecAll(t, data, func(r *rpn.RPN) { RegisterAll(r) })
}

func TestLog10(t *testing.T) {
	data := []rpn.UnitTestExecData{
		{
			Args:    []string{"log10"},
			WantErr: rpn.ErrStackEmpty,
		},
		{
			Args: []string{"4", "log10", "1000", "*", "int"},
			Want: []string{"602d"},
		},
		{
			Args: []string{"4d", "log10", "1000", "*", "int"},
			Want: []string{"602d"},
		},
		{
			Args:    []string{"true", "log"},
			WantErr: rpn.ErrExpectedANumber,
		},
	}
	rpn.UnitTestExecAll(t, data, func(r *rpn.RPN) { RegisterAll(r) })
}
