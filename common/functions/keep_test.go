package functions

import (
	"mattwach/rpngo/common/rpn"
	"testing"
)

func TestKeep(t *testing.T) {
	data := []rpn.UnitTestExecData{
		{
			Args:    []string{"keep"},
			WantErr: rpn.ErrStackEmpty,
		},
		{
			Args:    []string{"-1", "keep"},
			WantErr: rpn.ErrIllegalValue,
		},
		{
			Args: []string{"0", "keep"},
		},
		{
			Args: []string{"1", "keep"},
		},
		{
			Args:    []string{"1", "2", "3", "-1", "keep"},
			Want:    []string{"1", "2", "3"},
			WantErr: rpn.ErrIllegalValue,
		},
		{
			Args: []string{"1", "2", "3", "0", "keep"},
		},
		{
			Args: []string{"1", "2", "3", "1", "keep"},
			Want: []string{"3"},
		},
		{
			Args: []string{"1", "2", "3", "2", "keep"},
			Want: []string{"2", "3"},
		},
		{
			Args: []string{"1", "2", "3", "3", "keep"},
			Want: []string{"1", "2", "3"},
		},
		{
			Args: []string{"1", "2", "3", "4", "keep"},
			Want: []string{"1", "2", "3"},
		},
	}
	rpn.UnitTestExecAll(t, data, func(r *rpn.RPN) { RegisterAll(r) })
}
