package fs

import (
	"mattwach/rpngo/common/rpn"
	"os"
)

type FileOpsDriver struct{}

func (fo *FileOpsDriver) FileSize(path string) (int, error) {
	s, err := os.Stat(path)
	if err != nil {
		return 0, err
	}
	return int(s.Size()), nil
}

func (fo *FileOpsDriver) ReadFile(path string) ([]byte, error) {
	return os.ReadFile(path)
}

func (fo *FileOpsDriver) WriteFile(path string, data []byte) error {
	return os.WriteFile(path, data, 0644)
}

func (fo *FileOpsDriver) AppendToFile(path string, data []byte) error {
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.Write(data)
	return err
}

func (fo *FileOpsDriver) Chdir(path string) error {
	return os.Chdir(path)
}

func (fo *FileOpsDriver) Format() error {
	return rpn.ErrNotSupported
}

func (fo *FileOpsDriver) ListFiles(dir string, lst []string) ([]string, error) {
	// Read the directory contents
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		name := entry.Name()
		if entry.IsDir() {
			name += "/"
		}
		lst = append(lst, name)
	}

	return lst, nil
}
