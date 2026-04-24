package functions

import (
	"math/cmplx"
	"mattwach/rpngo/common/rpn"
)

const sinHelp = "takes the sine of a number"

func sin(r *rpn.RPN) error {
	af, err := r.PopFrame()
	if err != nil {
		return err
	}
	a, err := af.Complex()
	if err != nil {
		return err
	}
	return r.PushFrame(rpn.ComplexFrameWithType(cmplx.Sin(r.ToRadians(a)), af.Type()))
}

const asinHelp = "takes the inverse sine of a number"

func asin(r *rpn.RPN) error {
	af, err := r.PopFrame()
	if err != nil {
		return err
	}
	a, err := af.Complex()
	if err != nil {
		return err
	}
	return r.PushFrame(r.FromRadians(cmplx.Asin(a), af))
}

const cosHelp = "takes the cosine of a number"

func cos(r *rpn.RPN) error {
	af, err := r.PopFrame()
	if err != nil {
		return err
	}
	a, err := af.Complex()
	if err != nil {
		return err
	}
	return r.PushFrame(rpn.ComplexFrameWithType(cmplx.Cos(r.ToRadians(a)), af.Type()))
}

const acosHelp = "takes the inverse cosine of a number"

func acos(r *rpn.RPN) error {
	af, err := r.PopFrame()
	if err != nil {
		return err
	}
	a, err := af.Complex()
	if err != nil {
		return err
	}
	return r.PushFrame(r.FromRadians(cmplx.Acos(a), af))
}

const tanHelp = "takes the tangent of a number"

func tan(r *rpn.RPN) error {
	af, err := r.PopFrame()
	if err != nil {
		return err
	}
	a, err := af.Complex()
	if err != nil {
		return err
	}
	return r.PushFrame(rpn.ComplexFrameWithType(cmplx.Tan(r.ToRadians(a)), af.Type()))
}

const atanHelp = "takes the inverse tangent of a number"

func atan(r *rpn.RPN) error {
	af, err := r.PopFrame()
	if err != nil {
		return err
	}
	a, err := af.Complex()
	if err != nil {
		return err
	}
	return r.PushFrame(r.FromRadians(cmplx.Atan(a), af))
}
