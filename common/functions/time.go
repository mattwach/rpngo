package functions

import (
	"mattwach/rpngo/common/rpn"
	"time"
)

const delayHelp = "Pauses for the given number of seconds"

const delayInterruptCheckSeconds = 0.25
const delayInterruptDuration = time.Duration(delayInterruptCheckSeconds*1000) * time.Millisecond

var DelaySleepFn = func(r *rpn.RPN, t float64) error {
	if t > delayInterruptCheckSeconds {
		// long sleeps need to check the interrupt function so that the user can
		// break and the watchdog timer can be satisfied
		deadline := time.Now().Add(time.Duration(t*1000000) * time.Microsecond)
		for time.Now().Add(delayInterruptDuration).Before(deadline) {
			if r.Interrupt() {
				return rpn.ErrInterrupted
			}
			time.Sleep(delayInterruptDuration)
		}
		t = float64(deadline.Sub(time.Now()).Microseconds()) / 1000000
	}
	time.Sleep(time.Duration(t*1000) * time.Millisecond)
	return nil
}

func delay(r *rpn.RPN) error {
	f, err := r.PopFrame()
	if err != nil {
		return err
	}
	v, err := f.Real()
	if err != nil {
		return err
	}
	if v <= 0 {
		return nil
	}
	return DelaySleepFn(r, v)
}

const timeHelp = "Returns unix epoch time, assuming the clock on the hardware is calibrated."

func timeFn(r *rpn.RPN) error {
	t := time.Now()
	return r.PushFrame(rpn.RealFrame(float64(t.UnixMicro()) / 1000000))
}
