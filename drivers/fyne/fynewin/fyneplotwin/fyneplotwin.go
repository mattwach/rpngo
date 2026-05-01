package fyneplotwin

import (
	"image/color"
	"mattwach/rpngo/common/rpn"
	"mattwach/rpngo/common/window/plotwin"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"github.com/fogleman/gg"
)

// FynePlotWin holds the context for a stack window.
// Important, RPN is owned by the readline goroutine thus should be accessed
// with care.  This means that putting a pointer to it in this struct is
// probably starting down a bad path.
type FynePlotWin struct {
	win           fyne.Window
	canvas        interactiveImage
	autoYCheckbox *widget.Check
	color         color.RGBA
	clearFirst    bool
}

// New is expected to be called in the context of the main thread.
func New(win fyne.Window, r chan *rpn.RPN, parent *plotwin.PixelPlotWindow) *FynePlotWin {
	pw := &FynePlotWin{
		win:        win,
		clearFirst: true,
	}

	pw.win.Resize(fyne.NewSize(1024, 768))

	winSize := pw.win.Canvas().Size()
	pw.canvas.rpnInstance = r
	pw.canvas.ggimg = gg.NewContext(int(winSize.Width), int(winSize.Height))
	pw.canvas.image = canvas.NewImageFromImage(pw.canvas.ggimg.Image())
	pw.canvas.parent = parent
	pw.autoYCheckbox = widget.NewCheck("Auto Y", pw.autoYClicked)
	bottom := container.NewHBox(pw.autoYCheckbox)
	pw.win.SetContent(container.NewBorder(nil, bottom, nil, nil, &pw.canvas))
	return pw
}

// Color implements the window/PixelWindow interface.
func (pw *FynePlotWin) Color(c color.RGBA) {
	pw.color = c
	pw.canvas.ggimg.SetColor(c)
}

// FilledRect implements the window/PixelWindow interface.
func (pw *FynePlotWin) FilledRect(x, y, w, h int) {
	pw.clearIfNeeded()
	pw.canvas.ggimg.DrawRectangle(float64(x), float64(y), float64(w), float64(h))
	pw.canvas.ggimg.Fill()
}

// HLine implements the window/PixelWindow interface.
func (pw *FynePlotWin) HLine(x, y, w int) {
	pw.clearIfNeeded()
	pw.canvas.ggimg.DrawLine(float64(x), float64(y), float64(x+w), float64(y))
	pw.canvas.ggimg.Stroke()
}

// PixelSize implements the window/PixelWindow interface.
func (pw *FynePlotWin) PixelSize() (int, int) {
	return pw.WindowSize()
}

// Refresh implements the window/PixelWindow interface.
// It can be called from the fyne or readline gorroutine, with
// pw.canvas.inMainContext indicating which one.
func (pw *FynePlotWin) Refresh() {
	fn := func() {
		pw.canvas.image.Refresh()
		pw.autoYCheckbox.SetChecked(pw.canvas.parent.Common.AutoY)
	}
	if pw.canvas.inMainContext {
		fn()
	} else {
		fyne.DoAndWait(fn)
	}
	// clear the next time drawing is started
	pw.clearFirst = true
}

// ResizeWindow implements the window/WindowBase interface.
func (pw *FynePlotWin) ResizeWindow(x, y, w, h int) error {
	// not supported in fyne
	return nil
}

// SetPoint implements the window/PixelWindow interface.
func (pw *FynePlotWin) SetPoint(x, y int) {
	pw.clearIfNeeded()
	pw.canvas.ggimg.SetPixel(x, y)
}

// ShowBorder implements the window/WindowBase interface.  It is ignored by
// the fyne implementation.
func (pw *FynePlotWin) ShowBorder(sw, sh int) error {
	return nil
}

// Text implements the window/PixelWindow interface.
func (pw *FynePlotWin) Text(s string, x, y int) {
	pw.clearIfNeeded()
	pw.canvas.ggimg.DrawString(s, float64(x), float64(y))
}

// VLine implements the window/PixelWindow interface.
func (pw *FynePlotWin) VLine(x, y, h int) {
	pw.clearIfNeeded()
	pw.canvas.ggimg.DrawLine(float64(x), float64(y), float64(x), float64(y+h))
	pw.canvas.ggimg.Stroke()
}

// WindowXY implements the window/WindowBase interface.
func (pw *FynePlotWin) WindowXY() (int, int) {
	return 0, 0
}

// WindowSize implements the window/WindowBase interface.
func (pw *FynePlotWin) WindowSize() (int, int) {
	return pw.canvas.ggimg.Width(), pw.canvas.ggimg.Height()
}

// autoYClicked is called when the autoY checkbox is toggled.
func (pw *FynePlotWin) autoYClicked(on bool) {
	pw.canvas.inMainContext = true
	defer func() {
		pw.canvas.inMainContext = false
	}()
	select {
	case r := <-pw.canvas.rpnInstance:
		defer func() {
			pw.canvas.rpnInstance <- r
		}()
		pw.canvas.parent.Common.AutoY = on
		pw.canvas.parent.Update(r, true)
	default:
		// no action if the rpn instance is not available
	}
}

// clearIfNeeded clears the canvas if clearFirst is true. It is used as a
// workaround where we want to clear the canvas after a refresh.
func (pw *FynePlotWin) clearIfNeeded() {
	if !pw.clearFirst {
		return
	}
	pw.canvas.ggimg.SetRGB(0, 0, 0)
	pw.canvas.ggimg.Clear()
	pw.canvas.ggimg.SetColor(pw.color)
	pw.clearFirst = false
}
