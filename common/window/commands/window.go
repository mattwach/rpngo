// Package commands is creates window management commands
package commands

import (
	"mattwach/rpngo/common/elog"
	"mattwach/rpngo/common/rpn"
	"mattwach/rpngo/common/window"
	"mattwach/rpngo/common/window/common"
)

type WindowManagerCommands struct {
	manager window.WindowManager
}

func InitWindowManagerCommands(
	r *rpn.RPN,
	manager window.WindowManager) *WindowManagerCommands {
	common.RegisterConceptHelp(r, true)
	elog.Heap("alloc: /window/commands/window.go:50: wc := WindowCommands{root: root, screen: screen, newPlotWindowFn: newPlotWindowFn}")
	wc := WindowManagerCommands{manager: manager}
	r.Register("snapshot", wc.snapshot, rpn.CatIO, snapshotHelp)
	r.Register("w.del", wc.wDelete, rpn.CatWindow, wDeleteHelp)
	r.Register("w.dump", wc.wDump, rpn.CatWindow, wDumpHelp)
	r.Register("w.listp", wc.wListP, rpn.CatWindow, wListPHelp)
	r.Register("w.getp", wc.wGetP, rpn.CatWindow, wGetPHelp)
	r.Register("w.setp", wc.wSetP, rpn.CatWindow, wSetPHelp)
	r.Register("w.snapshot", wc.wSnapshot, rpn.CatWindow, wSnapshotHelp)
	r.Register("w.update", wc.wUpdate, rpn.CatWindow, wUpdateHelp)
	return &wc
}

const wUpdateHelp = "Updates the given window or window group"

func (wc *WindowManagerCommands) wUpdate(r *rpn.RPN) error {
	f, err := r.PopFrame()
	if err != nil {
		return err
	}
	if !f.IsString() {
		return rpn.ErrExpectedAString
	}
	return wc.manager.UpdateByName(r, f.UnsafeString(), true)
}

const wDumpHelp = "Dump the state of all created windows and groups"

func (wc *WindowManagerCommands) wDump(r *rpn.RPN) error {
	wc.manager.Dump(r)
	return nil
}

const wDeleteHelp = "Deletes a window or window group\n" +
	"Example: 'p1' w.del"

func (wc *WindowManagerCommands) wDelete(r *rpn.RPN) error {
	name, err := r.PopFrame()
	if err != nil {
		return err
	}
	if !name.IsString() {
		return rpn.ErrExpectedAString
	}
	return wc.manager.DeleteWindowOrGroup(name.UnsafeString())
}
