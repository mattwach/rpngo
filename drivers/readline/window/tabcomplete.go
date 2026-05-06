package window

import (
	"log"
	"mattwach/rpngo/common/fileops"
	"mattwach/rpngo/common/rpn"
	"mattwach/rpngo/common/window/tabcomplete"
)

type ReadlineTabComplete struct {
	tabc    tabcomplete.TabComplete
	rpnChan chan *rpn.RPN
}

func (tc *ReadlineTabComplete) Init(rpnChan chan *rpn.RPN, fs fileops.FileOpsDriver) {
	tc.tabc.Init(fs)
	tc.rpnChan = rpnChan
}

func (tc *ReadlineTabComplete) tabCompleteCallback(line string) []string {
	log.Printf("tab complete callback called with line: %s", line)
	select {
	case r := <-tc.rpnChan:
		defer func() {
			tc.rpnChan <- r
		}()
		return tc.tabc.FindAllWords(r, line)
	default:
		return nil
	}
}
