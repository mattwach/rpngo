package fyneplotwin

import (
	"fmt"
	"image/color"
	"mattwach/rpngo/common/rpn"
	"mattwach/rpngo/common/window/plotwin"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"

	"github.com/fogleman/gg"
)

type interactiveImage struct {
	widget.BaseWidget
	inMainContext bool
	ggimg         *gg.Context
	image         *canvas.Image
	parent        *plotwin.PixelPlotWindow
	rpnInstance   chan *rpn.RPN
}

// MouseMoved captures continuous movement over the image
func (i *interactiveImage) MouseMoved(ev *desktop.MouseEvent) {
	i.inMainContext = true
	defer func() {
		i.inMainContext = false
	}()
	s := fmt.Sprintf("(%.2f, %.2f)", ev.Position.X, ev.Position.Y)
	i.ggimg.SetRGB(0, 0, 1)
	i.ggimg.DrawRectangle(50, 30, 150, 25)
	i.ggimg.Fill()
	i.ggimg.SetRGB(1, 1, 1)
	i.ggimg.DrawString(s, 50, 50)
	i.image.Refresh()
}

func (i *interactiveImage) MouseIn(ev *desktop.MouseEvent) {}
func (i *interactiveImage) MouseOut()                      {}

func (i *interactiveImage) Resize(size fyne.Size) {
	i.inMainContext = true
	defer func() {
		i.inMainContext = false
	}()
	i.BaseWidget.Resize(size)
	if i.parent == nil {
		return
	}
	select {
	case r := <-i.rpnInstance:
		defer func() {
			i.rpnInstance <- r
		}()
		i.ggimg = gg.NewContext(int(size.Width), int(size.Height))
		i.image.Image = i.ggimg.Image()
		i.image.Resize(size)
		i.parent.ResizeWindow(0, 0, int(size.Width), int(size.Height))
		i.parent.Update(r, false)
	default:
		// no action if the rpn instance is not available
	}
}

func (i *interactiveImage) CreateRenderer() fyne.WidgetRenderer {
	i.inMainContext = true
	defer func() {
		i.inMainContext = false
	}()
	return widget.NewSimpleRenderer(i.image)
}

// FynePlotWin holds the context for a stack window.
// Important, RPN is owned by the readline goroutine thus should be accessed
// with care.  This means that putting a pointer to it in this struct is
// probably starting down a bad path.
type FynePlotWin struct {
	win        fyne.Window
	canvas     interactiveImage
	color      color.RGBA
	clearFirst bool
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
	pw.win.SetContent(&pw.canvas)
	pw.canvas.parent = parent
	return pw
}

func (pw *FynePlotWin) clearIfNeeded() {
	if !pw.clearFirst {
		return
	}
	pw.canvas.ggimg.SetRGB(0, 0, 0)
	pw.canvas.ggimg.Clear()
	pw.canvas.ggimg.SetColor(pw.color)
	pw.clearFirst = false
}

// We conform to the window/PixelWindow Interface so that common logic can
// update the plot
func (pw *FynePlotWin) ResizeWindow(x, y, w, h int) error {
	// not supported in fyne
	return nil
}

func (pw *FynePlotWin) ShowBorder(sw, sh int) error {
	// not supported in fyne
	return nil
}

func (pw *FynePlotWin) WindowXY() (int, int) {
	return 0, 0
}

func (pw *FynePlotWin) WindowSize() (int, int) {
	return pw.canvas.ggimg.Width(), pw.canvas.ggimg.Height()
}

// Refresh is expected to be called outside of the main fyne thread.
func (pw *FynePlotWin) Refresh() {
	if pw.canvas.inMainContext {
		// we got here from a fyne ui call
		pw.canvas.image.Refresh()
	} else {
		fyne.DoAndWait(func() {
			pw.canvas.image.Refresh()
		})
	}
	// clear the next time drawing is started
	pw.clearFirst = true
}

func (pw *FynePlotWin) PixelSize() (int, int) {
	return pw.WindowSize()
}

func (pw *FynePlotWin) Color(c color.RGBA) {
	pw.color = c
	pw.canvas.ggimg.SetColor(c)
}

func (pw *FynePlotWin) SetPoint(x, y int) {
	pw.clearIfNeeded()
	pw.canvas.ggimg.SetPixel(x, y)
}

func (pw *FynePlotWin) HLine(x, y, w int) {
	pw.clearIfNeeded()
	pw.canvas.ggimg.DrawLine(float64(x), float64(y), float64(x+w), float64(y))
	pw.canvas.ggimg.Stroke()
}

func (pw *FynePlotWin) VLine(x, y, h int) {
	pw.clearIfNeeded()
	pw.canvas.ggimg.DrawLine(float64(x), float64(y), float64(x), float64(y+h))
	pw.canvas.ggimg.Stroke()
}

func (pw *FynePlotWin) FilledRect(x, y, w, h int) {
	pw.clearIfNeeded()
	pw.canvas.ggimg.DrawRectangle(float64(x), float64(y), float64(w), float64(h))
	pw.canvas.ggimg.Fill()
}

func (pw *FynePlotWin) Text(s string, x, y int) {
	pw.clearIfNeeded()
	pw.canvas.ggimg.DrawString(s, float64(x), float64(y))
}
