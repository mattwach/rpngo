package functions

import "mattwach/rpngo/common/rpn"

const reverseHelp = "Reverses the stack"

func reverse(r *rpn.RPN) error {
	if len(r.Frames) <= 1 {
		return nil
	}
	mid := len(r.Frames) / 2
	for i := 0; i < mid; i++ {
		end := len(r.Frames) - i - 1
		r.Frames[i], r.Frames[end] = r.Frames[end], r.Frames[i]
	}
	return nil
}
