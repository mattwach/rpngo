// A simple console demonstration
package main

import (
	"fmt"
	"mattwach/rpngo/common/drivers/posix/fs"
	"mattwach/rpngo/common/rpn"
	"mattwach/rpngo/common/startup"
	"mattwach/rpngo/common/window/commands"
	"mattwach/rpngo/common/window/plotwin"
	"mattwach/rpngo/drivers/fyne/fynewin"
	"mattwach/rpngo/drivers/readline/window"
	"os"
)

func initRPN(r *rpn.RPN, fw *fynewin.FyneWin) error {
	_ = commands.InitWindowManagerCommands(r, fw)
	_ = plotwin.InitPlotCommands(r, fw, nil)
	if err := startup.Startup(r, &fs.FileOpsDriver{}); err != nil {
		return err
	}
	return nil
}

func interactive(r *rpn.RPN) error {
	var inter startup.Interrupt
	inter.Init()
	r.Interrupt = inter.Interrupt
	fw := fynewin.FyneWin{}
	fw.Init(r)
	var rlw window.ReadlineWindow
	if err := rlw.Init(fw.ExecRPN); err != nil {
		return err
	}
	defer rlw.Close()
	fw.AddChild("i", &rlw, nil)

	if err := initRPN(r, &fw); err != nil {
		return err
	}
	go func() {
		for {
			if err := rlw.ExecLine(); err != nil {
				break
			}
		}

		fw.Shutdown()
	}()

	fw.Run()
	return nil
}

func main() {
	if err := startup.Run(interactive); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
