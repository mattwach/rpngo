// Run rpngo on a microcontroller
//
// This is a "minimialist" implementation which can be thought of
// as a valdation stepping stone.
package main

import (
	"machine"
	"mattwach/rpngo/bin/tinygo"
	"mattwach/rpngo/common/functions"
	"mattwach/rpngo/common/parse"
	"mattwach/rpngo/common/rpn"
	"time"
)

// Globally allocate RPN to avoid Tijnygo heap usage
var r rpn.RPN

func main() {
	time.Sleep(2 * time.Second)
	r.Init(256)
	tinygo.Register(&r)
	functions.RegisterAll(&r)

	print("Type ? for help or topic? for more detailed help")

	for {
		line := readLine()
		var err error
		err = parse.Fields(line, r.Exec)
		if err == nil {
			for _, f := range r.Frames {
				for _, r := range f.String(true) {
					machine.Serial.WriteByte(byte(r))
				}
				machine.Serial.WriteByte('\n')
			}
		} else {
			print("Error: ", err.Error(), "\n")
		}
	}
}

func readLine() string {
	var msg []byte
	print("> ")
	for {
		c, err := machine.Serial.ReadByte()
		machine.Serial.WriteByte(c)
		if err == nil {
			switch c {
			case 0:
				break
			case 13:
				machine.Serial.WriteByte('\n')
				return string(msg)
			default:
				msg = append(msg, c)
			}
		}

		time.Sleep(time.Millisecond * 10)
	}
}
