// fynewin is the window manager for fyne.
//
// I decided that the fyne UI is different enough from ncurses and picocalc
// to use mostly it's own logic for window managerment (the ncurses/picocalc
// window management is located at common/window/window.go).  There is
// still some useful common logic there, however and it might be an
// iterative process to evolve to a good overall factoring.
package fynewin

import (
	"log"
	"mattwach/rpngo/common/rpn"
	"mattwach/rpngo/common/window"
	"mattwach/rpngo/common/window/common"
	"mattwach/rpngo/common/window/plotwin"
	"mattwach/rpngo/drivers/fyne/fynewin/fyneplotwin"
	"mattwach/rpngo/drivers/fyne/fynewin/stackwin"
	"sort"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/theme"
)

type child struct {
	wprops window.WindowWithProps
	fwin   fyne.Window
}

type FyneWin struct {
	wait     chan bool
	ready    chan bool
	rpnInst  chan *rpn.RPN
	fapp     fyne.App
	children map[string]child
}

func (f *FyneWin) Register(rpnInst chan *rpn.RPN) {
	if f.children == nil {
		f.children = make(map[string]child)
	}
	f.rpnInst = rpnInst
	r := <-f.rpnInst
	defer func() {
		f.rpnInst <- r
	}()
	common.RegisterConceptHelp(r, false)
	r.Register("w.new.plot", f.wNewPlot, rpn.CatWindow, common.WNewPlotHelp)
	r.Register("w.new.stack", f.wNewStack, rpn.CatWindow, common.WNewStackHelp)
}

func (f *FyneWin) AddChild(name string, wprops window.WindowWithProps, win fyne.Window) error {
	if f.children == nil {
		f.children = make(map[string]child)
	}
	_, ok := f.children[name]
	if ok {
		return rpn.ErrWindowAlreadyExists
	}
	f.children[name] = child{wprops, win}
	return nil
}

// Run starts fyne and blocks.
//
// Fyne only works when run from the main go routine.  Yet
// we often don't even need to use fyne at all.  The solution
// is to have the main routine block on a channel until the
// time an initial window is created.
func (f *FyneWin) Run() {
	f.wait = make(chan bool, 1)
	f.ready = make(chan bool, 1)
	log.Printf("fyne idle")
	start := <-f.wait
	if !start {
		// fyne was never needed
		return
	}

	log.Printf("fyne starting")
	f.fapp = app.New()
	f.fapp.Settings().SetTheme(theme.DarkTheme())
	f.fapp.SetIcon(resourceRpngoiconPng)
	// need to create a fyne window and hide it or it will kill rpngo
	// when all windows are cleared
	f.fapp.NewWindow("hidden")
	f.ready <- true
	f.fapp.Run()
}

// Note that sw, sh and updateInput are for interface compatibility with
// with plotwin.WindowManager interface
func (f *FyneWin) Update(r *rpn.RPN, sw, sh int, updateInput bool) error {
	if f.fapp == nil {
		return nil
	}
	var err error
	for _, c := range f.children {
		err = c.wprops.Update(r, true)
		if err != nil {
			break
		}
	}
	return err
}

func (f *FyneWin) UpdateByName(r *rpn.RPN, name string, force bool) error {
	w := f.children[name].wprops
	if w == nil {
		return rpn.ErrNotFound
	}
	return w.Update(r, force)
}

func (f *FyneWin) Snapshot(buff []byte, name string) ([]byte, error) {
	// TODO: implement
	return buff, nil
}

// Needed for compatibility with the plorwin.WindowManager interface
func (f *FyneWin) FindWindow(name string) window.WindowWithProps {
	return f.children[name].wprops
}

// Needed for compatibility with the plorwin.WindowManager interface
func (f *FyneWin) DeleteWindowOrGroup(name string) error {
	if name == "i" {
		return rpn.ErrCanNotDeleteInputWindow
	}
	w := f.children[name].fwin
	if w == nil {
		return rpn.ErrNotFound
	}
	w.Close()
	delete(f.children, name)
	return nil
}

// Needed for compatibility with the plorwin.WindowManager interface
func (f *FyneWin) Dump(r *rpn.RPN) {
	var names []string
	for name := range f.children {
		names = append(names, name)
	}
	sort.Strings(names)
	for _, name := range names {
		c := f.children[name]
		w, h := c.wprops.WindowSize()
		pad(r, 1)
		r.Print(name)
		r.Print("(type=")
		r.Print(c.wprops.Type())
		r.Print(", w=")
		r.Print(strconv.Itoa(w))
		r.Print(", h=")
		r.Print(strconv.Itoa(h))
		r.Print(")\n")
	}
}

func pad(r *rpn.RPN, indent int) {
	for range indent {
		r.Print("  ")
	}
}

func (f *FyneWin) makeReady() {
	if f.fapp == nil {
		log.Printf("signal fyne start")
		f.wait <- true
		<-f.ready
	}
}

func (f *FyneWin) Shutdown() {
	if f.fapp == nil {
		f.wait <- false
		return
	}
	fyne.DoAndWait(func() {
		f.fapp.Quit()
	})
	f.fapp = nil
}

func (f *FyneWin) wNew(r *rpn.RPN, prefix string, prepare func(fyne.Window) window.WindowWithProps) error {
	name, err := common.NewWindowNameFromStack(r)
	if err != nil {
		return err
	}
	_, ok := f.children[name]
	if ok {
		return rpn.ErrWindowAlreadyExists
	}
	f.makeReady()
	fyne.DoAndWait(func() {
		w := f.fapp.NewWindow(prefix + ": " + name)
		w.SetOnClosed(func() {
			delete(f.children, name)
		})
		f.children[name] = child{prepare(w), w}
		w.Show()
	})
	return nil
}

func (f *FyneWin) wNewStack(r *rpn.RPN) error {
	return f.wNew(r, "stack", func(w fyne.Window) window.WindowWithProps {
		return stackwin.New(w, r)
	})
}

func (f *FyneWin) wNewPlot(r *rpn.RPN) error {
	return f.wNew(r, "plot", func(w fyne.Window) window.WindowWithProps {
		ppw := &plotwin.PixelPlotWindow{}
		ppw.Init(fyneplotwin.New(w, f.rpnInst, ppw), 4096)
		return ppw
	})
}
