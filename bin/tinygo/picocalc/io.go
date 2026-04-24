package main

import (
	"mattwach/rpngo/drivers/tinygo/picocalc/i2ckbd"
	"mattwach/rpngo/drivers/tinygo/picocalc/ili948x"
	"mattwach/rpngo/drivers/tinygo/serial"
	"mattwach/rpngo/common/elog"
	"mattwach/rpngo/common/key"
	"mattwach/rpngo/common/rpn"
	"mattwach/rpngo/common/xmodem"
	"time"
)

var xmodemCommands xmodem.XmodemCommands

type picoCalcIO struct {
	serial        serial.Serial
	keyboard      i2ckbd.I2CKbd
	screen        ili948x.Ili948xScreen
	originalPrint func(string)
	rpnInst       *rpn.RPN
}

func (gi *picoCalcIO) Init(r *rpn.RPN) {
	gi.screen.Init()
	if err := gi.keyboard.Init(); err != nil {
		elog.Print("failed to init keyboard: ", err.Error())
	}
	gi.serial.Init(true)
	gi.rpnInst = r
	xmodemCommands.InitAndRegister(r, &gi.serial)
}

func (gi *picoCalcIO) RegisterPrint() {
	gi.originalPrint = gi.rpnInst.Print
	rpnInst.Print = gi.Print
}

func (gi *picoCalcIO) GetChar() (key.Key, error) {
	for {
		time.Sleep(20 * time.Millisecond)
		k := gi.serial.GetChar()
		if k != 0 {
			return k, nil
		}
		k, _ = gi.keyboard.GetChar()
		if k != 0 {
			return k, nil
		}
	}
}

func (gi *picoCalcIO) Print(str string) {
	gi.originalPrint(str)
	f, err := gi.rpnInst.GetVariable(".echo")
	if err != nil {
		return
	}
	enabled, err := f.Bool()
	if err != nil {
		return
	}
	if !enabled {
		return
	}
	for _, c := range str {
		_ = gi.serial.WriteByte(byte(c))
	}
}

func (gi *picoCalcIO) ctrlDown() bool {
	for i := 0; i < 5; i++ {
		gi.keyboard.GetChar()
		if gi.keyboard.CtrlDown {
			return true
		}
		time.Sleep(10 * time.Millisecond)
	}
	return false
}
