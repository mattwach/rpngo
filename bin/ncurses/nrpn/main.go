// A simple console demonstration
package main

import (
	"errors"
	"fmt"
	"mattwach/rpngo/common/drivers/posix/fs"
	"mattwach/rpngo/common/rpn"
	"mattwach/rpngo/common/startup"
	"mattwach/rpngo/common/window"
	"mattwach/rpngo/common/window/commands"
	"mattwach/rpngo/common/window/input"
	"mattwach/rpngo/common/window/plotwin"
	"mattwach/rpngo/drivers/curses"
	"os"
)

const scrollbytes = 256 * 1024

func interactive(r *rpn.RPN) error {
	var inter startup.Interrupt
	inter.Init()
	r.Interrupt = inter.Interrupt
	screen, err := curses.Init()
	if err != nil {
		return err
	}
	defer screen.End()
	var root window.WindowRoot
	err = buildUI(&root, screen, r)
	if err != nil {
		return err
	}
	newTextPlotWindow := func() (window.WindowWithProps, error) {
		var tpw plotwin.TxtPlotWindow
		pw, err := screen.NewTextWindow()
		if err != nil {
			return nil, err
		}
		tpw.Init(pw)
		return &tpw, nil
	}
	_ = commands.InitWindowManagerCommands(r, &root)
	_ = commands.InitWindowRootCommands(r, &root, screen, newTextPlotWindow)
	_ = plotwin.InitPlotCommands(r, &root, screen)
	if err := startup.Startup(r, &fs.FileOpsDriver{}); err != nil {
		return err
	}
	w, h := screen.ScreenSize()
	if err := root.Update(r, w, h, true); err != nil {
		return err
	}
	for {
		w, h = screen.ScreenSize()
		if err := root.Update(r, w, h, true); err != nil {
			if errors.Is(err, input.ErrExit) {
				return nil
			}
			return err
		}
	}
}

func buildUI(root *window.WindowRoot, screen *curses.Curses, r *rpn.RPN) error {
	w, h := screen.ScreenSize()
	root.Init(w, h)
	if err := addInputWindow(screen, root, r); err != nil {
		return err
	}
	return nil
}

func addInputWindow(screen window.Screen, root *window.WindowRoot, r *rpn.RPN) error {
	txtw, err := screen.NewTextWindow()
	if err != nil {
		return err
	}
	var iw input.InputWindow
	iw.Init(txtw.(*curses.Curses), txtw, r, &fs.FileOpsDriver{}, scrollbytes)
	if err != nil {
		return err
	}
	root.AddWindowChildToRoot(&iw, "i", 100)
	return nil
}

func main() {
	if err := startup.Run(interactive); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
