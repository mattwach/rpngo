package ili948x

import (
	"mattwach/rpngo/common/window"
)

const FontCharWidth = 8

type Ili948xScreen struct {
	// Control the LCD.
	Device *Ili948x
}

func (s *Ili948xScreen) Init() {
	s.Device = InitDisplay()
}

func (s *Ili948xScreen) NewTextWindow() (window.TextWindow, error) {
	tw := &Ili948xTxtW{}
	tw.Init(s.Device)
	return tw, nil
}

func (s *Ili948xScreen) NewPixelWindow() (window.PixelWindow, error) {
	pw := &Ili948xPixW{}
	pw.Init(s.Device)
	return pw, nil
}

func (s *Ili948xScreen) ScreenSize() (int, int) {
	return 320, 320
}
