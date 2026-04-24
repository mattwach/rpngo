package functions

import (
	"mattwach/rpngo/common/rpn"
	"testing"
)

func TestInt(t *testing.T) {
	data := []rpn.UnitTestExecData{
		{
			Args: []string{"-1", "int"},
			Want: []string{"-1d"},
		},
		{
			Args: []string{"0", "int"},
			Want: []string{"0d"},
		},
		{
			Args: []string{"1", "int"},
			Want: []string{"1d"},
		},
		{
			Args:    []string{"3+4i", "int"},
			WantErr: rpn.ErrComplexNumberNotSupported,
		},
		{
			Args: []string{"3.14", "int"},
			Want: []string{"3d"},
		},
		{
			Args: []string{"-3.14", "int"},
			Want: []string{"-3d"},
		},
		{
			Args: []string{"-2d", "int"},
			Want: []string{"-2d"},
		},
		{
			Args: []string{"0d", "int"},
			Want: []string{"0d"},
		},
		{
			Args: []string{"2d", "int"},
			Want: []string{"2d"},
		},
		{
			Args: []string{"-10o", "int"},
			Want: []string{"-8d"},
		},
		{
			Args: []string{"0o", "int"},
			Want: []string{"0d"},
		},
		{
			Args: []string{"10o", "int"},
			Want: []string{"8d"},
		},
		{
			Args: []string{"-10b", "int"},
			Want: []string{"-2d"},
		},
		{
			Args: []string{"0b", "int"},
			Want: []string{"0d"},
		},
		{
			Args: []string{"10b", "int"},
			Want: []string{"2d"},
		},
		{
			Args: []string{"-fx", "int"},
			Want: []string{"-15d"},
		},
		{
			Args: []string{"0x", "int"},
			Want: []string{"0d"},
		},
		{
			Args: []string{"fx", "int"},
			Want: []string{"15d"},
		},
		{
			Args: []string{"true", "int"},
			Want: []string{"1d"},
		},
		{
			Args: []string{"false", "int"},
			Want: []string{"0d"},
		},
		{
			Args: []string{"'5'", "int"},
			Want: []string{"5d"},
		},
		{
			Args: []string{"'-5'", "int"},
			Want: []string{"-5d"},
		},
		{
			Args:    []string{"'foo'", "int"},
			WantErr: rpn.ErrSyntax,
		},
	}
	rpn.UnitTestExecAll(t, data, func(r *rpn.RPN) { RegisterAll(r) })
}

func TestBin(t *testing.T) {
	data := []rpn.UnitTestExecData{
		{
			Args: []string{"-1", "bin"},
			Want: []string{"-1b"},
		},
		{
			Args: []string{"0", "bin"},
			Want: []string{"0b"},
		},
		{
			Args: []string{"1", "bin"},
			Want: []string{"1b"},
		},
		{
			Args:    []string{"3+4i", "bin"},
			WantErr: rpn.ErrComplexNumberNotSupported,
		},
		{
			Args: []string{"3.14", "bin"},
			Want: []string{"11b"},
		},
		{
			Args: []string{"-3.14", "bin"},
			Want: []string{"-11b"},
		},
		{
			Args: []string{"-2d", "bin"},
			Want: []string{"-10b"},
		},
		{
			Args: []string{"0d", "bin"},
			Want: []string{"0b"},
		},
		{
			Args: []string{"2d", "bin"},
			Want: []string{"10b"},
		},
		{
			Args: []string{"-10o", "bin"},
			Want: []string{"-1000b"},
		},
		{
			Args: []string{"0o", "bin"},
			Want: []string{"0b"},
		},
		{
			Args: []string{"10o", "bin"},
			Want: []string{"1000b"},
		},
		{
			Args: []string{"-10b", "bin"},
			Want: []string{"-10b"},
		},
		{
			Args: []string{"0b", "bin"},
			Want: []string{"0b"},
		},
		{
			Args: []string{"10b", "bin"},
			Want: []string{"10b"},
		},
		{
			Args: []string{"-fx", "bin"},
			Want: []string{"-1111b"},
		},
		{
			Args: []string{"0x", "bin"},
			Want: []string{"0b"},
		},
		{
			Args: []string{"fx", "bin"},
			Want: []string{"1111b"},
		},
		{
			Args: []string{"true", "bin"},
			Want: []string{"1b"},
		},
		{
			Args: []string{"false", "bin"},
			Want: []string{"0b"},
		},
		{
			Args: []string{"'5'", "bin"},
			Want: []string{"101b"},
		},
		{
			Args: []string{"'-5'", "bin"},
			Want: []string{"-101b"},
		},
		{
			Args:    []string{"'foo'", "bin"},
			WantErr: rpn.ErrSyntax,
		},
	}
	rpn.UnitTestExecAll(t, data, func(r *rpn.RPN) { RegisterAll(r) })
}

func TestOctal(t *testing.T) {
	data := []rpn.UnitTestExecData{
		{
			Args: []string{"-1", "oct"},
			Want: []string{"-1o"},
		},
		{
			Args: []string{"0", "oct"},
			Want: []string{"0o"},
		},
		{
			Args: []string{"1", "oct"},
			Want: []string{"1o"},
		},
		{
			Args:    []string{"3+4i", "oct"},
			WantErr: rpn.ErrComplexNumberNotSupported,
		},
		{
			Args: []string{"3.14", "oct"},
			Want: []string{"3o"},
		},
		{
			Args: []string{"-3.14", "oct"},
			Want: []string{"-3o"},
		},
		{
			Args: []string{"-2d", "oct"},
			Want: []string{"-2o"},
		},
		{
			Args: []string{"0d", "oct"},
			Want: []string{"0o"},
		},
		{
			Args: []string{"2d", "oct"},
			Want: []string{"2o"},
		},
		{
			Args: []string{"-10o", "oct"},
			Want: []string{"-10o"},
		},
		{
			Args: []string{"0o", "oct"},
			Want: []string{"0o"},
		},
		{
			Args: []string{"10o", "oct"},
			Want: []string{"10o"},
		},
		{
			Args: []string{"-10b", "oct"},
			Want: []string{"-2o"},
		},
		{
			Args: []string{"0b", "oct"},
			Want: []string{"0o"},
		},
		{
			Args: []string{"10b", "oct"},
			Want: []string{"2o"},
		},
		{
			Args: []string{"-fx", "oct"},
			Want: []string{"-17o"},
		},
		{
			Args: []string{"0x", "oct"},
			Want: []string{"0o"},
		},
		{
			Args: []string{"fx", "oct"},
			Want: []string{"17o"},
		},
		{
			Args: []string{"true", "oct"},
			Want: []string{"1o"},
		},
		{
			Args: []string{"false", "oct"},
			Want: []string{"0o"},
		},
		{
			Args: []string{"'5'", "oct"},
			Want: []string{"5o"},
		},
		{
			Args: []string{"'-5'", "oct"},
			Want: []string{"-5o"},
		},
		{
			Args:    []string{"'foo'", "oct"},
			WantErr: rpn.ErrSyntax,
		},
	}
	rpn.UnitTestExecAll(t, data, func(r *rpn.RPN) { RegisterAll(r) })
}

func TestHex(t *testing.T) {
	data := []rpn.UnitTestExecData{
		{
			Args: []string{"-1", "hex"},
			Want: []string{"-1x"},
		},
		{
			Args: []string{"0", "hex"},
			Want: []string{"0x"},
		},
		{
			Args: []string{"1", "hex"},
			Want: []string{"1x"},
		},
		{
			Args:    []string{"3+4i", "hex"},
			WantErr: rpn.ErrComplexNumberNotSupported,
		},
		{
			Args: []string{"3.14", "hex"},
			Want: []string{"3x"},
		},
		{
			Args: []string{"-3.14", "hex"},
			Want: []string{"-3x"},
		},
		{
			Args: []string{"-2d", "hex"},
			Want: []string{"-2x"},
		},
		{
			Args: []string{"0d", "hex"},
			Want: []string{"0x"},
		},
		{
			Args: []string{"2d", "hex"},
			Want: []string{"2x"},
		},
		{
			Args: []string{"-10o", "hex"},
			Want: []string{"-8x"},
		},
		{
			Args: []string{"0o", "hex"},
			Want: []string{"0x"},
		},
		{
			Args: []string{"10o", "hex"},
			Want: []string{"8x"},
		},
		{
			Args: []string{"-10b", "hex"},
			Want: []string{"-2x"},
		},
		{
			Args: []string{"0b", "hex"},
			Want: []string{"0x"},
		},
		{
			Args: []string{"10b", "hex"},
			Want: []string{"2x"},
		},
		{
			Args: []string{"-fx", "hex"},
			Want: []string{"-fx"},
		},
		{
			Args: []string{"0x", "hex"},
			Want: []string{"0x"},
		},
		{
			Args: []string{"fx", "hex"},
			Want: []string{"fx"},
		},
		{
			Args: []string{"true", "hex"},
			Want: []string{"1x"},
		},
		{
			Args: []string{"false", "hex"},
			Want: []string{"0x"},
		},
		{
			Args: []string{"'5'", "hex"},
			Want: []string{"5x"},
		},
		{
			Args: []string{"'-5'", "hex"},
			Want: []string{"-5x"},
		},
		{
			Args:    []string{"'foo'", "hex"},
			WantErr: rpn.ErrSyntax,
		},
	}
	rpn.UnitTestExecAll(t, data, func(r *rpn.RPN) { RegisterAll(r) })
}

func TestStr(t *testing.T) {
	data := []rpn.UnitTestExecData{
		{
			Args: []string{"-1", "str"},
			Want: []string{"{-1}"},
		},
		{
			Args: []string{"0", "str"},
			Want: []string{"{0}"},
		},
		{
			Args: []string{"1", "str"},
			Want: []string{"{1}"},
		},
		{
			Args: []string{"3+4i", "str"},
			Want: []string{"{3+4i}"},
		},
		{
			Args: []string{"3.14", "str"},
			Want: []string{"{3.14}"},
		},
		{
			Args: []string{"-3.14", "str"},
			Want: []string{"{-3.14}"},
		},
		{
			Args: []string{"-2d", "str"},
			Want: []string{"{-2d}"},
		},
		{
			Args: []string{"0d", "str"},
			Want: []string{"{0d}"},
		},
		{
			Args: []string{"2d", "str"},
			Want: []string{"{2d}"},
		},
		{
			Args: []string{"-10o", "str"},
			Want: []string{"{-10o}"},
		},
		{
			Args: []string{"0o", "str"},
			Want: []string{"{0o}"},
		},
		{
			Args: []string{"10o", "str"},
			Want: []string{"{10o}"},
		},
		{
			Args: []string{"-10b", "str"},
			Want: []string{"{-10b}"},
		},
		{
			Args: []string{"0b", "str"},
			Want: []string{"{0b}"},
		},
		{
			Args: []string{"10b", "str"},
			Want: []string{"{10b}"},
		},
		{
			Args: []string{"-fx", "str"},
			Want: []string{"{-fx}"},
		},
		{
			Args: []string{"0x", "str"},
			Want: []string{"{0x}"},
		},
		{
			Args: []string{"fx", "str"},
			Want: []string{"{fx}"},
		},
		{
			Args: []string{"true", "str"},
			Want: []string{"{true}"},
		},
		{
			Args: []string{"false", "str"},
			Want: []string{"{false}"},
		},
		{
			Args: []string{"'5'", "str"},
			Want: []string{"{5}"},
		},
		{
			Args: []string{"'-5'", "str"},
			Want: []string{"{-5}"},
		},
		{
			Args: []string{"'foo'", "str"},
			Want: []string{"{foo}"},
		},
	}
	rpn.UnitTestExecAll(t, data, func(r *rpn.RPN) { RegisterAll(r) })
}
