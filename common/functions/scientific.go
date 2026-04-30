package functions

import (
	"math"
	"math/cmplx"
	"mattwach/rpngo/common/rpn"
)

const powerHelp = "executes a to the power of b"

func power(r *rpn.RPN) error {
	af, bf, err := r.Pop2Frames()
	if err != nil {
		return err
	}
	if af.IsComplex() {
		b, err := bf.Complex()
		if err != nil {
			return err
		}

		if af.IsReal() && bf.IsReal() {
			// cmplx.Pow() can in some cases have unsightly rounding errors
			// (such as -6 2 **).  Catch the case where both a and b are real
			// numbers to avoid this common case.
			return r.PushFrame(
				rpn.RealFrame(math.Pow(real(af.UnsafeComplex()), real(b))))
		}

		return r.PushFrame(rpn.ComplexFrameWithType(cmplx.Pow(af.UnsafeComplex(), b), af.Type()))
	}
	if bf.IsComplex() {
		a, err := af.Complex()
		if err != nil {
			return err
		}
		return r.PushFrame(rpn.ComplexFrameWithType(cmplx.Pow(a, bf.UnsafeComplex()), bf.Type()))
	}
	a, err := af.Int()
	if err != nil {
		return err
	}
	b, err := bf.Int()
	if err != nil {
		return err
	}
	return r.PushFrame(rpn.IntFrameCloneType(powInts(a, b), af))
}

func powInts(x, n int64) int64 {
	if n < 0 {
		return 0
	}
	if n == 0 {
		return 1
	}
	if n == 1 {
		return x
	}
	y := powInts(x, n/2)
	if n%2 == 0 {
		return y * y
	}
	return x * y * y
}

const sqrtHelp = "takes the square root of a number"

func sqrt(r *rpn.RPN) error {
	af, err := r.PopFrame()
	if err != nil {
		return err
	}
	a, err := af.Complex()
	if err != nil {
		return err
	}
	a = cmplx.Sqrt(a)
	if af.IsComplex() {
		return r.PushFrame(rpn.ComplexFrameWithType(a, af.Type()))
	}
	return r.PushFrame(rpn.IntFrameCloneType(int64(real(a)), af))
}

const absHelp = "Takes the absolute value"

func abs(r *rpn.RPN) error {
	af, err := r.PopFrame()
	if err != nil {
		return err
	}
	if af.IsComplex() {
		a, _ := af.Complex()
		return r.PushFrame(rpn.RealFrame(cmplx.Abs(a)))
	}
	a, err := af.Int()
	if err != nil {
		return err
	}
	if a < 0 {
		a = -a
	}
	return r.PushFrame(rpn.IntFrameCloneType(a, af))
}

const sqHelp = "executes v * v"

func sq(r *rpn.RPN) error {
	a, err := r.PopFrame()
	if err != nil {
		return err
	}
	if a.IsComplex() {
		ac := a.UnsafeComplex()
		return r.PushFrame(rpn.ComplexFrameWithType(ac*ac, a.Type()))
	}
	ai, err := a.Int()
	if err != nil {
		return err
	}
	return r.PushFrame(rpn.IntFrameCloneType(ai*ai, a))
}

const logHelp = "executes natural logrithm"

func log(r *rpn.RPN) error {
	a, err := r.PopFrame()
	if err != nil {
		return err
	}
	if a.IsComplex() {
		return r.PushFrame(rpn.ComplexFrameWithType(cmplx.Log(a.UnsafeComplex()), a.Type()))
	}
	ac, err := a.Complex()
	if err != nil {
		return err
	}
	return r.PushFrame(rpn.ComplexFrame(cmplx.Log(ac)))
}

const log10Help = "executes base 10 logrithm"

func log10(r *rpn.RPN) error {
	a, err := r.PopFrame()
	if err != nil {
		return err
	}
	if a.IsComplex() {
		return r.PushFrame(rpn.ComplexFrameWithType(cmplx.Log10(a.UnsafeComplex()), a.Type()))
	}
	ac, err := a.Complex()
	if err != nil {
		return err
	}
	return r.PushFrame(rpn.ComplexFrame(cmplx.Log10(ac)))
}
