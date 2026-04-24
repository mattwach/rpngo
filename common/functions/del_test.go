package functions

import (
	"mattwach/rpngo/common/rpn"
	"testing"
)

func TestDel(t *testing.T) {
	data := []rpn.UnitTestExecData{
		{
			Args:    []string{"del"},
			WantErr: rpn.ErrStackEmpty,
		},
		{
			Args:    []string{"-1", "del"},
			WantErr: rpn.ErrIllegalValue,
		},
		{
			Args: []string{"0", "del"},
		},
		{
			Args: []string{"1", "del"},
		},
		{
			Args:    []string{"1", "2", "3", "-1", "del"},
			Want:    []string{"1", "2", "3"},
			WantErr: rpn.ErrIllegalValue,
		},
		{
			Args: []string{"1", "2", "3", "0", "del"},
			Want: []string{"1", "2", "3"},
		},
		{
			Args: []string{"1", "2", "3", "1", "del"},
			Want: []string{"1", "2"},
		},
		{
			Args: []string{"1", "2", "3", "2", "del"},
			Want: []string{"1"},
		},
		{
			Args: []string{"1", "2", "3", "3", "del"},
		},
		{
			Args: []string{"1", "2", "3", "4", "del"},
		},
	}
	rpn.UnitTestExecAll(t, data, func(r *rpn.RPN) { RegisterAll(r) })
}
