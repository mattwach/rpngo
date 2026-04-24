package rpn

import (
	"mattwach/rpngo/common/elog"
	"mattwach/rpngo/common/parse"
	"sort"
	"strconv"
)

const listVariablesHelp = "Prints all variable names"

func listVariables(r *RPN) error {
	elog.Heap("alloc: rpn/vars.go:12: names := make([]string, 0, len(r.variables))")
	names := make([]string, 0, len(r.variables)) // object allocated on the heap: size is not constant
	names = r.AppendAllVariableNames(names)
	sort.Strings(names)
	for _, n := range names {
		r.Print(n)
		r.Print("\n")
	}
	return nil
}

// Sets a variable
func (r *RPN) SetVariable(name string) error {
	f, err := r.PopFrame()
	if err != nil {
		return err
	}
	if err := checkVariableName(name); err != nil {
		r.PushFrame(f)
		return err
	}
	vlist := r.variables[name]
	if len(vlist) > 0 {
		vlist[len(vlist)-1] = f
	} else {
		elog.Heap("alloc: rpn/vars.go:36: r.variables[name] = []Frame{f}")
		r.variables[name] = []Frame{f} // object allocated on the heap: escapes at line 23
	}
	return nil
}

func checkVariableName(name string) error {
	if len(name) == 0 {
		return ErrIllegalName
	}
	if !isAlpha(rune(name[0])) {
		return ErrIllegalName
	}
	for _, r := range name {
		if !isAlphaNum(r) {
			return ErrIllegalName
		}
	}
	return nil
}

func isAlpha(r rune) bool {
	return (r == '.') || (r == '_') || ((r >= 'A') && (r <= 'Z')) || ((r >= 'a') && (r <= 'z'))
}

func isNum(r rune) bool {
	return (r >= '0') && (r <= '9')
}

func isAlphaNum(r rune) bool {
	return (r == '.') || (r == '_') || ((r >= '0') && (r <= '9')) || ((r >= 'A') && (r <= 'Z')) || ((r >= 'a') && (r <= 'z'))
}

// Clears a variable
func (r *RPN) ClearVariable(name string) error {
	if len(name) == 0 {
		return ErrIllegalName
	}
	if isNum(rune(name[0])) {
		return r.clearStackVariable(name)
	}
	_, ok := r.variables[name]
	if !ok {
		return ErrNotFound
	}
	delete(r.variables, name)
	return nil
}

// Gets a variable
func (r *RPN) GetVariable(name string) (Frame, error) {
	if len(name) == 0 {
		return Frame{}, ErrIllegalName
	}
	if isNum(rune(name[0])) {
		return r.getStackVariable(name)
	}
	vlist := r.variables[name]
	if len(vlist) == 0 {
		return Frame{}, ErrNotFound
	}
	return vlist[len(vlist)-1], nil
}

// Gets a variable from the stack
func (r *RPN) getStackVariable(name string) (Frame, error) {
	idx, err := strconv.Atoi(name)
	if err != nil {
		return Frame{}, ErrIllegalName
	}
	return r.PeekFrame(idx)
}

// Removes a variable from the stack
func (r *RPN) clearStackVariable(name string) error {
	idx, err := strconv.Atoi(name)
	if err != nil {
		return ErrIllegalName
	}
	_, err = r.DeleteFrame(idx)
	return err
}

func (r *RPN) moveVariableToHead(name string) error {
	if len(name) == 0 {
		return ErrIllegalName
	}
	if isNum(rune(name[0])) {
		// It's a stack variable, e.g. 2>
		idx, err := strconv.Atoi(name)
		if err != nil {
			return ErrIllegalName
		}
		f, err := r.DeleteFrame(idx)
		if err != nil {
			return err
		}
		return r.PushFrame(f)
	}
	// It's a named variable, e.g. x>
	vlist := r.variables[name]
	if len(vlist) == 0 {
		return ErrNotFound
	}
	f := vlist[len(vlist)-1]
	if len(vlist) == 1 {
		delete(r.variables, name)
	} else {
		r.variables[name] = r.variables[name][:len(vlist)-1]
	}
	return r.PushFrame(f)
}

func (r *RPN) moveHeadToVariable(name string) error {
	if len(name) == 0 {
		return ErrIllegalName
	}
	if isNum(rune(name[0])) {
		// It's a stack variable, e.g. 2<
		idx, err := strconv.Atoi(name)
		if err != nil {
			return ErrIllegalName
		}
		f, err := r.PopFrame()
		if err != nil {
			return err
		}
		err = r.InsertFrame(f, idx)
		if err != nil {
			r.PushFrame(f)
		}
		return err
	}
	// It's a named variable, e.g. x<
	if err := checkVariableName(name); err != nil {
		return err
	}
	f, err := r.PopFrame()
	if err != nil {
		return err
	}
	r.variables[name] = append(r.variables[name], f)
	return nil
}

func (r *RPN) moveAllStackToVariable(name string) error {
	if err := checkVariableName(name); err != nil {
		return err
	}
	if len(r.Frames) == 0 {
		return ErrStackEmpty
	}
	elog.Heap("alloc: rpn/vars.go:187: r.variables[name] = make([]Frame, len(r.Frames))")
	r.variables[name] = make([]Frame, len(r.Frames)) // object allocated on the heap: size is not constant
	copy(r.variables[name], r.Frames)
	r.Frames = r.Frames[:0]
	return nil
}

func (r *RPN) appendAllVariableToStack(name string) error {
	if err := checkVariableName(name); err != nil {
		return err
	}
	if len(r.variables[name]) == 0 {
		return ErrNotFound
	}
	r.Frames = append(r.Frames, r.variables[name]...)
	return nil
}

func (r *RPN) appendAllStackToVariable(name string) error {
	if err := checkVariableName(name); err != nil {
		return err
	}
	if (len(r.Frames) == 0) && (len(r.variables[name]) == 0) {
		return ErrStackEmpty
	}
	r.variables[name] = append(r.variables[name], r.Frames...)
	r.Frames = r.Frames[:0]
	return nil
}

func (r *RPN) moveAllVariableToStack(name string) error {
	if err := checkVariableName(name); err != nil {
		return err
	}
	if len(r.variables[name]) == 0 {
		return ErrNotFound
	}
	r.Frames = append(r.Frames, r.variables[name]...)
	delete(r.variables, name)
	return nil
}

// gets a variable as a string
func (r *RPN) GetStringVariable(name string) (string, error) {
	v, err := r.GetVariable(name)
	if err != nil {
		return "", err
	}
	return v.String(false), nil
}

// gets a variable as a complex
func (r *RPN) GetComplexVariable(name string) (complex128, error) {
	f, err := r.GetVariable(name)
	if err != nil {
		return 0, err
	}
	v, err := f.Complex()
	if err != nil {
		return 0, err
	}
	return v, nil
}

// Gets all variable names
func (r *RPN) AppendAllVariableNames(names []string) []string {
	for name := range r.variables {
		names = append(names, name)
	}
	return names
}

// Calls the callback for every defined variable.  The function can return
// false to abort the process
func (r *RPN) IterateAllVariables(fn func(string, []Frame) bool) {
	for k, v := range r.variables {
		if !fn(k, v) {
			break
		}
	}
}

// Gets all variable names
func (r *RPN) AllVariableNamesAndValues() map[string][]Frame {
	return r.variables
}

// Executes a Variables as a macro
func (r *RPN) execVariableAsMacro(name string) error {
	f, err := r.GetVariable(name)
	if err != nil {
		return err
	}
	if !f.IsString() {
		return ErrExpectedAString
	}
	// this call can be recursive so we need to allocate here
	if err := parse.Fields(f.UnsafeString(), r.Exec); err != nil {
		return err
	}
	return nil
}

const varSnapshotHelp = "Creates a string that defines current variables"

func varSnapshot(r *RPN) error {
	if len(r.variables) == 0 {
		return r.PushFrame(StringFrame("", STRING_BRACE_FRAME))
	}
	buff := make([]byte, 0, 5*len(r.Frames))
	buff = r.VarSnapshot(buff)
	return r.PushFrame(StringFrame(string(buff), STRING_BRACE_FRAME))
}

func (r *RPN) VarSnapshot(buff []byte) []byte {
	for name, vals := range r.variables {
		buff = r.snapshotOneVar(buff, name, vals)
	}
	return buff
}

func (r *RPN) snapshotOneVar(buff []byte, name string, vals []Frame) []byte {
	for i, v := range vals {
		buff = append(buff, []byte(v.String(true))...)
		buff = append(buff, ' ')
		buff = append(buff, []byte(name)...)
		if i == 0 {
			buff = append(buff, '=')
		} else {
			buff = append(buff, '<')
		}
		buff = append(buff, '\n')
	}
	return buff
}

const varClearHelp = "Clears all variables that do not start with a ."

func varClear(r *RPN) error {
	var delnames []string
	for name := range r.variables {
		if (len(name) > 0) && (name[0] != '.') {
			delnames = append(delnames, name)
		}
	}
	for _, name := range delnames {
		delete(r.variables, name)
	}
	return nil
}

const varClearAllHelp = "Clears all variables (including . ones)"

func varClearAll(r *RPN) error {
	r.variables = make(map[string][]Frame)
	return nil
}

const varExistsHelp = "Returns true if the variable name is defined."

func varExists(r *RPN) error {
	f, err := r.PopFrame()
	if err != nil {
		return err
	}
	if !f.IsString() {
		return ErrExpectedAString
	}
	_, ok := r.variables[f.UnsafeString()]
	return r.PushFrame(BoolFrame(ok))
}
