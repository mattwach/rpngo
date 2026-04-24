package functions

import (
	"mattwach/rpngo/common/parse"
	"mattwach/rpngo/common/rpn"
)

const IfHelp = "Pops action, then val. Executes val if cond is true.\n" +
	"Example: 4 3 > {\"a is greater than b\" printlnx} if"

func If(r *rpn.RPN) error {
	f, err := r.PopFrame()
	if err != nil {
		return err
	}
	cf, err := r.PopFrame()
	if err != nil {
		return err
	}
	cond, err := cf.Bool()
	if err != nil {
		return err
	}
	if cond {
		if err := parse.Fields(f.String(false), r.Exec); err != nil {
			return err
		}
	}
	return nil
}

const ifelseHelp = "Pops elsev, ifv, cond. Executes ifv if cond=true, else executes elsev.\n" +
	"Example: 4 3 > {'a is greater than b'} {'a is less than b'} ifelse printlnx"

func ifelse(r *rpn.RPN) error {
	elsev, err := r.PopFrame()
	if err != nil {
		return err
	}
	ifv, err := r.PopFrame()
	if err != nil {
		return err
	}
	cf, err := r.PopFrame()
	if err != nil {
		return err
	}
	cond, err := cf.Bool()
	if err != nil {
		return err
	}
	if cond {
		if err := parse.Fields(ifv.String(false), r.Exec); err != nil {
			return err
		}
	} else {
		if err := parse.Fields(elsev.String(false), r.Exec); err != nil {
			return err
		}
	}
	return nil
}

const forHelp = "Executes the head of the stack in a loop until a value < is found\n" +
	"Example: 1 'c 1 + c 50 <' for # put 1 to 50 on the stack"

// avoid allocating fields on the stack (for can be nested though)
var forFields []string
var forFieldsStart []int

func forFn(r *rpn.RPN) error {
	forFieldsStart = append(forFieldsStart, len(forFields))
	defer func() { forFieldsStart = forFieldsStart[:len(forFieldsStart)-1] }()
	mf, err := r.PopFrame()
	if err != nil {
		return err
	}
	if !mf.IsString() {
		return rpn.ErrExpectedAString
	}
	macro := mf.UnsafeString()
	addField := func(t string) error {
		forFields = append(forFields, t)
		return nil
	}
	if err := parse.Fields(macro, addField); err != nil {
		return err
	}
	for {
		if err := r.ExecSlice(forFields[forFieldsStart[len(forFieldsStart)-1]:]); err != nil {
			return err
		}
		cf, err := r.PopFrame()
		if err != nil {
			return err
		}
		cond, err := cf.Bool()
		if err != nil {
			return err
		}
		if !cond {
			break
		}
	}
	return nil
}
