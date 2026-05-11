package commandwin

import (
	"mattwach/rpngo/common/parse"
	"mattwach/rpngo/common/rpn"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type CommandWin struct {
	rpn    *rpn.RPN
	win    fyne.Window
	data   *widget.Entry
	result *widget.Entry
}

func New(win fyne.Window, r *rpn.RPN) *CommandWin {
	cw := &CommandWin{
		rpn:    r,
		win:    win,
		data:   widget.NewEntry(),
		result: widget.NewEntry(),
	}
	cw.data.TextStyle = fyne.TextStyle{Monospace: true}
	cw.data.OnSubmitted = cw.runCommandWithString
	submitButton := widget.NewButton("run", cw.runCommand)
	entryBox := container.NewBorder(nil, nil, nil, submitButton, cw.data)

	cw.result.SetPlaceHolder("<results appear here>")

	breakButton := widget.NewButton("break", cw.breakPressed)
	commandBox := container.NewHBox(breakButton)

	win.SetContent(container.NewVBox(cw.result, entryBox, commandBox))
	win.Resize(fyne.NewSize(640, 0))
	return cw
}

func (cw *CommandWin) Update(r *rpn.RPN, force bool) error {
	return nil
}

func (cw *CommandWin) ResizeWindow(x, y, w, h int) error {
	// Not used in fyne
	return nil
}

func (cw *CommandWin) ShowBorder(w, h int) error {
	// Not used in fyne
	return nil
}

func (cw *CommandWin) WindowXY() (int, int) {
	// Not used in fyne
	return 0, 0
}

func (cw *CommandWin) WindowSize() (int, int) {
	s := cw.win.Canvas().Size()
	return int(s.Width), int(s.Height)
}

func (cw *CommandWin) Type() string {
	return "stack"
}

func (cw *CommandWin) SetProp(name string, val rpn.Frame) error {
	return rpn.ErrUnknownProperty
}

func (cw *CommandWin) GetProp(name string) (rpn.Frame, error) {
	return rpn.Frame{}, rpn.ErrUnknownProperty
}

func (cw *CommandWin) ListProps() []string {
	return nil
}

func (cw *CommandWin) runCommandWithString(s string) {
	err := parse.Fields(s, cw.rpn.Exec)
	if err != nil {
		cw.result.SetText("error: " + err.Error())
	} else if len(cw.rpn.Frames) > 0 {
		cw.result.SetText(cw.rpn.Frames[len(cw.rpn.Frames)-1].String(true))
	} else {
		cw.result.SetText("<stack empty>")
	}
}

func (cw *CommandWin) runCommand() {
	cw.runCommandWithString(cw.data.Text)
}

func (cw *CommandWin) breakPressed() {

}
