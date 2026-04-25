// A simple console demonstration
package main

import (
	"fmt"
	"mattwach/rpngo/common/fileops"
	"mattwach/rpngo/common/parse"
	"mattwach/rpngo/common/rpn"
	"mattwach/rpngo/common/startup"
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

	for {
		if err := execLine(r, rl); err != nil {
			break
		}
	}

	return nil
}

func main() {
	if err := startup.Run(interactive); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
