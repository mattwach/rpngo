package main

import (
	"fmt"
	"image/color"
	"machine"
	"mattwach/rpngo/drivers/tinygo/fonts"
	"mattwach/rpngo/drivers/tinygo/picocalc/ili948x"
	"mattwach/rpngo/common/elog"
	"time"
)

// Since anything beow could fail depending on the state of the hardware,
// do the most reliable things first
func handlePanic(lcd *ili948x.Ili948x, r any) {
	led := machine.LED
	led.Configure(machine.PinConfig{Mode: machine.PinOutput})
	led.High()

	msg := fmt.Sprintf("Panic: %+v", r) // object allocated on the heap (OK)
	elog.Print(msg)
	lcd.FillRectangle(0, 0, 320, 100, ili948x.RED)
	var x int16 = 16
	var y int16 = 20
	c := color.RGBA{R: 255, G: 255, B: 255, A: 255}
	for _, r := range msg {
		fonts.Hack12pt.GetGlyph(rune(r&0xFF)).Draw(lcd, x, y, c)
		x += fonts.FontCharWidth
		if x > 280 {
			x = 16
			y += 16
		}
	}

	for {
		led.Low()
		time.Sleep(time.Millisecond * 500)
		led.High()
		time.Sleep(time.Millisecond * 500)
	}

}
