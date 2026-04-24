package functions

import (
	"mattwach/rpngo/common/rpn"
	"strconv"
)

const printHelp = "Prints the head element of the stack to the output window"

func printFn(r *rpn.RPN) error {
	f, err := r.PeekFrame(0)
	if err != nil {
		return err
	}
	r.Print(f.String(false))
	return nil
}

const printxHelp = "Pops head element of the stack and prints it"

func printx(r *rpn.RPN) error {
	f, err := r.PopFrame()
	if err != nil {
		return err
	}
	r.Print(f.String(false))
	return nil
}

const printsHelp = "Prints the head element of the stack plus a space"

func prints(r *rpn.RPN) error {
	f, err := r.PeekFrame(0)
	if err != nil {
		return err
	}
	r.Print(f.String(false))
	r.Print(" ")
	return nil
}

const printsxHelp = "Pops head element of the stack and prints it and a space"

func printsx(r *rpn.RPN) error {
	f, err := r.PopFrame()
	if err != nil {
		return err
	}
	r.Print(f.String(false))
	r.Print(" ")
	return nil
}

const printlnHelp = "Prints the head element of the stack plus a newline"

func printlnFn(r *rpn.RPN) error {
	f, err := r.PeekFrame(0)
	if err != nil {
		return err
	}
	r.Print(f.String(false))
	r.Print("\n")
	return nil
}

const printlnxHelp = "Pops head element of the stack and prints it and a newline"

func printlnx(r *rpn.RPN) error {
	f, err := r.PopFrame()
	if err != nil {
		return err
	}
	r.Print(f.String(false))
	r.Print("\n")
	return nil
}

const printallHelp = "Prints the whole stack"

func printall(r *rpn.RPN) error {
	i := len(r.Frames)
	for _, f := range r.Frames {
		i--
		r.Print(strconv.Itoa(i) + ": " + f.String(true) + "\n")
	}
	return nil
}

const inputHelp = "Pauses for user input and pushes the result to the stack as a string"

func input(r *rpn.RPN) error {
	str, err := r.Input(r)
	if err != nil {
		return err
	}
	return r.PushFrame(rpn.StringFrame(str, rpn.STRING_BRACE_FRAME))
}

const hexdumpHelp = "Hex dump of the top of the stack (converted to a string)"

func hexdump(r *rpn.RPN) error {
	f, err := r.PopFrame()
	if err != nil {
		return err
	}
	str := f.String(false)
	bytesPerRow := 1
	for {
		// Example if bytesPerRow = 4
		// 0000| 40 40 40 40  AAAA
		// (23 bytes needed)
		widthNeeded := (bytesPerRow * 4) + (bytesPerRow / 4) + 6
		if widthNeeded >= r.TextWidth {
			bytesPerRow /= 2
			break
		}
		bytesPerRow *= 2
	}

	rowStart := 0
	for rowStart < len(str) {
		s := strconv.FormatInt(int64(rowStart), 16)
		for i := len(s); i < 4; i++ {
			r.Print("0")
		}
		r.Print(s)
		r.Print("| ")
		bytesThisRow := len(str) - rowStart
		if bytesThisRow > bytesPerRow {
			bytesThisRow = bytesPerRow
		}
		for i := 0; i < bytesThisRow; i++ {
			s = strconv.FormatInt(int64(str[rowStart+i]), 16)
			if len(s) < 2 {
				r.Print("0")
			}
			r.Print(s)
			r.Print(" ")
			if ((rowStart + i + 1) % 4) == 0 {
				r.Print(" ")
			}
		}
		for i := bytesThisRow; i < bytesPerRow; i++ {
			r.Print("   ")
			if ((rowStart + i + 1) % 4) == 0 {
				r.Print(" ")
			}
		}
		for i := 0; i < bytesThisRow; i++ {
			c := byte(str[rowStart+i])
			if (c < 32) || (c > 127) {
				r.Print(".")
			} else {
				r.Print(string(str[rowStart+i]))
			}
		}
		r.Print("\n")
		rowStart += bytesPerRow
	}
	return nil
}
