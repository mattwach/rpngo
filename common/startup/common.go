package startup

import (
	"fmt"
	"mattwach/rpngo/common/elog"
	"mattwach/rpngo/common/fileops"
	"mattwach/rpngo/common/parse"
	"mattwach/rpngo/common/rpn"
	"path/filepath"
)

const commonStartup = `
# set some useful vars
3.141592653589793 pi=
2.718281828459045 e=

# some useful equations
# (-b +/- sqrt(b*b - 4*a*c)) / (2 * a)
{$2 * 4 * $1 sq - neg sqrt 1> neg $0 $2 - $3 2 * / 3< + 1> 2 * /} quad=

'w.new.group' exists
  {
  {w.reset @.init} .f1=
  {w.reset} .f2=
  {w.reset @.init @.plotinit} .f3=
  
  {
    w.reset
    'root' w.columns
    'v' w.new.var 'v'
    'showdot'
    true w.setp
  } .f7=
  
  # snapshot save and load
  {'' cd snapshot '.rpngo_snapshot' save 'snapshot saved to .rpngo_snapshot' printlnx} .f4=
  {'' cd '.rpngo_snapshot' . 'snapshot loaded from .rpngo_snapshot' printlnx} .f9=
} 
if

'p' .plotwin=
rad
`

const HistFile = ".rpngo_history"
const configName = ".rpngo"

// Startup tries to load .rpngo and tries to create a default
// file if one can not be loaded.
func Startup(r *rpn.RPN, fs fileops.FileOpsDriver) error {
	configPath, err := genConfigPath()
	if err != nil {
		return err
	}
	s := loadOrCreateConfigFile(fs, configPath)
	err = parse.Fields(s, r.Exec)
	if err != nil {
		return fmt.Errorf(
			"while parsing %s: %w\n\nYou may delete or rename %s to have it automatically regenerated.",
			configPath,
			err,
			configPath)
	}
	return nil
}

func genConfigPath() (string, error) {
	home, err := fileops.HomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, configName), nil
}

func loadOrCreateConfigFile(fs fileops.FileOpsDriver, configPath string) string {
	var s []byte
	var err error
	if fs != nil {
		s, err = fs.ReadFile(configPath)
	}
	if err != nil {
		elog.Print("while loading config ", configPath, ": ", err.Error())
	}
	if fs == nil || (err != nil) {
		s = createConfigFile(fs, configPath)
	}
	return string(s)
}

func createConfigFile(fs fileops.FileOpsDriver, configPath string) []byte {
	s := []byte(defaultConfig)
	if fs != nil {
		if err := fs.WriteFile(configPath, s); err != nil {
			elog.Print("while saving config ", configPath, ": ", err.Error())
		}
	}
	return s
}
