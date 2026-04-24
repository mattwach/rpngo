// Ili9341PixW implements the PixelWindow interface in
// window/window.go to the ili9341 LCD
package ili9341

import (
	"errors"
	"image/color"
	"mattwach/rpngo/drivers/tinygo/fonts"
	"mattwach/rpngo/common/window"

	"tinygo.org/x/drivers/ili9341"
)

type Ili9341PixW struct {
	// screen to send chars to
	device *ili9341.Device

	// window dimensions in pixels
	wx int16
	wy int16
	ww int16
	wh int16

	// pixel dimension, in pixels
	px int16
	py int16
	pw int16
	ph int16

	// current color
	col color.RGBA
}

// Init initializes a pixel window. x, y, w, and h are all in pixels
func (tw *Ili9341PixW) Init(d *ili9341.Device) {
	tw.device = d
	tw.ResizeWindow(0, 0, 5, 5)
}

func (tw *Ili9341PixW) Refresh() {}

func (tw *Ili9341PixW) ResizeWindow(x, y, w, h int) error {
	if (tw.wx == int16(x)) && (tw.wy == int16(y)) && (tw.ww == int16(w)) && (tw.wh == int16(h)) {
		return nil
	}
	if (w < 3) || (h < 3) {
		return errors.New("pixelwindow resize too small")
	}
	tw.wx = int16(x)
	tw.wy = int16(y)
	tw.ww = int16(w)
	tw.wh = int16(h)
	tw.px = tw.wx + 1
	tw.py = tw.wy + 1
	tw.pw = tw.ww - 2
	tw.ph = tw.wh - 2
	tw.device.FillRectangle(int16(x), int16(y), int16(w), int16(h), color.RGBA{})
	return nil
}

func (tw *Ili9341PixW) WindowXY() (int, int) {
	return int(tw.wx), int(tw.wy)
}

func (tw *Ili9341PixW) WindowSize() (int, int) {
	return int(tw.ww), int(tw.wh)
}

func (tw *Ili9341PixW) PixelSize() (int, int) {
	return int(tw.pw), int(tw.ph)
}

func (tw *Ili9341PixW) ShowBorder(screenw, screenh int) error {
	// Need to expand by one pixel as the main window does not include the border.
	x0 := tw.wx
	x1 := tw.wx + tw.ww - 1
	y0 := tw.wy
	y1 := tw.wy + tw.wh - 1
	/*
		tw.device.DrawFastHLine(x0, x1, y0, window.BorderColor)
		tw.device.DrawFastHLine(x0, x1, y1, window.BorderColor)
		tw.device.DrawFastVLine(x0, y0+1, y1-1, window.BorderColor)
		tw.device.DrawFastVLine(x1, y0+1, y1-1, window.BorderColor)
	*/
	slowHline(tw.device, x0, x1, y0, window.BorderColor)
	slowHline(tw.device, x0, x1, y1, window.BorderColor)
	slowVline(tw.device, x0, y0+1, y1-1, window.BorderColor)
	slowVline(tw.device, x1, y0+1, y1-1, window.BorderColor)
	return nil
}

func (tw *Ili9341PixW) Color(c color.RGBA) {
	tw.col = c
}

func (tw *Ili9341PixW) SetPoint(x, y int) {
	tw.device.SetPixel(tw.px+int16(x), tw.py+int16(y), tw.col)
}

func (tw *Ili9341PixW) HLine(x, y, w int) {
	//tw.device.DrawFastHLine(tw.px+int16(x), tw.px+int16(x+w)-1, tw.py+int16(y), tw.col)
	slowHline(tw.device, tw.px+int16(x), tw.px+int16(x+w)-1, tw.py+int16(y), tw.col)
}

func (tw *Ili9341PixW) VLine(x, y, h int) {
	//tw.device.DrawFastVLine(tw.px+int16(x), tw.py+int16(y), tw.py+int16(y+h)-1, tw.col)
	slowVline(tw.device, tw.px+int16(x), tw.py+int16(y), tw.py+int16(y+h)-1, tw.col)
}

func (tw *Ili9341PixW) FilledRect(x, y, w, h int) {
	tw.device.FillRectangle(tw.px+int16(x), tw.py+int16(y), int16(w), int16(h), tw.col)
}

func (tw *Ili9341PixW) Text(s string, x, y int) {
	// do it lower level to avoid importing a bunch of tinyfont code
	for _, r := range s {
		fonts.Hack12pt.GetGlyph(rune(r&0xFF)).Draw(
			tw.device, tw.px+int16(x), tw.py+int16(y), tw.col)
		x += fonts.FontCharWidth
	}
}
