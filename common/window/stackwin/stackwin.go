// Package stackwin shows a stack window
package stackwin

import (
	"mattwach/rpngo/common/rpn"
	"mattwach/rpngo/common/window"
	"strconv"
)

type StackWindow struct {
	txtb  window.TextBuffer
	round int8
	rsd   roundedStringData
}

type roundedStringData struct {
	buff      [32]byte
	dec       [12]byte
	inDecimal bool
	didx      int8
	idx       int
}

func (rsd *roundedStringData) reset() {
	rsd.inDecimal = false
	rsd.didx = 0
	rsd.idx = 0
}

func (sw *StackWindow) Init(txtw window.TextWindow) {
	sw.txtb.Init(txtw, 0)
	sw.round = -1
}

func (sw *StackWindow) ResizeWindow(x, y, w, h int) error {
	err := sw.txtb.Txtw.ResizeWindow(x, y, w, h)
	if err != nil {
		return err
	}
	sw.txtb.CheckSize()
	return nil
}

func (sw *StackWindow) ShowBorder(screenw, screenh int) error {
	return sw.txtb.Txtw.ShowBorder(screenw, screenh)
}

func (sw *StackWindow) WindowXY() (int, int) {
	return sw.txtb.Txtw.WindowXY()
}

func (sw *StackWindow) WindowSize() (int, int) {
	return sw.txtb.Txtw.WindowSize()
}

func (sw *StackWindow) Type() string {
	return "stack"
}

func (sw *StackWindow) SetProp(name string, val rpn.Frame) error {
	switch name {
	case "round":
		v, err := val.BoundedInt(-1, 10)
		if err != nil {
			return err
		}
		sw.round = int8(v)
		return nil
	default:
		return rpn.ErrUnknownProperty
	}
}

func (sw *StackWindow) GetProp(name string) (rpn.Frame, error) {
	switch name {
	case "round":
		return rpn.IntFrame(int64(sw.round), rpn.INTEGER_FRAME), nil
	default:
		return rpn.Frame{}, rpn.ErrUnknownProperty
	}
}

var props = []string{"round"}

// Lists props.  Do not nodify return value.
func (sw *StackWindow) ListProps() []string {
	return props
}

func (sw *StackWindow) Update(rpn *rpn.RPN, unusedForce bool) error {
	w, h := sw.txtb.Txtw.TextSize()
	sw.txtb.CheckSize()
	sw.txtb.Erase()
	framesBack := h
	if len(rpn.Frames) < framesBack {
		framesBack = len(rpn.Frames)
	}
	for i := 0; i < framesBack; i++ {
		f, err := rpn.PeekFrame(i)
		if err != nil {
			return err
		}
		sw.txtb.SetCursorXY(0, h-i-1)
		sw.txtb.TextColor(window.White)
		s := strconv.Itoa(i) + ": "
		if len(s) > w {
			s = s[:w]
		}
		sw.txtb.Print(s, false)
		lw := w - len(s)
		if lw > 0 {
			sw.txtb.TextColor(window.Cyan)
			s := sw.roundedString(f)
			if len(s) > lw {
				s = s[:lw]
			}
			sw.txtb.Print(s, false)
		}
	}
	if (len(rpn.Frames) == 0) && (h > 0) {
		sw.txtb.SetCursorXY(0, h-1)
		sw.txtb.TextColor(window.Cyan)
		s := "Stack Empty"
		if len(s) > w {
			s = s[:w]
		}
		sw.txtb.Print(s, false)
	}
	sw.txtb.Update()
	return nil
}

func (sw *StackWindow) roundedString(f rpn.Frame) string {
	s := f.String(true)
	if !f.IsComplex() || (sw.round < 0) {
		return s
	}
	sw.rsd.reset() // This is done to avoid heap allocations in tinygo

	leftDecimalFn := func() {
		if sw.round > 0 {
			iv, _ := strconv.Atoi(string(sw.rsd.dec[:sw.rsd.didx]))
			if sw.rsd.didx > sw.round {
				iv = (iv + 5) / 10
			}
			for _, b := range strconv.Itoa(int(iv)) {
				sw.rsd.buff[sw.rsd.idx] = byte(b)
				sw.rsd.idx++
			}
		}
		sw.rsd.inDecimal = false
	}

	for _, c := range s {
		if sw.rsd.inDecimal {
			if c == '.' {
				// skip
			} else if (c < '0') || (c > '9') {
				leftDecimalFn()
			} else if sw.rsd.didx <= sw.round {
				sw.rsd.dec[sw.rsd.didx] = byte(c)
				sw.rsd.didx++
			}
		}
		if !sw.rsd.inDecimal {
			sw.rsd.buff[sw.rsd.idx] = byte(c)
			sw.rsd.idx++
			if c == '.' {
				sw.rsd.didx = 0
				sw.rsd.inDecimal = true
				if sw.round == 0 {
					sw.rsd.idx--
				}
			}
		}
	}

	if sw.rsd.inDecimal {
		leftDecimalFn()
	}

	return string(sw.rsd.buff[:sw.rsd.idx])
}
