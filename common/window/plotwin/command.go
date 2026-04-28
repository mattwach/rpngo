package plotwin

import (
	"errors"
	"fmt"
	"mattwach/rpngo/common/elog"
	"mattwach/rpngo/common/parse"
	"mattwach/rpngo/common/rpn"
	"mattwach/rpngo/common/window"
	"strconv"
	"strings"
)

type WindowManager interface {
	FindWindow(name string) window.WindowWithProps
	Update(r *rpn.RPN, screenw, screenh int, updateInput bool) error
}

type PlotCommands struct {
	manager WindowManager
	screen  window.Screen
}

func InitPlotCommands(
	r *rpn.RPN,
	manager WindowManager,
	screen window.Screen) *PlotCommands {
	conceptHelp := map[string]string{
		"plot": "Plot functions using plot. Plot will push an 'x' value to the stack,\n" +
			"run the provided string, and pop the value as y value.\n" +
			"Examples:\n" +
			"    '2 *' plot # plots y = x * 2\n" +
			"    'sq' plot # plots y = x * x\n" +
			"    'sin' plot # plots y = sin(x)\n" +
			"Various properties can be set on the .plotwindow to change the number\n" +
			"of points and the boundaries of the plot.\n" +
			"There are some special variables that plot uses:\n" +
			"    .plotwin  : Name of the window to send plots to (there can be more than one)\n" +
			"                at a time.\n" +
			"    .plotinit : If no .plotwindow exists, this string is executed and is expected\n" +
			"                create one. Making this a variable allows for user customization.\n" +
			"    .t0       : Set to true if this is the first point of a given plot, false otherwise\n" +
			"                The intent to to allow plot functions that use tracking variables to\n" +
			"                known when it's time to intialize them.\n" +
			"See Also: window.props, plot.parametric",

		"plot.parametric": "Plot parametric functions using pplot. pplot will push a 't' value to\n" +
			"the stack, run the provided string then pop y, then x to determine the plot point x, y\n" +
			"Examples:\n" +
			"    '$0 cos 1> sin' pplot # draws an arc or full circle, depending on t range\n" +
			"    't= $t sin $t * $t cos $t *' pplot # draw a spiral\n" +
			"    '1 sw' draw a vertical line\n",
	}
	r.RegisterConceptHelp(conceptHelp)

	elog.Heap("alloc: /window/plotwin/command.go:48: pc := PlotCommands{root: root, screen: screen, addPlotFn: addPlotFn}")
	pc := PlotCommands{manager: manager, screen: screen} // object allocated on the heap: escapes at line 50
	r.Register("plot", pc.plot, rpn.CatPlot, plotHelp)
	r.Register("pplot", pc.pplot, rpn.CatPlot, pplotHelp)
	return &pc
}

const plotHelp = "Run a value from $plot.min to $plot.max with $plot.steps\n" +
	"executing the provided string and plotting to the window $.plotwin\n" +
	"Example (y = x * x): 'c *' plot"

func (pc *PlotCommands) plot(r *rpn.RPN) error {
	return pc.plotInternal(r, false)
}

const pplotHelp = "Run a value from $plot.min to $plot.max with $plot.steps\n" +
	"executing the provided string and plotting (x, y) to the window $.plotwin\n" +
	"Example (arc): 'c sin sw cos"

func (wc *PlotCommands) pplot(r *rpn.RPN) error {
	return wc.plotInternal(r, true)
}

func (pc *PlotCommands) plotInternal(r *rpn.RPN, isParametric bool) error {
	macro, err := r.PopFrame()
	if err != nil {
		return err
	}
	if !macro.IsString() {
		return rpn.ErrExpectedAString
	}
	// a quick check to make sure there will not be a parsing error after
	// adding the new plot.  Also checks for plto that already exist.
	elog.Heap("alloc: window/plotwin/command.go:84: var fields []string")
	var fields []string // object allocated on the heap: escapes at line 89
	addField := func(t string) error {
		fields = append(fields, t)
		return nil
	}
	if err := parse.Fields(macro.String(false), addField); err != nil {
		return err
	}
	wname, err := r.GetStringVariable(".plotwin")
	if err != nil {
		return err
	}
	pw := pc.manager.FindWindow(wname)
	if pw == nil {
		if err := pc.initPlot(r); err != nil {
			return err
		}
		pw = pc.manager.FindWindow(wname)
	}
	if pw == nil {
		return errors.New("executing $.plotinit did not result in a window: " + wname)
	}

	if pw.Type() != "plot" {
		return errors.New(wname + " has the wrong window type: " + pw.Type())
	}

	// From this point forward, we interface with the rpn in the same
	// low level manner avaialble to the user (e.g. using props). This
	// is done so there is only one "actual way" to crerate plots, which
	// means less code paths (less rarely-executed code, bugs)
	wwin := "'" + wname + "'"

	numplots, err := getNumPlots(r, wwin)
	if err != nil {
		return err
	}
	idx, err := findOpenPlotIndex(r, wwin, numplots, strings.Join(fields, "\n"))
	if err != nil {
		return err
	}

	if idx < 0 {
		// create a new one
		idx = numplots
		if err := r.ExecSlice([]string{wwin, "'numplots'"}); err != nil {
			return err
		}
		r.PushFrame(rpn.IntFrame(int64(numplots+1), rpn.INTEGER_FRAME))
		if err := r.Exec("w.setp"); err != nil {
			return err
		}
	}

	pidx := strconv.Itoa(int(idx))
	if err := setParametric(r, wwin, pidx, isParametric); err != nil {
		return err
	}

	if err := r.ExecSlice([]string{wwin, "'fn" + pidx + "'"}); err != nil {
		return err
	}
	if err := r.PushFrame(macro); err != nil {
		return err
	}
	if err := r.Exec("w.setp"); err != nil {
		return err
	}
	return nil
}

func getNumPlots(r *rpn.RPN, wwin string) (int, error) {
	if err := r.ExecSlice([]string{wwin, "'numplots'", "w.getp"}); err != nil {
		return 0, err
	}
	f, err := r.PopFrame()
	if err != nil {
		return 0, err
	}
	n, err := f.Int()
	return int(n), err
}

func setParametric(r *rpn.RPN, wwin, pidx string, isParametric bool) error {
	if err := r.ExecSlice([]string{wwin, "'parametric" + pidx + "'"}); err != nil {
		return err
	}
	if err := r.PushFrame(rpn.BoolFrame(isParametric)); err != nil {
		return err
	}
	if err := r.Exec("w.setp"); err != nil {
		return err
	}
	return nil
}

func findOpenPlotIndex(r *rpn.RPN, wwin string, numplots int, fn string) (int, error) {
	openIndex := -1
	for i := range numplots {
		if err := r.ExecSlice([]string{
			wwin,
			"'fn" + strconv.Itoa(i) + "'",
			"w.getp",
		}); err != nil {
			return 0, err
		}
		f, err := r.PopFrame()
		if err != nil {
			return 0, err
		}
		if fn == f.String(false) {
			return i, nil
		}
		if (openIndex < 0) && (f.String(false) == "") {
			openIndex = i
		}
	}
	return openIndex, nil
}

func (wc *PlotCommands) initPlot(r *rpn.RPN) error {
	macro, err := r.GetStringVariable(".plotinit")
	if err != nil {
		return err
	}
	if err := parse.Fields(macro, r.Exec); err != nil {
		return err
	}
	var w int
	var h int
	if wc.screen != nil {
		w, h = wc.screen.ScreenSize()
	}
	if uerr := wc.manager.Update(r, w, h, false); uerr != nil {
		elog.Heap("alloc: /window/plotwin/command.go:118: elog.Print('initPlot.Update error: ', uerr.Error())")
		elog.Print("initPlot.Update error: ", uerr.Error()) // object allocated on the heap: escapes at line 118
	}
	if err != nil {
		return fmt.Errorf("while executing $.plotinit: %v", err)
	}
	return nil
}
