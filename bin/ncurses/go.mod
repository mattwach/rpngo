module mattwach/rpngo/bin/ncurses

go 1.25.1

replace mattwach/rpngo/common => ../../common

replace mattwach/rpngo/drivers/curses => ../../drivers/curses

require (
	mattwach/rpngo/common v0.0.0-00010101000000-000000000000
	mattwach/rpngo/drivers/curses v0.0.0-00010101000000-000000000000
)

require github.com/gbin/goncurses v0.0.0-20251113135420-86371713952c // indirect
