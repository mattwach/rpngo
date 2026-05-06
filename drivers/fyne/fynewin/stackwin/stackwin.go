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
	rpnInst    chan *rpn.RPN
	win        fyne.Window
	data       *widget.Entry
	round      int8
	roundEntry *customwidget.CustomEntry
}

func New(win fyne.Window, rpnInst chan *rpn.RPN) *StackWin {
	sw := &StackWin{
		rpnInst: rpnInst,
		win:     win,
		data:    widget.NewMultiLineEntry(),
		round:   -1,
	}
	sw.data.TextStyle = fyne.TextStyle{Monospace: true}
	scroll := container.NewScroll(sw.data)
	roundLabel := widget.NewLabel("round")
	sw.roundEntry = customwidget.NewCustomEntry(sw.updateRoundEntry, 30)
	bottom := container.NewHBox(roundLabel, sw.roundEntry)
	win.SetContent(container.NewBorder(nil, bottom, nil, nil, scroll))
	win.Resize(fyne.NewSize(240, 640))
	return sw
}

func (sw *StackWin) Update(r *rpn.RPN, force bool) error {
	fyne.DoAndWait(func() { sw.updateMainContext(r) })
	return nil
}

func (sw *StackWin) updateMainContext(r *rpn.RPN) {
	if len(r.Frames) > 0 {
		lines := make([]string, 0)
		for i, f := range r.Frames {
			lines = append(
				lines,
				fmt.Sprintf("%3d: %v", len(r.Frames)-i-1, f.RoundedString(sw.round)))
		}
		sw.data.SetText(strings.Join(lines, "\n"))
	} else {
		sw.data.SetText("stack empty")
	}
	sw.data.CursorColumn = 0
	sw.data.CursorRow = len(r.Frames)
	sw.roundEntry.SetText(fmt.Sprintf("%d", sw.round))
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
	select {
	case r := <-sw.rpnInst:
		defer func() {
			sw.rpnInst <- r
		}()
		val, err := strconv.ParseInt(s, 10, 64)
		if err == nil {
			sw.SetProp("round", rpn.IntFrame(val, rpn.INTEGER_FRAME))
		}
		sw.updateMainContext(r)
	default:
		// do nothing
	}
}
