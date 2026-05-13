package rpn

import (
	"mattwach/rpngo/common/convert"
	"mattwach/rpngo/common/elog"
	"sort"
)

// RPN is the main structure
type RPN struct {
	Frames    []Frame
	variables map[string][]Frame
	functions map[string]func(*RPN) error
	// maps are category -> command -> help
	help          map[string]map[string]string
	Interrupt     func() bool
	Print         func(string)
	Input         func(*RPN) (string, error)
	TextWidth     int
	maxStackDepth int
	AngleUnit     FrameType
	conv          *convert.Conversion
	exiting       bool
}

// Init initializes an RPNCalc object
func (r *RPN) Init(maxStackDepth int) {
	r.Clear()
	elog.Heap("alloc: /rpn/rpn.go:26: r.Frames = make([]Frame, 0, maxStackDepth)")
	r.Frames = make([]Frame, 0, 16) // object allocated on the heap: object size 10240 exceeds maximum stack allocation size 256
	r.maxStackDepth = maxStackDepth
	r.functions = make(map[string]func(*RPN) error)
	elog.Heap("alloc: /rpn/rpn.go:28: r.variables = []map[string]Frame{make(map[string]Frame)}")
	r.variables = make(map[string][]Frame) // object allocated on the heap: escapes at line 28
	r.conv = convert.Init()                // must come before initHelp()
	r.initHelp()
	r.registerCore()
	r.Print = DefaultPrint
	r.Interrupt = DefaultInterrupt
	r.AngleUnit = POLAR_RAD_FRAME
	r.TextWidth = 80
}

func (r *RPN) registerCore() {
	r.Register("s.size", stackSize, CatStack, stackSizeHelp)
	r.Register("s.snapshot", stackSnapshot, CatStack, stackSnapshotHelp)
	r.Register("v.clear", varClear, CatVariables, varClearHelp)
	r.Register("v.clearall", varClearAll, CatVariables, varClearAllHelp)
	r.Register("v.exists", varExists, CatVariables, varExistsHelp)
	r.Register("v.list", listVariables, CatVariables, listVariablesHelp)
	r.Register("v.snapshot", varSnapshot, CatVariables, varSnapshotHelp)
	r.Register("deg", deg, CatEng, degHelp)
	r.Register("exists", exists, CatProg, existsHelp)
	r.Register("exit", exit, CatCore, exitHelp)
	r.Register("getangle", getAngle, CatEng, getAngleHelp)
	r.Register("grad", grad, CatEng, gradHelp)
	r.Register("rad", rad, CatEng, radHelp)
	r.Register("setangle", setAngle, CatEng, setAngleHelp)
}

// Register adds a new function
func (rpn *RPN) Register(name string, fn func(f *RPN) error, helpcat, helptxt string) {
	rpn.functions[name] = fn
	cat := rpn.help[helpcat]
	if cat == nil {
		rpn.help[helpcat] = map[string]string{name: helptxt}
	} else {
		cat[name] = helptxt
	}
}

func (rpn *RPN) AllFunctionNames() []string {
	var names []string
	for name := range rpn.functions {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

func DefaultPrint(msg string) {
	print(msg)
}

func DefaultInterrupt() bool {
	return false
}

func (r *RPN) Println(msg string) {
	r.Print(msg)
	r.Print("\n")
}

const existsHelp = "Returns true if the given string exists as a command. This " +
	"function can be useful when creating init script that will work in " +
	"different contexts, for example fyne vs ncurses.\n\n" +
	"example: 'w.new.group' exists"

func exists(r *RPN) error {
	f, err := r.PopFrame()
	if err != nil {
		return err
	}
	_, ok := r.functions[f.String(false)]
	r.PushFrame(BoolFrame(ok))
	return nil
}

const exitHelp = "exits RPN"

func exit(r *RPN) error {
	return ErrExit
}
