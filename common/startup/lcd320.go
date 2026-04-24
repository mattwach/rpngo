//go:build pico || pico2

package startup

const defaultConfig = commonStartup + `
{
  30 .wweight<
  's' w.new.stack
  .wweight> 0/
} .init=

@.init

{hists 'history saved' printlnx} .f5=

{reset} .f6=

{
  time t1=
  0 x=
  {$x 1 + x= $x 50000 <} for
  time $t1 - 50000 1> /
  ' loops/second' + printlnx
  heapstats
} benchmark=

{
  w.reset
  false .wend<
  250 .wweight<
  $.plotwin w.new.plot
  .wend> 0/ .wweight> 0/
} .plotinit=

# history load/save doesn't work on tinygo unless the media is formatted
{histl} {0/} try

true .echo=
`
