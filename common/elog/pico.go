//go:build pico || pico2

package elog

func Print(v ...any) {
	print("Log: ")
	for _, a := range v {
		print(a)
		print(" ")
	}
	print("\n")
}

// These are primarily for embedded logging.  We don't care on PC hardware.
func Heap(msg string) {
	// comment this out if not actively debugging something
	//print("Heap: ", msg, "\n")
}
