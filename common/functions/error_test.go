package functions

import (
	"mattwach/rpngo/common/rpn"
	"testing"
)

func TestTryAndError(t *testing.T) {
	data := []rpn.UnitTestExecData{
		{
			Args:    []string{"error"},
			WantErr: rpn.ErrStackEmpty,
		},
		{
			Args:    []string{"try"},
			WantErr: rpn.ErrStackEmpty,
		},
		{
			Args:    []string{"5", "try"},
			WantErr: rpn.ErrStackEmpty,
		},
		{
			Args: []string{"{2 3 +}", "{'foo'}", "try"},
			Want: []string{"5"},
		},
		{
			Args: []string{"{2 0 /}", "{'foo'}", "try"},
			Want: []string{"'2 0 ->/<-: divide by zero'", "'foo'"},
		},
		{
			Args: []string{"{'foo' error}", "{}", "try"},
			Want: []string{"''foo' ->error<-: foo'"},
		},
	}
	rpn.UnitTestExecAll(t, data, func(r *rpn.RPN) { RegisterAll(r) })
}
