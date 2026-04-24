package tinygo

import (
	"machine"
	"mattwach/rpngo/common/rpn"
)

const resetHelp = "Resets the calculator."

func reset(r *rpn.RPN) error {
	machine.CPUReset()
	return nil
}

func Register(r *rpn.RPN) {
	r.Register("reset", reset, rpn.CatCore, resetHelp)
}
