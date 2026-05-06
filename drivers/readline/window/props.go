package window

import (
	"mattwach/rpngo/common/parse"
	"mattwach/rpngo/common/rpn"
	"strings"
)

const MAX_SHOW_FRAMES = 1000

func (rlw *ReadlineWindow) SetProp(name string, val rpn.Frame) error {
	switch name {
	case "autofn":
		var fn []string
		addfn := func(s string) error {
			fn = append(fn, s)
			return nil
		}
		if err := parse.Fields(val.String(false), addfn); err != nil {
			return err
		}
		rlw.autofn = fn
		return nil
	case "prompt":
		if !val.IsString() {
			return rpn.ErrExpectedAString
		}
		rlw.inst.Config.Prompt = val.UnsafeString()
		rlw.inst.SetPrompt(val.UnsafeString())
		return nil
	case "showframes":
		v, err := val.BoundedInt(0, MAX_SHOW_FRAMES)
		if err != nil {
			return err
		}
		rlw.showFrames = int(v)
		return nil
	}
	return rpn.ErrNotSupported
}

func (rlw *ReadlineWindow) GetProp(name string) (rpn.Frame, error) {
	switch name {
	case "autofn":
		return rpn.StringFrame(strings.Join(rlw.autofn, " "), rpn.STRING_BRACE_FRAME), nil
	case "prompt":
		return rpn.StringFrame(rlw.inst.Config.Prompt, rpn.STRING_SINGLEQ_FRAME), nil
	case "showframes":
		return rpn.IntFrame(int64(rlw.showFrames), rpn.INTEGER_FRAME), nil
	}
	return rpn.Frame{}, rpn.ErrNotSupported
}

var inputProps = []string{"autofn", "prompt", "showframes"}

func (rlw *ReadlineWindow) ListProps() []string {
	return inputProps
}
