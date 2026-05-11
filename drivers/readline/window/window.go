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
	"github.com/fatih/color"
)

type ReadlineWindow struct {
	inst       *readline.Instance
	showFrames int
	autofn     []string
	rpnExec    func(func(*rpn.RPN) error) error
}

func (rlw *ReadlineWindow) Init(rpnExec func(func(*rpn.RPN) error) error) error {
	rlw.rpnExec = rpnExec
	rlw.showFrames = 1
	histFile := startup.HistFile
	home, err := fileops.HomeDir()
	if err == nil {
		histFile = filepath.Join(home, startup.HistFile)
	}
	var tabc ReadlineTabComplete
	tabc.Init(rpnExec, &fs.FileOpsDriver{})
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
	rlw.autoFn()
	line, err := rlw.inst.Readline()
	if err != nil {
		return err
	}
	rlw.rpnExec(func(r *rpn.RPN) error {
		color.Set(color.FgYellow)
		err = parse.Fields(line, r.Exec)
		color.Unset()
		if err != nil {
			color.Red("%v\n", err)
		} else {
			rlw.printFrames(r)
		}
		return nil
	})
	return nil
}

func (rlw *ReadlineWindow) autoFn() {
	if len(rlw.autofn) == 0 {
		return
	}
	err := rlw.rpnExec(func(r *rpn.RPN) error {
		return r.ExecSlice(rlw.autofn)
	})

	if err != nil {
		color.Red("autofn error: %v\n", err)
	}
}

func (rlw *ReadlineWindow) printFrames(r *rpn.RPN) {
	count := len(r.Frames)
	if rlw.showFrames < count {
		count = rlw.showFrames
	}
	if count == 0 {
		return
	}
	color.Set(color.FgCyan)
	for i := 0; i < count; i++ {
		f := r.Frames[len(r.Frames)-count+i]
		fmt.Println(f.String(true))
	}
	color.Unset()
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
