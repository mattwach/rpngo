package functions

import (
	"mattwach/rpngo/common/parse"
	"mattwach/rpngo/common/rpn"
	"strings"
)

const fieldsHelp = "Splits a string into fields and places all fields on the stack"

func fields(r *rpn.RPN) error {
	f, err := r.PopFrame()
	if err != nil {
		return err
	}
	if !f.IsString() {
		return rpn.ErrExpectedAString
	}
	fn := func(s string) error {
		return r.PushFrame(rpn.StringFrame(s, f.Type()))
	}
	return parse.Fields(f.String(false), fn)
}

const splitHelp = "Splits a string based on the given argument and places all " +
	"fields on the stack.\n\nexample: '25:59' ':' split"

func split(r *rpn.RPN) error {
	a, b, err := r.Pop2Frames()
	if err != nil {
		return err
	}
	if !a.IsString() || !b.IsString() {
		return rpn.ErrExpectedAString
	}
	parts := strings.Split(a.String(false), b.String(false))
	for _, part := range parts {
		if err := r.PushFrame(rpn.StringFrame(part, a.Type())); err != nil {
			return err
		}
	}
	return nil
}
