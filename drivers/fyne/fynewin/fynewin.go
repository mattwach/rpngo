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
	"mattwach/rpngo/common/window/common"
	"mattwach/rpngo/drivers/fyne/fynewin/plotwin"
	"mattwach/rpngo/drivers/fyne/fynewin/stackwin"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
)

type FyneWinChild interface {
	Update(r *rpn.RPN)
}

type FyneWin struct {
	wait     chan bool
	ready    chan bool
	fapp     fyne.App
	children map[string]FyneWinChild
}

func (f *FyneWin) Register(r *rpn.RPN) {
	common.RegisterConceptHelp(r, false)
	r.Register("w.new.plot", f.wNewPlot, rpn.CatWindow, common.WNewPlotHelp)
	r.Register("w.new.stack", f.wNewStack, rpn.CatWindow, common.WNewStackHelp)
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
	f.children = make(map[string]FyneWinChild)
	log.Printf("fyne idle")
	<-f.wait
	log.Printf("fyne starting")
	f.fapp = app.New()
	f.fapp.SetIcon(resourceRpngoiconPng)
	// need to create a fyne window and hide it or it will kill rpngo
	// when all windows are cleared
	f.fapp.NewWindow("hidden")
	f.ready <- true
	f.fapp.Run()
}

func (f *FyneWin) Update(r *rpn.RPN) {
	if f.fapp == nil {
		return
	}
	// Update is likely called from the readline goroutine.
	fyne.DoAndWait(func() {
		for _, c := range f.children {
			c.Update(r)
		}
	})
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
		return
	}
	f.fapp.Quit()
	f.fapp = nil
}

func (f *FyneWin) wNew(r *rpn.RPN, prefix string, prepare func(fyne.Window) FyneWinChild) error {
	name, err := common.NewWindowNameFromStack(r)
	if err != nil {
		return err
	}
	existing := f.children[name]
	if existing != nil {
		return rpn.ErrWindowAlreadyExists
	}
	f.makeReady()
	fyne.DoAndWait(func() {
		w := f.fapp.NewWindow(prefix + ": " + name)
		f.children[name] = prepare(w)
		w.Show()
	})
	return nil
}

func (f *FyneWin) wNewStack(r *rpn.RPN) error {
	return f.wNew(r, "stack", func(w fyne.Window) FyneWinChild {
		return stackwin.New(w, r)
	})
}

func (f *FyneWin) wNewPlot(r *rpn.RPN) error {
	return f.wNew(r, "plot", func(w fyne.Window) FyneWinChild {
		return plotwin.New(w, r)
	})
}
