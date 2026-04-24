package fileops

import (
	"errors"
	"fmt"
	"io"
	"mattwach/rpngo/common/elog"
	"mattwach/rpngo/common/parse"
	"mattwach/rpngo/common/rpn"
	"os/exec"
	"strings"
)

type FileOpsDriver interface {
	FileSize(path string) (int, error)
	Format() error
	ReadFile(path string) ([]byte, error)
	WriteFile(path string, data []byte) error
	AppendToFile(path string, data []byte) error
	ListFiles(path string, lst []string) ([]string, error)
	Chdir(path string) error
	Shell(args []string, stdin io.Reader) (string, error)
}

type FileOps struct {
	maxFileSize int
	driver      FileOpsDriver
}

func (fo *FileOps) InitAndRegister(r *rpn.RPN, maxFileSize int, driver FileOpsDriver) {
	fo.maxFileSize = maxFileSize
	fo.driver = driver
	fo.register(r)
}

const formatHelp = "Formats an SD card. For safety, must provide a single \"YES\" string argument. Not suported on PCs"

func (fo *FileOps) format(r *rpn.RPN) error {
	f, err := r.PopFrame()
	if err != nil {
		return err
	}
	if !f.IsString() || (f.UnsafeString() != "YES") {
		return errors.New("provide the argument 'YES' to confirm SD Card reformat")
	}
	return fo.driver.Format()
}

const loadHelp = "Loads the given filename and places it on the stack as a string variable"

func (fo *FileOps) load(r *rpn.RPN) error {
	f, err := r.PopFrame()
	if err != nil {
		return err
	}
	if !f.IsString() {
		return rpn.ErrExpectedAString
	}
	path := f.UnsafeString()
	sz, err := fo.driver.FileSize(path)
	if err != nil {
		return err
	}
	if sz > fo.maxFileSize {
		return fmt.Errorf("file is too large.  %v > %v max bytes", sz, fo.maxFileSize)
	}
	data, err := fo.driver.ReadFile(path)
	if err != nil {
		return err
	}
	return r.PushFrame(rpn.StringFrame(string(data), rpn.STRING_SINGLEQ_FRAME))
}

const saveHelp = "Uses $0 as the filename to save the contents of $1  Both are popped from the stack.  A '\\n' character is added.  Use append if this is not wanted."

func (fo *FileOps) save(r *rpn.RPN) error {
	f, err := r.PopFrame()
	if err != nil {
		return err
	}
	if !f.IsString() {
		return rpn.ErrExpectedAString
	}
	path := f.UnsafeString()
	data, err := r.PopFrame()
	if err != nil {
		return err
	}
	return fo.driver.WriteFile(path, []byte(data.String(false)))
}

const appendHelp = "Uses $0 as the filename to append the contents of $1  Both are popped from the stack.  No '\\n' character is added."

func (fo *FileOps) append(r *rpn.RPN) error {
	f, err := r.PopFrame()
	if err != nil {
		return err
	}
	if !f.IsString() {
		return rpn.ErrExpectedAString
	}
	path := f.UnsafeString()
	data, err := r.PopFrame()
	if err != nil {
		return err
	}
	return fo.driver.AppendToFile(path, []byte(data.String(false)))
}

const cdHelp = "Change the working directory"

func (fo *FileOps) cd(r *rpn.RPN) error {
	f, err := r.PopFrame()
	if err != nil {
		return err
	}
	if !f.IsString() {
		return rpn.ErrExpectedAString
	}
	path := f.UnsafeString()
	if err != nil {
		return err
	}
	if path == "" {
		path, _ = HomeDir()
	}
	return fo.driver.Chdir(path)
}

const shellHelp = `Executes a string as a shell command.

There are many ways to execute a shell command and the following special
variables control the behavior:

.stdin  - If set, the contents will be sent to stdin of the process
.stdout - If empty or false, stdout/stderr is simply printed.  
          If set to true, stdout/stderr is pushed to the stack
.env    - If set, environment variables will be set using KEY=VALUE with
          one variable per line

The exit code of the shell command is set to the variable $rc.
`

func (fo *FileOps) shell(r *rpn.RPN) error {
	f, err := r.PopFrame()
	if err != nil {
		return err
	}
	s := f.String(false)
	elog.Heap("alloc: /fileops/fileops.go:132: fields := make([]string, 0, 4)")
	fields := make([]string, 0, 4) // object allocated on the heap: escapes at line 132
	addField := func(t string) error {
		fields = append(fields, t)
		return nil
	}
	if err := parse.Fields(s, addField); err != nil {
		return err
	}
	if len(fields) == 0 {
		return rpn.ErrIllegalValue
	}

	stdin, err := checkStdinVar(r)
	if err != nil {
		return err
	}

	output, err := fo.driver.Shell(fields, stdin)

	rc := 0
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			rc = exitError.ExitCode()
		} else {
			rc = 1
		}
		r.Print("Error: " + err.Error() + " " + string(output) + "\n")
	}
	r.PushFrame(rpn.IntFrame(int64(rc), rpn.INTEGER_FRAME))
	r.SetVariable("rc")
	if err == nil {
		if err := setCmdOutput(r, string(output)); err != nil {
			return err
		}
	}
	return nil
}

func checkStdinVar(r *rpn.RPN) (io.Reader, error) {
	val, err := r.GetStringVariable(".stdin")
	if err != nil {
		return nil, nil
	}
	return strings.NewReader(val + "\n"), nil
}

func setCmdOutput(r *rpn.RPN, output string) error {
	stack := false
	stdout, err := r.GetVariable(".stdout")
	if err == nil {
		stack, err = stdout.Bool()
		if err != nil {
			return err
		}
	}
	if stack {
		return r.PushFrame(rpn.StringFrame(strings.TrimSpace(output), rpn.STRING_SINGLEQ_FRAME))
	}
	r.Print(output)
	return nil
}
