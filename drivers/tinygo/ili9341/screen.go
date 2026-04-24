package ili9341

import (
	"image/color"
	"mattwach/rpngo/common/elog"
	"mattwach/rpngo/common/window"

	"tinygo.org/x/drivers/ili9341"
)

type Ili9341Screen struct {
	// Control the LCD.
	Device *ili9341.Device
}

func (s *Ili9341Screen) Init() {
	s.Device = InitDisplay()
}

func (s *Ili9341Screen) NewTextWindow() (window.TextWindow, error) {
	elog.Heap("alloc: /drivers/tinygo/ili9341/screen.go:20: tw := &Ili9341TxtW{}")
	tw := &Ili9341TxtW{} // object allocated on the heap: escapes at line 20
	tw.Init(s.Device)
	return tw, nil
}

func (s *Ili9341Screen) NewPixelWindow() (window.PixelWindow, error) {
	elog.Heap("alloc: /drivers/tinygo/ili9341/screen.go:26: pw := &Ili9341PixW{}")
	pw := &Ili9341PixW{} // object allocated on the heap: escapes at line 26
	pw.Init(s.Device)
	return pw, nil
}

func (s *Ili9341Screen) ScreenSize() (int, int) {
	return 320, 240
}

// fastvline and fasthline are bugged for some reason, we we do it the slow way for now
func slowHline(d *ili9341.Device, x0 int16, x1 int16, y int16, c color.RGBA) error {
	for x0 <= x1 {
		d.SetPixel(x0, y, c)
		x0++
	}
	return nil
}

func slowVline(d *ili9341.Device, x int16, y0 int16, y1 int16, c color.RGBA) error {
	for y0 <= y1 {
		d.SetPixel(x, y0, c)
		y0++
	}
	return nil
}
