package fileops

import (
	"errors"
	"mattwach/rpngo/common/drivers/posix/fs"
	"mattwach/rpngo/common/functions"
	"mattwach/rpngo/common/rpn"
	"os"
	"path/filepath"
	"testing"
)

func TestSource(t *testing.T) {
	data := []struct {
		filedata string
		want     string
		wantErr  error
	}{
		{
			filedata: "2 3 + str",
			want:     "5",
		},
		{
			filedata: "'hi' foo",
			wantErr:  rpn.ErrSyntax,
			want:     "hi",
		},
		{
			wantErr: os.ErrNotExist,
		},
	}

	for _, d := range data {
		t.Run(d.filedata, func(t *testing.T) {
			dir := t.TempDir()
			defer os.RemoveAll(dir)
			filePath := filepath.Join(dir, "testfile.txt")

			if len(d.filedata) > 0 {
				err := os.WriteFile(filePath, []byte(d.filedata), 0644)
				if err != nil {
					t.Fatalf("error creating temp file: %v", err)
				}
			}
			var r rpn.RPN
			r.Init(256)
			functions.RegisterAll(&r)
			var fo FileOps
			fo.InitAndRegister(&r, 65536, &fs.FileOpsDriver{})
			err := r.ExecSlice([]string{"'" + filePath + "'", "source"})
			if !errors.Is(err, d.wantErr) {
				t.Errorf("err=%v, want: %v", err, d.wantErr)
			}
			if len(d.filedata) > 0 {
				f, err := r.PopFrame()
				if err != nil {
					t.Errorf("err=%v, want nil", err)
				}
				got := f.String(false)
				if got != d.want {
					t.Errorf("got: %v, want %v", got, d.want)
				}
			}
		})
	}
}
