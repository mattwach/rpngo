// A simple console demonstration
package main

import (
	"fmt"
	"mattwach/rpngo/common/drivers/posix/fs"
	"mattwach/rpngo/common/fileops"
	"mattwach/rpngo/common/parse"
	"mattwach/rpngo/common/rpn"
	"mattwach/rpngo/common/startup"
	"mattwach/rpngo/common/window/commands"
	"mattwach/rpngo/common/window/plotwin"
	"mattwach/rpngo/drivers/fyne/fynewin"
	"os"
	"path/filepath"

	"github.com/chzyer/readline"
)

func execLine(r *rpn.RPN, rl *readline.Instance) error {
	line, err := rl.Readline()
	if err != nil {
		return err
	}
	err = parse.Fields(line, r.Exec)
	if err != nil {
		fmt.Printf("error: %v\n", err)
	}

	for _, f := range r.Frames {
		fmt.Println(f.String(true))
	}

	return nil
}

func interactive(r *rpn.RPN) error {
	histFile := startup.HistFile
	home, err := fileops.HomeDir()
	if err == nil {
		histFile = filepath.Join(home, startup.HistFile)
	}
	rl, err := readline.NewEx(&readline.Config{
		Prompt:      "> ",
		HistoryFile: histFile,
	})
	if err != nil {
		panic(err)
	}
	defer rl.Close()
	fw := fynewin.FyneWin{}
	fw.Register(r)

	_ = commands.InitWindowManagerCommands(r, &fw)
	_ = plotwin.InitPlotCommands(r, &fw, nil)
	if err := startup.Startup(r, &fs.FileOpsDriver{}); err != nil {
		return err
	}

	go func() {
		for {
			if err := execLine(r, rl); err != nil {
				break
			}
			fw.Update(r, 0, 0, true)
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
