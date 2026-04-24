package functions

import (
	"mattwach/rpngo/common/rpn"
	"testing"
)

func TestBitwiseOps(t *testing.T) {
	data := []rpn.UnitTestExecData{
		{
			Args: []string{"0011b", "1010b", "&"},
			Want: []string{"10b"},
		},
		{
			Args:    []string{"'foo'", "1b", "&"},
			WantErr: rpn.ErrExpectedANumber,
		},
		{
			Args:    []string{"1b", "'foo'", "&"},
			WantErr: rpn.ErrExpectedANumber,
		},
		{
			Args: []string{"0011b", "1010b", "|"},
			Want: []string{"1011b"},
		},
		{
			Args:    []string{"'foo'", "1b", "|"},
			WantErr: rpn.ErrExpectedANumber,
		},
		{
			Args:    []string{"1b", "'foo'", "|"},
			WantErr: rpn.ErrExpectedANumber,
		},
		{
			Args: []string{"0011b", "1010b", "^"},
			Want: []string{"1001b"},
		},
		{
			Args:    []string{"'foo'", "1b", "^"},
			WantErr: rpn.ErrExpectedANumber,
		},
		{
			Args:    []string{"1b", "'foo'", "^"},
			WantErr: rpn.ErrExpectedANumber,
		},
		{
			Args: []string{"1000b", "2", ">>"},
			Want: []string{"10b"},
		},
		{
			Args:    []string{"'foo'", "1", ">>"},
			WantErr: rpn.ErrExpectedANumber,
		},
		{
			Args:    []string{"1b", "'foo'", ">>"},
			WantErr: rpn.ErrExpectedANumber,
		},
		{
			Args: []string{"10b", "2", "<<"},
			Want: []string{"1000b"},
		},
		{
			Args:    []string{"'foo'", "1", "<<"},
			WantErr: rpn.ErrExpectedANumber,
		},
		{
			Args:    []string{"1b", "'foo'", "<<"},
			WantErr: rpn.ErrExpectedANumber,
		},
	}
	rpn.UnitTestExecAll(t, data, func(r *rpn.RPN) { RegisterAll(r) })
}
