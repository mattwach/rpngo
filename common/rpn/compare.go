package rpn

const tolerance = 1e-9

func checkFloatEqual(a, b complex128) bool {
	d := real(a) - real(b)
	if d < -tolerance {
		return false
	}
	if d > tolerance {
		return false
	}
	d = imag(a) - imag(b)
	if d < -tolerance {
		return false
	}
	return d <= tolerance
}

// Ordering of non-comparable types
// bool < number < string

func (a Frame) IsLessThan(b Frame) bool {
	switch a.ftype & CLASS_MASK {
	case COMPLEX_CLASS:
		switch b.ftype & CLASS_MASK {
		case COMPLEX_CLASS:
			return real(a.cmplx) < real(b.cmplx)
		case INTEGER_CLASS:
			return real(a.cmplx) < float64(b.intv)
		case STRING_CLASS:
			return true
		case BOOL_CLASS:
			return false
		}
	case INTEGER_CLASS:
		switch b.ftype & CLASS_MASK {
		case COMPLEX_CLASS:
			return float64(a.intv) < real(b.cmplx)
		case INTEGER_CLASS:
			return a.intv < b.intv
		case STRING_CLASS:
			return true
		case BOOL_CLASS:
			return false
		}
	case STRING_CLASS:
		switch b.ftype & CLASS_MASK {
		case STRING_CLASS:
			return a.str < b.str
		default:
			return false
		}
	case BOOL_CLASS:
		switch b.ftype & CLASS_MASK {
		case BOOL_CLASS:
			return a.intv < b.intv
		default:
			return true
		}
	}

	return false
}

func (a Frame) IsLessThanOrEqual(b Frame) bool {
	switch a.ftype & CLASS_MASK {
	case COMPLEX_CLASS:
		switch b.ftype & CLASS_MASK {
		case COMPLEX_CLASS:
			return checkFloatEqual(a.cmplx, b.cmplx) || (real(a.cmplx) < real(b.cmplx))
		case INTEGER_CLASS:
			return checkFloatEqual(a.cmplx, complex(float64(b.intv), 0)) || (real(a.cmplx) < float64(b.intv))
		case STRING_CLASS:
			return true
		case BOOL_CLASS:
			return false
		}
	case INTEGER_CLASS:
		switch b.ftype & CLASS_MASK {
		case COMPLEX_CLASS:
			return checkFloatEqual(complex(float64(a.intv), 0), b.cmplx) || float64(a.intv) < real(b.cmplx)
		case INTEGER_CLASS:
			return a.intv <= b.intv
		case STRING_CLASS:
			return true
		case BOOL_CLASS:
			return false
		}
	case STRING_CLASS:
		switch b.ftype & CLASS_MASK {
		case STRING_CLASS:
			return a.str <= b.str
		default:
			return false
		}
	case BOOL_CLASS:
		switch b.ftype & CLASS_MASK {
		case BOOL_CLASS:
			return a.intv <= b.intv
		default:
			return true
		}
	}

	return false
}

func (a Frame) IsEqual(b Frame) bool {
	switch a.ftype & CLASS_MASK {
	case COMPLEX_CLASS:
		switch b.ftype & CLASS_MASK {
		case COMPLEX_CLASS:
			return checkFloatEqual(a.cmplx, b.cmplx)
		case INTEGER_CLASS:
			return checkFloatEqual(a.cmplx, complex(float64(b.intv), 0))
		case STRING_CLASS:
			return false
		case BOOL_FRAME:
			return false
		}
	case INTEGER_CLASS:
		switch b.ftype & CLASS_MASK {
		case COMPLEX_CLASS:
			return checkFloatEqual(complex(float64(a.intv), 0), b.cmplx)
		case INTEGER_CLASS:
			return a.intv == b.intv
		case STRING_CLASS:
			return false
		case BOOL_CLASS:
			return false
		}
	case STRING_CLASS:
		switch b.ftype & CLASS_MASK {
		case STRING_CLASS:
			return a.str == b.str
		default:
			return false
		}
	case BOOL_CLASS:
		switch b.ftype & CLASS_MASK {
		case BOOL_CLASS:
			return a.intv == b.intv
		default:
			return false
		}
	}

	return false
}
