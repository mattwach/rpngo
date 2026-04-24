package functions

import (
	"mattwach/rpngo/common/rpn"
	"testing"
)

func TestPrint(t *testing.T) {
	var got string
	prfn := func(msg string) {
		got += msg
	}
	data := []rpn.UnitTestExecData{
		{
			Args:    []string{"print"},
			WantErr: rpn.ErrNotEnoughStackFrames,
		},
		{
			Args: []string{"'foo'", "print"},
			Want: []string{"'foo'"},
		},
	}
	rpn.UnitTestExecAll(t, data, func(r *rpn.RPN) {
		r.Print = prfn
		RegisterAll(r)
	})
	want := "foo"
	if got != want {
		t.Errorf("got: %v, want %v", got, want)
	}
}

func TestPrintX(t *testing.T) {
	var got string
	prfn := func(msg string) {
		got += msg
	}
	data := []rpn.UnitTestExecData{
		{
			Args:    []string{"printx"},
			WantErr: rpn.ErrStackEmpty,
		},
		{
			Args: []string{"'foo'", "printx"},
		},
	}
	rpn.UnitTestExecAll(t, data, func(r *rpn.RPN) {
		r.Print = prfn
		RegisterAll(r)
	})
	want := "foo"
	if got != want {
		t.Errorf("got: %v, want %v", got, want)
	}
}

func TestPrintS(t *testing.T) {
	var got string
	prfn := func(msg string) {
		got += msg
	}
	data := []rpn.UnitTestExecData{
		{
			Args:    []string{"prints"},
			WantErr: rpn.ErrNotEnoughStackFrames,
		},
		{
			Args: []string{"'foo'", "prints"},
			Want: []string{"'foo'"},
		},
	}
	rpn.UnitTestExecAll(t, data, func(r *rpn.RPN) {
		r.Print = prfn
		RegisterAll(r)
	})
	want := "foo "
	if got != want {
		t.Errorf("got: %v, want %v", got, want)
	}
}

func TestPrintSX(t *testing.T) {
	var got string
	prfn := func(msg string) {
		got += msg
	}
	data := []rpn.UnitTestExecData{
		{
			Args:    []string{"printsx"},
			WantErr: rpn.ErrStackEmpty,
		},
		{
			Args: []string{"'foo'", "printsx"},
		},
	}
	rpn.UnitTestExecAll(t, data, func(r *rpn.RPN) {
		r.Print = prfn
		RegisterAll(r)
	})
	want := "foo "
	if got != want {
		t.Errorf("got: %v, want %v", got, want)
	}
}

func TestPrintLn(t *testing.T) {
	var got string
	prfn := func(msg string) {
		got += msg
	}
	data := []rpn.UnitTestExecData{
		{
			Args:    []string{"println"},
			WantErr: rpn.ErrNotEnoughStackFrames,
		},
		{
			Args: []string{"'foo'", "println"},
			Want: []string{"'foo'"},
		},
	}
	rpn.UnitTestExecAll(t, data, func(r *rpn.RPN) {
		r.Print = prfn
		RegisterAll(r)
	})
	want := "foo\n"
	if got != want {
		t.Errorf("got: %v, want %v", got, want)
	}
}

func TestPrintLnX(t *testing.T) {
	var got string
	prfn := func(msg string) {
		got += msg
	}
	data := []rpn.UnitTestExecData{
		{
			Args:    []string{"printlnx"},
			WantErr: rpn.ErrStackEmpty,
		},
		{
			Args: []string{"'foo'", "printlnx"},
		},
	}
	rpn.UnitTestExecAll(t, data, func(r *rpn.RPN) {
		r.Print = prfn
		RegisterAll(r)
	})
	want := "foo\n"
	if got != want {
		t.Errorf("got: %v, want %v", got, want)
	}
}

func TestPrintAll(t *testing.T) {
	var got string
	prfn := func(msg string) {
		got += msg
	}
	data := []rpn.UnitTestExecData{
		{
			Args: []string{"'foo'", "'bar'", "printall"},
			Want: []string{"'foo'", "'bar'"},
		},
	}
	rpn.UnitTestExecAll(t, data, func(r *rpn.RPN) {
		r.Print = prfn
		RegisterAll(r)
	})
	want := "1: 'foo'\n0: 'bar'\n"
	if got != want {
		t.Errorf("got: %v, want %v", got, want)
	}
}

func TestInput(t *testing.T) {
	infn := func(r *rpn.RPN) (string, error) {
		return "foo", nil
	}
	data := []rpn.UnitTestExecData{
		{
			Args: []string{"input"},
			Want: []string{"{foo}"},
		},
	}
	rpn.UnitTestExecAll(t, data, func(r *rpn.RPN) {
		r.Input = infn
		RegisterAll(r)
	})
}

func TestHexDump(t *testing.T) {
	var r rpn.RPN
	r.Init(16)
	r.TextWidth = 30
	var got string
	r.Print = func(s string) {
		got += s
	}
	err := r.PushFrame(rpn.StringFrame("Hello  World!", rpn.STRING_BRACE_FRAME))
	if err != nil {
		t.Errorf("err=%v, want nil", err)
	}
	err = hexdump(&r)
	if err != nil {
		t.Errorf("err=%v, want nil", err)
	}
	want := "0000| 48 65 6c 6c  Hell\n0004| 6f 20 20 57  o  W\n" +
		"0008| 6f 72 6c 64  orld\n000c| 21           !\n"
	if got != want {
		t.Errorf("want:\n%s\ngot:\n%s", want, got)
	}
}
