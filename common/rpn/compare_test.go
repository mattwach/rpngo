package rpn

import (
	"fmt"
	"testing"
)

func TestIsLessThan(t *testing.T) {
	data := []struct {
		a    Frame
		b    Frame
		want bool
	}{
		// a = complex
		{
			a:    ComplexFrame(complex(0, 1)),
			b:    ComplexFrame(complex(0, -1)),
			want: false,
		},
		{
			a:    ComplexFrame(complex(0, 1)),
			b:    ComplexFrame(complex(1, -1)),
			want: true,
		},
		{
			a:    ComplexFrame(complex(1, 1)),
			b:    ComplexFrame(complex(0, -1)),
			want: false,
		},
		{
			a:    ComplexFrame(complex(0, 1)),
			b:    StringFrame("foo", STRING_SINGLEQ_FRAME),
			want: true,
		},
		{
			a:    ComplexFrame(complex(0, 1)),
			b:    BoolFrame(true),
			want: false,
		},
		{
			a:    ComplexFrame(complex(0, 1)),
			b:    BoolFrame(false),
			want: false,
		},

		// a = int
		{
			a:    IntFrame(0, INTEGER_FRAME),
			b:    ComplexFrame(complex(0, 1)),
			want: false,
		},
		{
			a:    IntFrame(0, INTEGER_FRAME),
			b:    ComplexFrame(complex(1, 1)),
			want: true,
		},
		{
			a:    IntFrame(1, OCTAL_FRAME),
			b:    ComplexFrame(complex(0, 1)),
			want: false,
		},
		{
			a:    IntFrame(0, INTEGER_FRAME),
			b:    IntFrame(0, INTEGER_FRAME),
			want: false,
		},
		{
			a:    IntFrame(0, INTEGER_FRAME),
			b:    IntFrame(1, HEXIDECIMAL_FRAME),
			want: true,
		},
		{
			a:    IntFrame(1, OCTAL_FRAME),
			b:    IntFrame(0, HEXIDECIMAL_FRAME),
			want: false,
		},
		{
			a:    IntFrame(0, INTEGER_FRAME),
			b:    StringFrame("foo", STRING_SINGLEQ_FRAME),
			want: true,
		},
		{
			a:    IntFrame(0, INTEGER_FRAME),
			b:    BoolFrame(true),
			want: false,
		},
		{
			a:    IntFrame(0, INTEGER_FRAME),
			b:    BoolFrame(false),
			want: false,
		},

		// a = string
		{
			a:    StringFrame("foo", STRING_SINGLEQ_FRAME),
			b:    ComplexFrame(complex(0, 1)),
			want: false,
		},
		{
			a:    StringFrame("foo", STRING_SINGLEQ_FRAME),
			b:    IntFrame(0, INTEGER_FRAME),
			want: false,
		},
		{
			a:    StringFrame("foo", STRING_SINGLEQ_FRAME),
			b:    StringFrame("foo", STRING_SINGLEQ_FRAME),
			want: false,
		},
		{
			a:    StringFrame("bar", STRING_SINGLEQ_FRAME),
			b:    StringFrame("foo", STRING_SINGLEQ_FRAME),
			want: true,
		},
		{
			a:    StringFrame("foo", STRING_SINGLEQ_FRAME),
			b:    StringFrame("bar", STRING_SINGLEQ_FRAME),
			want: false,
		},
		{
			a:    StringFrame("foo", STRING_SINGLEQ_FRAME),
			b:    BoolFrame(true),
			want: false,
		},
		{
			a:    StringFrame("foo", STRING_SINGLEQ_FRAME),
			b:    BoolFrame(false),
			want: false,
		},

		// a = bool
		{
			a:    BoolFrame(true),
			b:    ComplexFrame(complex(0, 1)),
			want: true,
		},
		{
			a:    BoolFrame(false),
			b:    ComplexFrame(complex(0, 1)),
			want: true,
		},
		{
			a:    BoolFrame(true),
			b:    IntFrame(0, INTEGER_FRAME),
			want: true,
		},
		{
			a:    BoolFrame(false),
			b:    IntFrame(0, INTEGER_FRAME),
			want: true,
		},
		{
			a:    BoolFrame(true),
			b:    StringFrame("foo", STRING_SINGLEQ_FRAME),
			want: true,
		},
		{
			a:    BoolFrame(false),
			b:    StringFrame("foo", STRING_SINGLEQ_FRAME),
			want: true,
		},
		{
			a:    BoolFrame(false),
			b:    BoolFrame(false),
			want: false,
		},
		{
			a:    BoolFrame(false),
			b:    BoolFrame(true),
			want: true,
		},
		{
			a:    BoolFrame(true),
			b:    BoolFrame(false),
			want: false,
		},
		{
			a:    BoolFrame(true),
			b:    BoolFrame(true),
			want: false,
		},
	}

	for _, d := range data {
		name := fmt.Sprintf("%s < %s", d.a.String(false), d.b.String(false))
		t.Run(name, func(t *testing.T) {
			got := d.a.IsLessThan(d.b)
			if got != d.want {
				t.Errorf("got: %v, want %v", got, d.want)
			}
		})
	}
}

func TestIsLessThanOrEqual(t *testing.T) {
	data := []struct {
		a    Frame
		b    Frame
		want bool
	}{
		// a = complex
		{
			a:    ComplexFrame(complex(0, 1)),
			b:    ComplexFrame(complex(0, -1)),
			want: false,
		},
		{
			a:    ComplexFrame(complex(-1, 1)),
			b:    ComplexFrame(complex(0, -1)),
			want: true,
		},
		{
			a:    ComplexFrame(complex(0, 1)),
			b:    ComplexFrame(complex(1, -1)),
			want: true,
		},
		{
			a:    ComplexFrame(complex(1, 1)),
			b:    ComplexFrame(complex(0, -1)),
			want: false,
		},
		{
			a:    ComplexFrame(complex(0, 1)),
			b:    StringFrame("foo", STRING_SINGLEQ_FRAME),
			want: true,
		},
		{
			a:    ComplexFrame(complex(0, 1)),
			b:    BoolFrame(true),
			want: false,
		},
		{
			a:    ComplexFrame(complex(0, 1)),
			b:    BoolFrame(false),
			want: false,
		},

		// a = int
		{
			a:    IntFrame(0, INTEGER_FRAME),
			b:    ComplexFrame(complex(0, 1)),
			want: false,
		},
		{
			a:    IntFrame(0, INTEGER_FRAME),
			b:    ComplexFrame(complex(1, 1)),
			want: true,
		},
		{
			a:    IntFrame(1, OCTAL_FRAME),
			b:    ComplexFrame(complex(0, 1)),
			want: false,
		},
		{
			a:    IntFrame(0, INTEGER_FRAME),
			b:    IntFrame(0, INTEGER_FRAME),
			want: true,
		},
		{
			a:    IntFrame(0, INTEGER_FRAME),
			b:    IntFrame(1, HEXIDECIMAL_FRAME),
			want: true,
		},
		{
			a:    IntFrame(1, OCTAL_FRAME),
			b:    IntFrame(0, HEXIDECIMAL_FRAME),
			want: false,
		},
		{
			a:    IntFrame(0, INTEGER_FRAME),
			b:    StringFrame("foo", STRING_SINGLEQ_FRAME),
			want: true,
		},
		{
			a:    IntFrame(0, INTEGER_FRAME),
			b:    BoolFrame(true),
			want: false,
		},
		{
			a:    IntFrame(0, INTEGER_FRAME),
			b:    BoolFrame(false),
			want: false,
		},

		// a = string
		{
			a:    StringFrame("foo", STRING_SINGLEQ_FRAME),
			b:    ComplexFrame(complex(0, 1)),
			want: false,
		},
		{
			a:    StringFrame("foo", STRING_SINGLEQ_FRAME),
			b:    IntFrame(0, INTEGER_FRAME),
			want: false,
		},
		{
			a:    StringFrame("foo", STRING_SINGLEQ_FRAME),
			b:    StringFrame("foo", STRING_SINGLEQ_FRAME),
			want: true,
		},
		{
			a:    StringFrame("bar", STRING_SINGLEQ_FRAME),
			b:    StringFrame("foo", STRING_SINGLEQ_FRAME),
			want: true,
		},
		{
			a:    StringFrame("foo", STRING_SINGLEQ_FRAME),
			b:    StringFrame("bar", STRING_SINGLEQ_FRAME),
			want: false,
		},
		{
			a:    StringFrame("foo", STRING_SINGLEQ_FRAME),
			b:    BoolFrame(true),
			want: false,
		},
		{
			a:    StringFrame("foo", STRING_SINGLEQ_FRAME),
			b:    BoolFrame(false),
			want: false,
		},

		// a = bool
		{
			a:    BoolFrame(true),
			b:    ComplexFrame(complex(0, 1)),
			want: true,
		},
		{
			a:    BoolFrame(false),
			b:    ComplexFrame(complex(0, 1)),
			want: true,
		},
		{
			a:    BoolFrame(true),
			b:    IntFrame(0, INTEGER_FRAME),
			want: true,
		},
		{
			a:    BoolFrame(false),
			b:    IntFrame(0, INTEGER_FRAME),
			want: true,
		},
		{
			a:    BoolFrame(true),
			b:    StringFrame("foo", STRING_SINGLEQ_FRAME),
			want: true,
		},
		{
			a:    BoolFrame(false),
			b:    StringFrame("foo", STRING_SINGLEQ_FRAME),
			want: true,
		},
		{
			a:    BoolFrame(false),
			b:    BoolFrame(false),
			want: true,
		},
		{
			a:    BoolFrame(false),
			b:    BoolFrame(true),
			want: true,
		},
		{
			a:    BoolFrame(true),
			b:    BoolFrame(false),
			want: false,
		},
		{
			a:    BoolFrame(true),
			b:    BoolFrame(true),
			want: true,
		},
	}

	for _, d := range data {
		name := fmt.Sprintf("%s <= %s", d.a.String(false), d.b.String(false))
		t.Run(name, func(t *testing.T) {
			got := d.a.IsLessThanOrEqual(d.b)
			if got != d.want {
				t.Errorf("got: %v, want %v", got, d.want)
			}
		})
	}
}

func TestIsEqual(t *testing.T) {
	data := []struct {
		a    Frame
		b    Frame
		want bool
	}{
		// a = complex
		{
			a:    ComplexFrame(complex(0, 1)),
			b:    ComplexFrame(complex(0, -1)),
			want: false,
		},
		{
			a:    ComplexFrame(complex(1, 1)),
			b:    ComplexFrame(complex(1, 1)),
			want: true,
		},
		{
			a:    ComplexFrame(complex(0, 1)),
			b:    IntFrame(1, INTEGER_FRAME),
			want: false,
		},
		{
			a:    ComplexFrame(complex(1, 0)),
			b:    IntFrame(1, INTEGER_FRAME),
			want: true,
		},
		{
			a:    ComplexFrame(complex(0, 1)),
			b:    StringFrame("foo", STRING_SINGLEQ_FRAME),
			want: false,
		},
		{
			a:    ComplexFrame(complex(0, 1)),
			b:    BoolFrame(true),
			want: false,
		},
		{
			a:    ComplexFrame(complex(0, 1)),
			b:    BoolFrame(false),
			want: false,
		},

		// a = int
		{
			a:    IntFrame(0, INTEGER_FRAME),
			b:    ComplexFrame(complex(0, 1)),
			want: false,
		},
		{
			a:    IntFrame(1, INTEGER_FRAME),
			b:    ComplexFrame(complex(1, 0)),
			want: true,
		},
		{
			a:    IntFrame(0, INTEGER_FRAME),
			b:    ComplexFrame(complex(1, 1)),
			want: false,
		},
		{
			a:    IntFrame(0, INTEGER_FRAME),
			b:    IntFrame(0, INTEGER_FRAME),
			want: true,
		},
		{
			a:    IntFrame(0, INTEGER_FRAME),
			b:    IntFrame(1, HEXIDECIMAL_FRAME),
			want: false,
		},
		{
			a:    IntFrame(1, INTEGER_FRAME),
			b:    IntFrame(1, HEXIDECIMAL_FRAME),
			want: true,
		},
		{
			a:    IntFrame(1, OCTAL_FRAME),
			b:    IntFrame(1, HEXIDECIMAL_FRAME),
			want: true,
		},
		{
			a:    IntFrame(1, OCTAL_FRAME),
			b:    IntFrame(1, BINARY_FRAME),
			want: true,
		},
		{
			a:    IntFrame(0, INTEGER_FRAME),
			b:    StringFrame("foo", STRING_SINGLEQ_FRAME),
			want: false,
		},
		{
			a:    IntFrame(0, INTEGER_FRAME),
			b:    BoolFrame(true),
			want: false,
		},
		{
			a:    IntFrame(0, INTEGER_FRAME),
			b:    BoolFrame(false),
			want: false,
		},

		// a = string
		{
			a:    StringFrame("foo", STRING_SINGLEQ_FRAME),
			b:    ComplexFrame(complex(0, 1)),
			want: false,
		},
		{
			a:    StringFrame("foo", STRING_SINGLEQ_FRAME),
			b:    IntFrame(0, INTEGER_FRAME),
			want: false,
		},
		{
			a:    StringFrame("foo", STRING_SINGLEQ_FRAME),
			b:    StringFrame("foo", STRING_SINGLEQ_FRAME),
			want: true,
		},
		{
			a:    StringFrame("bar", STRING_SINGLEQ_FRAME),
			b:    StringFrame("foo", STRING_SINGLEQ_FRAME),
			want: false,
		},
		{
			a:    StringFrame("foo", STRING_SINGLEQ_FRAME),
			b:    StringFrame("bar", STRING_SINGLEQ_FRAME),
			want: false,
		},
		{
			a:    StringFrame("foo", STRING_SINGLEQ_FRAME),
			b:    BoolFrame(true),
			want: false,
		},
		{
			a:    StringFrame("foo", STRING_SINGLEQ_FRAME),
			b:    BoolFrame(false),
			want: false,
		},

		// a = bool
		{
			a:    BoolFrame(true),
			b:    ComplexFrame(complex(0, 1)),
			want: false,
		},
		{
			a:    BoolFrame(true),
			b:    IntFrame(0, INTEGER_FRAME),
			want: false,
		},
		{
			a:    BoolFrame(true),
			b:    StringFrame("foo", STRING_SINGLEQ_FRAME),
			want: false,
		},
		{
			a:    BoolFrame(false),
			b:    BoolFrame(false),
			want: true,
		},
		{
			a:    BoolFrame(false),
			b:    BoolFrame(true),
			want: false,
		},
		{
			a:    BoolFrame(true),
			b:    BoolFrame(false),
			want: false,
		},
		{
			a:    BoolFrame(true),
			b:    BoolFrame(true),
			want: true,
		},
	}

	for _, d := range data {
		name := fmt.Sprintf("%s == %s", d.a.String(false), d.b.String(false))
		t.Run(name, func(t *testing.T) {
			got := d.a.IsEqual(d.b)
			if got != d.want {
				t.Errorf("got: %v, want %v", got, d.want)
			}
		})
	}
}
