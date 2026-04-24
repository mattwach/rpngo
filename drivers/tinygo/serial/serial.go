//go:build pico || pico2

// Package serialconsile provides a common mechanism for processsing
// characters from the UART
package serial

import (
	"machine"
	"mattwach/rpngo/common/key"
	"mattwach/rpngo/common/rpn"
	"time"
)

type TermState int

const (
	NORMAL TermState = iota
	ESC
	ARROW
	SHIFTARROW
	PAGEUP
	PAGEDOWN
)

// Escape codes make it a bit trickey to determine when only ESC was pressed
// Using a timeout delay is the usual tactic.
const escWait = 100 * time.Millisecond

type Serial struct {
	state       TermState
	ignoreUntil time.Time
	escTimeout  time.Time
}

func (sc *Serial) Init(enable bool) {
	sc.state = NORMAL
	sc.ignoreUntil = time.Now().Add(time.Second)
}

func (sc *Serial) Open(r *rpn.RPN) error {
	return nil
}

func (sc *Serial) Close() error {
	return nil
}

func (sc *Serial) ReadByte() (byte, error) {
	return machine.Serial.ReadByte()
}

func (sc *Serial) WriteByte(c byte) error {
	return machine.Serial.WriteByte(c)
}

func (sc *Serial) GetChar() key.Key {
	c, err := machine.Serial.ReadByte()
	if err != nil {
		if (sc.state == ESC) && time.Now().After(sc.escTimeout) {
			sc.state = NORMAL
			return 27
		}
		// nothing available
		return 0
	}
	if time.Now().Before(sc.ignoreUntil) {
		// This might be some random junk from starting up cold
		return 0
	}
	switch sc.state {
	case NORMAL:
		switch c {
		case 13:
			return '\n'
		case 27:
			sc.escTimeout = time.Now().Add(escWait)
			sc.state = ESC
		case 127:
			return key.KEY_BACKSPACE
		default:
			return key.Key(c)
		}
	case ESC:
		sc.state = NORMAL
		switch c {
		case '[':
			sc.state = ARROW
		case 'c':
			return key.KEY_COPY
		case 'h':
			return key.KEY_HELP
		case 'q':
			return key.KEY_QUIT
		case 's':
			return key.KEY_SAVE
		case 'v':
			return key.KEY_PASTE
		case 'x':
			return key.KEY_CUT
		}
	case ARROW:
		sc.state = NORMAL
		switch c {
		case 27:
			sc.state = SHIFTARROW
		case 'A':
			return key.KEY_UP
		case 'B':
			return key.KEY_DOWN
		case 'C':
			return key.KEY_RIGHT
		case 'D':
			return key.KEY_LEFT
		case '5':
			sc.state = PAGEUP
		case '6':
			sc.state = PAGEDOWN
		}
	case SHIFTARROW:
		sc.state = NORMAL
		switch c {
		case 'A':
			return key.KEY_SUP
		case 'B':
			return key.KEY_SDOWN
		case 'C':
			return key.KEY_SRIGHT
		case 'D':
			return key.KEY_SLEFT
		}
	case PAGEUP:
		sc.state = NORMAL
		if c == 126 {
			return key.KEY_PAGEUP
		}
	case PAGEDOWN:
		sc.state = NORMAL
		if c == 126 {
			return key.KEY_PAGEDOWN
		}
	}
	// Return zero here to avoid messing up the state machine
	return 0
}
