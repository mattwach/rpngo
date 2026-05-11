package varwin

import (
	"fmt"
	"mattwach/rpngo/common/parse"
	"mattwach/rpngo/common/rpn"
	"sort"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// VarWin holds the context for a stack window.
// Important, RPN is owned by the readline goroutine thus should be accessed
// with care.  This means that putting a pointer to it in this struct is
// probably starting down a bad path.
type VarWin struct {
	rpn               *rpn.RPN
	win               fyne.Window
	data              *widget.Entry
	showdot           bool
	multiline         bool
	showDotCheckbox   *widget.Check
	multilineCheckbox *widget.Check
}

func New(win fyne.Window, r *rpn.RPN) *VarWin {
	vw := &VarWin{
		rpn:  r,
		win:  win,
		data: widget.NewMultiLineEntry(),
	}
	vw.data.TextStyle = fyne.TextStyle{Monospace: true}
	scroll := container.NewScroll(vw.data)
	vw.showDotCheckbox = widget.NewCheck("show dot", func(b bool) {
		vw.showdot = b
		vw.Update(vw.rpn, true)
	})
	vw.multilineCheckbox = widget.NewCheck("multiline", func(b bool) {
		vw.multiline = b
		vw.Update(vw.rpn, true)
	})
	bottom := container.NewHBox(vw.showDotCheckbox, vw.multilineCheckbox)
	win.SetContent(container.NewBorder(nil, bottom, nil, nil, scroll))
	win.Resize(fyne.NewSize(640, 800))
	vw.Update(r, true)
	return vw
}

func (vw *VarWin) ResizeWindow(x, y, w, h int) error {
	// Not used in fyne
	return nil
}

func (vw *VarWin) ShowBorder(w, h int) error {
	// Not used in fyne
	return nil
}

func (vw *VarWin) WindowXY() (int, int) {
	// Not used in fyne
	return 0, 0
}

func (vw *VarWin) WindowSize() (int, int) {
	s := vw.data.Size()
	return int(s.Width), int(s.Height)
}

func (vw *VarWin) Type() string {
	return "var"
}

func (vw *VarWin) SetProp(name string, val rpn.Frame) error {
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

func (vw *VarWin) GetProp(name string) (rpn.Frame, error) {
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

func (vw *VarWin) ListProps() []string {
	return props
}

const maxLineLength = 256

func (vw *VarWin) Update(r *rpn.RPN, force bool) error {
	names := r.AppendAllVariableNames(nil)
	if len(names) > 0 {
		sort.Strings(names)
		allValues := r.AllVariableNamesAndValues()
		var lines []string
		for _, name := range names {
			if !vw.showdot && (name[0] == '.') {
				continue
			}
			val := rpn.FramesToString(allValues[name])
			if !vw.multiline {
				val = parse.MakeSingleLine(val, maxLineLength)
			}
			lines = append(lines, fmt.Sprintf("%-15s | %s", name, val))
		}
		vw.data.SetText(strings.Join(lines, "\n"))
	} else {
		vw.data.SetText("no variables defined")
	}
	vw.multilineCheckbox.SetChecked(vw.multiline)
	vw.showDotCheckbox.SetChecked(vw.showdot)
	return nil
}
