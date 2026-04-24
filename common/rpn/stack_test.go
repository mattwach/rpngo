package rpn

import (
	"reflect"
	"testing"
)

func TestLengthAndClear(t *testing.T) {
	var r RPN
	r.Init(256)
	r.ExecSlice([]string{"1", "2", "3"})
	if r.StackLen() != 3 {
		t.Errorf("StackLen()=%v, want 3", r.StackLen())
	}
	r.Clear()
	if r.StackLen() != 0 {
		t.Errorf("StackLen()=%v, want 0", r.StackLen())
	}
}

func TestPush(t *testing.T) {
	var r RPN
	r.Init(256)
	f := StringFrame("foo", STRING_SINGLEQ_FRAME)
	err := r.PushFrame(f)
	if err != nil {
		t.Fatalf("err=%v, want nil", err)
	}
	if len(r.Frames) != 1 {
		t.Fatalf("len(frames)=%v, want 1", len(r.Frames))
	}
	if !reflect.DeepEqual(r.Frames[0], f) {
		t.Errorf("frame mismatch. got=%+v, want=%+v", r.Frames[0], f)
	}
}

func TestPopFrame(t *testing.T) {
	var r RPN
	r.Init(256)
	_, err := r.PopFrame()
	if err != ErrStackEmpty {
		t.Errorf("err: %v, want: ErrStackEmpty", err)
	}
	f := StringFrame("foo", STRING_SINGLEQ_FRAME)
	r.PushFrame(f)
	got, err := r.PopFrame()
	if err != nil {
		t.Errorf("err: %v, want: nil", err)
	}
	if !reflect.DeepEqual(got, f) {
		t.Errorf("frame mismatch. got=%+v, want=%+v", got, f)
	}
	if len(r.Frames) != 0 {
		t.Errorf("stack size is %v, want 0", len(r.Frames))
	}
}

func TestPop2Frames(t *testing.T) {
	var r RPN
	r.Init(256)
	_, _, err := r.Pop2Frames()
	if err != ErrStackEmpty {
		t.Errorf("err: %v, want: ErrStackEmpty", err)
	}
	wanta := StringFrame("foo", STRING_SINGLEQ_FRAME)
	r.PushFrame(wanta)
	wantb := StringFrame("bar", STRING_SINGLEQ_FRAME)
	r.PushFrame(wantb)
	gota, gotb, err := r.Pop2Frames()
	if err != nil {
		t.Errorf("err: %v, want: nil", err)
	}
	if !reflect.DeepEqual(gota, wanta) {
		t.Errorf("frame mismatch. got=%+v, want=%+v", gota, wanta)
	}
	if !reflect.DeepEqual(gotb, wantb) {
		t.Errorf("frame mismatch. got=%+v, want=%+v", gotb, wantb)
	}
	if len(r.Frames) != 0 {
		t.Errorf("stack size is %v, want 0", len(r.Frames))
	}
}

func TestPeekFrame(t *testing.T) {
	var r RPN
	r.Init(256)
	r.PushFrame(StringFrame("foo", STRING_SINGLEQ_FRAME))
	r.PushFrame(StringFrame("bar", STRING_SINGLEQ_FRAME))
	_, err := r.PeekFrame(-1)
	if err != ErrIllegalValue {
		t.Errorf("want ErrIllegalValue, got %v", err)
	}
	_, err = r.PeekFrame(2)
	if err != ErrNotEnoughStackFrames {
		t.Errorf("want ErrNotEnoughStackFrames, got %v", err)
	}

	want := StringFrame("bar", STRING_SINGLEQ_FRAME)
	got, err := r.PeekFrame(0)
	if err != nil {
		t.Errorf("want err=nil, got %v", err)
	}
	if !reflect.DeepEqual(want, got) {
		t.Errorf("want %v, for %v", want, got)
	}

	want = StringFrame("foo", STRING_SINGLEQ_FRAME)
	got, err = r.PeekFrame(1)
	if err != nil {
		t.Errorf("want err=nil, got %v", err)
	}
	if !reflect.DeepEqual(want, got) {
		t.Errorf("want %v, for %v", want, got)
	}
}

func TestDeleteFrame(t *testing.T) {
	var r RPN
	r.Init(256)
	r.PushFrame(StringFrame("foo", STRING_SINGLEQ_FRAME))
	r.PushFrame(StringFrame("bar", STRING_SINGLEQ_FRAME))
	r.PushFrame(StringFrame("baz", STRING_SINGLEQ_FRAME))
	_, err := r.DeleteFrame(-1)
	if err != ErrIllegalValue {
		t.Errorf("want ErrIllegalValue, got %v", err)
	}
	_, err = r.DeleteFrame(3)
	if err != ErrNotEnoughStackFrames {
		t.Errorf("want ErrNotEnoughStackFrames, got %v", err)
	}

	got, err := r.DeleteFrame(0)
	if err != nil {
		t.Errorf("want err=nil, got %v", err)
	}
	want := StringFrame("baz", STRING_SINGLEQ_FRAME)
	if !reflect.DeepEqual(want, got) {
		t.Errorf("want %v, got %v", want, got)
	}

	wants := []Frame{
		StringFrame("foo", STRING_SINGLEQ_FRAME),
		StringFrame("bar", STRING_SINGLEQ_FRAME),
	}
	if !reflect.DeepEqual(wants, r.Frames) {
		t.Errorf("want %v, get %v", wants, r.Frames)
	}

	got, err = r.DeleteFrame(1)
	if err != nil {
		t.Errorf("want err=nil, got %v", err)
	}
	want = StringFrame("foo", STRING_SINGLEQ_FRAME)
	if !reflect.DeepEqual(want, got) {
		t.Errorf("want %v, got %v", want, got)
	}

	wants = []Frame{
		StringFrame("bar", STRING_SINGLEQ_FRAME),
	}
	if !reflect.DeepEqual(wants, r.Frames) {
		t.Errorf("want %v, get %v", wants, r.Frames)
	}
}

func TestInsertFrame(t *testing.T) {
	var r RPN
	r.Init(256)
	err := r.InsertFrame(Frame{}, -1)
	wantErr := ErrIllegalValue
	if err != wantErr {
		t.Errorf("err=%v, want %v", err, wantErr)
	}
	err = r.InsertFrame(Frame{}, 1)
	wantErr = ErrNotEnoughStackFrames
	if err != wantErr {
		t.Errorf("err=%v, want %v", err, wantErr)
	}
	r.PushFrame(StringFrame("1", STRING_SINGLEQ_FRAME))
	err = r.InsertFrame(StringFrame("foo", STRING_SINGLEQ_FRAME), 1)
	wantErr = nil
	if err != wantErr {
		t.Errorf("err=%v, want %v", err, wantErr)
	}

	wants := []Frame{StringFrame("foo", STRING_SINGLEQ_FRAME), StringFrame("1", STRING_SINGLEQ_FRAME)}
	if !reflect.DeepEqual(wants, r.Frames) {
		t.Errorf("want %v, get %v", wants, r.Frames)
	}
}

func TestStackSize(t *testing.T) {
	data := []UnitTestExecData{
		{
			Args: []string{"s.size"},
			Want: []string{"0d"},
		},
		{
			Args: []string{"1", "s.size"},
			Want: []string{"1", "1d"},
		},
		{
			Args: []string{"1", "2", "s.size"},
			Want: []string{"1", "2", "2d"},
		},
	}
	UnitTestExecAll(t, data, nil)
}
