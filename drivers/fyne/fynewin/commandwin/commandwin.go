package commandwin

import (
	"errors"
	"mattwach/rpngo/common/parse"
	"mattwach/rpngo/common/rpn"
	"mattwach/rpngo/common/startup"
	"sort"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type CommandWin struct {
	rpnInst   chan *rpn.RPN
	win       fyne.Window
	data      *widget.Entry
	result    *widget.Entry
	interrupt *startup.Interrupt
	buttons   map[string]string
	vBox      *fyne.Container
	updateFn  func()
}

func New(
	win fyne.Window,
	rpnInst chan *rpn.RPN,
	interrupt *startup.Interrupt,
	updateFn func()) *CommandWin {
	cw := &CommandWin{
		rpnInst:   rpnInst,
		win:       win,
		data:      widget.NewEntry(),
		result:    widget.NewEntry(),
		interrupt: interrupt,
		buttons:   make(map[string]string),
		updateFn:  updateFn,
	}
	cw.data.TextStyle = fyne.TextStyle{Monospace: true}
	cw.data.OnSubmitted = cw.runCommandWithString
	submitButton := widget.NewButton("run", cw.runCommand)
	entryBox := container.NewBorder(nil, nil, nil, submitButton, cw.data)

	cw.result.SetText("<results appear here>")

	cw.vBox = container.NewVBox(cw.result, entryBox)
	cw.layoutButtons(false)
	win.SetContent(cw.vBox)
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
	return "command"
}

func (cw *CommandWin) SetProp(name string, val rpn.Frame) error {
	if !strings.HasPrefix(name, "b.") {
		return rpn.ErrUnknownProperty
	}
	vstr := val.String(false)
	if vstr == "" {
		delete(cw.buttons, name)
	} else {
		cw.buttons[name] = val.String(false)
	}
	cw.layoutButtons(true)
	return nil
}

func (cw *CommandWin) GetProp(name string) (rpn.Frame, error) {
	command, ok := cw.buttons[name]
	if ok {
		return rpn.StringFrame(command, rpn.STRING_BRACE_FRAME), nil
	}
	return rpn.Frame{}, rpn.ErrUnknownProperty
}

func (cw *CommandWin) layoutButtons(removeExisting bool) {
	names := append([]string{"break"}, cw.ListProps()...)
	layout := container.NewGridWrap(maxBounds(names), widget.NewButton("break", cw.breakPressed))
	for _, name := range names[1:] {
		layout.Add(widget.NewButton(buttonName(name), func() {
			cw.runCommandWithString(cw.buttons[name])
		}))
	}
	if removeExisting {
		cw.vBox.Remove(cw.vBox.Objects[len(cw.vBox.Objects)-1])
	}
	cw.vBox.Add(layout)
}

func buttonName(name string) string {
	parts := strings.Split(name, ".")
	if len(parts) <= 3 {
		return parts[len(parts)-1]
	}
	return strings.Join(parts[2:len(parts)], ".")
}

const hpad = 15
const vpad = 10

func maxBounds(names []string) fyne.Size {
	var maxs fyne.Size
	for _, name := range names {
		s := fyne.MeasureText(name, theme.TextSize(), fyne.TextStyle{})
		if s.Width > maxs.Width {
			maxs.Width = s.Width
		}
		if s.Height > maxs.Height {
			maxs.Height = s.Height
		}
	}
	maxs.Width += hpad * 2
	maxs.Height += vpad * 2
	return maxs
}

func (cw *CommandWin) ListProps() []string {
	var buttonNames []string
	for name := range cw.buttons {
		buttonNames = append(buttonNames, name)
	}
	sort.Strings(buttonNames)
	return buttonNames
}

func (cw *CommandWin) runCommandWithString(s string) {

	// This needs to be executed in a separate go routine because 's'
	// might contain fyne API calls and we are already running within
	// the main fyne thread.
	go func() {
		select {
		case r := <-cw.rpnInst:
			defer func() {
				cw.rpnInst <- r
			}()

			fyne.Do(func() {
				cw.data.SetText("")
			})

			err := parse.Fields(s, r.Exec)
			s := ""
			if err != nil {
				if errors.Is(err, rpn.ErrExit) {
					fyne.CurrentApp().Quit()
				}
				s = "error: " + err.Error()
			} else if len(r.Frames) > 0 {
				s = r.Frames[len(r.Frames)-1].String(true)
			}

			fyne.Do(func() {
				cw.result.SetText(s)
			})
			cw.updateFn()
		default:
			fyne.DoAndWait(func() {
				cw.result.SetText("RPN is busy (break to stop)")
			})
		}
	}()
}

func (cw *CommandWin) runCommand() {
	cw.runCommandWithString(cw.data.Text)
}

func (cw *CommandWin) breakPressed() {
	cw.interrupt.Signal()
}
