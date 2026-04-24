//go:build !pico && !pico2

package fileops

import "os"

func HomeDir() (string, error) {
	return os.UserHomeDir()
}
