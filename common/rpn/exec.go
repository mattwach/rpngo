package rpn

import (
	"fmt"
	"mattwach/rpngo/common/elog"
	"strconv"
	"strings"
)

// Exec executes a single instruction
func (rpn *RPN) Exec(arg string) error {
	if rpn.Interrupt() {
		return ErrInterrupted
	}
	if fn := rpn.functions[arg]; fn != nil {
		return fn(rpn)
	}
	if len(arg) > 1 {
		switch arg[len(arg)-1] {
		case '=':
			if (len(arg) >= 3) && arg[len(arg)-2] == '=' {
				return rpn.moveAllStackToVariable(arg[:len(arg)-2])
			}
			return rpn.SetVariable(arg[:len(arg)-1])
		case '/':
			return rpn.ClearVariable(arg[:len(arg)-1])
		case '>':
			if (len(arg) >= 3) && arg[len(arg)-2] == '>' {
				return rpn.moveAllVariableToStack(arg[:len(arg)-2])
			}
			return rpn.moveVariableToHead(arg[:len(arg)-1])
		case '<':
			if (len(arg) >= 3) && arg[len(arg)-2] == '<' {
				return rpn.appendAllStackToVariable(arg[:len(arg)-2])
			}
			return rpn.moveHeadToVariable(arg[:len(arg)-1])
		}
		switch arg[0] {
		case '$':
			if (len(arg) >= 3) && arg[1] == '$' {
				return rpn.appendAllVariableToStack(arg[2:])
			}
			f, err := rpn.GetVariable(arg[1:])
			if err != nil {
				return err
			}
			rpn.PushFrame(f)
			return nil
		case '@':
			return rpn.execVariableAsMacro(arg[1:])
		case '`':
			return rpn.addLabel(arg)
		}
	}
	if len(arg) >= 2 {
		last := arg[len(arg)-1]
		switch last {
		case '"':
			if arg[0] == '"' {
				return rpn.PushFrame(StringFrame(arg[1:len(arg)-1], STRING_DOUBLEQ_FRAME))
			}
		case '\'':
			if arg[0] == '\'' {
				return rpn.PushFrame(StringFrame(arg[1:len(arg)-1], STRING_SINGLEQ_FRAME))
			}
		case '}':
			if arg[0] == '{' {
				return rpn.PushFrame(StringFrame(arg[1:len(arg)-1], STRING_BRACE_FRAME))
			}
		case 'd':
			return rpn.parseAndPushInt(arg[:len(arg)-1], 10, INTEGER_FRAME)
		case 'x':
			return rpn.parseAndPushInt(arg[:len(arg)-1], 16, HEXIDECIMAL_FRAME)
		case 'o':
			return rpn.parseAndPushInt(arg[:len(arg)-1], 8, OCTAL_FRAME)
		case 'b':
			return rpn.parseAndPushInt(arg[:len(arg)-1], 2, BINARY_FRAME)
		}
	}
	if len(arg) > 0 && arg[len(arg)-1] == '?' {
		return rpn.printHelp(arg[:len(arg)-1])
	}
	if strings.Contains(arg, ">") {
		return rpn.convert(arg)
	}
	return rpn.parseAndPushComplex(arg)
}

func (rpn *RPN) ExecSlice(args []string) error {
	for i, arg := range args {
		if err := rpn.Exec(arg); err != nil {
			elog.Heap("alloc: /rpn/exec.go:73: return fmt.Errorf('exec %s: %w', highlightArg(args, i), err)")
			return fmt.Errorf("exec %s: %w", highlightArg(args, i), err) // object allocated on the heap: escapes at line 73
		}
	}
	return nil
}

func highlightArg(args []string, idx int) string {
	var parts []string
	for i, arg := range args {
		if i == idx {
			arg = "->" + arg + "<-"
		}
		parts = append(parts, arg)
	}
	return strings.Join(parts, " ")
}

func (rpn *RPN) parseAndPushInt(arg string, base int, t FrameType) error {
	v, err := strconv.ParseInt(arg, base, 64)
	if err != nil {
		return ErrSyntax
	}
	return rpn.PushFrame(IntFrame(v, t))
}

// Pushes a float onto the stack
func (rpn *RPN) parseAndPushComplex(arg string) error {
	var v complex128
	var err error

	if strings.HasSuffix(arg, "i") {
		v, err = parseComplexWithI(arg)
		if err != nil {
			return err
		}
	} else {
		fv, err := strconv.ParseFloat(arg, 64)
		if err != nil {
			return rpn.parseAndPushPolar(arg)
		}
		v = complex(fv, 0)
	}
	return rpn.PushFrame(ComplexFrame(v))
}

func (rpn *RPN) parseAndPushPolar(arg string) error {
	ltIdx := strings.IndexRune(arg, '<')
	if (ltIdx < 0) || (ltIdx >= (len(arg) - 1)) {
		return ErrSyntax
	}
	real := arg[:ltIdx]
	cmplx := arg[ltIdx+1:]
	r, err := strconv.ParseFloat(real, 64)
	if err != nil {
		return err
	}
	a, err := strconv.ParseFloat(cmplx, 64)
	if err != nil {
		return err
	}
	return rpn.PushFrame(PolarFrame(r, a, rpn.AngleUnit))

}

// parses a complex string that contains an i
func parseComplexWithI(arg string) (complex128, error) {
	// a is the "real" part and b is the "imag" part: a + bi
	var a string
	var b string
	for i, c := range arg {
		if c == '+' || (c == '-' && i > 0) {
			a = arg[:i]
			b = arg[i : len(arg)-1]
			break
		}
	}
	if a == "" {
		// no real part was given.  e.g. 5i
		a = "0"
		b = arg[:len(arg)-1]
	}
	b = strings.TrimPrefix(b, "+")
	switch b {
	case "":
		// the user specified just i
		b = "1"
	case "-":
		// the user specified just -i
		b = "-1"
	}
	fa, err := strconv.ParseFloat(a, 64)
	if err != nil {
		return 0, ErrSyntax
	}
	fb, err := strconv.ParseFloat(b, 64)
	if err != nil {
		return 0, ErrSyntax
	}
	return complex(fa, fb), nil
}

func (r *RPN) convert(arg string) error {
	f, err := r.PopFrame()
	if err != nil {
		return err
	}
	v, err := f.Real()
	if err != nil {
		return err
	}
	parts := strings.SplitN(arg, ">", 2)
	newv, err := r.conv.Convert(v, parts[0], parts[1])
	if err != nil {
		return err
	}
	newf := RealFrame(newv)
	newf.Annotate("`" + parts[1])
	return r.PushFrame(newf)
}

func (r *RPN) addLabel(label string) error {
	if len(r.Frames) == 0 {
		return ErrStackEmpty
	}
	f := &r.Frames[len(r.Frames)-1]
	if (f.ftype & CLASS_MASK) == STRING_CLASS {
		return ErrCanNotAddLabelToString
	}
	f.str = label
	return nil
}
