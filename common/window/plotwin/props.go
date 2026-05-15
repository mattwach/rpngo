package plotwin

import (
	"mattwach/rpngo/common/elog"
	"mattwach/rpngo/common/rpn"
	"strconv"
	"strings"
)

func (pw *plotWindowCommon) setProp(name string, val rpn.Frame) error {
	switch name {
	case "minx":
		v, err := val.Real()
		if err != nil {
			return err
		}
		pw.minx = v
		if pw.maxx <= pw.minx {
			pw.maxx = pw.minx + 1
		}
		pw.autox = false
		return nil
	case "maxx":
		v, err := val.Real()
		if err != nil {
			return err
		}
		pw.maxx = v
		if pw.maxx <= pw.minx {
			pw.minx = pw.maxx - 1
		}
		pw.autox = false
		return nil
	case "miny":
		v, err := val.Real()
		if err != nil {
			return err
		}
		pw.miny = v
		if pw.maxy <= pw.miny {
			pw.maxy = pw.miny + 1
		}
		pw.autoy = false
		return nil
	case "maxy":
		v, err := val.Real()
		if err != nil {
			return err
		}
		pw.maxy = v
		if pw.maxy <= pw.miny {
			pw.miny = pw.maxy - 1
		}
		pw.autoy = false
		return nil
	case "minv":
		v, err := val.Real()
		if err != nil {
			return err
		}
		pw.minv = v
		if pw.maxv <= pw.minv {
			pw.maxv = pw.minv + 1
		}
		return nil
	case "maxv":
		v, err := val.Real()
		if err != nil {
			return err
		}
		pw.maxv = v
		if pw.maxv <= pw.minv {
			pw.minv = pw.maxv - 1
		}
		return nil
	case "autox":
		v, err := val.Bool()
		if err != nil {
			return err
		}
		pw.autox = v
		return nil
	case "autoy":
		v, err := val.Bool()
		if err != nil {
			return err
		}
		pw.autoy = v
		return nil
	case "drawaxis":
		v, err := val.Bool()
		if err != nil {
			return err
		}
		pw.drawaxis = v
		return nil
	case "steps":
		v, err := val.BoundedInt(1, maxSteps)
		if err != nil {
			return err
		}
		pw.steps = uint32(v)
		return nil
	case "numplots":
		v, err := val.Int()
		if err != nil {
			return err
		}
		return pw.changePlotCount(int(v))
	}

	if strings.HasPrefix(name, "color") && (len(name) > 5) {
		idx, err := strconv.Atoi(name[5:])
		if (err == nil) && (idx >= 0) && (idx < len(pw.plots)) {
			v, err := val.BoundedInt(0, int64(pw.numColors)-1)
			if err != nil {
				return err
			}
			pw.plots[idx].coloridx = uint8(v)
			return nil
		}
	}

	if strings.HasPrefix(name, "parametric") && (len(name) > 10) {
		idx, err := strconv.Atoi(name[10:])
		if (err == nil) && (idx >= 0) && (idx < len(pw.plots)) {
			v, err := val.Bool()
			if err != nil {
				return err
			}
			pw.plots[idx].isParametric = v
			return nil
		}
	}

	if strings.HasPrefix(name, "fn") && (len(name) > 2) {
		idx, err := strconv.Atoi(name[2:])
		if (err == nil) && (idx >= 0) && (idx < len(pw.plots)) {
			v := val.String(false)
			return pw.setPlotFn(v, idx)
		}
	}

	return rpn.ErrUnknownProperty
}

func (pw *plotWindowCommon) getProp(name string) (rpn.Frame, error) {
	switch name {
	case "minx":
		return rpn.RealFrame(pw.minx), nil
	case "maxx":
		return rpn.RealFrame(pw.maxx), nil
	case "miny":
		return rpn.RealFrame(pw.miny), nil
	case "maxy":
		return rpn.RealFrame(pw.maxy), nil
	case "minv":
		return rpn.RealFrame(pw.minv), nil
	case "maxv":
		return rpn.RealFrame(pw.maxv), nil
	case "autox":
		return rpn.BoolFrame(pw.autox), nil
	case "autoy":
		return rpn.BoolFrame(pw.autoy), nil
	case "drawaxis":
		return rpn.BoolFrame(pw.drawaxis), nil
	case "numplots":
		return rpn.IntFrame(int64(len(pw.plots)), rpn.INTEGER_FRAME), nil
	case "steps":
		return rpn.IntFrame(int64(pw.steps), rpn.INTEGER_FRAME), nil
	}

	if strings.HasPrefix(name, "color") && (len(name) > 5) {
		idx, err := strconv.Atoi(name[5:])
		if (err == nil) && (idx >= 0) && (idx < len(pw.plots)) {
			return rpn.IntFrame(int64(pw.plots[idx].coloridx), rpn.INTEGER_FRAME), nil
		}
	}

	if strings.HasPrefix(name, "parametric") && (len(name) > 10) {
		idx, err := strconv.Atoi(name[10:])
		if (err == nil) && (idx >= 0) && (idx < len(pw.plots)) {
			return rpn.BoolFrame(pw.plots[idx].isParametric), nil
		}
	}

	if strings.HasPrefix(name, "fn") && (len(name) > 2) {
		idx, err := strconv.Atoi(name[2:])
		if (err == nil) && (idx >= 0) && (idx < len(pw.plots)) {
			return rpn.StringFrame(strings.Join(pw.plots[idx].fn, " "), rpn.STRING_BRACE_FRAME), nil
		}
	}

	return rpn.Frame{}, rpn.ErrUnknownProperty
}

var props = []string{"minv", "maxv", "minx", "maxx", "miny", "maxy", "numplots", "steps", "autox", "autoy"}

func (pw *plotWindowCommon) ListProps() []string {
	elog.Heap("alloc: window/plotwin/props.go:191: wprops := make([]string, len(props)+len(pw.plots)*3)")
	wprops := make([]string, len(props)+len(pw.plots)*3) // object allocated on the heap: size is not constant
	copy(wprops, props)
	j := len(props)
	for i := range pw.plots {
		plotid := strconv.Itoa(i)
		wprops[j] = "color" + plotid
		j++
		wprops[j] = "parametric" + plotid
		j++
		wprops[j] = "fn" + plotid
		j++
	}
	return wprops
}
