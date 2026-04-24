// Package varwin shows a variable window
package varwin

import (
	"mattwach/rpngo/common/elog"
	"mattwach/rpngo/common/rpn"
	"mattwach/rpngo/common/window"
	"sort"
	"strconv"
	"strings"
)

type VariableWindow struct {
	txtb      window.TextBuffer
	showdot   bool
	multiline bool
	names     []string
}

func (w *VariableWindow) Init(txtw window.TextWindow) {
	w.txtb.Init(txtw, 0)
	w.txtb.TextColor(window.White)
	elog.Heap("alloc: /window/varwin/varwin.go:21: w.namesAndValues = make([]rpn.NameAndValues, 0, 16)")
	w.names = make([]string, 0, 16) // object allocated on the heap: object size 320 exceeds maximum stack allocation size 256
}

func (vw *VariableWindow) ResizeWindow(x, y, w, h int) error {
	err := vw.txtb.Txtw.ResizeWindow(x, y, w, h)
	if err != nil {
		return err
	}
	vw.txtb.CheckSize()
	return nil
}

func (vw *VariableWindow) ShowBorder(screenw, screenh int) error {
	return vw.txtb.Txtw.ShowBorder(screenw, screenh)
}

func (vw *VariableWindow) WindowXY() (int, int) {
	return vw.txtb.Txtw.WindowXY()
}

func (vw *VariableWindow) WindowSize() (int, int) {
	return vw.txtb.Txtw.WindowSize()
}

func (vw *VariableWindow) Type() string {
	return "var"
}

func (vw *VariableWindow) SetProp(name string, val rpn.Frame) error {
	switch name {
	case "showdot":
		v, err := val.Bool()
		if err != nil {
			return err
		}
		vw.showdot = v
		return nil
	case "multiline":
		v, err := val.Bool()
		if err != nil {
			return err
		}
		vw.multiline = v
		return nil
	default:
		return rpn.ErrUnknownProperty
	}
}

func (vw *VariableWindow) GetProp(name string) (rpn.Frame, error) {
	switch name {
	case "showdot":
		return rpn.BoolFrame(vw.showdot), nil
	case "multiline":
		return rpn.BoolFrame(vw.multiline), nil
	default:
		return rpn.Frame{}, rpn.ErrUnknownProperty
	}
}

var props = []string{"showdot", "multiline"}

func (vw *VariableWindow) ListProps() []string {
	return props
}

func (vw *VariableWindow) Update(r *rpn.RPN, unusedForce bool) error {
	w, h := vw.txtb.Txtw.TextSize()
	vw.txtb.CheckSize()
	vw.txtb.Erase()
	vw.names = r.AppendAllVariableNames(vw.names[:0])
	sort.Strings(vw.names)
	vw.txtb.SetCursorXY(0, 0)
	hidden := 0
	row := 0
	allValues := r.AllVariableNamesAndValues()
	for _, name := range vw.names {
		if !vw.showdot && (name[0] == '.') {
			hidden++
			continue
		}
		if row < (h - 1) {
			val := framesToString(allValues[name])
			if !vw.multiline {
				val = makeSingleLine(val, w-len(name)-2)
			} else {
				row += countCRs(val)
			}
			vw.txtb.TextColor(window.White)
			vw.txtb.Print(name, false)
			vw.txtb.Print(": ", false)
			vw.txtb.TextColor(window.Cyan)
			vw.txtb.Print(val, false)
			vw.txtb.Write('\n', false)
			row++
		}
	}
	vw.txtb.TextColor(window.Blue)
	vw.txtb.SetCursorXY(0, h-1)
	vw.txtb.Print("num: "+strconv.Itoa(len(vw.names))+" hidden: "+strconv.Itoa(hidden), false)
	vw.txtb.Update()
	return nil
}

func countCRs(val string) int {
	n := 0
	for _, c := range val {
		if c == '\n' {
			n++
		}
	}
	return n
}

func framesToString(frames []rpn.Frame) string {
	var parts []string
	for _, f := range frames {
		parts = append(parts, f.String(true))
	}
	return strings.Join(parts, " -> ")
}

func makeSingleLine(line string, width int) string {
	if width < 0 {
		return ""
	}
	if strings.Contains(line, "\n") {
		line = removeCRsAndComments(line)
	}
	if len(line) < width {
		return line
	}
	if width < 4 {
		return line[:width]
	}
	return line[:width-4] + "..."
}

func removeCRsAndComments(line string) string {
	if !strings.Contains(line, "#") {
		// no comments
		return strings.ReplaceAll(line, "\n", " ")
	}
	var parts []string
	for _, part := range strings.Split(line, "\n") {
		commentIdx := strings.Index(part, "#")
		if commentIdx >= 0 {
			part = part[:commentIdx]
		}
		if len(part) == 0 {
			continue
		}
		parts = append(parts, strings.Fields(part)...)
	}
	return strings.Join(parts, " ")
}
