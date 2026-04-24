package window

import (
	"mattwach/rpngo/common/elog"
)

// TextBuffer serves multiple purposes and allows character drivers (of which
// we have several LCD and curses variants) have a reduced interface and
// implementation requirements.  This also increases consistency (by does
// involve reinventing already-solved issues for some targets - particurlarily
// the ncurses one).
//
// Base requirements are:
//  1. Support color text.  Maybe backgrounds too, although that had not yet
//     seen any use
//  2. Do not redraw characters with the same character.  Drawing characters
//     on a color LCD via SPI is a somehwat expensive process.  A low hanginng
//     fruit optimization is to not do it when possible.
//  3. Support scrolling.  It is reasonable for the user to need to look back at
//     history, especially when a command (such as help) has a multiscreen
//     result.
//  5. Keepthe door open for editing support.  Although that will come later.
//
// Implementation.
//
// We have a char, which is a 16-bit value that combines:
//   - An 8-bit ASCII byte (expansion to unicode might be needed later, if this
//     project takes on any international interet).
//   - A foreground and background color (4 bits each and why we might drop
//     or reduce the bg bits later.
//
// We have a chars slice that contains the entire buffer as a single block.
// The format is subsequent bytes tracing to the right and wrapping to the next
// line, which is a typical implementation.
//
// We have headidx which is a pointer to 0,0 of the visible area of the LCD
// Characters before headidx are what you would see if you scrollup and
// characters after headix+screensize are what you will see if you scroll down.
type TextBuffer struct {
	// holds the characters that make up the entire text area
	buffer []ColorChar
	// holds the characters the represent the visible screen
	screen []ColorChar

	// to make scrolling cheaper, we have a head index which starts
	// at 0 and moves forward when we scroll up
	headidx int
	// Scrollback bytes is used to determine how large to make chars and
	// h.  Set it to zero for no scrolling support
	scrollbytes int

	// character position, on the screen and not the buffer
	cx         int
	cy         int
	showCursor bool

	// height of chars buffer.  Note that bw should equal the TextWindow
	// width but ch could be larger.
	bw int
	bh int

	// text color
	col ColorChar

	// text area to update
	Txtw TextWindow
}

func (tb *TextBuffer) Init(txtw TextWindow, scrollbytes int) {
	tb.Txtw = txtw
	tb.scrollbytes = scrollbytes
	// This check is so that some unit tests (like stackwin) are not required
	// to provide a TextWindow.  Arguably there should be a shared fake one
	// that can be used but no tests in stackwin would use it explicitly yet.
	if txtw != nil {
		tb.CheckSize()
	}
}

func (tb *TextBuffer) Cursor(c bool) {
	if tb.showCursor != c {
		tb.showCursor = c
		if tb.showCursor {
			tb.drawCursor()
		} else {
			tb.eraseCursor()
		}
	}
}

func (tb *TextBuffer) Update() {
	var bi int = tb.headidx
	var si int = 0
	w, h := tb.Txtw.TextSize()
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			if tb.screen[si] != tb.buffer[bi] {
				tb.screen[si] = tb.buffer[bi]
				tb.Txtw.DrawChar(int(x), int(y), tb.screen[si])
			}
			bi = (bi + 1) % len(tb.buffer)
			si++
		}
	}
	tb.Txtw.Refresh() // ncurses needs this, LCDs do not
}

// Refreshes an area of the screen, even if buffer matches screen
// This is useful for recovering from otehr things drawing on the text window.
func (tb *TextBuffer) RefreshArea(tx, ty, w, h int) {
	tw, th := tb.Txtw.TextSize()
	for y := ty; y < ty+h; y++ {
		if (y < 0) || (y >= th) {
			continue
		}
		for x := tx; x < tw+w; x++ {
			if (x < 0) || (x >= tw) {
				continue
			}
			si := y*tb.bw + x
			bi := (tb.headidx + si) % len(tb.buffer)
			tb.screen[si] = tb.buffer[bi]
			tb.Txtw.DrawChar(x, y, tb.screen[si])
		}
	}
	tb.Txtw.Refresh() // ncurses needs this, LCDs do not
}

// Checks if the underlying window has resized
func (tb *TextBuffer) CheckSize() {
	tw, th := tb.Txtw.TextSize()
	if (tw == 0) || (th == 0) {
		return
	}
	scrollh := tb.scrollbytes / tw
	if (int(tb.bw) == tw) && (int(tb.bh) == (scrollh + th)) {
		// already the right size
		return
	}
	tb.SetCursorXY(0, 0)
	tb.bw = tw
	tb.bh = th + scrollh
	elog.Heap("alloc: /window/txtbuffer.go:141: tb.buffer = make([]ColorChar, tb.bw*tb.bh)")
	tb.buffer = make([]ColorChar, tb.bw*tb.bh) // object allocated on the heap: size is not constant
	elog.Heap("alloc: /window/txtbuffer.go:142: tb.screen = make([]ColorChar, tw*th)")
	tb.screen = make([]ColorChar, tw*th) // object allocated on the heap: size is not constant
	// A reflow would be nicer than an erase.  I attempted to do it in the
	// branch reflow_text_buffer but ran into complications with handling new
	// line characters properly.  I think it could be done if the definition
	// of buffer was not 1:1 to the screen but instead was just a block
	// of characters.  But this would require logic changes everywhere and
	// we would probably want to drop DrawChar() and rework the editor.
	tb.Erase()
}

func (tb *TextBuffer) Erase() {
	b := tb.col | ColorChar(' ')
	for i := range tb.buffer {
		tb.buffer[i] = b
	}
	tb.headidx = 0
	tb.Update()
}

// stagees a character write to x, y with no speial handling
// of characters or scrolling and no cursor position change.
func (tb *TextBuffer) DrawChar(x, y int, c ColorChar) {
	bidx := (tb.headidx + y*tb.bw + x) % len(tb.buffer)
	tb.buffer[bidx] = c
}

func (tb *TextBuffer) Write(b byte, updatenow bool) error {
	if len(tb.buffer) == 0 {
		// No underlying window
		return nil
	}
	bidx := (tb.headidx + int(tb.cy)*tb.bw + int(tb.cx)) % len(tb.buffer)
	if b != '\n' {
		tb.buffer[bidx] = tb.col | ColorChar(b)
		if updatenow {
			// just update this character
			sidx := tb.cy*tb.bw + tb.cx
			if tb.showCursor || (tb.screen[sidx] != tb.buffer[bidx]) {
				tb.screen[sidx] = tb.buffer[bidx]
				tb.Txtw.DrawChar(int(tb.cx), int(tb.cy), tb.screen[sidx])
			}
		}
		tb.cx++
	}
	tw, th := tb.Txtw.TextSize()
	if (b == '\n') || (int(tb.cx) >= tw) {
		// next line
		if tb.showCursor && (int(tb.cx) < tw) {
			tb.eraseCursor()
		}
		tb.cx = 0
		tb.cy++
	}
	if int(tb.cy) >= th {
		// erase rest of the line, which might contain old data
		tb.Scroll(1)

		tb.cy--
		lineidx := (tb.headidx + int(tb.cy*tb.bw)) % len(tb.buffer)
		for i := 0; i < tw; i++ {
			tb.buffer[(lineidx+i)%len(tb.buffer)] = tb.col | ColorChar(' ')
		}
		if updatenow {
			tb.Update()
		}
	}
	if tb.showCursor {
		tb.drawCursor()
	}
	return nil
}

// takes the character under the cursor and redraws it with black text
// and a white background
func (tb *TextBuffer) drawCursor() {
	c := tb.screen[tb.cy*tb.bw+tb.cx] & 0xFF
	tb.Txtw.DrawChar(int(tb.cx), int(tb.cy), CursorColor|c)
}

// restore the original properties of the character under the cursor
func (tb *TextBuffer) eraseCursor() {
	c := tb.screen[tb.cy*tb.bw+tb.cx]
	tb.Txtw.DrawChar(int(tb.cx), int(tb.cy), c)
}

func (tb *TextBuffer) CursorX() int {
	return int(tb.cx)
}

func (tb *TextBuffer) CursorY() int {
	return int(tb.cy)
}

func (tb *TextBuffer) CursorXY() (int, int) {
	return int(tb.cx), int(tb.cy)
}

func (tb *TextBuffer) SetCursorX(x int) {
	tb.SetCursorXY(x, int(tb.cy))
}

func (tb *TextBuffer) SetCursorY(y int) {
	tb.SetCursorXY(int(tb.cx), y)
}

func (tb *TextBuffer) SetCursorXY(x, y int) {
	if (int(tb.cx) == x) && (int(tb.cy) == y) {
		return
	}
	if tb.showCursor {
		tb.eraseCursor()
	}
	tb.cx = x
	tb.cy = y
	if tb.showCursor {
		tb.drawCursor()
	}
}

func (tb *TextBuffer) TextColor(col ColorChar) {
	tb.col = col
}

func (tb *TextBuffer) BufferLines() int {
	return len(tb.buffer) / tb.bw
}

func (tb *TextBuffer) Scroll(i int) {
	// scrolling is as easy as moving the headidx
	tb.headidx += i * int(tb.bw)
	for tb.headidx >= len(tb.buffer) {
		tb.headidx -= len(tb.buffer)
	}
	for tb.headidx < 0 {
		tb.headidx += len(tb.buffer)
	}
}

func (tb *TextBuffer) Print(msg string, updatenow bool) {
	for _, b := range msg {
		if err := tb.Write(byte(b), updatenow); err != nil {
			return
		}
	}
}

func (tb *TextBuffer) PrintBytes(bytes []byte, updatenow bool) {
	for _, b := range bytes {
		if err := tb.Write(b, updatenow); err != nil {
			return
		}
	}
}

func (tb *TextBuffer) PrintErr(err error, updatenow bool) {
	tb.TextColor(Red)
	tb.Print(err.Error(), updatenow)
	tb.Write('\n', updatenow)
	tb.TextColor(White)
}

func (tb *TextBuffer) Shift(n int) {
	x, y := tb.CursorXY()
	x += n
	for x >= tb.Txtw.TextWidth() {
		y += 1
		if y >= tb.Txtw.TextHeight() {
			tb.Scroll(1)
			y--
		}
		x -= tb.Txtw.TextWidth()
	}
	for x < 0 {
		x += tb.Txtw.TextWidth()
		y -= 1
		if y < 0 {
			tb.Scroll(1)
			y = 0
		}
	}
	tb.SetCursorXY(x, y)
}
