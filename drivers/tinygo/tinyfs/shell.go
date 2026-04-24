//go:build pico || pico2

package tinyfs

import (
	"errors"
	"fmt"
	"io"
	"mattwach/rpngo/common/elog"
	"os"
)

var (
	errArgsNotSupported       = errors.New("args not supported")
	errEmptyArgument          = errors.New("empty argument")
	errOnlyOnePathIsSupported = errors.New("only one path is supported")
	errPathIsRequired         = errors.New("path is required")
	errExpectedTwoPaths       = errors.New("expected two paths")
	errUnknownCommand         = errors.New("unknown command. cat, cp, ls, mv, pwd, and rm are supported")
)

// This file implements a very basic shell that implements a minimal
// interface of what could be considered "necessary" commands.
// It's done as shell emulation so that the PC and PicoCalc
// operations will use similar syntax (althogh the PC version is
// much mor robust with a real shell, of course).

func (fo *FileOpsDriver) Shell(args []string, stdin io.Reader) (string, error) {
	if err := fo.checkMount(); err != nil {
		return "", err
	}
	if len(args) == 0 {
		return "", nil
	}
	var val string
	var err error
	switch args[0] {
	case "cat":
		val, err = fo.cat(args[1:])
	case "cp":
		err = fo.cp(args[1:])
	case "ls":
		val, err = fo.ls(args[1:])
	case "mv":
		err = fo.mv(args[1:])
	case "pwd":
		val, err = fo.getpwd(args[1:])
	case "rm":
		err = fo.rm(args[1:])
	default:
		return "", errUnknownCommand
	}

	return val, err
}

func (fo *FileOpsDriver) ls(args []string) (string, error) {
	path := ""
	longMode := false
	for _, arg := range args {
		if arg == "" {
			return "", errEmptyArgument
		} else if arg[0] != '-' {
			if len(path) > 0 {
				return "", errOnlyOnePathIsSupported
			}
			path = arg
		} else if arg == "-l" {
			longMode = true
		} else {
			return "", fmt.Errorf("unknown flag: %v", arg)
		}
	}
	path = absPath(fo.pwd, path, true, false)
	dir, err := fo.fs.Open(path)
	if err != nil {
		return "", err
	}
	defer dir.Close()
	infos, err := dir.Readdir(0)
	if err != nil {
		return "", err
	}
	elog.Heap("alloc: /drivers/tinygo/tinyfs/shell.go:77: buff := make([]byte, 0, 128)")
	buff := make([]byte, 0, 128) // object allocated on the heap: escapes at line 77
	for _, info := range infos {
		if longMode {
			buff = appendLongInfo(buff, info)
		} else {
			buff = appendShortInfo(buff, info)
		}
	}
	return string(buff), nil
}

func (fo *FileOpsDriver) rm(args []string) error {
	if len(args) == 0 {
		return errPathIsRequired
	}
	for _, path := range args {
		if err := fo.rmPath(path); err != nil {
			return err
		}
	}
	return nil
}

func (fo *FileOpsDriver) cat(args []string) (string, error) {
	if len(args) != 1 {
		return "", errPathIsRequired
	}
	bytes, err := fo.ReadFile(args[0])
	return string(bytes), err
}

func (fo *FileOpsDriver) cp(args []string) error {
	if len(args) != 2 {
		return errExpectedTwoPaths
	}
	dst, err := fo.dstPathForCopy(args[1])
	if err != nil {
		return err
	}
	data, err := fo.ReadFile(args[0])
	if err != nil {
		return err
	}
	return fo.WriteFile(dst, data)
}

func (fo *FileOpsDriver) mv(args []string) error {
	if len(args) != 2 {
		return errExpectedTwoPaths
	}
	dst, err := fo.dstPathForCopy(args[1])
	if err != nil {
		return err
	}
	src := absPath(fo.pwd, args[0], true, false)
	return fo.fs.Rename(src, dst)
}

func (fo *FileOpsDriver) dstPathForCopy(p string) (string, error) {
	dst := absPath(fo.pwd, p, true, false)
	dstSt, err := fo.fs.Stat(dst)
	if err == nil {
		// the path exists
		if dstSt.IsDir() {
			dst = dst + "/" + basePath(p)
		}
	}
	return dst, nil
}

func (fo *FileOpsDriver) rmPath(path string) error { // object allocated on the heap: escapes at line 111
	if len(path) == 0 {
		return errEmptyArgument
	}
	if path[0] == '-' {
		return fmt.Errorf("unknown flag: %v", path)
	}
	path = absPath(fo.pwd, path, true, false)
	return fo.fs.Remove(path)
}

func appendLongInfo(buff []byte, info os.FileInfo) []byte {
	if info.IsDir() {
		buff = append(buff, byte('d'))
	} else {
		buff = append(buff, byte('-'))
	}
	elog.Heap("alloc: /drivers/tinygo/tinyfs/shell.go:95: return append(buff, []byte(fmt.Sprintf('rwxrwxrwx %5d %s\n', info.Size(), info.Name()))...)")
	return append(buff, []byte(fmt.Sprintf("rwxrwxrwx %5d %s\n", info.Size(), info.Name()))...) // object allocated on the heap: escapes at line 95
}

func appendShortInfo(buff []byte, info os.FileInfo) []byte {
	buff = append(buff, []byte(info.Name())...)
	if info.IsDir() {
		buff = append(buff, byte('/'))
	}
	return append(buff, byte('\n'))
}

func (fo *FileOpsDriver) getpwd(args []string) (string, error) {
	if len(args) != 0 {
		return "", errArgsNotSupported
	}
	return "/" + fo.pwd + "\n", nil
}
