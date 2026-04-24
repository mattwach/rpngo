// Package commands is creates window management commands
package commands

import (
	"mattwach/rpngo/common/elog"
	"mattwach/rpngo/common/rpn"
	"mattwach/rpngo/common/window"
	"mattwach/rpngo/common/window/stackwin"
	"mattwach/rpngo/common/window/varwin"
)

type WindowCommands struct {
	root            *window.WindowRoot
	screen          window.Screen
	newPlotWindowFn func() (window.WindowWithProps, error)
}

func InitWindowCommands(
	r *rpn.RPN,
	root *window.WindowRoot,
	screen window.Screen,
	newPlotWindowFn func() (window.WindowWithProps, error)) *WindowCommands {
	conceptHelp := map[string]string{
		"window.layout": "Windows are arranged with window groups.  There\n" +
			"is always a window group named 'root' which is the parent of all \n" +
			"windows and groups.\n" +
			"- Add a new window group to the root window with w.new.group.\n" +
			"- Move a window or group to a different window group with w.move.beg and w.move.end\n" +
			"- Change the weight of a window or group with w.weight (default weight is 100).\n" +
			"- Change the layout mode of a window group to columns with w.columns.\n" +
			"- Print info on all existing windows and groups with w.dump.\n" +
			"- You may also set .wtarget, .wend, and .wweight to direct how and\n" +
			"  where the next window/group will be create.  w.reset resets these\n" +
			"  to .wtarget=root, .wend=true, .wweight=100. Using illegal types or\n" +
			"  values for these variables will cause them to revert to the defaults\n" +
			"  as well.\n" +
			"See Also: windows, window.props",

		"window.props": "Each window supports properties that changes how the window operates\n" +
			"- Print all properties and values for a window with w.listp\n" +
			"- Get a single property with w.getp\n" +
			"- Set a single property with w.setp\n" +
			"See Also: windows, window.layout, plotting",

		"windows": "The display can be customized with different windows\n" +
			"- Add a window with a w.new.<type> command. Example: w.new.stack\n" +
			"- Reset to a single window with w.reset.\n" +
			"See Also: window.layout, window.props",
	}
	r.RegisterConceptHelp(conceptHelp)
	elog.Heap("alloc: /window/commands/window.go:50: wc := WindowCommands{root: root, screen: screen, newPlotWindowFn: newPlotWindowFn}")
	wc := WindowCommands{root: root, screen: screen, newPlotWindowFn: newPlotWindowFn} // object allocated on the heap: escapes at line 65
	r.Register("snapshot", wc.snapshot, rpn.CatIO, snapshotHelp)
	r.Register("w.columns", wc.wColumns, rpn.CatWindow, wColumnsHelp)
	r.Register("w.del", wc.wDelete, rpn.CatWindow, wDeleteHelp)
	r.Register("w.dump", wc.wDump, rpn.CatWindow, wDumpHelp)
	r.Register("w.move.beg", wc.wMoveBeg, rpn.CatWindow, wMoveBegHelp)
	r.Register("w.move.end", wc.wMoveEnd, rpn.CatWindow, wMoveEndHelp)
	r.Register("w.new.group", wc.wNewGroup, rpn.CatWindow, wNewGroupHelp)
	r.Register("w.new.plot", wc.wNewPlot, rpn.CatWindow, wNewPlotHelp)
	r.Register("w.new.stack", wc.wNewStack, rpn.CatWindow, wNewStackHelp)
	r.Register("w.new.var", wc.wNewVar, rpn.CatWindow, wNewVarHelp)
	r.Register("w.listp", wc.wListP, rpn.CatWindow, wListPHelp)
	r.Register("w.getp", wc.wGetP, rpn.CatWindow, wGetPHelp)
	r.Register("w.setp", wc.wSetP, rpn.CatWindow, wSetPHelp)
	r.Register("w.reset", wc.wReset, rpn.CatWindow, wResetHelp)
	r.Register("w.snapshot", wc.wSnapshot, rpn.CatWindow, wSnapshotHelp)
	r.Register("w.update", wc.wUpdate, rpn.CatWindow, wUpdateHelp)
	r.Register("w.weight", wc.wWeight, rpn.CatWindow, wWeightHelp)
	return &wc
}

const wUpdateHelp = "Updates the given window or window group"

func (wc *WindowCommands) wUpdate(r *rpn.RPN) error {
	f, err := r.PopFrame()
	if err != nil {
		return err
	}
	if !f.IsString() {
		return rpn.ErrExpectedAString
	}
	return wc.root.UpdateByName(r, f.UnsafeString(), true)
}

const wDumpHelp = "Dump the state of all created windows and groups"

func (wc *WindowCommands) wDump(r *rpn.RPN) error {
	wc.root.Dump(r)
	return nil
}

const wResetHelp = "Resets window configuration to just a single input window"

func (wc *WindowCommands) wReset(r *rpn.RPN) error {
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

func (wc *WindowCommands) wColumns(r *rpn.RPN) error {
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

const wDeleteHelp = "Deletes a window or window group\n" +
	"Example: 'p1' w.del"

func (wc *WindowCommands) wDelete(r *rpn.RPN) error {
	name, err := r.PopFrame()
	if err != nil {
		return err
	}
	if !name.IsString() {
		return rpn.ErrExpectedAString
	}
	return wc.root.DeleteWindowOrGroup(name.UnsafeString())
}

const wMoveBegHelp = "Moves a window or group to the beginning of a window group\n" +
	"Example: 's1' 'root' w.move.beg"

func (wc *WindowCommands) wMoveBeg(r *rpn.RPN) error {
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

func (wc *WindowCommands) wMoveEnd(r *rpn.RPN) error {
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

func (wc *WindowCommands) wNewGroup(r *rpn.RPN) error {
	name, err := wc.newWindowNameFromStack(r)
	if err != nil {
		return err
	}
	wc.root.AddNewWindowGroupChild(r, name)
	return nil
}

const wNewStackHelp = "Creates a new stack window with the given name and\n" +
	"adds it to the root window. Example: 's1' w.new.stack"

func (wc *WindowCommands) wNewStack(r *rpn.RPN) error {
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

const wNewPlotHelp = "Creates a new plot window with the given name and\n" +
	"adds it to the root window. Example: 'p1' w.new.plot"

func (wc *WindowCommands) wNewPlot(r *rpn.RPN) error {
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

const wNewVarHelp = "Creates a new variable window with the given name and\n" +
	"adds it to the root window. Example: 'v1' w.new.var"

func (wc *WindowCommands) wNewVar(r *rpn.RPN) error {
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

func (wc *WindowCommands) newTextWindow(r *rpn.RPN) (window.TextWindow, string, error) {
	name, err := wc.newWindowNameFromStack(r)
	if err != nil {
		return nil, "", err
	}
	txtw, err := wc.screen.NewTextWindow()
	return txtw, name, err
}

func (wc *WindowCommands) newWindowNameFromStack(r *rpn.RPN) (string, error) {
	name, err := r.PopFrame()
	if err != nil {
		return "", err
	}
	if !name.IsString() {
		return "", rpn.ErrExpectedAString
	}
	existing := wc.root.FindWindow(name.UnsafeString())
	if existing != nil {
		return "", rpn.ErrWindowAlreadyExists
	}
	return name.UnsafeString(), nil
}

const wWeightHelp = "Changes the weight of a window or window group causing it\n" +
	"to take more or less screen space. The default value is 100.\n" +
	"Example: 's1' 20 w.weight"

func (wc *WindowCommands) wWeight(r *rpn.RPN) error {
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
