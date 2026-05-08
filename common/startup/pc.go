//go:build !pico && !pico2

package startup

import (
	"fmt"
	"log"
	"mattwach/rpngo/common/drivers/posix/fs"
	"mattwach/rpngo/common/drivers/posix/serial"
	"mattwach/rpngo/common/fileops"
	"mattwach/rpngo/common/functions"
	"mattwach/rpngo/common/parse"
	"mattwach/rpngo/common/rpn"
	"mattwach/rpngo/common/xmodem"
	"os"
	"os/signal"
	"path/filepath"
)

const defaultConfig = commonStartup + `
# Create and layout windows
'w.new.group' exists
{
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
  }
  .init=
  @.init
}
if

{
  time t1=
  0 x= {$x 1 + x= $x 3000000 <} for
  time $t1 - 3000000 1> /
} benchmark=

{$.plotwin w.new.plot} .plotinit=

{histl} {0/} try
{
  hists
  'i' 'autohist' true w.setp
} {0/} try

'/dev/ttyACM0' .serial=
`

const maxStackDepth = 65536

func Run(interactive func(*rpn.RPN) error) error {
	os.RemoveAll("/tmp/rpngo.log")
	logFile, err := os.Create("/tmp/rpngo.log")
	if err != nil {
		return err
	}
	defer logFile.Close()
	log.SetOutput(logFile)
	log.Println("Application started")
	var r rpn.RPN
	r.Init(maxStackDepth)
	functions.RegisterAll(&r)
	var fo fileops.FileOps
	fo.InitAndRegister(&r, 65536, &fs.FileOpsDriver{})
	var xm xmodem.XmodemCommands
	xm.InitAndRegister(&r, &serial.Serial{})

	if len(os.Args) > 1 {
		return cli(&r)
	}

	return interactive(&r)
}

func cli(r *rpn.RPN) error {
	tryLoadState(r)
	defer trySaveState(r)
	if err := r.ExecSlice(os.Args[1:]); err != nil {
		return err
	}
	if len(r.Frames) > 0 {
		fmt.Println(r.Frames[len(r.Frames)-1].String(true))
	}
	return nil
}

const stateName = ".rpn_cli_state"

func genStatePath() string {
	home, err := fileops.HomeDir()
	if err != nil {
		log.Printf("Can not generate state path: %v", err)
		return ""
	}
	return filepath.Join(home, stateName)
}

func tryLoadState(r *rpn.RPN) {
	path := genStatePath()
	if len(path) == 0 {
		return
	}
	buff, err := os.ReadFile(path)
	if err != nil {
		log.Printf("Failed to load %s: %v", path, err)
		return
	}
	err = parse.Fields(string(buff), r.Exec)
	if err != nil {
		log.Printf("Failed to parse %s: %v", path, err)
	} else {
		log.Printf("Read state snapshot %s", path)
	}
}

func trySaveState(r *rpn.RPN) {
	path := genStatePath()
	if len(path) == 0 {
		return
	}
	buff := r.VarSnapshot(make([]byte, 0, 256))
	buff = r.StackSnapshot(buff)
	err := os.WriteFile(path, buff, 0644)
	if err != nil {
		log.Printf("Failed to write %s: %v", path, err)
	} else {
		log.Printf("Wrote state snapshot to %s", path)
	}
}

type Interrupt struct {
	sigc chan os.Signal
}

func (i *Interrupt) Init() {
	i.sigc = make(chan os.Signal, 1)
	signal.Notify(i.sigc, os.Interrupt)
}

func (i *Interrupt) Interrupt() bool {
	select {
	case <-i.sigc:
		return true
	default:
		return false
	}
}
