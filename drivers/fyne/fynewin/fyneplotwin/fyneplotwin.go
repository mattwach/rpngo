package fyneplotwin

import (
	"image/color"
	"mattwach/rpngo/common/rpn"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"

	"github.com/fogleman/gg"
)

// PlotWin holds the context for a stack window.
// Important, RPN is owned by the readline goroutine thus should be accessed
// with care.  This means that putting a pointer to it in this struct is
// probably starting down a bad path.
type PlotWin struct {
	win        fyne.Window
	img        *gg.Context
	canvas     *canvas.Image
	color      color.RGBA
	clearFirst bool
}

func New(win fyne.Window, r *rpn.RPN) *PlotWin {
	pw := &PlotWin{
		win:        win,
		img:        gg.NewContext(1024, 768),
		clearFirst: true,
	}

	pw.canvas = canvas.NewImageFromImage(pw.img.Image())
	pw.win.SetContent(container.NewStack(pw.canvas))
	pw.win.Resize(fyne.NewSize(1024, 768))
	return pw
}

func (pw *PlotWin) clearIfNeeded() {
	if !pw.clearFirst {
		return
	}
	pw.img.SetRGB(0, 0, 0)
	pw.img.Clear()
	pw.img.SetColor(pw.color)
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
	return pw.img.Width(), pw.img.Height()
}

func (pw *PlotWin) Refresh() {
	pw.win.SetContent(container.NewStack(pw.canvas))
	// clear the next time drawing is started
	pw.clearFirst = true
}

func (pw *PlotWin) PixelSize() (int, int) {
	return pw.WindowSize()
}

func (pw *PlotWin) Color(c color.RGBA) {
	pw.color = c
	pw.img.SetColor(c)
}

func (pw *PlotWin) SetPoint(x, y int) {
	pw.clearIfNeeded()
	pw.img.SetPixel(x, y)
}

func (pw *PlotWin) HLine(x, y, w int) {
	pw.clearIfNeeded()
	pw.img.DrawLine(float64(x), float64(y), float64(x+w), float64(y))
	pw.img.Stroke()
}

func (pw *PlotWin) VLine(x, y, h int) {
	pw.clearIfNeeded()
	pw.img.DrawLine(float64(x), float64(y), float64(x), float64(y+h))
	pw.img.Stroke()
}

func (pw *PlotWin) FilledRect(x, y, w, h int) {
	pw.clearIfNeeded()
	pw.img.DrawRectangle(float64(x), float64(y), float64(w), float64(h))
	pw.img.Fill()
}

func (pw *PlotWin) Text(s string, x, y int) {
	pw.clearIfNeeded()
	pw.img.DrawString(s, float64(x), float64(y))
}
