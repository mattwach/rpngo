//go:build !pico && !pico2

package startup

const defaultConfig = commonStartup + `
# Create and layout windows
{
  'g' w.new.group
  'g' w.columns
  'i' 'g' w.move.end
  'i' 25 w.weight
  'g2' w.new.group
  'g2' 25 w.weight
  'g2' w.columns
  'g2' .wtarget=
  's' w.new.stack
  'v' w.new.var
  'g' .wtarget=
} .init=

@.init

{
  time t1=
  0 x= {$x 1 + x= $x 3000000 <} for
  time $t1 - 3000000 1> /
} benchmark=

{$.plotwin w.new.plot} .plotinit=

{histl} {0/} try
hists
'i' 'autohist' true w.setp

'/dev/ttyACM0' .serial=
`
