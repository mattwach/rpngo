package functions

import (
	"mattwach/rpngo/common/rpn"
	"testing"
)

func TestFilter(t *testing.T) {
	data := []rpn.UnitTestExecData{
		{
			Args:    []string{"filter"},
			WantErr: rpn.ErrStackEmpty,
		},
		{
			Args: []string{"{sq}", "filter"},
		},
		{
			Args: []string{"1", "2", "3", "{sq}", "filter"},
			Want: []string{"1", "4", "9"},
		},
		{
			Args: []string{"10", "100", "30", "200", "{$0 100 >= {0/} if}", "filter"},
			Want: []string{"10", "30"},
		},
		{
			Args: []string{"1", "2", "3", "0", "sum=", "{$sum + sum=}", "filter", "$sum"},
			Want: []string{"6"},
		},
		{
			Args: []string{"10", "5", "7", "12", "$0", "min=", "{$0 $min < {min=} {0/} ifelse}", "filter", "$min"},
			Want: []string{"5"},
		},
		{
			Args: []string{"1", "2", "3", "{0/}", "filter"},
		},
		{
			Args:    []string{"1", "2", "3", "{0/ 0/}", "filter"},
			Want:    []string{"1"},
			WantErr: rpn.ErrNotEnoughStackFrames,
		},
	}
	rpn.UnitTestExecAll(t, data, func(r *rpn.RPN) { RegisterAll(r) })
}
