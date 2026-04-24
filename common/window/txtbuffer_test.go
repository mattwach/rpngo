package window

import "testing"

type FakeTextWindow struct {
	tw    int
	th    int
	chars []ColorChar
}

func (ftw *FakeTextWindow) SetCursorXY(x, y int) {}

func (ftw *FakeTextWindow) ResizeWindow(x, y, w, h int) error {
	return nil
}

func (ftw *FakeTextWindow) ShowBorder(sw, sh int) error {
	return nil
}

func (ftw *FakeTextWindow) Erase() {
	for i := range ftw.chars {
		ftw.chars[i] = ' '
	}
}

func (ftw *FakeTextWindow) WindowXY() (int, int) {
	return 0, 0
}

func (ftw *FakeTextWindow) WindowSize() (int, int) {
	return ftw.tw, ftw.th
}

func (ftw *FakeTextWindow) DrawChar(x, y int, char ColorChar) {
	ftw.chars[y*ftw.tw+x] = char
}

func (ftw *FakeTextWindow) TextWidth() int {
	return ftw.tw
}

func (ftw *FakeTextWindow) TextHeight() int {
	return ftw.th
}

func (ftw *FakeTextWindow) TextSize() (int, int) {
	return ftw.tw, ftw.th
}

func (ftw *FakeTextWindow) Refresh() {}

func (ftw *FakeTextWindow) Cursor(bool) {}

func makeChars(col ColorChar, s string) []ColorChar {
	cc := make([]ColorChar, len(s))
	for i, c := range s {
		cc[i] = col | ColorChar(c)
	}
	return cc
}

func checkChars(t *testing.T, name string, want []ColorChar, got []ColorChar) {
	t.Helper()
	if len(want) != len(got) {
		t.Fatalf("%s len mismatch want=%v, got=%v", name, len(want), len(got))
	}
	for i := range want {
		if want[i] != got[i] {
			t.Fatalf("%s want[%v]=%04x, got[%v]=%04x", name, i, want[i], i, got[i])
		}
	}
}

func TestInitResizeAndErase(t *testing.T) {
	tw := FakeTextWindow{
		tw:    3,
		th:    4,
		chars: makeChars(0, "xxxxxxxxxxxx"),
	}
	var tb TextBuffer
	// 8 bytes should create 2 blank lines
	tb.Init(&tw, 6)

	if tb.scrollbytes != 6 {
		t.Errorf("tb.scrollbytes=%v, want 8", tb.scrollbytes)
	}

	if tb.bw != 3 {
		t.Errorf("tb.bw=%v, want 3", tb.bw)
	}

	if tb.bh != 6 {
		t.Errorf("tb.bh=%v, want 6", tb.bw)
	}

	wantbuff := makeChars(0, "                  ")
	wantscreen := makeChars(0, "            ")

	checkChars(t, "tb.buffer", wantbuff, tb.buffer)
	checkChars(t, "tb.screen", wantscreen, tb.screen)
	checkChars(t, "tw.chars", wantscreen, tw.chars)
}

func TestWriteAndUpdateAndScroll(t *testing.T) {
	tw := FakeTextWindow{
		tw:    3,
		th:    4,
		chars: makeChars(0, "xxxxxxxxxxxx"),
	}
	var tb TextBuffer
	tb.Init(&tw, 6)
	tb.Scroll(1)
	tb.TextColor(White)

	// write x to the buffer, but not the screen
	//
	// buffer: screen
	// ...
	// x..     ...
	// ...     ...
	// ...     ...
	// ...     ...
	// ...
	tb.Write('x', false)
	wantbuff := makeChars(0, "                  ")
	wantbuff[3] = White | 'x'
	wantscreen := makeChars(0, "            ")

	checkChars(t, "1: tb.buffer", wantbuff, tb.buffer)
	checkChars(t, "1: tb.screen", wantscreen, tb.screen)
	checkChars(t, "1: tw.chars", wantscreen, tw.chars)

	// write y to the buffer and the screen
	//
	// buffer: screen
	// ...
	// xy.     .y.
	// ...     ...
	// ...     ...
	// ...     ...
	// ...
	tb.Write('y', true)
	wantbuff[4] = White | 'y'
	wantscreen[1] = White | 'y'

	checkChars(t, "2: tb.buffer", wantbuff, tb.buffer)
	checkChars(t, "2: tb.screen", wantscreen, tb.screen)
	checkChars(t, "2: tw.chars", wantscreen, tw.chars)

	// update to get the screen updated with the x
	//
	// buffer: screen
	// ...
	// xy.     xy.
	// ...     ...
	// ...     ...
	// ...     ...
	// ...
	tb.Update()
	wantscreen[0] = White | 'x'
	checkChars(t, "3: tb.buffer", wantbuff, tb.buffer)
	checkChars(t, "3: tb.screen", wantscreen, tb.screen)
	checkChars(t, "3: tw.chars", wantscreen, tw.chars)

	// next line, and another character
	//
	// buffer: screen
	// ...
	// xy.     xy.
	// z..     ...
	// ...     ...
	// ...     ...
	// ...
	tb.Write('\n', false)
	tb.Write('z', false)
	wantbuff[6] = White | 'z'
	checkChars(t, "4: tb.buffer", wantbuff, tb.buffer)
	checkChars(t, "4: tb.screen", wantscreen, tb.screen)
	checkChars(t, "4: tw.chars", wantscreen, tw.chars)

	// update the screen
	//
	// buffer: screen
	// ...
	// xy.     xy.
	// z..     z..
	// ...     ...
	// ...     ...
	// ...
	tb.Update()
	wantscreen[3] = White | 'z'
	checkChars(t, "5: tb.buffer", wantbuff, tb.buffer)
	checkChars(t, "5: tb.screen", wantscreen, tb.screen)
	checkChars(t, "5: tw.chars", wantscreen, tw.chars)

	// now we are going to change the x character to the edge and print
	// a couple of characters
	//
	// buffer: screen
	// ...
	// xy.     xy.
	// z.a     z..
	// b..     ...
	// ...     ...
	// ...
	tb.SetCursorX(2)
	tb.Write('a', false)
	tb.Write('b', false)
	wantbuff[8] = White | 'a'
	wantbuff[9] = White | 'b'
	checkChars(t, "6: tb.buffer", wantbuff, tb.buffer)
	checkChars(t, "6: tb.screen", wantscreen, tb.screen)
	checkChars(t, "6: tw.chars", wantscreen, tw.chars)

	// goto the lower corner of the screen and print two characters,
	// this should create a scroll event.  Becuase update is off, the
	// screen should not change yet.  The scroll event also changes
	// the foreground color of some of the space chars on the same row
	// as d.
	//
	// buffer: screen
	// ...
	// xy.
	// z.a     xy.
	// b..     z..
	// ..c     ...
	// d..     ...
	tb.SetCursorXY(2, 3)
	tb.Write('c', false)
	tb.Write('d', false)
	wantbuff[14] = White | 'c'
	wantbuff[15] = White | 'd'
	wantbuff[16] = White | ' '
	wantbuff[17] = White | ' '
	checkChars(t, "7: tb.buffer", wantbuff, tb.buffer)
	checkChars(t, "7: tb.screen", wantscreen, tb.screen)
	checkChars(t, "7: tw.chars", wantscreen, tw.chars)

	// goto the lower corner of the screen again and print two characters,
	// but with updates on.  This should scroll again and update the screen too
	//
	// buffer: screen
	// f..
	// xy.
	// z.a
	// b..     b..
	// ..c     ..c
	// d.e     d.e
	//         f..
	tb.SetCursorX(2)
	tb.Write('e', true)
	tb.Write('f', true)
	wantbuff[17] = White | 'e'
	wantbuff[0] = White | 'f'
	wantbuff[1] = White | ' '
	wantbuff[2] = White | ' '
	wantscreen = makeChars(White, "b    cd ef  ")
	wantscreen[1] = ' '
	wantscreen[2] = ' '
	wantscreen[3] = ' '
	wantscreen[4] = ' '
	checkChars(t, "8: tb.buffer", wantbuff, tb.buffer)
	checkChars(t, "8: tb.screen", wantscreen, tb.screen)
	checkChars(t, "8: tw.chars", wantscreen, tw.chars)

	// goto the lower corner of the screen again and print two characters,
	// with updates off.  This time, we are checking the the whole line
	// is cleared when autoscroll happens
	//
	// buffer: screen
	// f.g
	// h..
	// z.a
	// b..
	// ..c     b..
	// d.e     ..c
	//         d.e
	//         f..
	tb.SetCursorX(2)
	tb.Write('g', false)
	tb.Write('h', false)
	wantbuff[2] = White | 'g'
	wantbuff[3] = White | 'h'
	wantbuff[4] = White | ' '
	wantbuff[5] = White | ' '
	checkChars(t, "9: tb.buffer", wantbuff, tb.buffer)
	checkChars(t, "9: tb.screen", wantscreen, tb.screen)
	checkChars(t, "9: tw.chars", wantscreen, tw.chars)

	// go ahead and pass the update to the screen
	//
	// buffer: screen
	// f.g
	// h..
	// z.a
	// b..
	// ..c     ..c
	// d.e     d.e
	//         f.g
	//         h..
	tb.Update()
	wantscreen = makeChars(White, "  cd ef gh  ")
	wantscreen[0] = ' '
	wantscreen[1] = ' '
	checkChars(t, "10: tb.buffer", wantbuff, tb.buffer)
	checkChars(t, "10: tb.screen", wantscreen, tb.screen)
	checkChars(t, "10: tw.chars", wantscreen, tw.chars)

	// scroll up two lines and update the screen
	//
	// buffer: screen
	// f.g
	// h..
	// z.a     z.a
	// b..     b..
	// ..c     ..c
	// d.e     d.e
	//
	//
	tb.Scroll(-2)
	tb.Update()
	wantscreen = makeChars(White, "z ab    cd e")
	wantscreen[1] = ' '
	wantscreen[4] = ' '
	wantscreen[5] = ' '
	wantscreen[6] = ' '
	wantscreen[7] = ' '
	checkChars(t, "10: tb.buffer", wantbuff, tb.buffer)
	checkChars(t, "10: tb.screen", wantscreen, tb.screen)
	checkChars(t, "10: tw.chars", wantscreen, tw.chars)
}
