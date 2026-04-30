package fyneplotwin

import (
	"fmt"
	"image/color"
	"mattwach/rpngo/common/rpn"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"

	"github.com/fogleman/gg"
)

type interactiveImage struct {
	widget.BaseWidget
	ggimg *gg.Context
	image *canvas.Image
}

// MouseMoved captures continuous movement over the image
func (i *interactiveImage) MouseMoved(ev *desktop.MouseEvent) {
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

func (i *interactiveImage) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(i.image)
}

// PlotWin holds the context for a stack window.
// Important, RPN is owned by the readline goroutine thus should be accessed
// with care.  This means that putting a pointer to it in this struct is
// probably starting down a bad path.
type PlotWin struct {
	win        fyne.Window
	canvas     interactiveImage
	color      color.RGBA
	clearFirst bool
}

// New is expected to be called in the context of the main thread.
func New(win fyne.Window, r *rpn.RPN) *PlotWin {
	pw := &PlotWin{
		win:        win,
		clearFirst: true,
	}

	pw.canvas.ggimg = gg.NewContext(1024, 768)
	pw.canvas.image = canvas.NewImageFromImage(pw.canvas.ggimg.Image())
	pw.win.SetContent(container.NewStack(&pw.canvas))
	pw.win.Resize(fyne.NewSize(1024, 768))
	return pw
}

func (pw *PlotWin) clearIfNeeded() {
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
func (pw *PlotWin) ResizeWindow(x, y, w, h int) error {
	// not supported in fyne
	return nil
}

func (pw *PlotWin) ShowBorder(sw, sh int) error {
	// not supported in fyne
	return nil
}

func (pw *PlotWin) WindowXY() (int, int) {
	return 0, 0
}

func (pw *PlotWin) WindowSize() (int, int) {
	return pw.canvas.ggimg.Width(), pw.canvas.ggimg.Height()
}

// Refresh is expected to be called outside of the main fyne thread.
func (pw *PlotWin) Refresh() {
	fyne.DoAndWait(func() {
		pw.canvas.image.Refresh()
	})
	// clear the next time drawing is started
	pw.clearFirst = true
}

func (pw *PlotWin) PixelSize() (int, int) {
	return pw.WindowSize()
}

func (pw *PlotWin) Color(c color.RGBA) {
	pw.color = c
	pw.canvas.ggimg.SetColor(c)
}

func (pw *PlotWin) SetPoint(x, y int) {
	pw.clearIfNeeded()
	pw.canvas.ggimg.SetPixel(x, y)
}

func (pw *PlotWin) HLine(x, y, w int) {
	pw.clearIfNeeded()
	pw.canvas.ggimg.DrawLine(float64(x), float64(y), float64(x+w), float64(y))
	pw.canvas.ggimg.Stroke()
}

func (pw *PlotWin) VLine(x, y, h int) {
	pw.clearIfNeeded()
	pw.canvas.ggimg.DrawLine(float64(x), float64(y), float64(x), float64(y+h))
	pw.canvas.ggimg.Stroke()
}

func (pw *PlotWin) FilledRect(x, y, w, h int) {
	pw.clearIfNeeded()
	pw.canvas.ggimg.DrawRectangle(float64(x), float64(y), float64(w), float64(h))
	pw.canvas.ggimg.Fill()
}

func (pw *PlotWin) Text(s string, x, y int) {
	pw.clearIfNeeded()
	pw.canvas.ggimg.DrawString(s, float64(x), float64(y))
}
