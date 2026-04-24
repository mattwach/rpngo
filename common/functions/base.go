package functions

import (
	"mattwach/rpngo/common/rpn"
)

func convert(r *rpn.RPN, t rpn.FrameType) error {
	f, err := r.PopFrame()
	if err != nil {
		return err
	}
	if f.IsBool() {
		if f.UnsafeBool() {
			return r.PushFrame(rpn.IntFrame(1, t))
		} else {
			return r.PushFrame(rpn.IntFrame(0, t))
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
	v, err := f.Int()
	if err != nil {
		return err
	}
	return r.PushFrame(rpn.IntFrame(v, t))
}

const intHelp = "Converts head element to an integer number"

func intFn(r *rpn.RPN) error {
	return convert(r, rpn.INTEGER_FRAME)
}

const hexHelp = "Converts head element to a hexidecimal number"

func hex(r *rpn.RPN) error {
	return convert(r, rpn.HEXIDECIMAL_FRAME)
}

const octHelp = "Converts head element to an octal number"

func oct(r *rpn.RPN) error {
	return convert(r, rpn.OCTAL_FRAME)
}

const binHelp = "Converts head element to a binary number"

func bin(r *rpn.RPN) error {
	return convert(r, rpn.BINARY_FRAME)
}

const strHelp = "Converts head element to a string"

func str(r *rpn.RPN) error {
	f, err := r.PopFrame()
	if err != nil {
		return err
	}
	return r.PushFrame(rpn.StringFrame(f.String(false), rpn.STRING_BRACE_FRAME))
}
