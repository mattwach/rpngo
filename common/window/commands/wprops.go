package commands

import (
	"mattwach/rpngo/common/rpn"
)

const wListPHelp = "Prints all properties / values for a window\n" +
	"Example 'p1' w.listp"

func (wc *WindowCommands) wListP(r *rpn.RPN) error {
	wname, err := r.PopFrame()
	if err != nil {
		return err
	}
	if !wname.IsString() {
		return rpn.ErrExpectedAString
	}
	w := wc.root.FindWindow(wname.UnsafeString())
	if w == nil {
		return rpn.ErrNotFound
	}
	for _, p := range w.ListProps() {
		f, err := w.GetProp(p)
		if err != nil {
			// unexpected
			return err
		}
		r.Print(p)
		r.Print(": ")
		r.Println(f.String(true))
	}
	return nil
}

const wGetPHelp = "Pushes the value of the given property to the stack.\n" +
	"Example: 'p1' 'minx' w.getp"

func (wc *WindowCommands) wGetP(r *rpn.RPN) error {
	wname, pname, err := r.Pop2Frames()
	if err != nil {
		return err
	}
	if !wname.IsString() || !pname.IsString() {
		return rpn.ErrExpectedAString
	}
	w := wc.root.FindWindow(wname.UnsafeString())
	if w == nil {
		return rpn.ErrNotFound
	}
	f, err := w.GetProp(pname.UnsafeString())
	if err != nil {
		return err
	}
	return r.PushFrame(f)
}

const wSetPHelp = "Sets a property on a window.\n" +
	"Example: 'p1' 'minx' -1 w.setp"

func (wc *WindowCommands) wSetP(r *rpn.RPN) error {
	f, err := r.PopFrame()
	if err != nil {
		return err
	}
	wname, pname, err := r.Pop2Frames()
	if err != nil {
		return err
	}
	if !wname.IsString() || !pname.IsString() {
		return rpn.ErrExpectedAString
	}
	w := wc.root.FindWindow(wname.UnsafeString())
	if w == nil {
		return rpn.ErrNotFound
	}
	return w.SetProp(pname.UnsafeString(), f)
}

const wSnapshotHelp = "Creates a set of commands that will restore window properties"

func (wc *WindowCommands) wSnapshot(r *rpn.RPN) error {
	wname, err := r.PopFrame()
	if err != nil {
		return err
	}
	if !wname.IsString() {
		return rpn.ErrExpectedAString
	}
	buff := make([]byte, 0, 128)
	buff, err = wc.root.Snapshot(buff, wname.UnsafeString())
	if err != nil {
		return err
	}
	return r.PushFrame(rpn.StringFrame(string(buff), rpn.STRING_BRACE_FRAME))
}

const snapshotHelp = "Generates a string that will restore the current state of the calculator"

func (wc *WindowCommands) snapshot(r *rpn.RPN) error {
	buff := make([]byte, 0, 256)
	buff = append(buff, []byte("d\nv.clearall\nw.reset\n")...)
	buff, _ = wc.root.Snapshot(buff, "root")
	buff = r.VarSnapshot(buff)
	buff = r.StackSnapshot(buff)
	return r.PushFrame(rpn.StringFrame(string(buff), rpn.STRING_BRACE_FRAME))
}
