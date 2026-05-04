package plotwin

import (
	"errors"
	"fmt"
	"log"
	"mattwach/rpngo/common/elog"
	"mattwach/rpngo/common/parse"
	"mattwach/rpngo/common/rpn"
	"strings"
)

const maxSteps = 500000
const maxPlotCount = 32

type PointStats struct {
	minx        float64
	maxx        float64
	miny        float64
	maxy        float64
	initialized bool
}

func (ps *PointStats) reset() {
	ps.initialized = false
}

func (ps *PointStats) update(x, y float64, colidx uint8) error {
	if !ps.initialized {
		ps.initialized = true
		ps.minx = x
		ps.maxx = x
		ps.miny = y
		ps.maxy = y
	} else {
		if x < ps.minx {
			ps.minx = x
		}
		if x > ps.maxx {
			ps.maxx = x
		}
		if y < ps.miny {
			ps.miny = y
		}
		if y > ps.maxy {
			ps.maxy = y
		}
	}
	return nil
}

type Plot struct {
	fn           []string
	coloridx     uint8
	isParametric bool
}

type plotWindowCommon struct {
	MinX      float64
	MaxX      float64
	MinY      float64
	MaxY      float64
	coloridx  uint8
	numColors uint8
	AutoX     bool
	AutoY     bool
	MinV      float64
	MaxV      float64
	steps     uint32
	plots     []Plot
	stats     PointStats
}

func (pw *plotWindowCommon) init(numColors uint8, steps uint32) {
	pw.AutoX = true
	pw.AutoY = true
	pw.MinV = -1
	pw.MaxV = 1
	pw.steps = steps
	pw.numColors = numColors
}

func (pw *plotWindowCommon) HasAParametricPlot() bool {
	for _, plot := range pw.plots {
		if plot.isParametric {
			return true
		}
	}
	return false
}

func (pw *plotWindowCommon) nextColor(numColors uint8) {
	pw.coloridx++
	if pw.coloridx >= numColors {
		pw.coloridx = 0
	}
}

func (pw *plotWindowCommon) changePlotCount(n int) error {
	if (n < 0) || (n >= maxPlotCount) {
		return rpn.ErrIllegalValue
	}
	var col uint8 = uint8(len(pw.plots))
	for n > len(pw.plots) {
		pw.plots = append(pw.plots, Plot{coloridx: col % pw.numColors})
		col++
	}
	if n < len(pw.plots) {
		pw.plots = pw.plots[:n]
	}
	return nil
}

func (pw *plotWindowCommon) setPlotFn(fnstr string, idx int) error {
	if (idx < 0) || (idx >= len(pw.plots)) {
		return rpn.ErrIllegalValue
	}
	elog.Heap("alloc: window/plotwin/common.go:106: p := &pw.plots[idx]")
	p := &pw.plots[idx] // object allocated on the heap: escapes at line 112
	p.fn = p.fn[:0]
	addField := func(t string) error {
		p.fn = append(p.fn, t)
		return nil
	}
	if err := parse.Fields(fnstr, addField); err != nil {
		p.fn = p.fn[:0]
		return err
	}

	return nil
}

func (pw *plotWindowCommon) setAxisMinMax(r *rpn.RPN) {
	// emergency corrections to avoid hanging in other parts of the code.
	// This will only happen if other parts of the code fail to
	// keep these values proper.
	if !pw.AutoX && (pw.MinX >= pw.MaxX) {
		log.Printf("warning: MinX >= MaxX (%f >= %f), adjusting to avoid math issues", pw.MinX, pw.MaxX)
		pw.AutoX = true
	}
	if !pw.AutoY && (pw.MinY >= pw.MaxY) {
		log.Printf("warning: MinY >= MaxY (%f >= %f), adjusting to avoid math issues", pw.MinY, pw.MaxY)
		pw.AutoY = true
	}
	// first determine the ranges
	if pw.AutoX || pw.AutoY {
		pw.stats.reset()
		for _, plot := range pw.plots {
			if err := pw.addPoints(r, plot, pw.steps, pw.stats.update); err != nil {
				// this plot has some type of error, but there is nothing to be done
				// here outside of not contributing any more points from this point
				// to the stats
			}
		}
		if pw.AutoX {
			pw.adjustAutoX()
		}
		if pw.AutoY {
			pw.adjustAutoY()
		}
	}
}

func (pw *plotWindowCommon) createPoints(r *rpn.RPN, fn func(x, y float64, coloridx uint8) error) error {
	var finalErr error
	for i, plot := range pw.plots {
		if err := pw.addPoints(r, plot, pw.steps, fn); err != nil {
			r.Print("error plotting {")
			r.Print(strings.Join(plot.fn, " "))
			r.Print("}: ")
			r.Print(err.Error())
			r.Println(" (removing plot)")
			pw.plots[i].fn = pw.plots[i].fn[:0]
			finalErr = err
		}
	}
	return finalErr
}

func (pw *plotWindowCommon) addPoints(r *rpn.RPN, plot Plot, steps uint32, fn func(x, y float64, coloridx uint8) error) error {
	if len(plot.fn) == 0 {
		return nil
	}
	startlen := r.StackLen()
	step := (pw.MaxV - pw.MinV) / float64(steps)
	var x float64
	t0 := true
	for v := pw.MinV; v <= pw.MaxV; v += step {
		if t0 {
			if err := setT0(r, true); err != nil {
				return err
			}
		}
		if err := r.PushFrame(rpn.RealFrame(v)); err != nil {
			return err
		}
		if err := r.ExecSlice(plot.fn); err != nil {
			if errors.Is(err, rpn.ErrDivideByZero) {
				// just skip this point
				continue
			}
			return err
		}
		if t0 {
			if err := setT0(r, false); err != nil {
				return err
			}
			t0 = false
		}
		yf, err := r.PopFrame()
		if err != nil {
			return err
		}
		y, err := yf.Real()
		if err != nil {
			return err
		}
		if plot.isParametric {
			xf, err := r.PopFrame()
			if err != nil {
				return err
			}
			x, err = xf.Real()
			if err != nil {
				return err
			}
		} else {
			x = v
		}
		nowlen := r.StackLen()
		if nowlen != startlen {
			return fmt.Errorf(
				"stack changed size running plot string (old: %d, new %d)",
				startlen,
				nowlen)
		}
		if err := fn(x, y, plot.coloridx); err != nil {
			return err
		}
	}
	return nil
}

func setT0(r *rpn.RPN, t0 bool) error {
	if err := r.PushFrame(rpn.BoolFrame(t0)); err != nil {
		return err
	}
	return r.SetVariable(".t0")
}

func (pw *plotWindowCommon) adjustAutoX() {
	pw.MinX = pw.stats.minx
	pw.MaxX = pw.stats.maxx
	if pw.MinX == pw.MaxX {
		// create a little spread to avoid math issues
		pw.MinX -= 1.0
		pw.MaxX += 1.0
	}
}

func (pw *plotWindowCommon) adjustAutoY() {
	pw.MinY = pw.stats.miny
	pw.MaxY = pw.stats.maxy
	if pw.MinY == pw.MaxY {
		// create a little spread to avoid math issues
		pw.MinY -= 1.0
		pw.MaxY += 1.0
	}
	// open up the y a bit
	delta := (pw.MaxY - pw.MinY) / 5
	pw.MaxY += delta
	pw.MinY -= delta
}

func (pw *plotWindowCommon) pixelToCoordX(x, w int) float64 {
	return pw.MinX + (pw.MaxX-pw.MinX)*(float64(x)/float64(w))
}

func (pw *plotWindowCommon) pixelToCoordY(y, h int) float64 {
	return pw.MinY + (pw.MaxY-pw.MinY)*(float64(h-y)/float64(h))
}

func (pw *plotWindowCommon) transformX(x float64, w int) (int, bool) {
	x = (x - pw.MinX) / (pw.MaxX - pw.MinX)
	if x < 0 {
		return 0, false
	}
	xi := int(float64(w)*x + 0.5)
	if xi < 0 || xi >= w {
		return 0, false
	}
	return xi, true
}

func (pw *plotWindowCommon) transformY(y float64, h int) (int, bool) {
	y = (y - pw.MinY) / (pw.MaxY - pw.MinY)
	if y < 0 {
		return 0, false
	}
	py := h - int(float64(h)*y+0.5) - 1
	if py < 0 || py > h {
		return 0, false
	}
	return py, true
}

// use a nice-looking scale. 1, 0.5, 0.25, 0.1, 0.05, 0.025, 0.01, etc
func searchScaleDownward(cpu, minSpacing float64) float64 {
	tens := 1.0
	partial := 1
	te := 1.0

	for {
		switch partial {
		case 1:
			partial = 2
		case 2:
			partial = 4
		case 4:
			partial = 1
			tens *= 10
		}

		newte := 1.0 / (tens * float64(partial))
		if (cpu * newte) < minSpacing {
			// too far
			break
		}
		te = newte
	}

	return te
}

func searchScaleUpward(cpu, maxSpacing float64) float64 {
	tens := 1.0
	partialDeci := 10
	te := 1.0

	for {
		switch partialDeci {
		case 10:
			partialDeci = 25
		case 25:
			partialDeci = 50
		case 50:
			partialDeci = 10
			tens *= 10
		}

		newte := tens * float64(partialDeci) / 10.0
		if (cpu * newte) > maxSpacing {
			// too far
			break
		}
		te = newte
	}

	return te
}
