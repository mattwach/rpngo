package stackwin

import (
	"fmt"
	"mattwach/rpngo/common/rpn"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// StackWin holds teh context for a stack window.
// Important, RPN is owned by the readline goroutine thus should be accessed
// with care.  This means that putting a pointer to it in this struct is
// probably starting down a bad path.
type StackWin struct {
	win  fyne.Window
	data *widget.Entry
}

func New(win fyne.Window, r *rpn.RPN) *StackWin {
	sw := &StackWin{
		win:  win,
		data: widget.NewMultiLineEntry(),
	}
	sw.data.TextStyle = fyne.TextStyle{Monospace: true}
	scroll := container.NewScroll(sw.data)
	win.SetContent(scroll)
	sw.Update(r)
	win.Resize(fyne.NewSize(240, 640))
	return sw
}

func (sw *StackWin) Update(r *rpn.RPN) {
	if len(r.Frames) > 0 {
		lines := make([]string, 0)
		for i, f := range r.Frames {
			lines = append(
				lines,
				fmt.Sprintf("%3d: %v", len(r.Frames)-i-1, f.String(true)))
		}
		sw.data.SetText(strings.Join(lines, "\n"))
	} else {
		sw.data.SetText("stack empty")
	}
	sw.data.CursorColumn = 0
	sw.data.CursorRow = len(r.Frames)
}
