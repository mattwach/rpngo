//go:build !pico && !pico2

package elog

import "log"

// Logging in tinygo is expensive, both because log.Printf() takes a bit of
// memory/resources and becuase the code can get slowed down while the
// log messages try to stream over the UART (generally at 115 kbps)
//
// This file provides a simple wrapper solution where cheap calls can be put
// in the embedded version.  Also there are some calls (like heap logging)
// can are useful to selectively comment in and out when they are needed.

func Print(v ...any) {
	log.Print(v...)
}

// These are primarily for embedded logging.  We don't care on PC hardware.
func Heap(msg string) {}
