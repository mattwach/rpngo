// Implements textwindow for the tinygo env.
//
// Currently this targets the ili948x.  If/when more
// devices are supported, some refactoring may need to occcur.
package ili948x

import (
	"image/color"
	"mattwach/rpngo/drivers/tinygo/fonts"
	"mattwach/rpngo/common/window"
)

const textPad = 3

type Ili948xTxtW struct {
	// screen to send chars to
	device *Ili948x

	// window dimensions in pixels
	wx int16
	wy int16
	ww int16
	wh int16

	// character dimension in pixels
	cw int16
	ch int16

	// character y offset
	cyoffset int16

	// width and height as text cells
	textw int16
	texth int16

	// A cell to draw the characters in
	bitmap Bitmap
	// saves a little performance in drawing
	lastr window.ColorChar
}

// Init initializes a text window. x, y, w, and h are all in pixels
func (tw *Ili948xTxtW) Init(d *Ili948x) {
	tw.cw = FontCharWidth
	tw.ch = 14
	tw.cyoffset = 11
	tw.bitmap.Init(tw.cw, tw.ch)
	tw.device = d
	tw.ResizeWindow(0, 0, 1, 1)
}

func (tw *Ili948xTxtW) ResizeWindow(x, y, w, h int) error {
	if (tw.wx == int16(x)) && (tw.wy == int16(y)) && (tw.ww == int16(w)) && (tw.wh == int16(h)) {
		return nil
	}
	tw.wx = int16(x)
	tw.wy = int16(y)

	if (tw.ww != int16(w)) || (tw.wh != int16(h)) {
		tw.textw = (int16(w) - textPad*2) / tw.cw
		if tw.textw <= 0 {
			tw.textw = 1
		}
		tw.texth = (int16(h) - textPad*2) / tw.ch
		if tw.texth <= 0 {
			tw.texth = 1
		}
		tw.ww = int16(w)
		tw.wh = int16(h)
	}
	tw.Erase()
	return nil
}

func (tw *Ili948xTxtW) Refresh() {
	// maybe no need to do this?
}

func fgColor(c window.ColorChar) color.RGBA {
	r, g, b := c.FGColor8()
	return color.RGBA{R: r, G: g, B: b}
}

func bgColor(c window.ColorChar) RGB565 {
	r, g, b := c.BGColor8()
	return NewRGB565(r>>3, g>>2, b>>3)
}

func (tw *Ili948xTxtW) DrawChar(tx, ty int, r window.ColorChar) {
	if r != tw.lastr {
		tw.lastr = r
		tw.bitmap.FillWith(bgColor(r))
		fonts.Hack12pt.GetGlyph(rune(r&0xFF)).Draw(&tw.bitmap, 0, tw.cyoffset, fgColor(r))
	}
	tw.device.DrawBitmap(tw.wx+int16(tx)*tw.cw+textPad, tw.wy+int16(ty)*tw.ch+textPad, &tw.bitmap)
}

func (tw *Ili948xTxtW) Erase() {
	tw.device.FillRectangle(tw.wx, tw.wy, tw.ww, tw.wh, 0) // fill with black
}

func (tw *Ili948xTxtW) ShowBorder(screenw, screenh int) error {
	tw.device.DrawHLine(tw.wx, tw.wx+tw.ww-1, tw.wy, BLUE)
	tw.device.DrawHLine(tw.wx, tw.wx+tw.ww-1, tw.wy+tw.wh-1, BLUE)
	tw.device.DrawVLine(tw.wx, tw.wy, tw.wy+tw.wh-1, BLUE)
	tw.device.DrawVLine(tw.wx+tw.ww-1, tw.wy, tw.wy+tw.wh-1, BLUE)
	return nil
}

func (tw *Ili948xTxtW) TextWidth() int {
	return int(tw.textw)
}

func (tw *Ili948xTxtW) TextHeight() int {
	return int(tw.texth)
}

func (tw *Ili948xTxtW) TextSize() (int, int) {
	return int(tw.textw), int(tw.texth)
}

func (tw *Ili948xTxtW) WindowXY() (int, int) {
	return int(tw.wx), int(tw.wy)
}

func (tw *Ili948xTxtW) WindowSize() (int, int) {
	return int(tw.ww), int(tw.wh)
}
