//go:build pico || pico2

package tinyfs

import (
	"errors"
	"machine"
	"mattwach/rpngo/common/elog"
	"os"

	"tinygo.org/x/tinyfs"
)

var errNotADirectory = errors.New("not a directory")

// Setp for a PicoCalc or hardware with the same SD Card Connection
var spi = machine.SPI0

const sckPin = machine.GP18
const sdoPin = machine.GP19
const sdiPin = machine.GP16
const csPin = machine.GP17

type FileOpsDriver struct {
	initErr error
	mounted bool
	fs      tinyfs.Filesystem
	// present working directory.  '/' should only be used
	// between directories
	pwd string
}

func (fo *FileOpsDriver) Init() error {
	fo.mounted = false
	var bd tinyfs.BlockDevice
	bd, fo.initErr = blockDevice()
	if fo.initErr != nil {
		return fo.initErr
	}
	fo.initFS(bd)
	return fo.initErr
}

func (fo *FileOpsDriver) checkMount() error {
	if fo.initErr != nil {
		return fo.initErr
	}
	if fo.mounted {
		return nil
	}
	if err := fo.fs.Mount(); err != nil {
		return err
	}
	fo.mounted = true
	return nil
}

func (fo *FileOpsDriver) FileSize(path string) (int, error) {
	if err := fo.checkMount(); err != nil {
		return 0, err
	}
	s, err := fo.fs.Stat(absPath(fo.pwd, path, true, false))
	if err != nil {
		return 0, err
	}
	return int(s.Size()), nil
}

func (fo *FileOpsDriver) Format() error {
	if fo.initErr != nil {
		return fo.initErr
	}
	if fo.mounted {
		if err := fo.fs.Unmount(); err != nil {
			return err
		}
		fo.mounted = false
	}
	return fo.fs.Format()
}

func (fo *FileOpsDriver) ReadFile(path string) ([]byte, error) {
	if err := fo.checkMount(); err != nil {
		return nil, err
	}
	apath := absPath(fo.pwd, path, true, false)
	sz, err := fo.FileSize(apath)
	if err != nil {
		return nil, err
	}
	f, err := fo.fs.Open(apath)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	elog.Heap("alloc: /drivers/tinygo/tinyfs/fileops.go:71: data := make([]byte, sz)")
	data := make([]byte, sz) // object allocated on the heap: size is not constant
	totalRead := 0
	for totalRead < int(sz) {
		read, err := f.Read(data[totalRead:])
		if err != nil {
			return nil, err
		}
		totalRead += read
	}
	return data, nil
}

func (fo *FileOpsDriver) WriteFile(path string, data []byte) error {
	// Note, the tiyfs fat driver is very picky that all three of these are
	// set or it falls back to read mode.
	return fo.writeOrAppend(path, data, os.O_CREATE|os.O_WRONLY|os.O_TRUNC)
}

func (fo *FileOpsDriver) AppendToFile(path string, data []byte) error {
	// Note, the tiyfs fat driver is very picky that all three of these are
	// set or it falls back to read mode.
	return fo.writeOrAppend(path, data, os.O_WRONLY|os.O_CREATE|os.O_APPEND)
}

func (fo *FileOpsDriver) writeOrAppend(path string, data []byte, flags int) error {
	if err := fo.checkMount(); err != nil {
		return err
	}
	for len(path) > 0 && path[0] == '/' {
		path = path[1:]
	}
	f, err := fo.fs.OpenFile(absPath(fo.pwd, path, true, false), flags)
	if err != nil {
		return err
	}
	defer f.Close()
	totalWritten := 0
	for totalWritten < len(data) {
		written, err := f.Write(data[totalWritten:])
		if err != nil {
			return err
		}
		totalWritten += written
	}
	return nil
}

func (fo *FileOpsDriver) Chdir(path string) error {
	newPath := absPath(fo.pwd, path, false, false)
	s, err := fo.fs.Stat(newPath)
	if err != nil {
		return err
	}
	if !s.IsDir() {
		return errNotADirectory
	}
	fo.pwd = newPath
	return nil
}

func (fo *FileOpsDriver) ListFiles(path string, lst []string) ([]string, error) {
	path = absPath(fo.pwd, path, true, false)
	dir, err := fo.fs.Open(path)
	if err != nil {
		return lst, err
	}
	defer dir.Close()
	infos, err := dir.Readdir(0)
	if err != nil {
		return lst, err
	}
	for _, info := range infos {
		name := info.Name()
		if info.IsDir() {
			name += "/"
		}
		lst = append(lst, name)
	}
	return lst, nil
}
