//go:build pico || pico2

package fileops

func HomeDir() (string, error) {
	return "/", nil
}
