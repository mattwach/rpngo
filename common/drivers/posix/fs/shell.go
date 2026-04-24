package fs

import (
	"io"
	"os/exec"
)

func (f *FileOpsDriver) Shell(args []string, stdin io.Reader) (string, error) {
	cmd := exec.Command(args[0], args[1:]...)
	if stdin != nil {
		cmd.Stdin = stdin
	}
	val, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}
	return string(val), nil
}
