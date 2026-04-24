package rpn

import "testing"

func TestInit(t *testing.T) {
	var r RPN
	r.Init(256)
	err := r.ExecSlice([]string{"11", "22", "33", "s.size"})
	if err != nil {
		t.Fatalf("r.Exec() err=%v", err)
	}
	if len(r.Frames) != 4 {
		t.Fatalf("want len(r.frames) = 4, got %v", len(r.Frames))
	}
	f := r.Frames[3]
	if f.ftype != INTEGER_FRAME {
		t.Errorf("want frame.Type = INTEGER_FRAME, got %v", f.ftype)
	}
	if f.intv != 3 {
		t.Errorf("want frame.Int = 3, got %v", f.intv)
	}
}

func TestRegister(t *testing.T) {
	var r RPN
	r.Init(256)
	fn := func(r *RPN) error {
		return r.PushFrame(IntFrame(55, INTEGER_FRAME))
	}
	r.Register("fiftyfive", fn, "helpcat", "helptxt")
	err := r.ExecSlice([]string{"fiftyfive"})
	if err != nil {
		t.Fatalf("r.Exec() err=%v", err)
	}
	if len(r.Frames) != 1 {
		t.Fatalf("want len(r.frames) = 1, got %v", len(r.Frames))
	}
	f := r.Frames[0]
	if f.ftype != INTEGER_FRAME {
		t.Errorf("want frame.Type = INTEGER_FRAME, got %v", f.ftype)
	}
	if f.intv != 55 {
		t.Errorf("want frame.Int = 55, got %v", f.intv)
	}
	if r.help["helpcat"]["fiftyfive"] != "helptxt" {
		t.Errorf("want 'helptxt', got %v", r.help["helpcat"]["fiftyfive"])
	}
}

func TestAllFunctionNames(t *testing.T) {
	var r RPN
	r.Init(256)
	names := r.AllFunctionNames()
	var found bool
	for _, n := range names {
		if n == "s.size" {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("AllFunctionNames() did not contain s.size")
	}
}

func TestPrintln(t *testing.T) {
	var r RPN
	r.Init(256)
	var got string
	r.Print = func(msg string) {
		got = got + msg
	}
	r.Println("hello")
	if got != "hello\n" {
		t.Errorf("want hello\\n, got %v", got)
	}
}
