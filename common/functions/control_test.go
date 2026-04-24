package functions

import (
	"mattwach/rpngo/common/rpn"
	"testing"
)

func TestIf(t *testing.T) {
	data := []rpn.UnitTestExecData{
		{
			Args:    []string{"if"},
			WantErr: rpn.ErrStackEmpty,
		},
		{
			Args:    []string{"1", "if"},
			WantErr: rpn.ErrStackEmpty,
		},
		{
			Args:    []string{"1", "2", "if"},
			WantErr: rpn.ErrExpectedABoolean,
		},
		{
			Args: []string{"false", "1", "if"},
		},
		{
			Args: []string{"true", "1", "if"},
			Want: []string{"1"},
		},
	}
	rpn.UnitTestExecAll(t, data, func(r *rpn.RPN) { RegisterAll(r) })
}

func TestIfElse(t *testing.T) {
	data := []rpn.UnitTestExecData{
		{
			Args:    []string{"ifelse"},
			WantErr: rpn.ErrStackEmpty,
		},
		{
			Args:    []string{"1", "ifelse"},
			WantErr: rpn.ErrStackEmpty,
		},
		{
			Args:    []string{"1", "2", "ifelse"},
			WantErr: rpn.ErrStackEmpty,
		},
		{
			Args:    []string{"1", "2", "3", "ifelse"},
			WantErr: rpn.ErrExpectedABoolean,
		},
		{
			Args: []string{"false", "1", "2", "ifelse"},
			Want: []string{"2"},
		},
		{
			Args: []string{"true", "1", "2", "ifelse"},
			Want: []string{"1"},
		},
	}
	rpn.UnitTestExecAll(t, data, func(r *rpn.RPN) { RegisterAll(r) })
}

func TestFor(t *testing.T) {
	data := []rpn.UnitTestExecData{
		{
			Args:    []string{"for"},
			WantErr: rpn.ErrStackEmpty,
		},
		{
			Args:    []string{"'", "for"},
			WantErr: rpn.ErrSyntax,
		},
		{
			Args:    []string{"1", "for"},
			WantErr: rpn.ErrExpectedAString,
		},
		{
			Args: []string{"'false'", "for"},
		},
		{
			Args: []string{"0d", "x=", "'$x 1d + x= 10d $x !='", "for", "$x"},
			Want: []string{"10d"},
		},
	}
	rpn.UnitTestExecAll(t, data, func(r *rpn.RPN) { RegisterAll(r) })
}
