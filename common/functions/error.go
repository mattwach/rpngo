package functions

import (
	"errors"
	"mattwach/rpngo/common/parse"
	"mattwach/rpngo/common/rpn"
)

const tryHelp = "Executes first arg.  If there is an error, pushes the error " +
	"as a string, then executes the second arg.\n" +
	"Example: {$x $y /} { 0/ 'please try again' println} try"

func try(r *rpn.RPN) error {
	errorv, err := r.PopFrame()
	if err != nil {
		return err
	}
	tryv, err := r.PopFrame()
	if err != nil {
		return err
	}
	err = parse.Fields(tryv.String(false), r.Exec)
	if err == nil {
		return nil
	}
	err = r.PushFrame(rpn.StringFrame(err.Error(), rpn.STRING_SINGLEQ_FRAME))
	if err != nil {
		return err
	}
	return parse.Fields(errorv.String(false), r.Exec)
}

const errorHelp = "Pops a frame and returns it as an error"

func errorFn(r *rpn.RPN) error {
	f, err := r.PopFrame()
	if err != nil {
		return err
	}
	return errors.New(f.String(false))
}
