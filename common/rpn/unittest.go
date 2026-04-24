package rpn

import (
	"errors"
	"reflect"
	"testing"
)

type UnitTestExecData struct {
	Name    string
	Args    []string
	Want    []string
	WantErr error
}

// A runner to dave on boilerplate testing code
func UnitTestExecAll(t *testing.T, data []UnitTestExecData, prepfn func(*RPN)) {
	t.Helper()
	for _, d := range data {
		var r RPN
		r.Init(256)
		if prepfn != nil {
			prepfn(&r)
		}
		t.Run(d.Name, func(t *testing.T) {
			UnitTestExec(t, &r, d.Args, d.Want, d.WantErr)
		})
	}
}

// Exec creates sends args to rpn, collects the results and compares
// them to want.
func UnitTestExec(t *testing.T, r *RPN, args, want []string, wantErr error) {
	t.Helper()
	err := r.ExecSlice(args)
	if !errors.Is(err, wantErr) {
		t.Fatalf("args=%v err=%v, want=%v", args, err, wantErr)
	}
	var got []string
	for _, f := range r.Frames {
		got = append(got, f.String(true))
	}
	if !reflect.DeepEqual(want, got) {
		t.Errorf("stack mismatch.  args=%v got=%+v, want=%+v", args, got, want)
	}
}
