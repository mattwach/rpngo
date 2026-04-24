package xmodem

import (
	"errors"
	"mattwach/rpngo/common/rpn"
	"strconv"
	"time"
)

var (
	errCRCMismatch            = errors.New("crc mismatch")
	errEarlyReceiverResponse  = errors.New("receiver responded mid-packet")
	errEndOfTransmission      = errors.New("end of transmission")
	errIncorrectPacketId      = errors.New("incorrect packet id")
	errIncorrectPacketIdSum   = errors.New("incorrect packet id sum")
	errInitSerial             = errors.New("please init serial with: true serial")
	errInvalidInitialPacket   = errors.New("invalid initial packet type")
	errReadTimeout            = errors.New("read timeout")
	errNakReceived            = errors.New("nak received")
	errNakSent                = errors.New("nak sent")
	errResponseTimeout        = errors.New("response timeout")
	errUnexpectedHeaderByte   = errors.New("unexpected header byte")
	errTimeoutWaitingForStart = errors.New("timeout waiting for start message")
	errUnexpectedByteReceived = errors.New("unexpected byte received (not ack/nak)")
)

type Serial interface {
	Open(r *rpn.RPN) error
	Close() error
	ReadByte() (byte, error) // Non-blocking
	WriteByte(c byte) error
}

type XmodemCommands struct {
	serial       Serial
	attemptsLeft int
	state        readState
	deadline     time.Time
	nextPacketId uint8
	idx          int
	packet       [133]uint8
	buff         []byte
}

type echostate int

const (
	ECHO_ERROR echostate = iota
	ECHO_FALSE
	ECHO_TRUE
)

const handshakeAttempt = 10

func (sc *XmodemCommands) InitAndRegister(r *rpn.RPN, serial Serial) {
	sc.serial = serial
	r.Register("rx", sc.xmodemRead, rpn.CatIO, xmodemReadHelp)
	r.Register("sx", sc.xmodemWrite, rpn.CatIO, xmodemWriteHelp)
}

func (sc *XmodemCommands) SetSerial(serial Serial) {
	sc.serial = serial
}

const xmodemReadHelp = "Attempts to read data using the xmodem protocal to the top of the stack."

type readState uint8

const (
	FlushReadBuffer readState = iota
	InitialHandshake
	XferPackets
	Finished
)

const (
	SOH uint8 = 0x01
	EOT       = 0x04
	ACK       = 0x06
	NAK       = 0x15
	ETB       = 0x17
	CAN       = 0x18
	C         = 0x43
)

func (sc *XmodemCommands) xmodemRead(r *rpn.RPN) error {
	if sc.serial == nil {
		return errInitSerial
	}
	if err := sc.serial.Open(r); err != nil {
		return err
	}
	defer sc.serial.Close()
	st := disableEcho(r)
	defer restoreEcho(r, st)
	sc.state = FlushReadBuffer
	sc.attemptsLeft = handshakeAttempt
	sc.nextPacketId = 0x01
	sc.buff = make([]byte, 0, 128)
	var err error
	r.Println("rx: init")
	for err == nil {
		switch sc.state {
		case FlushReadBuffer:
			sc.flushReadBuffer()
		case InitialHandshake:
			err = sc.serial.WriteByte(C)
			if err == nil {
				err = sc.readPacket(r)
				if err != nil {
					err = sc.nakWith(r, "first packet", err)
				}
			}
		case XferPackets:
			err = sc.readPacket(r)
			if err != nil {
				err = sc.nakWith(r, "packet", err)
			}
		case Finished:
			return r.PushFrame(rpn.StringFrame(sc.trimExtra(), rpn.STRING_BRACE_FRAME))
		}
	}
	return err
}

func (sc *XmodemCommands) xmodemWrite(r *rpn.RPN) error {
	if sc.serial == nil {
		return errInitSerial
	}
	if err := sc.serial.Open(r); err != nil {
		return err
	}
	defer sc.serial.Close()
	s, err := r.PopFrame()
	if err != nil {
		return err
	}
	st := disableEcho(r)
	defer restoreEcho(r, st)
	sc.buff = []byte(s.String(false))
	sc.state = FlushReadBuffer
	sc.attemptsLeft = handshakeAttempt
	sc.nextPacketId = 0x01
	sc.deadline = time.Now().Add(60 * time.Second)
	r.Println("sx: wait for receiver")
	for err == nil {
		switch sc.state {
		case FlushReadBuffer:
			sc.flushReadBuffer()
		case InitialHandshake:
			b, err := sc.serial.ReadByte()
			if err == nil {
				if b == C {
					r.Println("sx: xfer start")
					sc.state = XferPackets
				}
			}
			if time.Now().After(sc.deadline) {
				return errTimeoutWaitingForStart
			}
		case XferPackets:
			err = sc.writePacket(r)
		case Finished:
			return nil
		}
	}
	return err
}

func disableEcho(r *rpn.RPN) echostate {
	f, err := r.GetVariable(".echo")
	if err != nil {
		return ECHO_ERROR
	}
	en, err := f.Bool()
	if err != nil {
		return ECHO_ERROR
	}
	restoreEcho(r, ECHO_FALSE)
	if en {
		return ECHO_TRUE
	}
	return ECHO_FALSE
}

func restoreEcho(r *rpn.RPN, es echostate) {
	if es == ECHO_ERROR {
		return
	}
	if err := r.PushFrame(rpn.BoolFrame(es == ECHO_TRUE)); err != nil {
		return
	}
	r.SetVariable(".echo")
}

// Try to pull characters from the read buffer until it goes silent
// for 1 second.
func (sc *XmodemCommands) flushReadBuffer() {
	deadline := time.Now().Add(1 * time.Second)
	for {
		b, err := sc.serial.ReadByte()
		if (err == nil) && (b != 'C') {
			deadline = time.Now().Add(1 * time.Second)
		} else if time.Now().After(deadline) {
			sc.state = InitialHandshake
			return
		}
	}
}

func (sc *XmodemCommands) trimExtra() string {
	// I see a string of 0x1A (SUB) characters in the last
	// block.  A bit surprising.
	end := len(sc.buff) - 1
	for ; end > 0; end-- {
		if (sc.buff[end] == '\n') || (sc.buff[end] >= 32) {
			break
		}
	}
	s := string(sc.buff[:end+1])
	sc.buff = nil
	return s
}

func (sc *XmodemCommands) writePacket(r *rpn.RPN) error {
	r.Print("sc write (")
	r.Print(strconv.Itoa(int(sc.nextPacketId)))
	r.Print(") ")
	lastPacket := sc.buildWritePacket()
	for _, b := range sc.packet {
		err := sc.serial.WriteByte(b)
		if err != nil {
			r.Println(err.Error())
			return err
		}
		b, err = sc.serial.ReadByte()
		if err == nil {
			return sc.reduceAttempts(r, "read byte "+string(b), errEarlyReceiverResponse)
		}
	}

	if err := sc.waitForAck(r); err != nil {
		return sc.reduceAttempts(r, "wait ACK", err)
	}

	if lastPacket {
		return sc.sendEOT(r)
	}

	return nil
}

func (sc *XmodemCommands) sendEOT(r *rpn.RPN) error {
	r.Print("sx: EOT ")
	for {
		err := sc.serial.WriteByte(EOT)
		if err != nil {
			return err
		}
		err = sc.waitForAck(r)
		if err == nil {
			sc.state = Finished
			return nil
		}
		err = sc.reduceAttempts(r, "EOT", err)
		if err != nil {
			return err
		}
	}
}

func (sc *XmodemCommands) waitForAck(r *rpn.RPN) error {
	deadline := time.Now().Add(3 * time.Second)
	for {
		c, err := sc.serial.ReadByte()
		if err == nil {
			switch c {
			case ACK:
				sc.attemptsLeft = handshakeAttempt
				sc.nextPacketId++
				r.Println("ACK")
				return nil
			case NAK:
				return errNakReceived
			default:
				return errors.New("got " + string(c))
			}
		}
		if time.Now().After(deadline) {
			return errResponseTimeout
		}
	}
}

func (sc *XmodemCommands) buildWritePacket() bool {
	sc.packet[0] = SOH
	sc.packet[1] = sc.nextPacketId
	sc.packet[2] = 0xFF - sc.packet[1]
	start := (int(sc.nextPacketId) - 1) * 128
	end := start + 128
	lastPacket := end > len(sc.buff)
	if lastPacket {
		end = len(sc.buff)
		// fill in some pad
		for i := 3; i < (3 + 128); i++ {
			sc.packet[i] = 0x1A
		}
	}
	copy(sc.packet[3:], sc.buff[start:end])
	crc := sc.calcPacketChecksum()
	sc.packet[131] = uint8(crc >> 8)
	sc.packet[132] = uint8(crc & 0xFF)
	return lastPacket
}

func (sc *XmodemCommands) readPacket(r *rpn.RPN) error {
	for {
		if err := sc.readPacketData(r); err != nil {
			return err
		}
		if sc.state == Finished {
			r.Println("rx: finished")
			return nil
		}

		if err := sc.validatePacket(); err != nil {
			return err
		}

		r.Print("rx: recv (")
		r.Print(strconv.Itoa(int(sc.nextPacketId)))
		r.Print(") ")

		if sc.packet[1] == sc.nextPacketId {
			sc.nextPacketId++
			sc.attemptsLeft = handshakeAttempt
			sc.buff = append(sc.buff, sc.packet[3:131]...)
			r.Println(strconv.Itoa(len(sc.buff)))
			if sc.state == InitialHandshake {
				if sc.packet[0] != SOH {
					return errInvalidInitialPacket
				}
				sc.state = XferPackets
			}
			return sc.serial.WriteByte(ACK)
		} else if sc.packet[1] == (sc.nextPacketId - 1) {
			// sender might have not received our previous ack
			r.Println("repeat")
			if err := sc.serial.WriteByte(ACK); err != nil {
				return err
			}
		} else {
			r.Println("bad id")
			return errIncorrectPacketId
		}
	}

}

func (sc *XmodemCommands) readPacketData(r *rpn.RPN) error {
	if sc.state == InitialHandshake {
		sc.deadline = time.Now().Add(3 * time.Second)
	} else {
		sc.deadline = time.Now().Add(1 * time.Second)
	}

	sc.idx = 0
	for sc.idx < len(sc.packet) {
		if time.Now().After(sc.deadline) {
			return errReadTimeout
		}
		// read in bursts so we are not checking the timer on every byte
		for n := 0; n < 1024; n++ {
			b, err := sc.serial.ReadByte()
			if err != nil {
				// We get errors when the buffer is empty, which is almost
				// guaranteed to happen.  Just rely on the timeout.
				continue
			}
			if (sc.idx == 0) && (sc.state == XferPackets) {
				if b == EOT {
					if err := sc.serial.WriteByte(ACK); err != nil {
						return err
					}
					sc.state = Finished
					return nil
				}
			}
			sc.packet[sc.idx] = b
			sc.idx++
			if sc.idx >= len(sc.packet) {
				break
			}
		}
	}

	return nil
}

func (sc *XmodemCommands) reduceAttempts(r *rpn.RPN, ctx string, err error) error {
	sc.attemptsLeft--
	r.Print(ctx)
	r.Print(": ")
	r.Print(err.Error())
	r.Print(" attemptsLeft: ")
	r.Println(strconv.Itoa(sc.attemptsLeft))
	if sc.attemptsLeft <= 0 {
		return err
	}
	return nil
}

func (sc *XmodemCommands) nakWith(r *rpn.RPN, ctx string, err error) error {
	sc.serial.WriteByte(NAK)
	return sc.reduceAttempts(r, ctx+" NAK ", err)
}

func (sc *XmodemCommands) validatePacket() error {
	switch sc.packet[0] {
	case SOH, EOT, ETB, CAN:
		// ok
	default:
		return errUnexpectedHeaderByte
	}

	if uint16(sc.packet[1])+uint16(sc.packet[2]) != 0xFF {
		return errIncorrectPacketIdSum
	}

	crc := sc.calcPacketChecksum()
	if (sc.packet[131] != uint8(crc>>8)) || (sc.packet[132] != uint8(crc&0xFF)) {
		return errCRCMismatch
	}

	return nil
}

func (sc *XmodemCommands) calcPacketChecksum() uint16 {
	var crc uint16
	for _, b := range sc.packet[3:131] {
		crc ^= uint16(b) << 8
		for i := 0; i < 8; i++ {
			if crc&0x8000 != 0 {
				crc = (crc << 1) ^ 0x1021
			} else {
				crc = crc << 1
			}
		}
	}
	return crc
}

const xmodemWriteHelp = "Attempts to send the data at the top of the stack using the xmodem protocol."
