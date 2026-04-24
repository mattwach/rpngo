package functions

import (
	"mattwach/rpngo/common/rpn"
	"sort"
)

const sortHelp = "Sorts values on the stack.  Uses the " +
	"< conditional as the baseline for sorting."

type sortInterface struct {
	data []rpn.Frame
}

func (si *sortInterface) Len() int {
	return len(si.data)
}

func (si *sortInterface) Less(i, j int) bool {
	return si.data[i].IsLessThan(si.data[j])
}

func (si *sortInterface) Swap(i, j int) {
	si.data[i], si.data[j] = si.data[j], si.data[i]
}

func sortFn(r *rpn.RPN) error {
	if len(r.Frames) == 0 {
		return nil
	}
	si := sortInterface{data: r.Frames}
	sort.Sort(&si)
	return nil
}
