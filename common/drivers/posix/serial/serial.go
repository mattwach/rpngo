package serial

import (
	"errors"
	"fmt"
	"log"
	"mattwach/rpngo/common/rpn"
	"os"
)

var (
	errNoData                 = errors.New("no data")
	errSerialPortNotOpen      = errors.New("serial port not open")
	errSerialPortNeedsToBeSet = errors.New("$.serial needs to be set (e.g. /dev/ttyAMC0)")
)

type readData struct {
	read chan byte
	err  chan error
	done chan bool
}

type Serial struct {
	f     *os.File
	rdata *readData
}

func (sc *Serial) Open(r *rpn.RPN) error {
	if sc.f != nil {
		return fmt.Errorf("serial file is already open: v", sc.f.Name())
	}
	f, err := r.GetVariable(".serial")
	if err != nil {
		return err
	}
	if !f.IsString() {
		return rpn.ErrExpectedAString
	}
	sc.f, err = os.OpenFile(f.UnsafeString(), os.O_RDWR|os.O_SYNC, 0666)
	if err != nil {
		log.Printf("sopen err %v", err)
		return err
	}
	log.Print("sopen ok")
	return nil
}

func (sc *Serial) Close() error {
	log.Printf("sclose")
	if sc.f == nil {
		return errSerialPortNotOpen
	}
	err := sc.f.Close()
	if sc.rdata != nil {
		<-sc.rdata.done
		sc.rdata = nil
	}
	sc.f = nil
	return err
}

var readbuff = make([]byte, 1)

func (sc *Serial) ReadByte() (byte, error) {
	if sc.f == nil {
		return 0, errSerialPortNotOpen
	}

	if sc.rdata == nil {
		sc.rdata = &readData{
			read: make(chan byte, 16),
			err:  make(chan error, 1),
			done: make(chan bool),
		}
		go func() {
			log.Print("rdata loop enter")
			for {
				count, err := sc.f.Read(readbuff)
				if err != nil {
					log.Printf("rb err=%v", err)
					sc.rdata.err <- err
					break
				}
				if count > 0 {
					log.Printf("rb got=%v %c", readbuff[0], readbuff[0])
					sc.rdata.read <- readbuff[0]
				}
			}
			log.Print("rdata loop exit")
			sc.rdata.done <- true
		}()
	}

	select {
	case rd := <-sc.rdata.read:
		log.Printf("rbr byte: %v %c", rd, rd)
		return rd, nil
	case err := <-sc.rdata.err:
		log.Printf("rbr err=%v", err)
		return 0, err
	default:
		return 0, errNoData
	}
}

func (sc *Serial) WriteByte(c byte) error {
	if sc.f == nil {
		return errSerialPortNotOpen
	}
	_, err := sc.f.Write([]byte{c})
	return err
}
