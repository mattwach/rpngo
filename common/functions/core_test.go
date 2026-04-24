package functions

import (
	"mattwach/rpngo/common/rpn"
	"testing"
)

func TestNOOP(t *testing.T) {
	data := []rpn.UnitTestExecData{
		{
			Args: []string{"noop"},
		},
	}
	rpn.UnitTestExecAll(t, data, func(r *rpn.RPN) { RegisterAll(r) })
}

func TestAdd(t *testing.T) {
	data := []rpn.UnitTestExecData{
		{
			Args:    []string{"+"},
			WantErr: rpn.ErrStackEmpty,
		},
		{
			Args:    []string{"1", "+"},
			WantErr: rpn.ErrStackEmpty,
		},
		{
			Args:    []string{"1", "true", "+"},
			WantErr: rpn.ErrExpectedANumber,
		},
		{
			Args:    []string{"true", "1", "+"},
			WantErr: rpn.ErrExpectedANumber,
		},
		{
			Args: []string{"1", "2", "+"},
			Want: []string{"3"},
		},
		{
			Args: []string{"1", "2d", "+"},
			Want: []string{"3"},
		},
		{
			Args: []string{"1d", "2", "+"},
			Want: []string{"3"},
		},
		{
			Args: []string{"1d", "2d", "+"},
			Want: []string{"3d"},
		},
		{
			Args: []string{"'foo'", "7", "+"},
			Want: []string{"'foo7'"},
		},
		{
			Args: []string{"7", "'foo'", "+"},
			Want: []string{"'7foo'"},
		},
		{
			Args: []string{"'foo'", "7d", "+"},
			Want: []string{"'foo7d'"},
		},
		{
			Args: []string{"7d", "'foo'", "+"},
			Want: []string{"'7dfoo'"},
		},
		{
			Args: []string{"'foo'", "true", "+"},
			Want: []string{"'footrue'"},
		},
		{
			Args: []string{"true", "'foo'", "+"},
			Want: []string{"'truefoo'"},
		},
		{
			Args: []string{"\"foo\"", "'bar'", "+"},
			Want: []string{"\"foobar\""},
		},
	}
	rpn.UnitTestExecAll(t, data, func(r *rpn.RPN) { RegisterAll(r) })
}

func TestSubtract(t *testing.T) {
	data := []rpn.UnitTestExecData{
		{
			Args:    []string{"-"},
			WantErr: rpn.ErrStackEmpty,
		},
		{
			Args:    []string{"1", "-"},
			WantErr: rpn.ErrStackEmpty,
		},
		{
			Args:    []string{"1", "true", "-"},
			WantErr: rpn.ErrExpectedANumber,
		},
		{
			Args:    []string{"true", "1", "-"},
			WantErr: rpn.ErrExpectedANumber,
		},
		{
			Args: []string{"1", "2", "-"},
			Want: []string{"-1"},
		},
		{
			Args: []string{"1", "2d", "-"},
			Want: []string{"-1"},
		},
		{
			Args: []string{"1d", "2", "-"},
			Want: []string{"-1"},
		},
		{
			Args: []string{"1d", "2d", "-"},
			Want: []string{"-1d"},
		},
		{
			Args:    []string{"'foo'", "7", "-"},
			WantErr: rpn.ErrExpectedANumber,
		},
		{
			Args:    []string{"7", "'foo'", "-"},
			WantErr: rpn.ErrExpectedANumber,
		},
		{
			Args:    []string{"\"foo\"", "'bar'", "-"},
			WantErr: rpn.ErrExpectedANumber,
		},
	}
	rpn.UnitTestExecAll(t, data, func(r *rpn.RPN) { RegisterAll(r) })
}

func TestMultiply(t *testing.T) {
	data := []rpn.UnitTestExecData{
		{
			Args:    []string{"*"},
			WantErr: rpn.ErrStackEmpty,
		},
		{
			Args:    []string{"1", "*"},
			WantErr: rpn.ErrStackEmpty,
		},
		{
			Args:    []string{"1", "true", "*"},
			WantErr: rpn.ErrExpectedANumber,
		},
		{
			Args:    []string{"true", "1", "*"},
			WantErr: rpn.ErrExpectedANumber,
		},
		{
			Args: []string{"2", "3", "*"},
			Want: []string{"6"},
		},
		{
			Args: []string{"2", "3d", "*"},
			Want: []string{"6"},
		},
		{
			Args: []string{"2d", "3", "*"},
			Want: []string{"6"},
		},
		{
			Args: []string{"2d", "3d", "*"},
			Want: []string{"6d"},
		},
		{
			Args:    []string{"'foo'", "7", "*"},
			WantErr: rpn.ErrExpectedANumber,
		},
		{
			Args:    []string{"7", "'foo'", "*"},
			WantErr: rpn.ErrExpectedANumber,
		},
		{
			Args:    []string{"\"foo\"", "'bar'", "*"},
			WantErr: rpn.ErrExpectedANumber,
		},
	}
	rpn.UnitTestExecAll(t, data, func(r *rpn.RPN) { RegisterAll(r) })
}

func TestDivide(t *testing.T) {
	data := []rpn.UnitTestExecData{
		{
			Args:    []string{"/"},
			WantErr: rpn.ErrStackEmpty,
		},
		{
			Args:    []string{"1", "/"},
			WantErr: rpn.ErrStackEmpty,
		},
		{
			Args:    []string{"1", "true", "/"},
			WantErr: rpn.ErrExpectedANumber,
		},
		{
			Args:    []string{"true", "1", "/"},
			WantErr: rpn.ErrExpectedANumber,
		},
		{
			Args: []string{"5", "2", "/"},
			Want: []string{"2.5"},
		},
		{
			Args: []string{"5", "2d", "/"},
			Want: []string{"2.5"},
		},
		{
			Args: []string{"5d", "2", "/"},
			Want: []string{"2.5"},
		},
		{
			Args: []string{"5d", "2d", "/"},
			Want: []string{"2d"},
		},
		{
			Args:    []string{"5", "0", "/"},
			WantErr: rpn.ErrDivideByZero,
		},
		{
			Args:    []string{"5d", "0d", "/"},
			WantErr: rpn.ErrDivideByZero,
		},
		{
			Args:    []string{"5", "0d", "/"},
			WantErr: rpn.ErrDivideByZero,
		},
		{
			Args:    []string{"5d", "0", "/"},
			WantErr: rpn.ErrDivideByZero,
		},
		{
			Args:    []string{"'foo'", "7", "/"},
			WantErr: rpn.ErrExpectedANumber,
		},
		{
			Args:    []string{"7", "'foo'", "/"},
			WantErr: rpn.ErrExpectedANumber,
		},
		{
			Args:    []string{"\"foo\"", "'bar'", "/"},
			WantErr: rpn.ErrExpectedANumber,
		},
	}
	rpn.UnitTestExecAll(t, data, func(r *rpn.RPN) { RegisterAll(r) })
}

func TestMod(t *testing.T) {
	data := []rpn.UnitTestExecData{
		{
			Args:    []string{"%"},
			WantErr: rpn.ErrStackEmpty,
		},
		{
			Args:    []string{"1", "%"},
			WantErr: rpn.ErrStackEmpty,
		},
		{
			Args:    []string{"1", "true", "%"},
			WantErr: rpn.ErrExpectedANumber,
		},
		{
			Args:    []string{"true", "1", "%"},
			WantErr: rpn.ErrExpectedANumber,
		},
		{
			Args: []string{"5", "2", "%"},
			Want: []string{"1d"},
		},
		{
			Args: []string{"5", "2d", "%"},
			Want: []string{"1d"},
		},
		{
			Args: []string{"5d", "2", "%"},
			Want: []string{"1d"},
		},
		{
			Args: []string{"5d", "2d", "%"},
			Want: []string{"1d"},
		},
		{
			Args:    []string{"5", "0", "%"},
			WantErr: rpn.ErrDivideByZero,
		},
		{
			Args:    []string{"5d", "0d", "%"},
			WantErr: rpn.ErrDivideByZero,
		},
		{
			Args:    []string{"5", "0d", "%"},
			WantErr: rpn.ErrDivideByZero,
		},
		{
			Args:    []string{"5d", "0", "%"},
			WantErr: rpn.ErrDivideByZero,
		},
		{
			Args:    []string{"'foo'", "7", "%"},
			WantErr: rpn.ErrExpectedANumber,
		},
		{
			Args:    []string{"7", "'foo'", "%"},
			WantErr: rpn.ErrExpectedANumber,
		},
		{
			Args:    []string{"\"foo\"", "'bar'", "%"},
			WantErr: rpn.ErrExpectedANumber,
		},
	}
	rpn.UnitTestExecAll(t, data, func(r *rpn.RPN) { RegisterAll(r) })
}

func TestNegate(t *testing.T) {
	data := []rpn.UnitTestExecData{
		{
			Args:    []string{"neg"},
			WantErr: rpn.ErrStackEmpty,
		},
		{
			Args:    []string{"'foo'", "neg"},
			WantErr: rpn.ErrIllegalValue,
		},
		{
			Args: []string{"1", "neg"},
			Want: []string{"-1"},
		},
		{
			Args: []string{"-1", "neg"},
			Want: []string{"1"},
		},
		{
			Args: []string{"0", "neg"},
			Want: []string{"-0"},
		},
		{
			Args: []string{"1d", "neg"},
			Want: []string{"-1d"},
		},
		{
			Args: []string{"0d", "neg"},
			Want: []string{"0d"},
		},
		{
			Args: []string{"true", "neg"},
			Want: []string{"false"},
		},
		{
			Args: []string{"false", "neg"},
			Want: []string{"true"},
		},
	}
	rpn.UnitTestExecAll(t, data, func(r *rpn.RPN) { RegisterAll(r) })
}

func TestExec(t *testing.T) {
	data := []rpn.UnitTestExecData{
		{
			Args:    []string{"@"},
			WantErr: rpn.ErrStackEmpty,
		},
		{
			Args: []string{"''", "@"},
		},
		{
			Args: []string{"'3'", "@"},
			Want: []string{"3"},
		},
		{
			Args: []string{"'2 3 +'", "@"},
			Want: []string{"5"},
		},
		{
			Args:    []string{"foo", "@"},
			WantErr: rpn.ErrSyntax,
		},
	}
	rpn.UnitTestExecAll(t, data, func(r *rpn.RPN) { RegisterAll(r) })
}

func TestRand(t *testing.T) {
	data := []rpn.UnitTestExecData{
		{
			Args: []string{"rand", "$0", "0", ">=", "1>", "1", "<="},
			Want: []string{"true", "true"},
		},
	}
	rpn.UnitTestExecAll(t, data, func(r *rpn.RPN) { RegisterAll(r) })
}

func TestReal(t *testing.T) {
	data := []rpn.UnitTestExecData{
		{
			Args:    []string{"real"},
			WantErr: rpn.ErrStackEmpty,
		},
		{
			Args: []string{"1d", "real"},
			Want: []string{"1"},
		},
		{
			Args:    []string{"'foo'", "real"},
			WantErr: rpn.ErrExpectedANumber,
		},
		{
			Args: []string{"1", "real"},
			Want: []string{"1"},
		},
		{
			Args: []string{"i", "real"},
			Want: []string{"0"},
		},
		{
			Args: []string{"2+3i", "real"},
			Want: []string{"2"},
		},
	}
	rpn.UnitTestExecAll(t, data, func(r *rpn.RPN) { RegisterAll(r) })
}

func TestImag(t *testing.T) {
	data := []rpn.UnitTestExecData{
		{
			Args:    []string{"imag"},
			WantErr: rpn.ErrStackEmpty,
		},
		{
			Args: []string{"1d", "imag"},
			Want: []string{"0"},
		},
		{
			Args:    []string{"'foo'", "imag"},
			WantErr: rpn.ErrExpectedANumber,
		},
		{
			Args: []string{"1", "imag"},
			Want: []string{"0"},
		},
		{
			Args: []string{"i", "imag"},
			Want: []string{"1"},
		},
		{
			Args: []string{"2+3i", "imag"},
			Want: []string{"3"},
		},
	}
	rpn.UnitTestExecAll(t, data, func(r *rpn.RPN) { RegisterAll(r) })
}

func TestTrueFalse(t *testing.T) {
	data := []rpn.UnitTestExecData{
		{
			Args: []string{"true"},
			Want: []string{"true"},
		},
		{
			Args: []string{"false"},
			Want: []string{"false"},
		},
	}
	rpn.UnitTestExecAll(t, data, func(r *rpn.RPN) { RegisterAll(r) })
}

func TestRound(t *testing.T) {
	data := []rpn.UnitTestExecData{
		{
			Args:    []string{"round"},
			WantErr: rpn.ErrStackEmpty,
		},
		{
			Args:    []string{"1", "round"},
			WantErr: rpn.ErrStackEmpty,
		},
		{
			Args: []string{"1.2345", "2", "round"},
			Want: []string{"1.23"},
		},
		{
			Args: []string{"1.235", "2", "round"},
			Want: []string{"1.24"},
		},
		{
			Args: []string{"-1.2345", "2", "round"},
			Want: []string{"-1.23"},
		},
		{
			Args: []string{"-1.235", "2", "round"},
			Want: []string{"-1.24"},
		},
		{
			Args: []string{"1.234<2.348", "2", "round"},
			Want: []string{"1.23<2.35 `rad"},
		},
		{
			Args: []string{"deg", "1.234<2.348", "2", "round"},
			Want: []string{"1.23<2.35 `deg"},
		},
		{
			Args: []string{"grad", "1.234<2.34", "2", "round"},
			Want: []string{"1.23<2.34 `grad"},
		},
		{
			Args:    []string{"-1.235", "i", "round"},
			WantErr: rpn.ErrComplexNumberNotSupported,
		},
		{
			Args: []string{"1.235", "2.1", "round"},
			Want: []string{"1.24"},
		},
		{
			Args:    []string{"1.235", "-1", "round"},
			WantErr: rpn.ErrIllegalValue,
		},
		{
			Args:    []string{"1.235", "17", "round"},
			WantErr: rpn.ErrIllegalValue,
		},
	}
	rpn.UnitTestExecAll(t, data, func(r *rpn.RPN) { RegisterAll(r) })
}

func TestFloat(t *testing.T) {
	data := []rpn.UnitTestExecData{
		{
			Args: []string{"-1", "float"},
			Want: []string{"-1"},
		},
		{
			Args: []string{"0", "float"},
			Want: []string{"0"},
		},
		{
			Args: []string{"1", "float"},
			Want: []string{"1"},
		},
		{
			Args: []string{"3+4i", "float"},
			Want: []string{"3+4i"},
		},
		{
			Args: []string{"3.14", "float"},
			Want: []string{"3.14"},
		},
		{
			Args: []string{"-3.14", "float"},
			Want: []string{"-3.14"},
		},
		{
			Args: []string{"-2d", "float"},
			Want: []string{"-2"},
		},
		{
			Args: []string{"0d", "float"},
			Want: []string{"0"},
		},
		{
			Args: []string{"2d", "float"},
			Want: []string{"2"},
		},
		{
			Args: []string{"-10o", "float"},
			Want: []string{"-8"},
		},
		{
			Args: []string{"0o", "float"},
			Want: []string{"0"},
		},
		{
			Args: []string{"10o", "float"},
			Want: []string{"8"},
		},
		{
			Args: []string{"-10b", "float"},
			Want: []string{"-2"},
		},
		{
			Args: []string{"0b", "float"},
			Want: []string{"0"},
		},
		{
			Args: []string{"10b", "float"},
			Want: []string{"2"},
		},
		{
			Args: []string{"-fx", "float"},
			Want: []string{"-15"},
		},
		{
			Args: []string{"0x", "float"},
			Want: []string{"0"},
		},
		{
			Args: []string{"fx", "float"},
			Want: []string{"15"},
		},
		{
			Args: []string{"true", "float"},
			Want: []string{"1"},
		},
		{
			Args: []string{"false", "float"},
			Want: []string{"0"},
		},
		{
			Args: []string{"'5'", "float"},
			Want: []string{"5"},
		},
		{
			Args: []string{"'-5'", "float"},
			Want: []string{"-5"},
		},
		{
			Args:    []string{"'foo'", "float"},
			WantErr: rpn.ErrSyntax,
		},
	}
	rpn.UnitTestExecAll(t, data, func(r *rpn.RPN) { RegisterAll(r) })
}

func TestPolar(t *testing.T) {
	data := []rpn.UnitTestExecData{
		{
			Args: []string{"-1", "polar"},
			Want: []string{"1<3.141592653589793 `rad"},
		},
		{
			Args: []string{"0", "polar"},
			Want: []string{"0<0 `rad"},
		},
		{
			Args: []string{"1", "polar"},
			Want: []string{"1<0 `rad"},
		},
		{
			Args: []string{"3+4i", "polar"},
			Want: []string{"5<0.9272952180016122 `rad"},
		},
		{
			Args: []string{"3.14", "polar"},
			Want: []string{"3.14<0 `rad"},
		},
		{
			Args: []string{"-3.14", "polar"},
			Want: []string{"3.14<3.141592653589793 `rad"},
		},
		{
			Args: []string{"-2d", "polar"},
			Want: []string{"2<3.141592653589793 `rad"},
		},
		{
			Args: []string{"0d", "polar"},
			Want: []string{"0<0 `rad"},
		},
		{
			Args: []string{"2d", "polar"},
			Want: []string{"2<0 `rad"},
		},
		{
			Args: []string{"-10o", "polar"},
			Want: []string{"8<3.141592653589793 `rad"},
		},
		{
			Args: []string{"0o", "polar"},
			Want: []string{"0<0 `rad"},
		},
		{
			Args: []string{"10o", "polar"},
			Want: []string{"8<0 `rad"},
		},
		{
			Args: []string{"-10b", "polar"},
			Want: []string{"2<3.141592653589793 `rad"},
		},
		{
			Args: []string{"0b", "polar"},
			Want: []string{"0<0 `rad"},
		},
		{
			Args: []string{"10b", "polar"},
			Want: []string{"2<0 `rad"},
		},
		{
			Args: []string{"-fx", "polar"},
			Want: []string{"15<3.141592653589793 `rad"},
		},
		{
			Args: []string{"0x", "polar"},
			Want: []string{"0<0 `rad"},
		},
		{
			Args: []string{"fx", "polar"},
			Want: []string{"15<0 `rad"},
		},
		{
			Args: []string{"true", "polar"},
			Want: []string{"1<0 `rad"},
		},
		{
			Args: []string{"false", "polar"},
			Want: []string{"0<0 `rad"},
		},
		{
			Args: []string{"'5'", "polar"},
			Want: []string{"5<0 `rad"},
		},
		{
			Args: []string{"'-5'", "polar"},
			Want: []string{"5<3.141592653589793 `rad"},
		},
		{
			Args:    []string{"'foo'", "polar"},
			WantErr: rpn.ErrSyntax,
		},
	}
	rpn.UnitTestExecAll(t, data, func(r *rpn.RPN) { RegisterAll(r) })
}

func TestPhase(t *testing.T) {
	data := []rpn.UnitTestExecData{
		{
			Args:    []string{"phase"},
			WantErr: rpn.ErrStackEmpty,
		},
		{
			Args: []string{"1d", "phase"},
			Want: []string{"0 `rad"},
		},
		{
			Args:    []string{"'foo'", "imag"},
			WantErr: rpn.ErrExpectedANumber,
		},
		{
			Args: []string{"1", "phase"},
			Want: []string{"0 `rad"},
		},
		{
			Args: []string{"deg", "i", "phase"},
			Want: []string{"90 `deg"},
		},
		{
			Args: []string{"grad", "i", "phase"},
			Want: []string{"100 `grad"},
		},
	}
	rpn.UnitTestExecAll(t, data, func(r *rpn.RPN) { RegisterAll(r) })
}

func TestFrac(t *testing.T) {
	data := []rpn.UnitTestExecData{
		{
			Args:    []string{"frac"},
			WantErr: rpn.ErrStackEmpty,
		},
		{
			Args:    []string{"'foo'", "frac"},
			WantErr: rpn.ErrExpectedANumber,
		},
		{
			Args: []string{"1", "frac"},
			Want: []string{"0"},
		},
		{
			Args: []string{"1.2", "frac"},
			Want: []string{"0.2"},
		},
		{
			Args: []string{"1.8", "frac"},
			Want: []string{"0.8"},
		},
		{
			Args: []string{"-1", "frac"},
			Want: []string{"0"},
		},
		{
			Args: []string{"-1.2", "frac"},
			Want: []string{"-0.2"},
		},
		{
			Args: []string{"-1.8", "frac"},
			Want: []string{"-0.8"},
		},
		{
			Args: []string{"1+i", "frac"},
			Want: []string{"0"},
		},
		{
			Args: []string{"1.2+1.5i", "frac"},
			Want: []string{"0.2+0.5i"},
		},
		{
			Args: []string{"1.8+1.8i", "frac"},
			Want: []string{"0.8+0.8i"},
		},
		{
			Args: []string{"-1-i", "frac"},
			Want: []string{"0"},
		},
		{
			Args: []string{"-1.2-1.5i", "frac"},
			Want: []string{"-0.2-0.5i"},
		},
		{
			Args: []string{"-1.8-1.8i", "frac"},
			Want: []string{"-0.8-0.8i"},
		},
		{
			Args: []string{"1x", "frac"},
			Want: []string{"0x"},
		},
	}
	rpn.UnitTestExecAll(t, data, func(r *rpn.RPN) { RegisterAll(r) })
}
