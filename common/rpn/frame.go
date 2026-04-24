package rpn

import (
	"math/cmplx"
	"strconv"
)

// FrameType is the combination of bitfields
// bit 0-3 is the display type (not usable in isolation)
// bit 4 is 1 if the type is a number, 0 otherwise (can use NUMBER_MASK)
// bit 4-8 represent the class (can use CLASS_MASK)
// bit 0-7 represent th type (no mask)

const (
	CLASS_MASK  = 0xF0
	NUMBER_MASK = 0x10
)

const (
	// & NUMBER_MASK = NUMBER_MASK
	INTEGER_CLASS = 0x10
	COMPLEX_CLASS = 0x30

	// & NUMBER_MASK = 0x00
	EMPTY_CLASS  = 0x00
	STRING_CLASS = 0x20
	BOOL_CLASS   = 0x40
)

type FrameType uint8

const (
	EMPTY_FRAME FrameType = FrameType(EMPTY_CLASS)

	STRING_SINGLEQ_FRAME = STRING_CLASS
	STRING_DOUBLEQ_FRAME = STRING_CLASS | 0x01
	STRING_BRACE_FRAME   = STRING_CLASS | 0x02

	COMPLEX_FRAME    = COMPLEX_CLASS
	POLAR_RAD_FRAME  = COMPLEX_CLASS | 0x01
	POLAR_DEG_FRAME  = COMPLEX_CLASS | 0x02
	POLAR_GRAD_FRAME = COMPLEX_CLASS | 0x03

	INTEGER_FRAME     = INTEGER_CLASS
	HEXIDECIMAL_FRAME = INTEGER_CLASS | 0x01
	OCTAL_FRAME       = INTEGER_CLASS | 0x02
	BINARY_FRAME      = INTEGER_CLASS | 0x03

	BOOL_FRAME = FrameType(BOOL_CLASS)
)

// Frame Defines a single stack frame
type Frame struct {
	ftype FrameType
	str   string
	cmplx complex128
	// If ftype == BOOL_FRAME, intv holds 1 or 0
	intv int64
}

// Annotates a frame.  Don't call this on string frames
// or the string will be replaced.
func (f *Frame) Annotate(s string) {
	f.str = s
}

func (f *Frame) IsInt() bool {
	return (f.ftype & CLASS_MASK) == INTEGER_CLASS
}

func (f *Frame) IsComplex() bool {
	return (f.ftype & CLASS_MASK) == COMPLEX_CLASS
}

func (f *Frame) IsNumber() bool {
	return (f.ftype & NUMBER_MASK) == NUMBER_MASK
}

func (f *Frame) IsBool() bool {
	return f.ftype == BOOL_FRAME
}

func (f *Frame) IsString() bool {
	return (f.ftype & CLASS_MASK) == STRING_CLASS
}

func (f *Frame) Type() FrameType {
	return f.ftype
}

func (f *Frame) Complex() (complex128, error) {
	if (f.ftype & CLASS_MASK) == COMPLEX_CLASS {
		return f.cmplx, nil
	}
	if f.IsInt() {
		return complex(float64(f.intv), 0), nil
	}
	return 0, ErrExpectedANumber
}

func (f *Frame) UnsafeComplex() complex128 {
	if (f.ftype & CLASS_MASK) == COMPLEX_CLASS {
		return f.cmplx
	}
	return complex(float64(f.intv), 0)
}

const complexTolerance = 1e-10

func (f *Frame) Real() (float64, error) {
	if (f.ftype & CLASS_MASK) == COMPLEX_CLASS {
		i := imag(f.cmplx)
		if (i != 0) && ((i < -complexTolerance) || (i > complexTolerance)) {
			return 0, ErrComplexNumberNotSupported
		}
		return real(f.cmplx), nil
	}
	if f.IsInt() {
		return float64(f.intv), nil
	}
	return 0, ErrExpectedANumber
}

func (f *Frame) Int() (int64, error) {
	if (f.ftype & CLASS_MASK) == INTEGER_CLASS {
		return f.intv, nil
	}
	if (f.ftype & CLASS_MASK) == COMPLEX_CLASS {
		if imag(f.cmplx) != 0 {
			return 0, ErrComplexNumberNotSupported
		}
		return int64(real(f.cmplx)), nil
	}
	return 0, ErrExpectedANumber
}

func (f *Frame) UnsafeInt() int64 {
	return f.intv
}

func (f *Frame) BoundedInt(min, max int64) (int64, error) {
	v, err := f.Int()
	if err != nil {
		return 0, err
	}
	if (v < min) || (v > max) {
		return 0, ErrIllegalValue
	}
	return v, nil
}

func (f *Frame) Bool() (bool, error) {
	if f.ftype == BOOL_FRAME {
		return f.intv != 0, nil
	}
	return false, ErrExpectedABoolean
}

func (f *Frame) UnsafeBool() bool {
	return f.intv != 0
}

func (f *Frame) String(quote bool) string {
	var s string
	switch f.ftype {
	case EMPTY_FRAME:
		return "nil"
	case STRING_SINGLEQ_FRAME:
		if quote {
			return "'" + f.str + "'"
		}
		return f.str
	case STRING_DOUBLEQ_FRAME:
		if quote {
			return "\"" + f.str + "\""
		}
		return f.str
	case STRING_BRACE_FRAME:
		if quote {
			return "{" + f.str + "}"
		}
		return f.str
	case COMPLEX_FRAME:
		s = f.complexString()
	case POLAR_RAD_FRAME, POLAR_DEG_FRAME, POLAR_GRAD_FRAME:
		s = f.polarString()
	case BOOL_FRAME:
		if f.intv != 0 {
			s = "true"
		} else {
			s = "false"
		}
	case INTEGER_FRAME:
		s = strconv.FormatInt(f.intv, 10) + "d"
	case HEXIDECIMAL_FRAME:
		s = strconv.FormatInt(f.intv, 16) + "x"
	case OCTAL_FRAME:
		s = strconv.FormatInt(f.intv, 8) + "o"
	case BINARY_FRAME:
		s = strconv.FormatInt(f.intv, 2) + "b"
	default:
		return "BAD_TYPE"
	}
	if len(f.str) > 0 {
		s += " " + f.str
	}
	return s
}

func (f *Frame) UnsafeString() string {
	return f.str
}

func BoolFrame(v bool) Frame {
	if v {
		return Frame{intv: 1, ftype: BOOL_FRAME}
	}
	return Frame{intv: 0, ftype: BOOL_FRAME}
}

func IntFrame(v int64, t FrameType) Frame {
	return Frame{intv: v, ftype: t}
}

func IntFrameCloneType(v int64, f Frame) Frame {
	return Frame{intv: v, ftype: f.ftype}
}

func ComplexFrame(v complex128) Frame {
	return Frame{cmplx: v, ftype: COMPLEX_FRAME}
}

func ComplexFrameWithType(v complex128, t FrameType) Frame {
	f := Frame{cmplx: v, ftype: t}
	switch t {
	case POLAR_DEG_FRAME:
		f.str = "`deg"
	case POLAR_RAD_FRAME:
		f.str = "`rad"
	case POLAR_GRAD_FRAME:
		f.str = "`grad"
	}
	return f
}

func PolarFrame(r, a float64, t FrameType) Frame {
	return ComplexFrameWithType(cmplx.Rect(r, toRadiansFloat(a, t)), t)
}

func RealFrame(v float64) Frame {
	return Frame{cmplx: complex(v, 0), ftype: COMPLEX_FRAME}
}

func StringFrame(v string, t FrameType) Frame {
	return Frame{str: v, ftype: t}
}

func EmptyFrame() Frame {
	return Frame{ftype: EMPTY_FRAME}
}

func (f *Frame) complexString() string {
	if imag(f.cmplx) == 0 {
		return strconv.FormatFloat(real(f.cmplx), 'g', 16, 64)
	}
	if real(f.cmplx) == 0 {
		return complexString(imag(f.cmplx))
	}
	r := strconv.FormatFloat(real(f.cmplx), 'g', 16, 64)
	if imag(f.cmplx) < 0 {
		return r + complexString(imag(f.cmplx))
	}
	return r + "+" + complexString(imag(f.cmplx))
}

func (f *Frame) polarString() string {
	r, a := cmplx.Polar(f.cmplx)
	return strconv.FormatFloat(r, 'g', 16, 64) + "<" + strconv.FormatFloat(FromRadiansFloat(a, f.ftype), 'g', 16, 64)
}

func complexString(v float64) string {
	if v == 1 {
		return "i"
	}
	if v == -1 {
		return "-i"
	}
	return strconv.FormatFloat(v, 'g', 16, 64) + "i"
}
