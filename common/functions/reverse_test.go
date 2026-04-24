package functions

import (
	"mattwach/rpngo/common/rpn"
	"testing"
)

func TestReverse(t *testing.T) {
	data := []rpn.UnitTestExecData{
		{
			Args: []string{"reverse"},
		},
		{
			Args: []string{"1", "reverse"},
			Want: []string{"1"},
		},
		{
			Args: []string{"1", "2", "reverse"},
			Want: []string{"2", "1"},
		},
		{
			Args: []string{"1", "2", "3", "reverse"},
			Want: []string{"3", "2", "1"},
		},
		{
			Args: []string{"1", "2", "3", "4", "reverse"},
			Want: []string{"4", "3", "2", "1"},
		},
		{
			Args: []string{"1", "2", "3", "4", "5", "reverse"},
			Want: []string{"5", "4", "3", "2", "1"},
		},
	}
	rpn.UnitTestExecAll(t, data, func(r *rpn.RPN) { RegisterAll(r) })
}
