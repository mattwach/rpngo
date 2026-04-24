// Package input contains dirrent io interfaces that abstract the actual implmentation
// from the API.
package input

import (
	"errors"
	"mattwach/rpngo/common/fileops"
	"mattwach/rpngo/common/key"
	"mattwach/rpngo/common/parse"
	"mattwach/rpngo/common/rpn"
	"mattwach/rpngo/common/window"
	"strings"
)

var ErrExit = errors.New("exit")

// Input gets input from the keyboard/keypad
type Input interface {
	GetChar() (key.Key, error)
}

type InputWindow struct {
	input      Input
	txtb       window.TextBuffer
	gl         *getLine
	firstInput bool
	showFrames int
	autofn     []string
}

func (iw *InputWindow) Init(input Input, txtw window.TextWindow, r *rpn.RPN, fs fileops.FileOpsDriver, scrollbytes int) {
	iw.input = input
	// Make this >0 when we are ready to try scrolling.
	iw.txtb.Init(txtw, scrollbytes)
	iw.gl = initGetLine(input, &iw.txtb, fs)
	iw.firstInput = true
	iw.showFrames = 1
	r.Print = iw.Print
	r.Input = iw.Input
	r.Register("edit", iw.edit, rpn.CatProg, editHelp)
	r.Register("editf", iw.editFile, rpn.CatProg, editFileHelp)
	r.Register("histl", iw.gl.histLoad, rpn.CatStatus, histLoadHelp)
	r.Register("hists", iw.gl.histSave, rpn.CatStatus, histSaveHelp)
}

func (iw *InputWindow) Print(msg string) {
	iw.txtb.Print(msg, true)
}

func (iw *InputWindow) Input(r *rpn.RPN) (string, error) {
	return iw.gl.get(r)
}

func (iw *InputWindow) Update(r *rpn.RPN, unusedForce bool) error {
	iw.txtb.CheckSize()
	if iw.firstInput {
		if err := iw.firstRun(r); err != nil {
			return err
		}
	}
	iw.txtb.TextColor(window.White)
	r.TextWidth = iw.txtb.Txtw.TextWidth()
	if len(iw.autofn) != 0 {
		if err := r.ExecSlice(iw.autofn); err != nil {
			iw.txtb.TextColor(window.Red)
			r.Println(err.Error())
			iw.txtb.TextColor(window.White)
		}
	}
	iw.txtb.Print("> ", true)
	line, err := iw.gl.get(r)
	iw.txtb.TextColor(window.Yellow)
	if err != nil {
		iw.txtb.PrintErr(err, true)
		return nil
	}
	line = strings.TrimSpace(line)
	if len(line) == 0 {
		return nil
	}
	if line == "exit" {
		return ErrExit
	}
	action, err := parseLine(r, line)
	if err != nil {
		iw.txtb.PrintErr(err, true)
		return nil
	}
	if action {
		iw.printFrames(r)
	}
	iw.txtb.Update()
	return nil
}

func (iw *InputWindow) printFrames(r *rpn.RPN) {
	count := len(r.Frames)
	if iw.showFrames < count {
		count = iw.showFrames
	}
	if count == 0 {
		return
	}
	iw.txtb.TextColor(window.Cyan)
	for i := 0; i < count; i++ {
		f := r.Frames[len(r.Frames)-count+i]
		iw.txtb.Print(f.String(true), false)
		iw.txtb.Write('\n', false)
	}
}

func (iw *InputWindow) firstRun(r *rpn.RPN) error {
	iw.firstInput = false
	r.RegisterConceptHelp(map[string]string{"exiting": "Enter exit or type Ctrl-D to exit the program"})
	iw.txtb.TextColor(window.Cyan)
	iw.txtb.Print("Enter ? to list help topics, topic? for more details\n", true)
	iw.txtb.TextColor(window.White)
	return nil
}

func (iw *InputWindow) ResizeWindow(x, y, w, h int) error {
	err := iw.txtb.Txtw.ResizeWindow(x, y, w, h)
	if err != nil {
		return err
	}
	iw.txtb.CheckSize()
	return nil
}

func (iw *InputWindow) ShowBorder(screenw, screenh int) error {
	return iw.txtb.Txtw.ShowBorder(screenw, screenh)
}

func (iw *InputWindow) WindowXY() (int, int) {
	return iw.txtb.Txtw.WindowXY()
}

func (iw *InputWindow) WindowSize() (int, int) {
	return iw.txtb.Txtw.WindowSize()
}

func (iw *InputWindow) Type() string {
	return "input"
}

func parseLine(r *rpn.RPN, line string) (bool, error) {
	if err := parse.Fields(line, r.Exec); err != nil {
		return false, err
	}
	return true, nil
}
