package plotwin

import (
	"fmt"
	"mattwach/rpngo/common/window"
)

func (pw *TxtPlotWindow) drawAxis() error {
	x, xok := pw.common.transformX(0, pw.txtw.TextWidth())
	if xok {
		if err := pw.drawVerticalAxis(x); err != nil {
			return err
		}
	}

	y, yok := pw.common.transformY(0, pw.txtw.TextHeight())
	if yok {
		if err := pw.drawHorizontalAxis(y); err != nil {
			return err
		}
	}

	if xok && yok {
		pw.txtw.DrawChar(x, y, window.White|'+')
	}

	return nil
}

func (pw *TxtPlotWindow) drawVerticalAxis(x int) error {
	h := pw.txtw.TextHeight()
	for y := 0; y < h; y++ {
		pw.txtw.DrawChar(x, y, window.White|'|')
	}
	pw.drawVerticalTickMarks(x)
	return nil
}

// About axis drawing
// the goal is to make it look nice while also being useful
// We are going for rounded tick marks (e.g. 0.25, 1.0, ...)
// and reasonable pixels per mark

const minTxtVerticalSpacing = 5.0
const maxTxtVerticalSpacing = 10.0

func (pw *TxtPlotWindow) drawVerticalTickMarks(wx int) {
	wh := pw.txtw.TextHeight()            // height in characters
	yr := pw.common.MaxY - pw.common.MinY // units
	cpu := float64(wh) / yr               // characters / unit
	var te float64 = 1                    // ticks every (0.5, 1, etc)
	if cpu > maxTxtVerticalSpacing {
		te = searchScaleDownward(cpu, minTxtVerticalSpacing)
	} else if cpu < minTxtVerticalSpacing {
		te = searchScaleUpward(cpu, maxTxtVerticalSpacing)
	}

	// becuase te was carefully selected to provide a limited number
	// of tick marks, we can use a simple loop to determine the min te
	stepsBack := 0
	for (float64(stepsBack-1) * te) > pw.common.MinY {
		stepsBack--
	}

	for {
		y := float64(stepsBack) * te
		if y >= pw.common.MaxY {
			break
		}
		if stepsBack != 0 {
			pw.drawVerticalTick(wx, y)
		}
		stepsBack++
	}
}

func (pw *TxtPlotWindow) drawVerticalTick(wx int, y float64) {
	ww := pw.txtw.TextWidth()
	wy, _ := pw.common.transformY(y, pw.txtw.TextHeight())
	pw.txtw.DrawChar(wx, wy, window.White|'+')
	if (wx + 10) < ww {
		window.DrawStr(pw.txtw, wx+3, wy, fmt.Sprintf("%.2f", y), window.White)
	}
}

func (pw *TxtPlotWindow) drawHorizontalAxis(y int) error {
	w := pw.txtw.TextWidth()
	for x := 0; x < w; x++ {
		pw.txtw.DrawChar(x, y, window.White|'-')
	}
	pw.drawHorizontalTickMarks(y)
	return nil
}

const minTxtHorizontalSpacing = 9.0
const maxTxtHorizontalSpacing = 18.0

func (pw *TxtPlotWindow) drawHorizontalTickMarks(wy int) {
	ww := pw.txtw.TextWidth()             // width in characters
	xr := pw.common.MaxX - pw.common.MinX // units
	cpu := float64(ww) / xr               // characters / unit
	var te float64 = 1                    // ticks every (0.5, 1, etc)
	if cpu > maxTxtHorizontalSpacing {
		te = searchScaleDownward(cpu, minTxtHorizontalSpacing)
	} else if cpu < minTxtHorizontalSpacing {
		te = searchScaleUpward(cpu, maxTxtHorizontalSpacing)
	}

	stepsBack := 0
	for (float64(stepsBack-1) * te) > pw.common.MinX {
		stepsBack--
	}

	for {
		x := float64(stepsBack) * te
		if x >= pw.common.MaxX {
			break
		}
		if stepsBack != 0 {
			pw.drawHorizontalTick(x, wy)
		}
		stepsBack++
	}
}

const horizontalTxtNumberPad = 5

func (pw *TxtPlotWindow) drawHorizontalTick(x float64, wy int) {
	ww, wh := pw.txtw.TextSize()
	wx, _ := pw.common.transformX(x, pw.txtw.TextWidth())
	if wy >= 0 {
		pw.txtw.DrawChar(wx, wy, window.White|'+')
	}
	if (wx > horizontalTxtNumberPad) && (wx < (ww - horizontalTxtNumberPad)) && (wy+1) < wh {
		s := fmt.Sprintf("%.2f", x)
		window.DrawStr(pw.txtw, wx-len(s)/2, wy+1, s, window.White)
	}
}
