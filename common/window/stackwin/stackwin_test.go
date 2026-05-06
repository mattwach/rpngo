package stackwin

import (
	"errors"
	"mattwach/rpngo/common/rpn"
	"testing"
)

func TestSetGetProp(t *testing.T) {
	data := []struct {
		name        string
		set         rpn.Frame
		wantDefault string
		wantGetErr  error
		wantSetErr  error
	}{
		{
			name:        "round",
			set:         rpn.IntFrame(2, rpn.INTEGER_FRAME),
			wantDefault: "-1d",
		},
		{
			name:       "round",
			set:        rpn.StringFrame("hello", rpn.STRING_SINGLEQ_FRAME),
			wantSetErr: rpn.ErrExpectedANumber,
		},
		{
			name: "round",
			set:  rpn.IntFrame(-1, rpn.INTEGER_FRAME),
		},
		{
			name:       "foo",
			wantGetErr: rpn.ErrUnknownProperty,
			wantSetErr: rpn.ErrUnknownProperty,
		},
	}

	for _, d := range data {
		t.Run(d.name+"."+d.set.String(false), func(t *testing.T) {
			var sw StackWindow
			sw.Init(nil)
			f, err := sw.GetProp(d.name)
			if !errors.Is(err, d.wantGetErr) {
				t.Fatalf("err=%v, wantGetErr=%v", err, d.wantGetErr)
			}
			if (d.wantDefault != "") && (f.String(true) != d.wantDefault) {
				t.Errorf("default get=%v, want=%v", f.String(true), d.wantDefault)
			}
			err = sw.SetProp(d.name, d.set)
			if !errors.Is(err, d.wantSetErr) {
				t.Fatalf("err=%v, wantSetErr=%v", err, d.wantSetErr)
			}
			f, err = sw.GetProp(d.name)
			if !errors.Is(err, d.wantGetErr) {
				t.Fatalf("err=%v, wantGetErr=%v", err, d.wantGetErr)
			}
			if (d.wantGetErr != nil) && (f.String(true) != d.set.String(true)) {
				t.Errorf("got: %v, want %v", f.String(true), d.set.String(true))
			}
			inAllProps := false
			for _, p := range sw.ListProps() {
				if p == d.name {
					inAllProps = true
					break
				}
			}
			if (d.wantGetErr != nil) && inAllProps {
				t.Errorf("wantGetErr=%v inAllProps=%v", d.wantGetErr, inAllProps)
			}
			if (d.wantGetErr == nil) && !inAllProps {
				t.Errorf("wantGetErr=%v inAllProps=%v", d.wantGetErr, inAllProps)
			}
		})
	}
}

func TestCountListProps(t *testing.T) {
	// did we add props and forget to change them?
	wantCount := 1
	var sw StackWindow
	sw.Init(nil)
	props := sw.ListProps()
	if len(props) != wantCount {
		t.Errorf("got props %+v, want a count of %v", props, wantCount)
	}
	for _, p := range props {
		_, err := sw.GetProp(p)
		if err != nil {
			t.Errorf("getprop err=%v, want nil", err)
		}
	}
}
