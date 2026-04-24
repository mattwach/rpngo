package fileops

import (
	"mattwach/rpngo/common/parse"
	"mattwach/rpngo/common/rpn"
)

const sourceHelp = "Loads the given path and executes commands within it.\n" +
	"Example: 'myfile.txt' source"

func (fo *FileOps) source(r *rpn.RPN) error {
	f, err := r.PopFrame()
	if err != nil {
		return err
	}
	if !f.IsString() {
		return rpn.ErrExpectedAString
	}
	path := f.UnsafeString()
	data, err := fo.driver.ReadFile(path)
	if err != nil {
		return err
	}
	if err := parse.Fields(string(data), r.Exec); err != nil {
		return err
	}
	return nil
}
