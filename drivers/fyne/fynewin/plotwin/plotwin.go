package plotwin

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
	win    fyne.Window
	color  color.RGBA
	img    *gg.Context
	canvas *canvas.Image
}

func New(win fyne.Window, r *rpn.RPN) *PlotWin {
	pw := &PlotWin{
		win: win,
		img: gg.NewContext(1024, 768),
	}

	pw.canvas = canvas.NewImageFromImage(pw.img.Image())
	pw.canvas.FillMode = canvas.ImageFillContain
	pw.win.SetContent(container.NewStack(pw.canvas))
	pw.win.Resize(fyne.NewSize(1024, 768))
	pw.Update(r)
	return pw
}

func (pw *PlotWin) Update(r *rpn.RPN) {
	pw.img.SetRGB(0, 0, 0)
	pw.img.Clear()
	pw.img.SetRGB(1, 1, 0)
	pw.img.DrawCircle(512, 512, 100)
	pw.img.Fill()
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

func (pw *PlotWin) Refresh() {}

func (pw *PlotWin) PixelSize() (int, int) {
	return pw.WindowSize()
}

func (pw *PlotWin) Color(c color.RGBA) {
	pw.color = c
}

func (pw *PlotWin) SetPoint(x, y int) {
	pw.img.SetPixel(x, y)
}

func (pw *PlotWin) HLine(x, y, w int) {
	pw.img.DrawLine(float64(x), float64(y), float64(x+w), float64(y))
}

func (pw *PlotWin) VLine(x, y, h int) {
	pw.img.DrawLine(float64(x), float64(y), float64(x), float64(y+h))
}

func (pw *PlotWin) FilledRect(x, y, w, h int) {
	pw.img.DrawRectangle(float64(x), float64(y), float64(w), float64(h))
	pw.img.Fill()
}

func (pw *PlotWin) Text(s string, x, y int) {
	pw.img.DrawString(s, float64(x), float64(y))
}
