// Implements textwindow for the tinygo env.
//
// Currently this targets the ili9341.  If/when more
// devices are supported, some refactoring may need to occcur.
package ili9341

import (
	"image/color"
	"mattwach/rpngo/drivers/tinygo/fonts"
	"mattwach/rpngo/drivers/tinygo/pixel565"
	"mattwach/rpngo/common/window"

	"tinygo.org/x/drivers/ili9341"
	"tinygo.org/x/drivers/pixel"
)

const textPad = 3

type Ili9341TxtW struct {
	// screen to send chars to
	device *ili9341.Device

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
	image pixel565.Pixel565
	// saves a little performance in drawing
	lastr window.ColorChar
}

// Init initializes a text window. x, y, w, and h are all in pixels
func (tw *Ili9341TxtW) Init(d *ili9341.Device) {
	tw.cw = fonts.FontCharWidth
	tw.ch = 14
	tw.cyoffset = 11
	tw.image.Init(tw.cw, tw.ch)
	tw.device = d
	tw.ResizeWindow(0, 0, 1, 1)
}

func (tw *Ili9341TxtW) ResizeWindow(x, y, w, h int) error {
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

func (tw *Ili9341TxtW) Refresh() {
	// maybe no need to do this?
}

func fgColor(c window.ColorChar) color.RGBA {
	r, g, b := c.FGColor8()
	return color.RGBA{R: r, G: g, B: b}
}

func bgColor(c window.ColorChar) pixel.RGB565BE {
	r, g, b := c.BGColor8()
	return pixel.NewRGB565BE(r, g, b)
}

func (tw *Ili9341TxtW) DrawChar(tx, ty int, r window.ColorChar) {
	c := r & 0xFF
	if (c < 32) || (c > 127) {
		// a possibly-unprintable character that might cause a crash in Draw().
		r = window.Magenta | '.'
	}
	if r != tw.lastr {
		tw.lastr = r
		tw.image.Image.FillSolidColor(bgColor(r))
		fonts.Hack12pt.GetGlyph(rune(r&0xFF)).Draw(&tw.image, 0, tw.cyoffset, fgColor(r))
	}
	tw.device.DrawBitmap(tw.wx+int16(tx)*tw.cw+textPad, tw.wy+int16(ty)*tw.ch+textPad, tw.image.Image)
}

func (tw *Ili9341TxtW) Erase() {
	tw.device.FillRectangle(tw.wx, tw.wy, tw.ww, tw.wh, color.RGBA{}) // fill with black
}

func (tw *Ili9341TxtW) ShowBorder(screenw, screenh int) error {
	c := color.RGBA{R: 0, G: 0, B: 100}
	slowHline(tw.device, tw.wx, tw.wx+tw.ww-1, tw.wy, c)
	slowHline(tw.device, tw.wx, tw.wx+tw.ww-1, tw.wy+tw.wh-1, c)
	slowVline(tw.device, tw.wx, tw.wy, tw.wy+tw.wh-1, c)
	slowVline(tw.device, tw.wx+tw.ww-1, tw.wy, tw.wy+tw.wh-1, c)
	return nil
}

func (tw *Ili9341TxtW) TextWidth() int {
	return int(tw.textw)
}

func (tw *Ili9341TxtW) TextHeight() int {
	return int(tw.texth)
}

func (tw *Ili9341TxtW) TextSize() (int, int) {
	return int(tw.textw), int(tw.texth)
}

func (tw *Ili9341TxtW) WindowXY() (int, int) {
	return int(tw.wx), int(tw.wy)
}

func (tw *Ili9341TxtW) WindowSize() (int, int) {
	return int(tw.ww), int(tw.wh)
}
