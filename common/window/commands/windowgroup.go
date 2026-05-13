// Package commands is creates window management commands
package commands

import (
	"mattwach/rpngo/common/elog"
	"mattwach/rpngo/common/rpn"
	"mattwach/rpngo/common/window"
	"mattwach/rpngo/common/window/common"
	"mattwach/rpngo/common/window/stackwin"
	"mattwach/rpngo/common/window/varwin"
)

type WindowRootCommands struct {
	root            *window.WindowRoot
	screen          window.Screen
	newPlotWindowFn func() (window.WindowWithProps, error)
}

func InitWindowRootCommands(
	r *rpn.RPN,
	root *window.WindowRoot,
	screen window.Screen,
	newPlotWindowFn func() (window.WindowWithProps, error)) *WindowRootCommands {
	wc := WindowRootCommands{root: root, screen: screen, newPlotWindowFn: newPlotWindowFn}
	r.Register("w.columns", wc.wColumns, rpn.CatWindow, wColumnsHelp)
	r.Register("w.move.beg", wc.wMoveBeg, rpn.CatWindow, wMoveBegHelp)
	r.Register("w.move.end", wc.wMoveEnd, rpn.CatWindow, wMoveEndHelp)
	r.Register("w.new.group", wc.wNewGroup, rpn.CatWindow, wNewGroupHelp)
	r.Register("w.new.plot", wc.wNewPlot, rpn.CatWindow, common.WNewPlotHelp)
	r.Register("w.new.stack", wc.wNewStack, rpn.CatWindow, common.WNewStackHelp)
	r.Register("w.new.var", wc.wNewVar, rpn.CatWindow, common.WNewVarHelp)
	r.Register("w.reset", wc.wReset, rpn.CatWindow, common.WResetHelp)
	r.Register("w.weight", wc.wWeight, rpn.CatWindow, wWeightHelp)
	return &wc
}

func (wc *WindowRootCommands) wReset(r *rpn.RPN) error {
	r.PushFrame(rpn.StringFrame("root", rpn.STRING_SINGLEQ_FRAME))
	r.SetVariable(".wtarget")
	r.PushFrame(rpn.BoolFrame(true))
	r.SetVariable(".wend")
	r.PushFrame(rpn.IntFrame(100, rpn.INTEGER_FRAME))
	r.SetVariable(".wweight")
	iw := wc.root.FindWindow("i")
	if iw == nil {
		return rpn.ErrInputWindowNotFound
	}
	wc.root.RemoveAllChildren()
	wc.root.UseColumnLayout("root", false)
	wc.root.AddWindowChildToRoot(iw, "i", 100)
	return nil
}

const wColumnsHelp = "Sets a window group layout to column mode\n" +
	"Example: 'g1' w.columns"

func (wc *WindowRootCommands) wColumns(r *rpn.RPN) error {
	name, err := r.PopFrame()
	if err != nil {
		return err
	}
	if !name.IsString() {
		return rpn.ErrExpectedAString
	}
	if err := wc.root.UseColumnLayout(name.UnsafeString(), true); err != nil {
		return err
	}
	return nil
}

const wMoveBegHelp = "Moves a window or group to the beginning of a window group\n" +
	"Example: 's1' 'root' w.move.beg"

func (wc *WindowRootCommands) wMoveBeg(r *rpn.RPN) error {
	src, dst, err := r.Pop2Frames()
	if err != nil {
		return err
	}
	if !src.IsString() || !dst.IsString() {
		return rpn.ErrExpectedAString
	}
	return wc.root.MoveWindowOrGroup(src.UnsafeString(), dst.UnsafeString(), true)
}

const wMoveEndHelp = "Moves a window or group to the end of a window group\n" +
	"Example: 's1' 'root' w.move.end"

func (wc *WindowRootCommands) wMoveEnd(r *rpn.RPN) error {
	src, dst, err := r.Pop2Frames()
	if err != nil {
		return err
	}
	if !src.IsString() || !dst.IsString() {
		return rpn.ErrExpectedAString
	}
	return wc.root.MoveWindowOrGroup(src.UnsafeString(), dst.UnsafeString(), false)
}

const wNewGroupHelp = "Creates a new window group with the given name and\n" +
	"adds it to the root window. Example: 'g1' w.new.group"

func (wc *WindowRootCommands) wNewGroup(r *rpn.RPN) error {
	name, err := wc.newWindowNameFromStack(r)
	if err != nil {
		return err
	}
	wc.root.AddNewWindowGroupChild(r, name)
	return nil
}

func (wc *WindowRootCommands) wNewStack(r *rpn.RPN) error {
	txtw, name, err := wc.newTextWindow(r)
	if err != nil {
		return err
	}
	elog.Heap("alloc: /window/commands/window.go:187: var sw stackwin.StackWindow")
	var sw stackwin.StackWindow // object allocated on the heap: escapes at line 189
	sw.Init(txtw)
	wc.root.AddWindowChild(r, &sw, name)
	return nil
}

func (wc *WindowRootCommands) wNewPlot(r *rpn.RPN) error {
	name, err := wc.newWindowNameFromStack(r)
	if err != nil {
		return err
	}
	pw, err := wc.newPlotWindowFn()
	if err != nil {
		return err
	}
	wc.root.AddWindowChild(r, pw, name)
	return nil
}

func (wc *WindowRootCommands) wNewVar(r *rpn.RPN) error {
	txtw, name, err := wc.newTextWindow(r)
	if err != nil {
		return err
	}
	elog.Heap("alloc: /window/commands/window.go:217: var vw varwin.VariableWindow")
	var vw varwin.VariableWindow // object allocated on the heap: escapes at line 219
	vw.Init(txtw)
	wc.root.AddWindowChild(r, &vw, name)
	return nil
}

func (wc *WindowRootCommands) newTextWindow(r *rpn.RPN) (window.TextWindow, string, error) {
	name, err := wc.newWindowNameFromStack(r)
	if err != nil {
		return nil, "", err
	}
	txtw, err := wc.screen.NewTextWindow()
	return txtw, name, err
}

func (wc *WindowRootCommands) newWindowNameFromStack(r *rpn.RPN) (string, error) {
	name, err := common.NewWindowNameFromStack(r)
	if err != nil {
		return "", err
	}
	existing := wc.root.FindWindow(name)
	if existing != nil {
		return "", rpn.ErrWindowAlreadyExists
	}
	return name, nil
}

const wWeightHelp = "Changes the weight of a window or window group causing it\n" +
	"to take more or less screen space. The default value is 100.\n" +
	"Example: 's1' 20 w.weight"

func (wc *WindowRootCommands) wWeight(r *rpn.RPN) error {
	cw, err := r.PopFrame()
	if err != nil {
		return err
	}
	w, err := cw.Int()
	if err != nil {
		return err
	}
	name, err := r.PopFrame()
	if err != nil {
		return err
	}
	if !name.IsString() {
		return rpn.ErrExpectedAString
	}
	if err := wc.root.SetWindowWeight(name.UnsafeString(), int(w)); err != nil {
		return err
	}
	return nil
}
