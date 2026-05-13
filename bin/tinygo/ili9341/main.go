// Run rpngo on a microcontroller
//
// This is a "minimialist" implementation which can be thought of
// as a valdation stepping stone.
package main

import (
	"errors"
	"mattwach/rpngo/common/elog"
	"mattwach/rpngo/common/fileops"
	"mattwach/rpngo/common/functions"
	"mattwach/rpngo/common/key"
	"mattwach/rpngo/common/rpn"
	"mattwach/rpngo/common/startup"
	"mattwach/rpngo/common/window"
	"mattwach/rpngo/common/window/commands"
	"mattwach/rpngo/common/window/input"
	"mattwach/rpngo/common/window/plotwin"
	"mattwach/rpngo/common/xmodem"
	"mattwach/rpngo/drivers/tinygo"
	"mattwach/rpngo/drivers/tinygo/ili9341"
	"mattwach/rpngo/drivers/tinygo/pixelwinbuffer"
	"mattwach/rpngo/drivers/tinygo/serial"
	"mattwach/rpngo/drivers/tinygo/tinyfs"
	"time"
)

const scrollbytes = 8 * 1024
const maxStackDepth = 256

// Persistant globals.  Tinygo curently can't move heap pointers,
// which leads to heap fragmentation and thus the heap should
// be used sparingly
var rpnInst rpn.RPN
var getInputInst getInput
var inputWin input.InputWindow
var screen ili9341.Ili9341Screen
var root window.WindowRoot
var fileOps fileops.FileOps
var fileOpsDriver tinyfs.FileOpsDriver
var xmodemCommands xmodem.XmodemCommands

type getInput struct {
	lcd    *ili9341.Ili9341TxtW
	serial serial.Serial
}

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	time.Sleep(200 * time.Millisecond)

	elog.Print("Started")
	rpnInst.Init(maxStackDepth)
	tinygo.Register(&rpnInst)
	functions.RegisterAll(&rpnInst)
	_ = fileOpsDriver.Init()
	fileOps.InitAndRegister(&rpnInst, 65536, &fileOpsDriver)

	screen.Init()
	err := buildUI(&root, &screen, &rpnInst)
	if err != nil {
		return err
	}
	newPixelPlotWindow := func() (window.WindowWithProps, error) {
		elog.Heap("alloc: /bin/tinygo/ili9341/main.go:67: var ppw plotwin.PixelPlotWindow")
		var ppw plotwin.PixelPlotWindow // object allocated on the heap: escapes at line 67
		pw, err := screen.NewPixelWindow()
		if err != nil {
			return nil, err
		}
		elog.Heap("alloc: /bin/tinygo/ili9341/main.go:72: var pb pixelwinbuffer.PixelBuffer")
		var pb pixelwinbuffer.PixelBuffer // object allocated on the heap: escapes at line 74
		pb.Init(pw)
		ppw.Init(&pb, 256)
		return &ppw, nil
	}
	_ = commands.InitWindowManagerCommands(&rpnInst, &root)
	_ = commands.InitWindowRootCommands(&rpnInst, &root, &screen, newPixelPlotWindow)
	_ = plotwin.InitPlotCommands(&rpnInst, &root, &screen)
	starterr := startup.Startup(&rpnInst, &fileOpsDriver)
	w, h := screen.ScreenSize()
	if err := root.Update(&rpnInst, w, h, true); err != nil {
		return err
	}
	if starterr != nil {
		rpnInst.Println(starterr.Error())
	}
	for {
		w, h = screen.ScreenSize()
		if err := root.Update(&rpnInst, w, h, true); err != nil {
			if errors.Is(err, rpn.ErrExit) {
				return nil
			}
			return err
		}
	}
}

func buildUI(root *window.WindowRoot, screen window.Screen, r *rpn.RPN) error {
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
	inputWin.Init(&getInputInst, txtw, r, &fileOpsDriver, scrollbytes)
	getInputInst.Init(txtw.(*ili9341.Ili9341TxtW))
	xmodemCommands.InitAndRegister(&rpnInst, &getInputInst.serial)
	root.AddWindowChildToRoot(&inputWin, "i", 100)
	return nil
}

func (gi *getInput) Init(lcd *ili9341.Ili9341TxtW) {
	gi.lcd = lcd
	gi.serial.Init(true)
}

func (gi *getInput) GetChar() (key.Key, error) {
	for {
		time.Sleep(20 * time.Millisecond)
		k := gi.serial.GetChar()
		if k != 0 {
			return k, nil
		}
	}
}
