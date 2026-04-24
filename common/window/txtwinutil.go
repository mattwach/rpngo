package window

func DrawStr(tw TextWindow, x, y int, msg string, col ColorChar) {
	w, h := tw.TextSize()
	for _, c := range msg {
		if (x >= 0) && (x < w) && (y >= 0) && (y < h) {
			tw.DrawChar(x, y, col|ColorChar(c))
		}
		x++
		if (x >= w) || c == '\n' {
			x = 0
			y++
		}
	}
}
