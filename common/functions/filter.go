package functions

import (
	"mattwach/rpngo/common/elog"
	"mattwach/rpngo/common/parse"
	"mattwach/rpngo/common/rpn"
)

const filterHelp = "Copies each value from the bottom to the top of the stack " +
	"to the top of the stack.  Calling filter function after each one.  Discards the " +
	"original values, leaving only the new ones.\n" +
	"\n" +
	"Examples:\n" +
	"{2 *} filter  # multiply every stack value by 2\n" +
	"{$0 100 >= {0/} if} filter  # keep only values < 100\n" +
	"\n" +
	"See Also: filterm, filtern, filtermn"

func filter(r *rpn.RPN) error {
	fn, err := r.PopFrame()
	if err != nil {
		return err
	}
	endIdx := len(r.Frames)
	if endIdx == 0 {
		return nil
	}
	elog.Heap("alloc: functions/filter.go:117: fields := make([]string, 0, 16)")
	fields := make([]string, 0, 16) // object allocated on the heap: escapes at line 117
	addField := func(t string) error {
		fields = append(fields, t)
		return nil
	}
	if err := parse.Fields(fn.String(false), addField); err != nil {
		return err
	}
	for i := 0; i < endIdx; i++ {
		if i >= len(r.Frames) {
			return rpn.ErrNotEnoughStackFrames
		}
		if err := r.PushFrame(r.Frames[i]); err != nil {
			return err
		}
		if err := r.ExecSlice(fields); err != nil {
			return err
		}
	}
	newSize := len(r.Frames) - endIdx
	if newSize <= 0 {
		r.Frames = r.Frames[:0]
	} else {
		copy(r.Frames, r.Frames[endIdx:])
		r.Frames = r.Frames[:newSize]
	}
	return nil
}
