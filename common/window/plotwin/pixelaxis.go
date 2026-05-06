package plotwin

import (
	"image/color"
	"strconv"
)

var white = color.RGBA{R: 255, G: 255, B: 255, A: 255}

func (pw *PixelPlotWindow) drawAxis() {
	pw.pixw.Color(white)
	w, h := pw.pixw.PixelSize()
	pw.drawVerticalAxis(w, h)
	pw.drawHorizontalAxis(w, h)
}

func (pw *PixelPlotWindow) drawVerticalAxis(w, h int) {
	x, xok := pw.common.transformX(0, w)
	if xok {
		pw.pixw.VLine(x, 0, h-1)
		pw.drawVerticalTickMarks(x, w, h)
	}
}

// About axis drawing
// the goal is to make it look nice while also being useful
// We are going for rounded tick marks (e.g. 0.25, 1.0, ...)
// and reasonable pixels per mark

const fontCharWidth = 8
const fontCharHeight = 12
const minPixelVerticalSpacing = 2 * fontCharHeight
const maxPixelVerticalSpacing = 4 * fontCharHeight
const tickLength = 4
const horizTextOffset = 10
const vertTextOffset = 20

func (pw *PixelPlotWindow) drawVerticalTickMarks(wx, ww, wh int) {
	yr := pw.common.maxy - pw.common.miny // units
	cpu := float64(wh) / yr               // characters / unit
	var te float64 = 1                    // ticks every (0.5, 1, etc)
	if cpu > maxPixelVerticalSpacing {
		te = searchScaleDownward(cpu, minPixelVerticalSpacing)
	} else if cpu < minPixelVerticalSpacing {
		te = searchScaleUpward(cpu, maxPixelVerticalSpacing)
	}

	// becuase te was carefully selected to provide a limited number
	// of tick marks, we can use a simple loop to determine the min te
	stepsBack := 0
	for (float64(stepsBack-1) * te) > pw.common.miny {
		stepsBack--
	}

	for {
		y := float64(stepsBack) * te
		if y >= pw.common.maxy {
			break
		}
		if stepsBack != 0 {
			pw.drawVerticalTick(wx, ww, wh, y)
		}
		stepsBack++
	}
}

func (pw *PixelPlotWindow) drawVerticalTick(wx, ww, wh int, y float64) {
	wy, _ := pw.common.transformY(y, wh)
	x0 := wx - tickLength
	if x0 < 0 {
		x0 = 0
	}
	x1 := wx + tickLength
	if x1 >= ww {
		x1 = ww - 1
	}
	pw.pixw.HLine(x0, wy, x1-x0)
	s := strconv.FormatFloat(y, 'g', 2, 64)
	lx := wx + horizTextOffset
	ly := wy + fontCharHeight/2
	if (lx > 0) && (lx < (ww - len(s)*fontCharWidth)) && (ly > fontCharHeight) && (ly < wh) {
		pw.pixw.Text(s, lx, ly)
	}
}

func (pw *PixelPlotWindow) drawHorizontalAxis(w, h int) {
	y, yok := pw.common.transformY(0, h)
	if yok {
		pw.pixw.HLine(0, y, w)
		pw.drawHorizontalTickMarks(y, w, h)
	}
}

const minPixelHorizontalSpacing = 5 * fontCharWidth
const maxPixelHorizontalSpacing = 8 * fontCharWidth

func (pw *PixelPlotWindow) drawHorizontalTickMarks(wy, ww, wh int) {
	xr := pw.common.maxx - pw.common.minx // units
	cpu := float64(ww) / xr               // characters / unit
	var te float64 = 1                    // ticks every (0.5, 1, etc)
	if cpu > maxPixelHorizontalSpacing {
		te = searchScaleDownward(cpu, minPixelHorizontalSpacing)
	} else if cpu < minPixelHorizontalSpacing {
		te = searchScaleUpward(cpu, maxPixelHorizontalSpacing)
	}

	stepsBack := 0
	for (float64(stepsBack-1) * te) > pw.common.minx {
		stepsBack--
	}

	for {
		x := float64(stepsBack) * te
		if x >= pw.common.maxx {
			break
		}
		if stepsBack != 0 {
			pw.drawHorizontalTick(x, wy, ww, wh)
		}
		stepsBack++
	}
}

func (pw *PixelPlotWindow) drawHorizontalTick(x float64, wy, ww, wh int) {
	wx, _ := pw.common.transformX(x, ww)
	y0 := wy - tickLength
	if y0 < 0 {
		y0 = 0
	}
	y1 := wy + tickLength
	if y1 >= wh {
		y1 = wh - 1
	}
	pw.pixw.VLine(wx, y0, y1-y0)
	s := strconv.FormatFloat(x, 'g', 2, 64)
	lx := wx - len(s)*fontCharWidth/2
	ly := wy + vertTextOffset
	if (lx > 0) && (lx < (ww - len(s)*fontCharWidth)) && (ly > fontCharHeight) && (ly < wh) {
		pw.pixw.Text(s, lx, ly)
	}
}
