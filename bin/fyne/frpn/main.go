// A simple console demonstration
package main

import (
	"fmt"
	"log"
	"mattwach/rpngo/common/drivers/posix/fs"
	"mattwach/rpngo/common/rpn"
	"mattwach/rpngo/common/startup"
	"mattwach/rpngo/common/window/commands"
	"mattwach/rpngo/common/window/plotwin"
	"mattwach/rpngo/drivers/fyne/fynewin"
	"mattwach/rpngo/drivers/readline/window"
	"os"
)

func updateFyne(rpnInst chan *rpn.RPN, fw *fynewin.FyneWin) {
	r := <-rpnInst
	defer func() {
		rpnInst <- r
	}()
	fw.Update(r, 0, 0, true)
}

func interactive(r *rpn.RPN) error {
	var inter startup.Interrupt
	inter.Init()
	log.Printf("inter is %d", &inter)
	r.Interrupt = inter.Interrupt
	// pack rpn into a channel so it can be shared between the readline and fyne
	// goroutines.
	rpnInst := make(chan *rpn.RPN, 1)
	rpnInst <- r
	return interactiveChannel(rpnInst, &inter)
}

func initRPN(rpnInst chan *rpn.RPN, fw *fynewin.FyneWin) error {
	r := <-rpnInst
	defer func() {
		rpnInst <- r
	}()
	_ = commands.InitWindowManagerCommands(r, fw)
	_ = plotwin.InitPlotCommands(r, fw, nil)
	if err := startup.Startup(r, &fs.FileOpsDriver{}); err != nil {
		return err
	}
	return nil
}

func interactiveChannel(rpnInst chan *rpn.RPN, interrupt *startup.Interrupt) error {
	var rlw window.ReadlineWindow
	if err := rlw.Init(rpnInst); err != nil {
		return err
	}
	defer rlw.Close()
	fw := fynewin.FyneWin{}
	fw.Init(rpnInst, interrupt)
	fw.AddChild("i", &rlw, nil)

	if err := initRPN(rpnInst, &fw); err != nil {
		return err
	}
	go func() {
		for {
			if err := rlw.ExecLine(); err != nil {
				break
			}
			updateFyne(rpnInst, &fw)
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
