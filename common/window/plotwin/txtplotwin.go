// Package plotwin shows a 2d plot
package plotwin

import (
	"mattwach/rpngo/common/rpn"
	"mattwach/rpngo/common/window"
)

// colors to autorotate through
var colorWheelTxt = []window.ColorChar{
	window.Red,
	window.Green,
	window.Blue,
	window.Yellow,
	window.Magenta,
	window.Cyan,
}

type TxtPlotWindow struct {
	txtw   window.TextWindow
	common plotWindowCommon
}

func (pw *TxtPlotWindow) Init(txtw window.TextWindow) {
	pw.txtw = txtw
	pw.common.init(uint8(len(colorWheelTxt)))
}

func (pw *TxtPlotWindow) ResizeWindow(x, y, w, h int) error {
	return pw.txtw.ResizeWindow(x, y, w, h)
}

func (pw *TxtPlotWindow) ShowBorder(screenw, screenh int) error {
	return pw.txtw.ShowBorder(screenw, screenh)
}

func (pw *TxtPlotWindow) WindowXY() (int, int) {
	return pw.txtw.WindowXY()
}

func (pw *TxtPlotWindow) WindowSize() (int, int) {
	return pw.txtw.WindowSize()
}

func (pw *TxtPlotWindow) Type() string {
	return "plot"
}

func (pw *TxtPlotWindow) Update(r *rpn.RPN, unusedForce bool) error {
	pw.txtw.Erase()
	defer pw.txtw.Refresh()
	pw.common.setAxisMinMax(r)
	// do not exit the program if either of these fail
	_ = pw.drawAxis()
	_ = pw.common.createPoints(r, pw.plotPoint)
	return nil
}

func (pw *TxtPlotWindow) SetProp(name string, val rpn.Frame) error {
	return pw.common.setProp(name, val)
}

func (pw *TxtPlotWindow) GetProp(name string) (rpn.Frame, error) {
	return pw.common.getProp(name)
}

func (pw *TxtPlotWindow) ListProps() []string {
	return pw.common.ListProps()
}

func (pw *TxtPlotWindow) plotPoint(x, y float64, colidx uint8) error {
	w, h := pw.txtw.TextSize()
	tx, xok := pw.common.transformX(x, w)
	if !xok {
		return nil
	}
	ty, yok := pw.common.transformY(y, h)
	if !yok {
		return nil
	}
	pw.txtw.DrawChar(tx, ty, colorWheelTxt[colidx]|'*')
	return nil
}
