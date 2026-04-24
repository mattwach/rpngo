package functions

import "mattwach/rpngo/common/rpn"

const keepHelp = "Keeps up to n elements from the head of the stack."

func keep(r *rpn.RPN) error {
	f, err := r.PopFrame()
	if err != nil {
		return err
	}
	n, err := f.Int()
	if err != nil {
		return err
	}
	if n < 0 {
		return rpn.ErrIllegalValue
	}
	if n == 0 {
		r.Frames = r.Frames[:0]
		return nil
	}
	if int(n) >= len(r.Frames) {
		return nil
	}
	copy(r.Frames, r.Frames[len(r.Frames)-int(n):])
	r.Frames = r.Frames[:n]
	return nil
}
