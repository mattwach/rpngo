package plotwin

import (
	"image/color"
	"mattwach/rpngo/common/rpn"
	"mattwach/rpngo/common/window"
)

// colors to autorotate through
var colorWheelPixel = []color.RGBA{
	{R: 255, G: 0, B: 0, A: 255},
	{R: 0, G: 255, B: 0, A: 255},
	{R: 0, G: 0, B: 255, A: 255},
	{R: 255, G: 255, B: 0, A: 255},
	{R: 255, G: 0, B: 255, A: 255},
	{R: 0, G: 255, B: 255, A: 255},
}

type PixelPlotWindow struct {
	pixw       window.PixelWindow
	Common     plotWindowCommon
	lastcolidx uint8
	needUpdate bool
}

func (pw *PixelPlotWindow) Init(pixw window.PixelWindow, steps uint32) {
	pw.pixw = pixw
	pw.Common.init(uint8(len(colorWheelPixel)), steps)
}

func (pw *PixelPlotWindow) ResizeWindow(x, y, w, h int) error {
	pw.needUpdate = true
	if err := pw.pixw.ResizeWindow(x, y, w, h); err != nil {
		return err
	}
	pw.pixw.Color(color.RGBA{})
	psw, psh := pw.pixw.PixelSize()
	pw.pixw.FilledRect(0, 0, psw, psh)
	return nil
}

func (pw *PixelPlotWindow) ShowBorder(screenw, screenh int) error {
	return pw.pixw.ShowBorder(screenw, screenh)
}

func (pw *PixelPlotWindow) WindowXY() (int, int) {
	return pw.pixw.WindowXY()
}

func (pw *PixelPlotWindow) WindowSize() (int, int) {
	return pw.pixw.WindowSize()
}

func (pw *PixelPlotWindow) Type() string {
	return "plot"
}

func (pw *PixelPlotWindow) Update(r *rpn.RPN, force bool) error {
	// Updates are expensive so don't do them if not needed
	if !force && !pw.needUpdate {
		return nil
	}
	pw.Common.setAxisMinMax(r)
	pw.drawAxis()
	pw.lastcolidx = 255
	// Do not exit the program if this fails
	_ = pw.Common.createPoints(r, pw.plotPoint)
	pw.pixw.Refresh()
	pw.needUpdate = false
	return nil
}

func (pw *PixelPlotWindow) SetProp(name string, val rpn.Frame) error {
	pw.needUpdate = true
	return pw.Common.setProp(name, val)
}

func (pw *PixelPlotWindow) GetProp(name string) (rpn.Frame, error) {
	return pw.Common.getProp(name)
}

func (pw *PixelPlotWindow) ListProps() []string {
	return pw.Common.ListProps()
}

func (pw *PixelPlotWindow) plotPoint(x, y float64, colidx uint8) error {
	w, h := pw.pixw.PixelSize()
	if colidx != pw.lastcolidx {
		pw.lastcolidx = colidx
		pw.pixw.Color(colorWheelPixel[colidx])
	}
	wx, xok := pw.Common.transformX(x, w)
	if !xok {
		return nil
	}
	wy, yok := pw.Common.transformY(y, h)
	if !yok {
		return nil
	}
	pw.pixw.SetPoint(wx, wy)
	return nil
}

func (pw *PixelPlotWindow) PixelToCoord(x, y int) (float64, float64) {
	w, h := pw.pixw.PixelSize()
	return pw.Common.pixelToCoordX(x, w), pw.Common.pixelToCoordY(y, h)
}
