package window

import (
	"mattwach/rpngo/common/fileops"
	"mattwach/rpngo/common/rpn"
	"mattwach/rpngo/common/window/tabcomplete"
)

type ReadlineTabComplete struct {
	tabc    tabcomplete.TabComplete
	rpnExec func(func(*rpn.RPN) error) error
}

func (tc *ReadlineTabComplete) Init(rpnExec func(func(*rpn.RPN) error) error, fs fileops.FileOpsDriver) {
	tc.tabc.Init(fs)
	tc.rpnExec = rpnExec
}

func (tc *ReadlineTabComplete) tabCompleteCallback(line string) []string {

	var words []string
	tc.rpnExec(func(r *rpn.RPN) error {
		words = tc.tabc.FindAllWords(r, line)
		return nil
	})
	return words
}
