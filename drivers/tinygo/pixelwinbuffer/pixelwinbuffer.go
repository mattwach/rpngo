// Package pixelwinbuffer provides an alternative to
// "erase all pixels, redraw".  Which either creates on-screen
// flashing or takes up memory (significant amount for the rp2040).
// It works like this:
//
// All drawing calls fill in a 1bpp buffer that the pixel is updated
// when it's time to refresh, a previous version of this
// buffer is compared to the latest one.
// Every pixel that is set for the old one and not set
// for the new one is cleared.
//
// This has the effect of a slight double image as the old
// and new will be onscreen at the same time for a very small
// amount of time until the old pixels can be cleared.
package pixelwinbuffer

import (
	"errors"
	"image/color"
	"mattwach/rpngo/drivers/tinygo/fonts"
	"mattwach/rpngo/common/elog"
	"mattwach/rpngo/common/window"
)

var bgcol = color.RGBA{}

// pixelBufferDisplayer is an adapter that will allow tinyfont to route pixels
// through the buffer
type pixelBufferDisplayer struct {
	buffer window.PixelWindow
}

func (pb *pixelBufferDisplayer) Display() error {
	return nil
}

func (pb *pixelBufferDisplayer) SetPixel(x, y int16, c color.RGBA) {
	pb.buffer.Color(c)
	pb.buffer.SetPoint(int(x), int(y))
}

func (pb *pixelBufferDisplayer) Size() (int16, int16) {
	w, h := pb.buffer.PixelSize()
	return int16(w), int16(h)
}

type PixelBuffer struct {
	previous  []uint8
	current   []uint8
	target    window.PixelWindow
	pw        int
	ph        int
	col       color.RGBA
	skipbit   bool
	displayer pixelBufferDisplayer
}

func (pb *PixelBuffer) Init(target window.PixelWindow) error {
	pb.target = target
	pb.displayer.buffer = pb
	return pb.ResizeWindow(0, 0, 5, 5)
}

func (pb *PixelBuffer) Refresh() {
	var mask uint8 = 0x80
	boffset := 0
	pb.target.Color(bgcol)
	for y := 0; y < pb.ph; y++ {
		for x := 0; x < pb.pw; x++ {
			if ((pb.previous[boffset] & mask) != 0) && ((pb.current[boffset] & mask) == 0) {
				pb.target.SetPoint(x, y)
			}
			if mask == 0x01 {
				boffset++
				mask = 0x80
			} else {
				mask >>= 1
			}
		}
	}
	// get ready for the next round
	for i := range pb.current {
		pb.previous[i] = pb.current[i]
		pb.current[i] = 0
	}
}

func (pb *PixelBuffer) ResizeWindow(x, y, w, h int) error {
	pb.target.ResizeWindow(x, y, w, h)
	size := ((w + 7) * h) / 8
	if size <= 0 {
		return errors.New("pixelbufffer, size <= 0")
	}
	if len(pb.previous) < size {
		elog.Heap("alloc: /drivers/pixelwinbuffer/pixelwinbuffer.go:94: pb.previous = make([]uint8, size)")
		pb.previous = make([]uint8, size) // object allocated on the heap: size is not constant
	}
	if len(pb.current) < size {
		elog.Heap("alloc: /drivers/pixelwinbuffer/pixelwinbuffer.go:97: pb.current = make([]uint8, size)")
		pb.current = make([]uint8, size) // object allocated on the heap: size is not constant
	}
	pb.pw, pb.ph = pb.target.PixelSize() // for less overhead
	return nil
}

func (pb *PixelBuffer) ShowBorder(sw, sh int) error {
	pb.target.ShowBorder(sw, sh)
	return nil
}

func (pb *PixelBuffer) WindowXY() (int, int) {
	return pb.target.WindowXY()
}

func (pb *PixelBuffer) WindowSize() (int, int) {
	return pb.target.WindowSize()
}

func (pb *PixelBuffer) PixelSize() (int, int) {
	return pb.target.PixelSize()
}

func (pb *PixelBuffer) Color(c color.RGBA) {
	pb.skipbit = c == bgcol
	pb.col = c
	pb.target.Color(c)
}

func (pb *PixelBuffer) setBit(x, y int) {
	poffset := y*pb.pw + x
	pb.current[poffset>>3] |= (0x80 >> (poffset & 7))
}

func (pb *PixelBuffer) setBitHline(x, y, w int) {
	offset := y*pb.pw + x
	maxoffset := offset + w - 1
	for offset <= maxoffset {
		if ((offset & 0x07) == 0) && ((maxoffset - offset) >= 8) {
			pb.current[offset>>3] = 0xFF
			offset += 8
		} else {
			pb.current[offset>>3] |= (0x80 >> (offset & 7))
			offset++
		}
	}
}

func (pb *PixelBuffer) SetPoint(x, y int) {
	pb.target.SetPoint(x, y)
	if pb.skipbit {
		return
	}
	pb.setBit(x, y)
}

func (pb *PixelBuffer) HLine(x, y, w int) {
	pb.target.HLine(x, y, w)
	if pb.skipbit {
		return
	}
	pb.setBitHline(x, y, w)
}

func (pb *PixelBuffer) VLine(x, y, h int) {
	pb.target.VLine(x, y, h)
	if pb.skipbit {
		return
	}
	ymax := y + h
	for y < ymax {
		pb.setBit(x, y)
		y++
	}
}

func (pb *PixelBuffer) FilledRect(x, y, w, h int) {
	pb.target.FilledRect(x, y, w, h)
	if pb.skipbit {
		return
	}
	ymax := y + h - 1
	for y < ymax {
		pb.setBitHline(x, y, w)
		y++
	}
}

func (pb *PixelBuffer) Text(s string, x, y int) {
	// do it lower level to avoid importing a bunch of tinyfont code
	for _, r := range s {
		fonts.Hack12pt.GetGlyph(rune(r&0xFF)).Draw(
			&pb.displayer, int16(x), int16(y), pb.col)
		x += fonts.FontCharWidth
	}
}
