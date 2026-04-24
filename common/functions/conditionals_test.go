package functions

import (
	"mattwach/rpngo/common/rpn"
	"testing"
)

func TestConditionals(t *testing.T) {
	data := []rpn.UnitTestExecData{
		{
			Args: []string{"1", "2", ">"},
			Want: []string{"false"},
		},
		{
			Args: []string{"2", "2", ">"},
			Want: []string{"false"},
		},
		{
			Args: []string{"3", "2", ">"},
			Want: []string{"true"},
		},
		{
			Args: []string{"'foo'", "2", ">"},
			Want: []string{"true"},
		},
		{
			Args: []string{"2", "'foo'", ">"},
			Want: []string{"false"},
		},

		{
			Args: []string{"1", "2", ">="},
			Want: []string{"false"},
		},
		{
			Args: []string{"2", "2", ">="},
			Want: []string{"true"},
		},
		{
			Args: []string{"3", "2", ">="},
			Want: []string{"true"},
		},
		{
			Args: []string{"'foo'", "2", ">="},
			Want: []string{"true"},
		},
		{
			Args: []string{"2", "'foo'", ">="},
			Want: []string{"false"},
		},

		{
			Args: []string{"1", "2", "<"},
			Want: []string{"true"},
		},
		{
			Args: []string{"2", "2", "<"},
			Want: []string{"false"},
		},
		{
			Args: []string{"3", "2", "<"},
			Want: []string{"false"},
		},
		{
			Args: []string{"'foo'", "2", "<"},
			Want: []string{"false"},
		},
		{
			Args: []string{"2", "'foo'", "<"},
			Want: []string{"true"},
		},

		{
			Args: []string{"1", "2", "<="},
			Want: []string{"true"},
		},
		{
			Args: []string{"2", "2", "<="},
			Want: []string{"true"},
		},
		{
			Args: []string{"3", "2", "<="},
			Want: []string{"false"},
		},
		{
			Args: []string{"'foo'", "2", "<="},
			Want: []string{"false"},
		},
		{
			Args: []string{"2", "'foo'", "<="},
			Want: []string{"true"},
		},

		{
			Args: []string{"1", "2", "="},
			Want: []string{"false"},
		},
		{
			Args: []string{"1", "2d", "="},
			Want: []string{"false"},
		},
		{
			Args: []string{"1d", "2", "="},
			Want: []string{"false"},
		},
		{
			Args: []string{"1d", "2d", "="},
			Want: []string{"false"},
		},
		{
			Args: []string{"1", "'foo'", "="},
			Want: []string{"false"},
		},
		{
			Args: []string{"'foo'", "2", "="},
			Want: []string{"false"},
		},
		{
			Args: []string{"'foo'", "'bar'", "="},
			Want: []string{"false"},
		},
		{
			Args: []string{"1", "false", "="},
			Want: []string{"false"},
		},
		{
			Args: []string{"false", "2", "="},
			Want: []string{"false"},
		},
		{
			Args: []string{"false", "true", "="},
			Want: []string{"false"},
		},
		{
			Args: []string{"1", "1", "="},
			Want: []string{"true"},
		},
		{
			Args: []string{"1", "1d", "="},
			Want: []string{"true"},
		},
		{
			Args: []string{"1d", "1", "="},
			Want: []string{"true"},
		},
		{
			Args: []string{"1d", "1d", "="},
			Want: []string{"true"},
		},
		{
			Args: []string{"'foo'", "'foo'", "="},
			Want: []string{"true"},
		},
		{
			Args: []string{"false", "false", "="},
			Want: []string{"true"},
		},
		{
			Args: []string{"true", "true", "="},
			Want: []string{"true"},
		},

		{
			Args: []string{"1", "2", "!="},
			Want: []string{"true"},
		},
		{
			Args: []string{"1", "2d", "!="},
			Want: []string{"true"},
		},
		{
			Args: []string{"1d", "2", "!="},
			Want: []string{"true"},
		},
		{
			Args: []string{"1d", "2d", "!="},
			Want: []string{"true"},
		},
		{
			Args: []string{"1", "'foo'", "!="},
			Want: []string{"true"},
		},
		{
			Args: []string{"'foo'", "2", "!="},
			Want: []string{"true"},
		},
		{
			Args: []string{"'foo'", "'bar'", "!="},
			Want: []string{"true"},
		},
		{
			Args: []string{"1", "false", "!="},
			Want: []string{"true"},
		},
		{
			Args: []string{"false", "2", "!="},
			Want: []string{"true"},
		},
		{
			Args: []string{"false", "true", "!="},
			Want: []string{"true"},
		},
		{
			Args: []string{"1", "1", "!="},
			Want: []string{"false"},
		},
		{
			Args: []string{"1", "1d", "!="},
			Want: []string{"false"},
		},
		{
			Args: []string{"1d", "1", "!="},
			Want: []string{"false"},
		},
		{
			Args: []string{"1d", "1d", "!="},
			Want: []string{"false"},
		},
		{
			Args: []string{"'foo'", "'foo'", "!="},
			Want: []string{"false"},
		},
		{
			Args: []string{"false", "false", "!="},
			Want: []string{"false"},
		},
		{
			Args: []string{"true", "true", "!="},
			Want: []string{"false"},
		},
		{
			Args:    []string{"min"},
			WantErr: rpn.ErrStackEmpty,
		},
		{
			Args:    []string{"0", "min"},
			WantErr: rpn.ErrStackEmpty,
		},
		{
			Args: []string{"0x", "0d", "min"},
			Want: []string{"0x"},
		},
		{
			Args: []string{"0x", "1d", "min"},
			Want: []string{"0x"},
		},
		{
			Args: []string{"1x", "0d", "min"},
			Want: []string{"0d"},
		},
		{
			Args:    []string{"max"},
			WantErr: rpn.ErrStackEmpty,
		},
		{
			Args:    []string{"0", "max"},
			WantErr: rpn.ErrStackEmpty,
		},
		{
			Args: []string{"0x", "0d", "max"},
			Want: []string{"0x"},
		},
		{
			Args: []string{"0x", "1d", "max"},
			Want: []string{"1d"},
		},
		{
			Args: []string{"1x", "0d", "max"},
			Want: []string{"1x"},
		},
	}
	rpn.UnitTestExecAll(t, data, func(r *rpn.RPN) { RegisterAll(r) })
}
