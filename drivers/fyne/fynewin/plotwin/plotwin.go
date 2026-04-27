package plotwin

import (
	"mattwach/rpngo/common/rpn"

	"fyne.io/fyne/v2"
)

// PlotWin holds the context for a stack window.
// Important, RPN is owned by the readline goroutine thus should be accessed
// with care.  This means that putting a pointer to it in this struct is
// probably starting down a bad path.
type PlotWin struct {
	win fyne.Window
}

func New(win fyne.Window) *PlotWin {
	pw := &PlotWin{
		win: win,
	}
	win.Resize(fyne.NewSize(1280, 720))
	return pw
}

func (sw *PlotWin) Update(r *rpn.RPN) {
}
