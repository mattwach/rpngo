// A simple console demonstration
package main

import (
	"fmt"
	"mattwach/rpngo/common/functions"
	"mattwach/rpngo/common/rpn"
	"os"
)

func run() error {
	var r rpn.RPN
	r.Init(256)
	functions.RegisterAll(&r)

	if err := r.ExecSlice(os.Args[1:]); err != nil {
		return err
	}

	for _, f := range r.Frames {
		fmt.Println(f.String(true))
	}

	return nil
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
