package rpn

import (
	"errors"
	"testing"
)

func TestSetGetAndClear(t *testing.T) {
	data := []UnitTestExecData{
		{
			Args:    []string{"$x"},
			WantErr: ErrNotFound,
		},
		{
			Args: []string{"5", "x.y_1=", "$x.y_1", "$x.y_1", "6", "x.y_1=", "$x.y_1"},
			Want: []string{"5", "5", "6"},
		},
		{
			Args:    []string{"5", "=="},
			WantErr: ErrIllegalName,
			Want:    []string{"5"},
		},
		{
			Args:    []string{"5", "$="},
			WantErr: ErrIllegalName,
			Want:    []string{"5"},
		},
		{
			Args:    []string{"5", "1="},
			WantErr: ErrIllegalName,
			Want:    []string{"5"},
		},
		{
			Args:    []string{"5", "-="},
			WantErr: ErrIllegalName,
			Want:    []string{"5"},
		},
		{
			Args:    []string{"x/"},
			WantErr: ErrNotFound,
		},
		{
			Args: []string{"5", "x=", "x/"},
		},
		{
			Args:    []string{"5", "x=", "x/", "$x"},
			WantErr: ErrNotFound,
		},
		{
			Args: []string{"1", "2", "$0"},
			Want: []string{"1", "2", "2"},
		},
		{
			Args:    []string{"$0"},
			WantErr: ErrNotEnoughStackFrames,
		},
		{
			Args: []string{"1", "2", "$1"},
			Want: []string{"1", "2", "1"},
		},
		{
			Args:    []string{"1", "2", "$2"},
			WantErr: ErrNotEnoughStackFrames,
			Want:    []string{"1", "2"},
		},
		{
			Args:    []string{"1", "$2x"},
			WantErr: ErrIllegalName,
			Want:    []string{"1"},
		},
		{
			Args:    []string{"1", "$-2"},
			WantErr: ErrNotFound,
			Want:    []string{"1"},
		},
		{
			Args: []string{"1", "2", "0/"},
			Want: []string{"1"},
		},
		{
			Args: []string{"1", "2", "1/"},
			Want: []string{"2"},
		},
		{
			Args:    []string{"0/"},
			WantErr: ErrNotEnoughStackFrames,
		},
		{
			Args:    []string{"1", "1/"},
			Want:    []string{"1"},
			WantErr: ErrNotEnoughStackFrames,
		},
		{
			Args:    []string{"1", "1x/"},
			Want:    []string{"1"},
			WantErr: ErrIllegalName,
		},
		{
			Args: []string{"1", "2", "1>"},
			Want: []string{"2", "1"},
		},
		{
			Args: []string{"1", "2", "0>"},
			Want: []string{"1", "2"},
		},
		{
			Args:    []string{"1", "2", "2>"},
			Want:    []string{"1", "2"},
			WantErr: ErrNotEnoughStackFrames,
		},
		{
			Args: []string{"1", "2", "1<"},
			Want: []string{"2", "1"},
		},
		{
			Args: []string{"1", "2", "0<"},
			Want: []string{"1", "2"},
		},
		{
			Args:    []string{"1", "2", "2<"},
			Want:    []string{"1", "2"},
			WantErr: ErrNotEnoughStackFrames,
		},
	}
	UnitTestExecAll(t, data, nil)
}

func TestPushAndPop(t *testing.T) {
	data := []UnitTestExecData{
		{
			Args:    []string{"x<"},
			WantErr: ErrStackEmpty,
		},
		{
			Args:    []string{"x<<"},
			WantErr: ErrStackEmpty,
		},
		{
			Args:    []string{"x=="},
			WantErr: ErrStackEmpty,
		},
		{
			Args:    []string{"x>>"},
			WantErr: ErrNotFound,
		},
		{
			Args:    []string{"$$x"},
			WantErr: ErrNotFound,
		},
		{
			Args: []string{"5", "x<", "$x"},
			Want: []string{"5"},
		},
		{
			Args: []string{"5", "x<", "x>"},
			Want: []string{"5"},
		},
		{
			Args:    []string{"5", "x<", "x>", "x>"},
			Want:    []string{"5"},
			WantErr: ErrNotFound,
		},
		{
			Args: []string{"4", "5", "x<<", "$x", "$x"},
			Want: []string{"5", "5"},
		},
		{
			Args: []string{"4", "5", "x<<", "$$x", "$$x"},
			Want: []string{"4", "5", "4", "5"},
		},
		{
			Args: []string{"4", "5", "x<<", "x>>"},
			Want: []string{"4", "5"},
		},
		{
			Args:    []string{"4", "5", "x<<", "x>>", "x>>"},
			Want:    []string{"4", "5"},
			WantErr: ErrNotFound,
		},
		{
			Args: []string{"4", "5", "x<<", "6", "7", "x<<", "$$x"},
			Want: []string{"4", "5", "6", "7"},
		},
		{
			Args: []string{"4", "5", "x==", "6", "7", "x==", "$$x"},
			Want: []string{"6", "7"},
		},
	}
	UnitTestExecAll(t, data, nil)
}

func TestGetStringVariable(t *testing.T) {
	data := []struct {
		name    string
		args    []string
		vname   string
		want    string
		wantErr error
	}{
		{
			name:  "basic",
			args:  []string{"'foo'", "x="},
			vname: "x",
			want:  "foo",
		},
		{
			name:    "not found",
			vname:   "x",
			wantErr: ErrNotFound,
		},
		{
			name:  "not a string",
			args:  []string{"5", "x="},
			vname: "x",
			want:  "5",
		},
	}
	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			var r RPN
			r.Init(256)
			err := r.ExecSlice(d.args)
			if err != nil {
				t.Fatalf("err=%v, want nil", err)
			}
			got, err := r.GetStringVariable(d.vname)
			if !errors.Is(err, d.wantErr) {
				t.Errorf("err=%v, want %v", err, d.wantErr)
			}
			if got != d.want {
				t.Errorf("got=%v, want=%v", got, d.want)
			}
		})
	}
}

func TestGetComplexVariable(t *testing.T) {
	data := []struct {
		name    string
		args    []string
		vname   string
		want    complex128
		wantErr error
	}{
		{
			name:  "basic",
			args:  []string{"5+i", "x="},
			vname: "x",
			want:  complex(5, 1),
		},
		{
			name:    "not found",
			vname:   "x",
			wantErr: ErrNotFound,
		},
		{
			name:  "integer",
			args:  []string{"5d", "x="},
			vname: "x",
			want:  complex(5, 0),
		},
		{
			name:    "string",
			args:    []string{"'foo'", "x="},
			vname:   "x",
			wantErr: ErrExpectedANumber,
		},
	}
	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			var r RPN
			r.Init(256)
			err := r.ExecSlice(d.args)
			if err != nil {
				t.Fatalf("err=%v, want nil", err)
			}
			got, err := r.GetComplexVariable(d.vname)
			if !errors.Is(err, d.wantErr) {
				t.Fatalf("err=%v, want %v", err, d.wantErr)
			}
			if got != d.want {
				t.Errorf("got=%v, want=%v", got, d.want)
			}
		})
	}
}

func TestExecVariableAsMacro(t *testing.T) {
	data := []UnitTestExecData{
		{
			Args:    []string{"@x"},
			WantErr: ErrNotFound,
		},
		{
			Args: []string{"'1 2'", "x=", "@x"},
			Want: []string{"1", "2"},
		},
		{
			Args: []string{"'1 2'", "@0"},
			Want: []string{"'1 2'", "1", "2"},
		},
		{
			Args: []string{"'1 2'", "x=", "'@x @x'", "y=", "@y"},
			Want: []string{"1", "2", "1", "2"},
		},
	}
	UnitTestExecAll(t, data, nil)
}

func TestVarExists(t *testing.T) {
	data := []UnitTestExecData{
		{
			Args:    []string{"v.exists"},
			WantErr: ErrStackEmpty,
		},
		{
			Args:    []string{"5", "v.exists"},
			WantErr: ErrExpectedAString,
		},
		{
			Args: []string{"'x'", "v.exists"},
			Want: []string{"false"},
		},
		{
			Args: []string{"5", "x=", "'x'", "v.exists"},
			Want: []string{"true"},
		},
	}
	UnitTestExecAll(t, data, nil)
}
