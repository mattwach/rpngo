package functions

import "mattwach/rpngo/common/rpn"

const delHelp = "Removes up to n elements from the head of the stack."

func del(r *rpn.RPN) error {
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
		return nil
	}
	if int(n) >= len(r.Frames) {
		n = int64(len(r.Frames))
	}
	r.Frames = r.Frames[:len(r.Frames)-int(n)]
	return nil
}
