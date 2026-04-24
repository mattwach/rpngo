package functions

import (
	"mattwach/rpngo/common/rpn"
	"testing"
)

func TestFields(t *testing.T) {
	data := []rpn.UnitTestExecData{
		{
			Args:    []string{"fields"},
			WantErr: rpn.ErrStackEmpty,
		},
		{
			Args: []string{"''", "fields"},
		},
		{
			Args: []string{"'  '", "fields"},
		},
		{
			Args: []string{"'hello'", "fields"},
			Want: []string{"'hello'"},
		},
		{
			Args: []string{"'hello world'", "fields"},
			Want: []string{"'hello'", "'world'"},
		},
		{
			Args:    []string{"5", "fields"},
			WantErr: rpn.ErrExpectedAString,
		},
	}
	rpn.UnitTestExecAll(t, data, func(r *rpn.RPN) { RegisterAll(r) })
}
