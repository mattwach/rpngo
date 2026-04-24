module mattwach/rpngo/bin/tinygo/ili9341

go 1.25.1

replace mattwach/rpngo/common => ../../../common

replace mattwach/rpngo/drivers/tinygo => ../../../drivers/tinygo

require (
	mattwach/rpngo/common v0.0.0-00010101000000-000000000000
	mattwach/rpngo/drivers/tinygo v0.0.0-00010101000000-000000000000
)

require (
	github.com/google/shlex v0.0.0-20191202100458-e7afc7fbc510 // indirect
	tinygo.org/x/drivers v0.35.0 // indirect
	tinygo.org/x/tinyfont v0.7.0 // indirect
	tinygo.org/x/tinyfs v0.5.0 // indirect
)
