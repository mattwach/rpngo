// Package ncurses implementation for IO
package curses

import (
	"mattwach/rpngo/common/key"
	"mattwach/rpngo/common/window"

	"github.com/gbin/goncurses"
)

type Curses struct {
	root      *goncurses.Window
	border    *goncurses.Window
	window    *goncurses.Window
	rgbToPair map[uint32]int16 // maps color,color values to a pair.
	col       window.ColorChar
	// redrawing the border causes the contents of the input window
	// to be wiped, thus we only want to do it if needed.
	borderChanged bool
	// set to true if ESC was detected.  This is needed to support
	// KEY_CUT, KEY_COPY, etc
	escPressed bool
}

func Init() (*Curses, error) {
	root, err := goncurses.Init()
	if err != nil {
		return nil, err
	}
	if goncurses.HasColors() {
		goncurses.StartColor()
	}
	goncurses.Echo(false)
	goncurses.Cursor(0)
	tw := &Curses{
		root:      root,
		window:    root,
		rgbToPair: make(map[uint32]int16),
	}
	return tw, nil
}

func (c *Curses) NewTextWindow() (window.TextWindow, error) {
	tw := &Curses{
		root:      c.root,
		rgbToPair: c.rgbToPair,
	}
	return tw, nil
}

func (c *Curses) NewPixelWindow() (window.PixelWindow, error) {
	panic("ncurses does not support pixel-based windows")
}

func (c *Curses) ShowBorder(screenw, screenh int) error {
	if !c.borderChanged {
		return nil
	}
	ch, err := c.colorPairFor(window.Blue)
	if err != nil {
		return err
	}
	c.border.AttrSet(ch)
	// only draw lines if these lines would not be on the edge
	bh, bw := c.border.MaxYX()
	wy, wx := c.window.YX()

	if wy > 0 {
		c.border.HLine(0, 0, '-', bw)
	}

	if wx > 0 {
		c.border.VLine(0, 0, '|', bh)
	}

	if (wx > 0) && (wy > 0) {
		c.border.MoveAddChar(0, 0, '+')
	}

	c.border.Refresh()
	ch, err = c.colorPairFor(window.White)
	if err != nil {
		return err
	}
	c.border.AttrSet(ch)
	c.borderChanged = false
	return nil
}

func (c *Curses) Refresh() {
	if c.window != nil {
		c.window.Refresh()
	}
}

func (c *Curses) End() {
	goncurses.End()
}

func (c *Curses) ResizeWindow(x, y, w, h int) error {
	if c.border != nil {
		by, bx := c.border.YX()
		bh, bw := c.border.MaxYX()
		if (x == bx) && (y == by) && (w == bw) && (h == bh) {
			return nil
		}
	}
	c.borderChanged = true
	if c.border != nil {
		// erase the contents so that artifacts do no collect on the screen
		c.border.Erase()
		c.border.Refresh()
		if err := c.border.Delete(); err != nil {
			return err
		}
	}
	if c.window != nil {
		if err := c.window.Delete(); err != nil {
			return err
		}
	}
	var err error
	screenh, screenw := c.root.MaxYX()
	if (x + w) < screenw {
		w--
	}
	if (y + h) < screenh {
		h--
	}
	c.border, err = goncurses.NewWindow(h, w, y, x)
	if err != nil {
		return err
	}

	// If the window touches the screen edge, there is no reason to
	// have a border there
	if x > 0 {
		x++
		w--
	}
	if y > 0 {
		y++
		h--
	}

	c.window, err = goncurses.NewWindow(h, w, y, x)
	if err != nil {
		return err
	}
	if err := c.window.Keypad(true); err != nil {
		return err
	}
	return nil
}

var charMap = map[goncurses.Key]key.Key{
	goncurses.KEY_LEFT:      key.KEY_LEFT,
	goncurses.KEY_RIGHT:     key.KEY_RIGHT,
	goncurses.KEY_UP:        key.KEY_UP,
	goncurses.KEY_DOWN:      key.KEY_DOWN,
	goncurses.KEY_BACKSPACE: key.KEY_BACKSPACE,
	goncurses.KEY_DC:        key.KEY_DEL,
	goncurses.KEY_IC:        key.KEY_INS,
	goncurses.KEY_END:       key.KEY_END,
	goncurses.KEY_HOME:      key.KEY_HOME,
	4:                       key.KEY_EOF,
	goncurses.KEY_PAGEUP:    key.KEY_PAGEUP,
	goncurses.KEY_PAGEDOWN:  key.KEY_PAGEDOWN,
	goncurses.KEY_F1:        key.KEY_F1,
	goncurses.KEY_F2:        key.KEY_F2,
	goncurses.KEY_F3:        key.KEY_F3,
	goncurses.KEY_F4:        key.KEY_F4,
	goncurses.KEY_F5:        key.KEY_F5,
	goncurses.KEY_F6:        key.KEY_F6,
	goncurses.KEY_F7:        key.KEY_F7,
	goncurses.KEY_F8:        key.KEY_F8,
	goncurses.KEY_F9:        key.KEY_F9,
	goncurses.KEY_F10:       key.KEY_F10,
	goncurses.KEY_F11:       key.KEY_F11,
	goncurses.KEY_F12:       key.KEY_F11,
	goncurses.KEY_SF:        key.KEY_SDOWN,
	goncurses.KEY_SR:        key.KEY_SUP,
	goncurses.KEY_SRIGHT:    key.KEY_SRIGHT,
	goncurses.KEY_SLEFT:     key.KEY_SLEFT,
	goncurses.KEY_SHOME:     key.KEY_SHOME,
	goncurses.KEY_SEND:      key.KEY_SEND,
}

func (c *Curses) GetChar() (key.Key, error) {
	if c.window == nil {
		return 0, nil
	}
	ch := c.window.GetChar()
	if c.escPressed {
		return c.getCharWithEscPressed(byte(ch))
	}
	if ch == 27 {
		c.escPressed = true
		return 0, nil
	}
	k, ok := charMap[ch]
	if ok {
		return k, nil
	}
	return key.Key(ch), nil
}

func (c *Curses) getCharWithEscPressed(ch byte) (key.Key, error) {
	c.escPressed = false
	switch ch {
	case 'c':
		return key.KEY_COPY, nil
	case 'h':
		return key.KEY_HELP, nil
	case 'q':
		return key.KEY_QUIT, nil
	case 'x':
		return key.KEY_CUT, nil
	case 'v':
		return key.KEY_PASTE, nil
	case 's', 'w':
		return key.KEY_SAVE, nil
	}

	return 0, nil
}

func (c *Curses) Erase() {
	if c.window != nil {
		c.window.Erase()
	}
}

func (c *Curses) TextWidth() int {
	if c.window == nil {
		return 0
	}
	_, x := c.window.MaxYX()
	return x
}

func (c *Curses) TextHeight() int {
	if c.window == nil {
		return 0
	}
	y, _ := c.window.MaxYX()
	return y
}

func (c *Curses) TextSize() (int, int) {
	if c.window == nil {
		return 0, 0
	}
	y, x := c.window.MaxYX()
	return x, y
}

func (c *Curses) WindowSize() (int, int) {
	if c.window == nil {
		return 0, 0
	}
	y, x := c.window.MaxYX()
	return x, y
}

func (c *Curses) ScreenSize() (int, int) {
	if c.window == nil {
		return 0, 0
	}
	y, x := c.window.MaxYX()
	return x, y
}

func (c *Curses) WindowXY() (int, int) {
	if c.window == nil {
		return 0, 0
	}
	y, x := c.window.YX()
	return x, y
}

func (c *Curses) DrawChar(x, y int, ch window.ColorChar) {
	if c.window == nil {
		return
	}
	newcol := ch & 0xFF00
	if newcol != c.col {
		c.textColor(newcol)
	}
	c.window.Move(y, x)
	b := byte(ch & 0xFF)
	c.window.AddChar(goncurses.Char(b))
}

func (c *Curses) textColor(col window.ColorChar) {
	if c.window == nil {
		return
	}
	ch, err := c.colorPairFor(col)
	if err != nil {
		return
	}
	c.col = col
	c.window.AttrSet(ch)
}

func (c *Curses) colorPairFor(col window.ColorChar) (goncurses.Char, error) {
	if !goncurses.HasColors() {
		return 0, nil
	}
	fc := idxToCol[uint8(col>>12)]
	bc := idxToCol[uint8((col&0x0F00)>>8)]
	pc := (uint32(fc) << 15) | uint32(bc)
	pidx, ok := c.rgbToPair[pc]
	if !ok {
		pidx = int16(len(c.rgbToPair) + 1) // zero is the default so start at 1
		if err := goncurses.InitPair(pidx, fc, bc); err != nil {
			return 0, err
		}
		c.rgbToPair[pc] = pidx
	}
	return goncurses.ColorPair(pidx), nil
}

var idxToCol = map[uint8]int16{
	0b0000: goncurses.C_BLACK,
	0b0001: goncurses.C_BLUE,
	0b0010: goncurses.C_GREEN,
	0b0011: goncurses.C_BLUE,
	0b0100: goncurses.C_GREEN,
	0b0101: goncurses.C_CYAN,
	0b0110: goncurses.C_GREEN,
	0b0111: goncurses.C_CYAN,
	0b1000: goncurses.C_RED,
	0b1001: goncurses.C_MAGENTA,
	0b1010: goncurses.C_RED,
	0b1011: goncurses.C_MAGENTA,
	0b1100: goncurses.C_YELLOW,
	0b1101: goncurses.C_WHITE,
	0b1110: goncurses.C_YELLOW,
	0b1111: goncurses.C_WHITE,
}
