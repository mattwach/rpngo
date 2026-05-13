package common

import "mattwach/rpngo/common/rpn"

func NewWindowNameFromStack(r *rpn.RPN) (string, error) {
	name, err := r.PopFrame()
	if err != nil {
		return "", err
	}
	if !name.IsString() {
		return "", rpn.ErrExpectedAString
	}
	return name.UnsafeString(), nil
}
