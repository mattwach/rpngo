// Handles input heystrokes from the user. e.g.an "editor"
//
// Many of the implementation choices below are to avoid heap allocations
// on microcontrollers.  Notably the use of []byte to hold the result
// as well as the extra code to keep from converting between bytes
// and strings (which forces a heap allocation).  Heap allocations are
// not forbidden of course, as this is an interactive session, but since
// line entry is character-based anyway, using a []byte as a baseline
// makes sense even if heap allocations were not a concern.
package input

import (
	"bytes"
	"errors"
	"mattwach/rpngo/common/elog"
	"mattwach/rpngo/common/fileops"
	"mattwach/rpngo/common/key"
	"mattwach/rpngo/common/rpn"
	"mattwach/rpngo/common/window"
	"path/filepath"
	"strings"
)

var errDisableAutoHistoryFirst = errors.New("disable auto history first")

const MAX_HISTORY_LINES = 500

type getLine struct {
	insertMode   bool
	input        Input
	txtb         *window.TextBuffer
	history      [MAX_HISTORY_LINES][]byte
	historyCount int
	histpath     string
	autoHistory  bool
	fs           fileops.FileOpsDriver
	// line is the current line.  It's kept here to support entering
	// scrolling mode without losing the current line contents.
	names []string
	line  []byte
	// cached list of files for tab complete
	fileList []string
}

const histFile = ".rpngo_history"

func initGetLine(input Input, txtb *window.TextBuffer, fs fileops.FileOpsDriver) *getLine {
	elog.Heap("alloc: /window/input/getline.go:41: gl := &getLine{")
	elog.Heap("alloc: /window/input/getline.go:46: namesAndValues: make([]rpn.NameAndValues, 0, 8),")
	var histpath string
	home, err := fileops.HomeDir()
	if err == nil {
		histpath = filepath.Join(home, histFile)
	}
	gl := &getLine{ // object allocated on the heap: object size 6048 exceeds maximum stack allocation size 256
		insertMode:   true,
		input:        input,
		txtb:         txtb,
		historyCount: 0,
		histpath:     histpath,
		fs:           fs,
		fileList:     make([]string, 0, 16),
		names:        make([]string, 0, 16), // object allocated on the heap: escapes at line 48
	}
	return gl
}

const histLoadHelp = "Loads in-memory history from 'histpath' (an input window property)"

func (gl *getLine) histLoad(r *rpn.RPN) error {
	if gl.autoHistory {
		return errDisableAutoHistoryFirst
	}
	return gl.loadHistory()
}

func (gl *getLine) loadHistory() error {
	data, err := gl.fs.ReadFile(gl.histpath)
	if err != nil {
		elog.Heap("alloc: /window/input/getline.go:69: elog.Print('Could not read hitory file: ', err.Error())")
		elog.Print("Could not read hitory file: ", err.Error()) // object allocated on the heap: escapes at line 69
		return err
	}
	for _, line := range strings.Split(string(data), "\n") {
		line := strings.TrimSpace(line)
		if len(line) == 0 {
			continue
		}
		hidx := gl.historyCount % MAX_HISTORY_LINES
		if len(gl.history[hidx]) > 0 {
			// history has wrapped the internal buffer
			gl.history[hidx] = gl.history[hidx][:0]
		}
		for _, c := range line {
			gl.history[hidx] = append(gl.history[hidx], byte(c))
		}
		gl.historyCount++
	}
	return nil
}

const histSaveHelp = "Saves in-memory history to 'histpath' (an input window property)"

func (gl *getLine) histSave(r *rpn.RPN) error {
	if gl.autoHistory {
		return errDisableAutoHistoryFirst
	}
	return gl.writeAllHistory()
}

func (gl *getLine) writeAllHistory() error {
	mini := gl.historyCount - MAX_HISTORY_LINES
	if mini < 0 {
		mini = 0
	}
	var buff []byte
	for i := mini; i < gl.historyCount; i++ {
		line := gl.history[i%MAX_HISTORY_LINES]
		buff = append(buff, line...)
		buff = append(buff, '\n')
	}
	err := gl.fs.WriteFile(gl.histpath, buff)
	if err != nil {
		elog.Print("failed to save history file: ", err.Error())
	}
	return err
}

func (gl *getLine) get(r *rpn.RPN) (string, error) {
	gl.txtb.Cursor(true)
	defer gl.txtb.Cursor(false)
	gl.line = gl.line[:0]
	gl.fileList = gl.fileList[:0]
	idx := 0
	// how many steps back into history, with 0 being not in history
	arrowKeyIdx := 0
	for {
		c, err := gl.input.GetChar()
		if err != nil {
			return "", err
		}
		switch c {
		case key.KEY_LEFT:
			if idx > 0 {
				idx--
				gl.txtb.Shift(-1)
			}
		case key.KEY_RIGHT:
			if idx < len(gl.line) {
				idx++
				gl.txtb.Shift(1)
			}
		case key.KEY_UP:
			arrowKeyIdx, idx = gl.replaceLineWithHistory(arrowKeyIdx, 1, idx)
		case key.KEY_DOWN:
			arrowKeyIdx, idx = gl.replaceLineWithHistory(arrowKeyIdx, -1, idx)
		case key.KEY_BACKSPACE:
			if idx > 0 {
				idx--
				gl.line = delete(gl.line, idx)
				gl.txtb.Shift(-1)
				gl.txtb.PrintBytes(gl.line[idx:], true)
				gl.txtb.Write(' ', true)
				gl.txtb.Shift(-(len(gl.line) - idx + 1))
			}
		case key.KEY_DEL:
			if idx < len(gl.line) {
				gl.line = delete(gl.line, idx)
				gl.txtb.PrintBytes(gl.line[idx:], true)
				gl.txtb.Write(' ', true)
				gl.txtb.Shift(-(len(gl.line) - idx + 1))
			}
		case key.KEY_INS:
			gl.insertMode = !gl.insertMode
		case key.KEY_END:
			gl.txtb.Shift(len(gl.line) - idx)
			idx = len(gl.line)
		case key.KEY_HOME:
			gl.txtb.Shift(-idx)
			idx = 0
		case '\t':
			idx = gl.tabComplete(r, idx)
		case 27: // ESCAPE key
			gl.enterScrollingMode(0)
		case key.KEY_PAGEUP:
			gl.enterScrollingMode(-gl.pageDelta())
		case key.KEY_EOF:
			return "exit", nil
		case key.KEY_F1:
			return gl.execMacro(r, idx, "@.f1")
		case key.KEY_F2:
			return gl.execMacro(r, idx, "@.f2")
		case key.KEY_F3:
			return gl.execMacro(r, idx, "@.f3")
		case key.KEY_F4:
			return gl.execMacro(r, idx, "@.f4")
		case key.KEY_F5:
			return gl.execMacro(r, idx, "@.f5")
		case key.KEY_F6:
			return gl.execMacro(r, idx, "@.f6")
		case key.KEY_F7:
			return gl.execMacro(r, idx, "@.f7")
		case key.KEY_F8:
			return gl.execMacro(r, idx, "@.f8")
		case key.KEY_F9:
			return gl.execMacro(r, idx, "@.f9")
		case key.KEY_F10:
			return gl.execMacro(r, idx, "@.f10")
		case key.KEY_F11:
			return gl.execMacro(r, idx, "@.f11")
		case key.KEY_F12:
			return gl.execMacro(r, idx, "@.f12")
		case key.KEY_BREAK:
			gl.txtb.Shift(len(gl.line) - idx)
			gl.txtb.Write('\n', true)
			gl.addToHistoryIfAutoEnabled()
			return "", nil
		default:
			b := byte(c)
			if b == '\n' {
				gl.txtb.Shift(len(gl.line) - idx)
				gl.txtb.Write(b, true)
				gl.addToHistoryIfAutoEnabled()
				return string(gl.line), nil
			}
			gl.addChar(idx, b)
			idx++
		}
	}
}

func (gl *getLine) execMacro(r *rpn.RPN, idx int, name string) (string, error) {
	gl.txtb.Shift(-idx)
	for idx > 0 {
		gl.txtb.Write(' ', true)
		idx--
	}
	return name, nil
}

func (gl *getLine) addChar(idx int, b byte) {
	if idx >= len(gl.line) {
		gl.line = append(gl.line, b)
		gl.txtb.Write(b, true)
	} else if gl.insertMode {
		gl.line = append(gl.line, 0) // grow the buffer
		copy(gl.line[idx+1:], gl.line[idx:])
		gl.line[idx] = b
		gl.txtb.Print(string(gl.line[idx:]), true)
		gl.txtb.Shift(-(len(gl.line) - idx - 1))
	} else {
		gl.line[idx] = b
		gl.txtb.Write(b, true)
	}
}

func (gl *getLine) addToHistoryIfAutoEnabled() {
	// if the last history element is the same as line, don't repeat it
	if gl.historyCount > 0 && bytes.Equal(gl.history[(gl.historyCount-1)%MAX_HISTORY_LINES], gl.line) {
		return
	}
	if len(gl.line) == 0 {
		// line is empty
		return
	}
	if gl.fs == nil {
		return
	}
	hidx := gl.historyCount % MAX_HISTORY_LINES
	if len(gl.history[hidx]) > 0 {
		gl.history[hidx] = gl.history[hidx][:0]
	}
	gl.history[hidx] = append(gl.history[hidx], gl.line...)
	gl.historyCount++
	if gl.autoHistory {
		line := append(gl.line, '\n')
		err := gl.fs.AppendToFile(gl.histpath, line)
		if err != nil {
			elog.Print("could not write history line: ", err.Error()) // object allocated on the heap: escapes at line 259
			gl.fs = nil
		}
	}
}

func (gl *getLine) replaceLineWithHistory(arrowKeyIdx int, delta int, idx int) (int, int) {
	// copy the current line to the current index so we can go back to it
	oldidx := (gl.historyCount - arrowKeyIdx) % MAX_HISTORY_LINES
	if len(gl.history[oldidx]) > 0 {
		gl.history[oldidx] = gl.history[oldidx][:0]
	}
	gl.history[oldidx] = append(gl.history[oldidx], gl.line...)
	arrowKeyIdx += delta
	if (arrowKeyIdx > gl.historyCount) || (arrowKeyIdx >= MAX_HISTORY_LINES) || (arrowKeyIdx < 0) {
		arrowKeyIdx -= delta
		return arrowKeyIdx, idx
	}
	oldlen := len(gl.line)
	newl := gl.history[(gl.historyCount-arrowKeyIdx)%MAX_HISTORY_LINES]
	// remove the existing line
	gl.txtb.Shift(-idx)
	for i := 0; i < oldlen; i++ {
		gl.txtb.Write(' ', true)
	}
	gl.txtb.Shift(-oldlen)
	gl.txtb.PrintBytes(newl, true)
	gl.line = gl.line[:0]
	gl.line = append(gl.line, newl...)
	return arrowKeyIdx, len(gl.line)
}

func (gl *getLine) pageDelta() int {
	return gl.txtb.Txtw.TextHeight() * 3 / 4
}

func (gl *getLine) enterScrollingMode(delta int) {
	gl.txtb.Cursor(false)
	scrollDelta, _ := gl.maybeScroll(0, delta)
	if scrollDelta == 0 {
		gl.drawScrollingBanner(true)
	}
	var exit bool
	for {
		c, err := gl.input.GetChar()
		if err != nil {
			return
		}
		switch c {
		case 'k':
			fallthrough
		case key.KEY_UP:
			scrollDelta, exit = gl.maybeScroll(scrollDelta, -1)
		case 'j':
			fallthrough
		case key.KEY_DOWN:
			scrollDelta, exit = gl.maybeScroll(scrollDelta, 1)
		case key.KEY_PAGEUP:
			scrollDelta, exit = gl.maybeScroll(scrollDelta, -gl.pageDelta())
		case ' ':
			fallthrough
		case key.KEY_PAGEDOWN:
			scrollDelta, exit = gl.maybeScroll(scrollDelta, gl.pageDelta())
		case 27:
			fallthrough
		case '\n':
			fallthrough
		case 'q':
			exit = true
		}
		if exit {
			gl.txtb.Scroll(-scrollDelta)
			gl.txtb.Update()
			gl.txtb.Cursor(true)
			gl.drawScrollingBanner(false)
			return
		}
	}
}

func (gl *getLine) maybeScroll(scrollDelta int, delta int) (int, bool) {
	newDelta := scrollDelta + delta
	maxDelta := gl.txtb.BufferLines() - gl.txtb.Txtw.TextHeight()
	if newDelta < -maxDelta {
		newDelta = -maxDelta
	}
	if newDelta > 0 {
		return 0, true
	}
	if newDelta != scrollDelta {
		gl.txtb.Scroll(newDelta - scrollDelta)
		gl.txtb.Update()
		gl.drawScrollingBanner(true)
	}
	return newDelta, false
}

// Here we draw directly on the text window (white text, blue background)
var scrollCol = window.White | (window.Blue >> 4)

const scrollMsg = "Scrolling Mode"

func (gl *getLine) drawScrollingBanner(enable bool) {
	w, h := gl.txtb.Txtw.TextSize()
	// Both DrawStr and RefreshArea can handle negative values
	x := (w - len(scrollMsg)) / 2
	if enable {
		window.DrawStr(gl.txtb.Txtw, x, h-1, scrollMsg, scrollCol)
	} else {
		gl.txtb.RefreshArea(x, h-1, len(scrollMsg), 1)
	}
}

func delete(line []byte, idx int) []byte {
	return append(line[:idx], line[idx+1:]...)
}
