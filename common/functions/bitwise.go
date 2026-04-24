package functions

import (
	"mattwach/rpngo/common/rpn"
)

const andHelp = "Performs a logical AND operation"

func and(r *rpn.RPN) error {
	return binaryOp(r, func(a, b int64) int64 { return a & b })
}

const orHelp = "Performs a logical OR operation"

func or(r *rpn.RPN) error {
	return binaryOp(r, func(a, b int64) int64 { return a | b })
}

const xorHelp = "Performs a logical XOR operation"

func xor(r *rpn.RPN) error {
	return binaryOp(r, func(a, b int64) int64 { return a ^ b })
}

const shiftLeftHelp = "Performs a logical shift left operation"

func shiftLeft(r *rpn.RPN) error {
	return binaryOp(r, func(a, b int64) int64 { return a << b })
}

const shiftRightHelp = "Performs a logical shift right operation"

func shiftRight(r *rpn.RPN) error {
	return binaryOp(r, func(a, b int64) int64 { return a >> b })
}

func binaryOp(r *rpn.RPN, fn func(a, b int64) int64) error {
	af, bf, err := r.Pop2Frames()
	if af.IsBool() {
		return binaryBoolOp(r, fn, af, bf)
	}
	a, err := af.Int()
	if err != nil {
		return err
	}
	b, err := bf.Int()
	if err != nil {
		return err
	}
	return r.PushFrame(rpn.IntFrameCloneType(fn(a, b), af))
}

func binaryBoolOp(r *rpn.RPN, fn func(a, b int64) int64, af, bf rpn.Frame) error {
	ab, err := af.Bool()
	if err != nil {
		return err
	}
	bb, err := bf.Bool()
	if err != nil {
		return err
	}
	var a int64 = 0
	var b int64 = 0
	if ab {
		a = 1
	}
	if bb {
		b = 1
	}
	v := fn(a, b)
	if v != 0 {
		return r.PushFrame(rpn.BoolFrame(true))
	}
	return r.PushFrame(rpn.BoolFrame(false))
}
