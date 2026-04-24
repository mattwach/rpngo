package functions

import (
	"mattwach/rpngo/common/rpn"
	"runtime"
)

const heapstatsHelp = "Pushes annotated memory heap stats to the stack"

func heapstats(r *rpn.RPN) error {
	var ms runtime.MemStats
	runtime.ReadMemStats(&ms)

	f := rpn.IntFrame(int64(ms.Sys), rpn.INTEGER_FRAME)
	f.Annotate("`sys")
	if err := r.PushFrame(f); err != nil {
		return err
	}

	f = rpn.IntFrame(int64(ms.TotalAlloc), rpn.INTEGER_FRAME)
	f.Annotate("`total alloc")
	if err := r.PushFrame(f); err != nil {
		return err
	}

	f = rpn.IntFrame(int64(ms.HeapIdle), rpn.INTEGER_FRAME)
	f.Annotate("`free")
	if err := r.PushFrame(f); err != nil {
		return err
	}

	f = rpn.IntFrame(int64(ms.Frees), rpn.INTEGER_FRAME)
	f.Annotate("`frees")
	if err := r.PushFrame(f); err != nil {
		return err
	}

	return nil
}
