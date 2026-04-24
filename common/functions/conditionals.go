package functions

import (
	"mattwach/rpngo/common/rpn"
)

const greaterThanHelp = "Returns true if a > b, false otherwise"

func greaterThan(r *rpn.RPN) error {
	af, bf, err := r.Pop2Frames()
	if err != nil {
		return err
	}
	return r.PushFrame(rpn.BoolFrame(!af.IsLessThanOrEqual(bf)))
}

const greaterThanEqualHelp = "Returns true if a >= b, false otherwise"

func greaterThanEqual(r *rpn.RPN) error {
	af, bf, err := r.Pop2Frames()
	if err != nil {
		return err
	}
	return r.PushFrame(rpn.BoolFrame(!af.IsLessThan(bf)))
}

const lessThanHelp = "Returns true if a < b, false otherwise"

func lessThan(r *rpn.RPN) error {
	af, bf, err := r.Pop2Frames()
	if err != nil {
		return err
	}
	return r.PushFrame(rpn.BoolFrame(af.IsLessThan(bf)))
}

const lessThanEqualHelp = "Returns true if a <= b, false otherwise"

func lessThanEqual(r *rpn.RPN) error {
	af, bf, err := r.Pop2Frames()
	if err != nil {
		return err
	}
	return r.PushFrame(rpn.BoolFrame(af.IsLessThanOrEqual(bf)))
}

const equalHelp = "Returns true if a = b, false otherwise (approximate)"

func equal(r *rpn.RPN) error {
	af, bf, err := r.Pop2Frames()
	if err != nil {
		return err
	}
	return r.PushFrame(rpn.BoolFrame(af.IsEqual(bf)))
}

const notEqualHelp = "Returns true if a != b, false otherwise (approximate)"

func notEqual(r *rpn.RPN) error {
	af, bf, err := r.Pop2Frames()
	if err != nil {
		return err
	}
	return r.PushFrame(rpn.BoolFrame(!af.IsEqual(bf)))
}

const minHelp = "Pops two frames and repushes the minimum value.  Pushes $1 if the frames were equal."

func min(r *rpn.RPN) error {
	af, bf, err := r.Pop2Frames()
	if err != nil {
		return err
	}
	if af.IsLessThanOrEqual(bf) {
		return r.PushFrame(af)
	}
	return r.PushFrame(bf)
}

const maxHelp = "Pops two frames and repushes the maximum value.  Pushes $1 if the frames were equal."

func max(r *rpn.RPN) error {
	af, bf, err := r.Pop2Frames()
	if err != nil {
		return err
	}
	if af.IsLessThan(bf) {
		return r.PushFrame(bf)
	}
	return r.PushFrame(af)
}
