package functions

import (
	"mattwach/rpngo/common/rpn"
	"testing"
)

func TestSort(t *testing.T) {
	data := []rpn.UnitTestExecData{
		{
			Args: []string{"sort"},
		},
		{
			Args: []string{"1", "sort"},
			Want: []string{"1"},
		},
		{
			Args: []string{"1", "2", "sort"},
			Want: []string{"1", "2"},
		},
		{
			Args: []string{"2", "1", "sort"},
			Want: []string{"1", "2"},
		},
		{
			Args: []string{"10", "5d", "'foo'", "true", "100", "6+i", "sort"},
			Want: []string{"true", "5d", "6+i", "10", "100", "'foo'"},
		},
	}
	rpn.UnitTestExecAll(t, data, func(r *rpn.RPN) { RegisterAll(r) })
}
