// Provides a window.WindowWithProps interface with most of the methods stubbed out.
// This is done to provie a similar user experience with nrpn and picocalc
package window

import (
	"fmt"
	"mattwach/rpngo/common/drivers/posix/fs"
	"mattwach/rpngo/common/fileops"
	"mattwach/rpngo/common/parse"
	"mattwach/rpngo/common/rpn"
	"mattwach/rpngo/common/startup"
	"path/filepath"

	"github.com/chzyer/readline"
)

type ReadlineWindow struct {
	inst    *readline.Instance
	rpnInst chan *rpn.RPN
}

func (rlw *ReadlineWindow) Init(rpnInst chan *rpn.RPN) error {
	rlw.rpnInst = rpnInst
	histFile := startup.HistFile
	home, err := fileops.HomeDir()
	if err == nil {
		histFile = filepath.Join(home, startup.HistFile)
	}
	var tabc ReadlineTabComplete
	tabc.Init(rpnInst, &fs.FileOpsDriver{})
	completer := readline.NewPrefixCompleter(
		readline.PcItemDynamic(tabc.tabCompleteCallback),
	)
	rlw.inst, err = readline.NewEx(&readline.Config{
		Prompt:       "> ",
		HistoryFile:  histFile,
		AutoComplete: completer,
	})
	return nil
}

func (rlw *ReadlineWindow) Close() {
	rlw.inst.Close()
}

func (rlw *ReadlineWindow) ExecLine() error {
	line, err := rlw.inst.Readline()
	if err != nil {
		return err
	}
	r := <-rlw.rpnInst
	defer func() {
		rlw.rpnInst <- r
	}()
	err = parse.Fields(line, r.Exec)
	if err != nil {
		fmt.Printf("error: %v\n", err)
	}

	if len(r.Frames) > 0 {
		fmt.Println(r.Frames[len(r.Frames)-1].String(true))
	}
	return nil
}

func (rlw *ReadlineWindow) ResizeWindow(x, y, w, h int) error {
	return nil
}

func (rlw *ReadlineWindow) ShowBorder(screenw, screenh int) error {
	return nil
}

func (rlw *ReadlineWindow) WindowXY() (int, int) {
	return 0, 0
}

func (rlw *ReadlineWindow) WindowSize() (int, int) {
	return 0, 0
}

func (w *ReadlineWindow) Update(r *rpn.RPN, force bool) error {
	return nil
}

func (rlw *ReadlineWindow) Type() string {
	return "input"
}

func (rlw *ReadlineWindow) SetProp(name string, val rpn.Frame) error {
	return nil
}

func (rlw *ReadlineWindow) GetProp(name string) (rpn.Frame, error) {
	return rpn.Frame{}, nil
}

func (rlw *ReadlineWindow) ListProps() []string {
	return nil
}
