package stackwin

import (
	"errors"
	"fmt"
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

const pi = 3.141592653589793

func TestRoundedString(t *testing.T) {
	data := []struct {
		round int8
		val   rpn.Frame
		want  string
	}{
		{
			round: 2,
			val:   rpn.StringFrame("hello", rpn.STRING_DOUBLEQ_FRAME),
			want:  "\"hello\"",
		},
		{
			round: 2,
			val:   rpn.IntFrame(2, rpn.INTEGER_FRAME),
			want:  "2d",
		},
		{
			round: 0,
			val:   rpn.RealFrame(pi),
			want:  "3",
		},
		{
			round: 1,
			val:   rpn.RealFrame(pi),
			want:  "3.1",
		},
		{
			round: 2,
			val:   rpn.RealFrame(pi),
			want:  "3.14",
		},
		{
			round: 3,
			val:   rpn.RealFrame(pi),
			want:  "3.142",
		},
		{
			round: 0,
			val:   rpn.RealFrame(-pi),
			want:  "-3",
		},
		{
			round: 1,
			val:   rpn.RealFrame(-pi),
			want:  "-3.1",
		},
		{
			round: 2,
			val:   rpn.RealFrame(-pi),
			want:  "-3.14",
		},
		{
			round: 3,
			val:   rpn.RealFrame(-pi),
			want:  "-3.142",
		},
		{
			round: 0,
			val:   rpn.ComplexFrame(complex(0, pi)),
			want:  "3i",
		},
		{
			round: 1,
			val:   rpn.ComplexFrame(complex(0, pi)),
			want:  "3.1i",
		},
		{
			round: 2,
			val:   rpn.ComplexFrame(complex(0, pi)),
			want:  "3.14i",
		},
		{
			round: 3,
			val:   rpn.ComplexFrame(complex(0, pi)),
			want:  "3.142i",
		},
		{
			round: 0,
			val:   rpn.ComplexFrame(complex(0, -pi)),
			want:  "-3i",
		},
		{
			round: 1,
			val:   rpn.ComplexFrame(complex(0, -pi)),
			want:  "-3.1i",
		},
		{
			round: 2,
			val:   rpn.ComplexFrame(complex(0, -pi)),
			want:  "-3.14i",
		},
		{
			round: 3,
			val:   rpn.ComplexFrame(complex(0, -pi)),
			want:  "-3.142i",
		},
		{
			round: 0,
			val:   rpn.ComplexFrame(complex(pi, pi)),
			want:  "3+3i",
		},
		{
			round: 1,
			val:   rpn.ComplexFrame(complex(pi, pi)),
			want:  "3.1+3.1i",
		},
		{
			round: 2,
			val:   rpn.ComplexFrame(complex(pi, pi)),
			want:  "3.14+3.14i",
		},
		{
			round: 3,
			val:   rpn.ComplexFrame(complex(pi, pi)),
			want:  "3.142+3.142i",
		},
		{
			round: 0,
			val:   rpn.ComplexFrame(complex(-pi, -pi)),
			want:  "-3-3i",
		},
		{
			round: 1,
			val:   rpn.ComplexFrame(complex(-pi, -pi)),
			want:  "-3.1-3.1i",
		},
		{
			round: 2,
			val:   rpn.ComplexFrame(complex(-pi, -pi)),
			want:  "-3.14-3.14i",
		},
		{
			round: 3,
			val:   rpn.ComplexFrame(complex(-pi, -pi)),
			want:  "-3.142-3.142i",
		},
	}

	for _, d := range data {
		t.Run(fmt.Sprintf("%v:%v", d.val.String(false), d.round), func(t *testing.T) {
			var sw StackWindow
			sw.Init(nil)
			sw.round = d.round
			got := sw.roundedString(d.val)
			if got != d.want {
				t.Errorf("got %v, want %v", got, d.want)
			}
		})
	}
}
