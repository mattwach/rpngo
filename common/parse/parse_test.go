package parse

import (
	"errors"
	"testing"
)

func TestFields(t *testing.T) {
	data := []struct {
		val     string
		want    []string
		wantErr error
	}{
		{
			val: "",
		},
		{
			val:  "1",
			want: []string{"1"},
		},
		{
			val:  "\n   a bc  def\n    ",
			want: []string{"a", "bc", "def"},
		},
		{
			val:  "   \"a bc  def\"    ",
			want: []string{"\"a bc  def\""},
		},
		{
			val:  "'a bc  def'",
			want: []string{"'a bc  def'"},
		},
		{
			val:  "'a\\nb'",
			want: []string{"'a\nb'"},
		},
		{
			val:  "{hello there}",
			want: []string{"{hello there}"},
		},
		{
			val:  "{hello {there}}",
			want: []string{"{hello {there}}"},
		},
		{
			val:     "'a",
			wantErr: ErrUnterminatedSingleQuote,
		},
		{
			val:     "a '",
			wantErr: ErrUnterminatedSingleQuote,
			want:    []string{"a"},
		},
		{
			val:     "'",
			wantErr: ErrUnterminatedSingleQuote,
		},
		{
			val:     "' # comment",
			wantErr: ErrUnterminatedSingleQuote,
		},
		{
			val:     "\"a",
			wantErr: ErrUnterminatedDouble,
		},
		{
			val:     "a \"",
			wantErr: ErrUnterminatedDouble,
			want:    []string{"a"},
		},
		{
			val:     "\"",
			wantErr: ErrUnterminatedDouble,
		},
		{
			val:     "\" # comment",
			wantErr: ErrUnterminatedDouble,
		},
		{
			val:     "{a",
			wantErr: ErrUnterminatedBrace,
		},
		{
			val:     "{{a}",
			wantErr: ErrUnterminatedBrace,
		},
		{
			val: "# comment",
		},
		{
			val: "\n# comment \"",
		},
		{
			val:  "'\"nested quote\"'",
			want: []string{"'\"nested quote\"'"},
		},
		{
			val:  "\"'nested quote'\"",
			want: []string{"\"'nested quote'\""},
		},
		{
			val:  "\\\"",
			want: []string{"\""},
		},
		{
			val:  "\\'",
			want: []string{"'"},
		},
		{
			val:  "\\\\",
			want: []string{"\\"},
		},
		{
			val: "\\",
		},
	}

	fields := make([]string, 16)

	for _, d := range data {
		t.Run(d.val, func(t *testing.T) {
			fields = fields[:0]
			err := Fields(d.val, func(t string) error {
				fields = append(fields, t)
				return nil
			})
			if !errors.Is(err, d.wantErr) {
				t.Errorf("err = %v, want %v", err, d.wantErr)
			}
			if len(fields) != len(d.want) {
				t.Fatalf("\n got: %+v\nwant: %+v", fields, d.want)
			}
			for i := range fields {
				if fields[i] != d.want[i] {
					t.Fatalf("\n got: %+v\nwant: %+v", fields, d.want)
				}
			}
		})
	}
}

func TestTruncateString(t *testing.T) {
	data := []struct {
		name       string
		str        string
		startIdx   int
		endIdx     int
		wantStarti int
		wantEndi   int
		wantStr    string
	}{
		{
			name:       "start",
			str:        "0123456789",
			startIdx:   0,
			endIdx:     4,
			wantStarti: 0,
			wantEndi:   4,
			wantStr:    "0123456",
		},
		{
			name:       "one",
			str:        "0123456789",
			startIdx:   1,
			endIdx:     5,
			wantStarti: 1,
			wantEndi:   5,
			wantStr:    "0123456",
		},
		{
			name:       "two",
			str:        "0123456789",
			startIdx:   2,
			endIdx:     6,
			wantStarti: 1,
			wantEndi:   5,
			wantStr:    "1234567",
		},
		{
			name:       "three",
			str:        "0123456789",
			startIdx:   3,
			endIdx:     7,
			wantStarti: 1,
			wantEndi:   5,
			wantStr:    "2345678",
		},
		{
			name:       "four",
			str:        "0123456789",
			startIdx:   4,
			endIdx:     8,
			wantStarti: 1,
			wantEndi:   5,
			wantStr:    "3456789",
		},
		{
			name:       "five",
			str:        "0123456789",
			startIdx:   5,
			endIdx:     9,
			wantStarti: 2,
			wantEndi:   6,
			wantStr:    "3456789",
		},
		{
			name:       "too big 1",
			str:        "0123456789",
			startIdx:   0,
			endIdx:     8,
			wantStarti: 0,
			wantEndi:   7,
			wantStr:    "0123456",
		},
		{
			name:       "too big 2",
			str:        "0123456789",
			startIdx:   1,
			endIdx:     9,
			wantStarti: 0,
			wantEndi:   7,
			wantStr:    "1234567",
		},
	}

	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			gots, gotsi, gotei := truncateString(
				d.str, d.startIdx, d.endIdx, 7)
			if gots != d.wantStr {
				t.Errorf("gots=%v, want=%v", gots, d.wantStr)
			}
			if gotsi != d.wantStarti {
				t.Errorf("got sidx=%v, want=%v", gotsi, d.wantStarti)
			}
			if gotei != d.wantEndi {
				t.Errorf("got eidx=%v, want=%v", gotei, d.wantEndi)
			}
		})
	}
}
