package stackwin

import (
	"fmt"
	"mattwach/rpngo/common/rpn"
	"mattwach/rpngo/drivers/fyne/fynewin/customwidget"
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// StackWin holds the context for a stack window.
// Important, RPN is owned by the readline goroutine thus should be accessed
// with care.  This means that putting a pointer to it in this struct is
// probably starting down a bad path.
type StackWin struct {
	rpn        *rpn.RPN
	win        fyne.Window
	clipboard  fyne.Clipboard
	data       *widget.Entry
	round      int8
	roundEntry *customwidget.CustomEntry
	copyButton *widget.Button
	fullUpdate func()
}

func New(win fyne.Window, clipboard fyne.Clipboard, r *rpn.RPN, fullUpdate func()) *StackWin {
	sw := &StackWin{
		rpn:        r,
		win:        win,
		clipboard:  clipboard,
		data:       widget.NewMultiLineEntry(),
		round:      -1,
		fullUpdate: fullUpdate,
	}
	sw.data.TextStyle = fyne.TextStyle{Monospace: true}
	scroll := container.NewScroll(sw.data)
	roundLabel := widget.NewLabel("round")
	sw.roundEntry = customwidget.NewCustomEntry(sw.updateRoundEntry, 30)
	clearButton := widget.NewButton("clear", sw.clearStack)
	sw.copyButton = widget.NewButton("copy", sw.copyToClipboard)
	bottom := container.NewHBox(roundLabel, sw.roundEntry, clearButton, sw.copyButton)
	win.SetContent(container.NewBorder(nil, bottom, nil, nil, scroll))
	win.Resize(fyne.NewSize(240, 640))
	sw.Update(r, true)
	return sw
}

func (sw *StackWin) Update(r *rpn.RPN, force bool) error {
	if len(r.Frames) > 0 {
		lines := make([]string, 0)
		for i, f := range r.Frames {
			lines = append(
				lines,
				fmt.Sprintf("%3d: %v", len(r.Frames)-i-1, f.RoundedString(sw.round, true)))
		}
		sw.data.SetText(strings.Join(lines, "\n"))
	} else {
		sw.data.SetText("stack empty")
	}
	sw.data.CursorColumn = 0
	sw.data.CursorRow = len(r.Frames)
	sw.roundEntry.SetText(fmt.Sprintf("%d", sw.round))
	if len(r.Frames) > 0 {
		sw.copyButton.Enable()
	} else {
		sw.copyButton.Disable()
	}
	return nil
}

func (sw *StackWin) ResizeWindow(x, y, w, h int) error {
	// Not used in fyne
	return nil
}

func (sw *StackWin) ShowBorder(w, h int) error {
	// Not used in fyne
	return nil
}

func (sw *StackWin) WindowXY() (int, int) {
	// Not used in fyne
	return 0, 0
}

func (sw *StackWin) WindowSize() (int, int) {
	s := sw.data.Size()
	return int(s.Width), int(s.Height)
}

func (sw *StackWin) Type() string {
	return "stack"
}

func (sw *StackWin) SetProp(name string, val rpn.Frame) error {
	switch name {
	case "round":
		v, err := val.BoundedInt(-1, 10)
		if err != nil {
			return err
		}
		sw.round = int8(v)
		return nil
	default:
		return rpn.ErrUnknownProperty
	}
}

func (sw *StackWin) GetProp(name string) (rpn.Frame, error) {
	switch name {
	case "round":
		return rpn.IntFrame(int64(sw.round), rpn.INTEGER_FRAME), nil
	default:
		return rpn.Frame{}, rpn.ErrUnknownProperty
	}
}

var props = []string{"round"}

// Lists props.  Do not nodify return value.
func (sw *StackWin) ListProps() []string {
	return props
}

func (sw *StackWin) updateRoundEntry(s string) {
	val, err := strconv.ParseInt(s, 10, 64)
	if err == nil {
		sw.SetProp("round", rpn.IntFrame(val, rpn.INTEGER_FRAME))
	}
	sw.Update(sw.rpn, true)
}

func (sw *StackWin) clearStack() {
	sw.rpn.Frames = sw.rpn.Frames[:0]
	sw.fullUpdate()
}

func (sw *StackWin) copyToClipboard() {
	if len(sw.rpn.Frames) == 0 {
		return
	}
	val := sw.rpn.Frames[len(sw.rpn.Frames)-1].RoundedString(sw.round, false)
	sw.clipboard.SetContent(val)
}
