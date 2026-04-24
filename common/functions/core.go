// package functions defines core functions
package functions

import (
	"math"
	"math/cmplx"
	"math/rand"
	"mattwach/rpngo/common/parse"
	"mattwach/rpngo/common/rpn"
)

const noopHelp = "No operation. e.g. 'noop' plot will plot y = x"

func noop(r *rpn.RPN) error {
	return nil
}

const addHelp = "Adds two numbers"

func add(r *rpn.RPN) error {
	a, b, err := r.Pop2Frames()
	if err != nil {
		return err
	}
	if a.IsString() {
		return r.PushFrame(rpn.StringFrame(a.String(false)+b.String(false), a.Type()))
	}
	if b.IsString() {
		return r.PushFrame(rpn.StringFrame(a.String(false)+b.String(false), b.Type()))
	}
	if a.IsComplex() {
		bc, err := b.Complex()
		if err != nil {
			return err
		}
		return r.PushFrame(rpn.ComplexFrameWithType(a.UnsafeComplex()+bc, a.Type()))
	}
	if b.IsComplex() {
		ac, err := a.Complex()
		if err != nil {
			return err
		}
		return r.PushFrame(rpn.ComplexFrameWithType(ac+b.UnsafeComplex(), b.Type()))
	}
	ab, err := a.Int()
	if err != nil {
		return err
	}
	bb, err := b.Int()
	if err != nil {
		return err
	}
	return r.PushFrame(rpn.IntFrameCloneType(ab+bb, a))
}

const subtractHelp = "Subtracts two numbers"

func subtract(r *rpn.RPN) error {
	a, b, err := r.Pop2Frames()
	if err != nil {
		return err
	}
	if a.IsComplex() {
		bc, err := b.Complex()
		if err != nil {
			return err
		}
		return r.PushFrame(rpn.ComplexFrameWithType(a.UnsafeComplex()-bc, a.Type()))
	}
	if b.IsComplex() {
		ac, err := a.Complex()
		if err != nil {
			return err
		}
		return r.PushFrame(rpn.ComplexFrameWithType(ac-b.UnsafeComplex(), b.Type()))
	}
	ab, err := a.Int()
	if err != nil {
		return err
	}
	bb, err := b.Int()
	if err != nil {
		return err
	}
	return r.PushFrame(rpn.IntFrameCloneType(ab-bb, a))
}

const multiplyHelp = "Multiplies two numbers"

func multiply(r *rpn.RPN) error {
	a, b, err := r.Pop2Frames()
	if err != nil {
		return err
	}
	if a.IsComplex() {
		bc, err := b.Complex()
		if err != nil {
			return err
		}
		return r.PushFrame(rpn.ComplexFrameWithType(a.UnsafeComplex()*bc, a.Type()))
	}
	if b.IsComplex() {
		ac, err := a.Complex()
		if err != nil {
			return err
		}
		return r.PushFrame(rpn.ComplexFrameWithType(ac*b.UnsafeComplex(), b.Type()))
	}
	ab, err := a.Int()
	if err != nil {
		return err
	}
	bb, err := b.Int()
	if err != nil {
		return err
	}
	return r.PushFrame(rpn.IntFrameCloneType(ab*bb, a))
}

const divideHelp = "Divides two numbers"

func divide(r *rpn.RPN) error {
	a, b, err := r.Pop2Frames()
	if err != nil {
		return err
	}
	if a.IsComplex() {
		bc, err := b.Complex()
		if err != nil {
			return err
		}
		if bc == 0 {
			return rpn.ErrDivideByZero
		}
		return r.PushFrame(rpn.ComplexFrameWithType(a.UnsafeComplex()/bc, a.Type()))
	}
	if b.IsComplex() {
		ac, err := a.Complex()
		if err != nil {
			return err
		}
		bc := b.UnsafeComplex()
		if bc == 0 {
			return rpn.ErrDivideByZero
		}
		return r.PushFrame(rpn.ComplexFrameWithType(ac/bc, b.Type()))
	}
	ab, err := a.Int()
	if err != nil {
		return err
	}
	bb, err := b.Int()
	if err != nil {
		return err
	}
	if bb == 0 {
		return rpn.ErrDivideByZero
	}
	return r.PushFrame(rpn.IntFrameCloneType(ab/bb, a))
}

const negateHelp = "Negates the top number"

func negate(r *rpn.RPN) error {
	f, err := r.PopFrame()
	if err != nil {
		return err
	}
	if f.IsComplex() {
		c, _ := f.Complex()
		return r.PushFrame(rpn.ComplexFrameWithType(-c, f.Type()))
	}
	if f.IsBool() {
		b, _ := f.Bool()
		return r.PushFrame(rpn.BoolFrame(!b))
	}
	if f.IsInt() {
		i, _ := f.Int()
		return r.PushFrame(rpn.IntFrameCloneType(-i, f))
	}
	return rpn.ErrIllegalValue
}

const execHelp = "Executes a string\n" +
	"Example: '4 5 +' @"

func exec(r *rpn.RPN) error {
	f, err := r.PopFrame()
	if err != nil {
		return err
	}
	if !f.IsString() {
		return rpn.ErrExpectedAString
	}
	return parse.Fields(f.UnsafeString(), r.Exec)
}

const randHelp = "Pushes a random number between 0 and 1"

func randFn(r *rpn.RPN) error {
	return r.PushFrame(rpn.RealFrame(rand.Float64()))
}

const polarHelp = "Converts head element to a complex polar"

func polar(r *rpn.RPN) error {
	f, err := r.PopFrame()
	if err != nil {
		return err
	}
	if f.IsBool() {
		if f.UnsafeBool() {
			return r.PushFrame(rpn.PolarFrame(1, 0, r.AngleUnit))
		} else {
			return r.PushFrame(rpn.PolarFrame(0, 0, r.AngleUnit))
		}
	}
	if f.IsString() {
		err := r.Exec(f.UnsafeString())
		if err != nil {
			return err
		}
		f, err = r.PopFrame()
		if err != nil {
			return err
		}
	}
	v, err := f.Complex()
	if err != nil {
		return err
	}
	return r.PushFrame(rpn.ComplexFrameWithType(v, r.AngleUnit))
}

const phaseHelp = "Returns the polar angle of a number as a real (use abs for magnitude)"

func phase(r *rpn.RPN) error {
	f, err := r.PopFrame()
	if err != nil {
		return err
	}
	c, err := f.Complex()
	if err != nil {
		return err
	}
	p := rpn.RealFrame(rpn.FromRadiansFloat(cmplx.Phase(c), r.AngleUnit))
	switch r.AngleUnit {
	case rpn.POLAR_RAD_FRAME:
		p.Annotate("`rad")
	case rpn.POLAR_DEG_FRAME:
		p.Annotate("`deg")
	case rpn.POLAR_GRAD_FRAME:
		p.Annotate("`grad")
	}
	return r.PushFrame(p)
}

const floatHelp = "Converts head element to a complex float"

func floatFn(r *rpn.RPN) error {
	f, err := r.PopFrame()
	if err != nil {
		return err
	}
	if f.IsBool() {
		if f.UnsafeBool() {
			return r.PushFrame(rpn.RealFrame(1))
		} else {
			return r.PushFrame(rpn.RealFrame(0))
		}
	}
	if f.IsString() {
		err := r.Exec(f.UnsafeString())
		if err != nil {
			return err
		}
		f, err = r.PopFrame()
		if err != nil {
			return err
		}
	}
	v, err := f.Complex()
	if err != nil {
		return err
	}
	return r.PushFrame(rpn.ComplexFrame(v))
}

const realHelp = "Takes the real portion of a complex number"

func realFn(r *rpn.RPN) error {
	f, err := r.PopFrame()
	if err != nil {
		return err
	}
	c, err := f.Complex()
	if err != nil {
		return err
	}
	return r.PushFrame(rpn.RealFrame(real(c)))
}

const imagHelp = "Takes the imaginary portion of a complex number (as a real number)"

func imagFn(r *rpn.RPN) error {
	f, err := r.PopFrame()
	if err != nil {
		return err
	}
	c, err := f.Complex()
	if err != nil {
		return err
	}
	return r.PushFrame(rpn.RealFrame(imag(c)))
}

const trueHelp = "Pushes a boolean true"

func trueFn(r *rpn.RPN) error {
	return r.PushFrame(rpn.BoolFrame(true))
}

const falseHelp = "Pushes a boolean false"

func falseFn(r *rpn.RPN) error {
	return r.PushFrame(rpn.BoolFrame(false))
}

const roundHelp = "Rounds a number to the given number of places"

func round(r *rpn.RPN) error {
	af, bf, err := r.Pop2Frames()
	if err != nil {
		return err
	}
	a, err := af.Complex()
	if err != nil {
		return err
	}
	b, err := bf.Int()
	if err != nil {
		return err
	}
	if (b < 0) || (b > 16) {
		return rpn.ErrIllegalValue
	}
	if af.IsComplex() && (af.Type() != rpn.COMPLEX_FRAME) {
		rl, an := cmplx.Polar(a)
		an = rpn.FromRadiansFloat(an, af.Type())
		for i := 0; i < int(b); i++ {
			rl *= 10
			an *= 10
		}
		rl = math.Round(rl)
		an = math.Round(an)
		for i := 0; i < int(b); i++ {
			rl /= 10
			an /= 10
		}
		return r.PushFrame(rpn.PolarFrame(rl, an, af.Type()))
	}
	rl := real(a)
	im := imag(a)
	for i := 0; i < int(b); i++ {
		rl *= 10
		im *= 10
	}
	rl = math.Round(rl)
	im = math.Round(im)
	for i := 0; i < int(b); i++ {
		rl /= 10
		im /= 10
	}
	return r.PushFrame(rpn.ComplexFrame(complex(rl, im)))
}

const fracHelp = "Returns the fractional parts of a number"

func frac(r *rpn.RPN) error {
	f, err := r.PopFrame()
	if err != nil {
		return err
	}
	if f.IsComplex() {
		c := f.UnsafeComplex()
		rl := real(c)
		rl = rl - math.Trunc(rl)
		i := imag(c)
		i = i - math.Trunc(i)
		return r.PushFrame(rpn.ComplexFrameWithType(complex(rl, i), f.Type()))
	}
	if f.IsNumber() {
		return r.PushFrame(rpn.IntFrameCloneType(0, f))
	}
	return rpn.ErrExpectedANumber
}

const modHelp = "Returns the modulus of two integers"

func mod(r *rpn.RPN) error {
	af, bf, err := r.Pop2Frames()
	if err != nil {
		return err
	}
	a, err := af.Int()
	if err != nil {
		return err
	}
	b, err := bf.Int()
	if err != nil {
		return err
	}
	if b == 0 {
		return rpn.ErrDivideByZero
	}
	if af.IsInt() {
		return r.PushFrame(rpn.IntFrameCloneType(a%b, af))
	}
	return r.PushFrame(rpn.IntFrame(a%b, rpn.INTEGER_FRAME))
}
