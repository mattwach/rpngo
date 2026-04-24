package window

import (
	"mattwach/rpngo/common/elog"
	"mattwach/rpngo/common/rpn"
)

type WindowRoot struct {
	group        windowGroup
	adjustNeeded bool
}

func (wr *WindowRoot) Init(w, h int) {
	wr.group.resize(0, 0, w, h)
}

func (wr *WindowRoot) FindWindow(name string) WindowWithProps {
	wge := wr.group.findwindowGroupEntry(name)
	if wge == nil {
		return nil
	}
	return wge.window
}

func (wr *WindowRoot) FindwindowGroup(name string) (*windowGroup, error) {
	if name == "root" {
		return &wr.group, nil
	}
	wge := wr.group.findwindowGroupEntry(name)
	if wge == nil {
		return nil, rpn.ErrNotFound
	}
	if wge.group == nil {
		return nil, rpn.ErrNotAWindowGroup
	}
	return wge.group, nil
}

func (wr *WindowRoot) DeleteWindowOrGroup(name string) error {
	if name == "root" {
		return rpn.ErrCanNotDeleteRootWindow
	}
	if name == "i" {
		return rpn.ErrCanNotDeleteInputWindow
	}
	pwge, wge := wr.group.findwindowGroupEntryAndParent(name)
	if wge == nil {
		return rpn.ErrNotFound
	}
	pwge.removeChild(wge)
	wr.adjustNeeded = true
	return nil
}

func (wr *WindowRoot) MoveWindowOrGroup(src string, dst string, beginning bool) error {
	if src == "root" {
		return rpn.ErrIllegalWindowOperation
	}
	srcpg, srcwge := wr.group.findwindowGroupEntryAndParent(src)
	if srcwge == nil {
		return rpn.ErrNotFound
	}
	if srcwge.group != nil {
		check := srcwge.group.findwindowGroupEntry(dst)
		if check != nil {
			return rpn.ErrIllegalWindowOperation
		}
	}
	dstpg, err := wr.FindwindowGroup(dst)
	if err != nil {
		return err
	}
	srcpg.removeChild(srcwge)
	if beginning {
		elog.Heap("alloc: /window/windowroot.go:74: dstpg.children = append([]*windowGroupEntry{srcwge}, dstpg.children...)")
		dstpg.children = append([]*windowGroupEntry{srcwge}, dstpg.children...) // object allocated on the heap: escapes at line 74
	} else {
		dstpg.children = append(dstpg.children, srcwge)
	}
	wr.adjustNeeded = true
	return nil
}

func (wr *WindowRoot) SetWindowWeight(name string, w int) error {
	if w < 1 {
		return rpn.ErrIllegalWindowOperation
	}
	if w > 10000 {
		return rpn.ErrIllegalWindowOperation
	}
	wge := wr.group.findwindowGroupEntry(name)
	if wge == nil {
		return rpn.ErrNotFound
	}
	wge.weight = w
	wr.adjustNeeded = true
	return nil
}

func (wr *WindowRoot) AddNewWindowGroupChild(r *rpn.RPN, name string) {
	elog.Heap("alloc: /window/windowroot.go:99: wr.addWindowGroupEntry(r, &windowGroupEntry{name: name, group: &windowGroup{}})")
	wr.addWindowGroupEntry(r, &windowGroupEntry{name: name, group: &windowGroup{}}) // object allocated on the heap: escapes at line 99
}

func (wr *WindowRoot) AddWindowChildToRoot(window WindowWithProps, name string, weight int) {
	elog.Heap("alloc: /window/windowroot.go:103: wr.group.children = append(wr.group.children, &windowGroupEntry{name: name, weight: weight, window: window})")
	wr.group.children = append(wr.group.children, &windowGroupEntry{name: name, weight: weight, window: window}) // object allocated on the heap: escapes at line 103
	wr.adjustNeeded = true
}

func (wr *WindowRoot) AddWindowChild(r *rpn.RPN, window WindowWithProps, name string) {
	elog.Heap("alloc: /window/windowroot.go:108: wr.addWindowGroupEntry(r, &windowGroupEntry{name: name, window: window})")
	wr.addWindowGroupEntry(r, &windowGroupEntry{name: name, window: window}) // object allocated on the heap: escapes at line 108
}

func (wr *WindowRoot) addWindowGroupEntry(r *rpn.RPN, wge *windowGroupEntry) {
	parentname, addend, weight := determineWindowParms(r)
	parent, err := wr.FindwindowGroup(parentname)
	if err != nil {
		// window not found, use the root group
		fr := rpn.StringFrame("root", rpn.STRING_SINGLEQ_FRAME)
		r.PushFrame(fr)
		r.SetVariable(".wtarget")
		parent = &wr.group
	}
	wge.weight = weight
	if addend {
		parent.children = append(parent.children, wge)
	} else {
		elog.Heap("alloc: /window/windowroot.go:125: parent.children = append([]*windowGroupEntry{wge}, parent.children...)")
		parent.children = append([]*windowGroupEntry{wge}, parent.children...) // object allocated on the heap: escapes at line 125
	}
	wr.adjustNeeded = true
}

func determineWindowParms(r *rpn.RPN) (parentname string, addend bool, weight int) {
	// save any parentname errors for addWindowGroupEntry
	parentname, _ = r.GetStringVariable(".wtarget")

	fr, err := r.GetVariable(".wend")
	if err != nil {
		fr = rpn.BoolFrame(true)
		r.PushFrame(fr)
		r.SetVariable(".wend")
	}
	addend, err = fr.Bool()
	if err != nil {
		addend = true
	}

	fr, err = r.GetVariable(".wweight")
	if err != nil {
		fr = rpn.IntFrame(100, rpn.INTEGER_FRAME)
		r.PushFrame(fr)
		r.SetVariable(".wweight")
	}
	weight64, _ := fr.Int()
	weight = int(weight64)
	if weight < 10 {
		weight = 10
	} else if weight > 10000 {
		weight = 10000
	}
	return
}

func (wr *WindowRoot) UseColumnLayout(name string, v bool) error {
	wg, err := wr.FindwindowGroup(name)
	if err != nil {
		return err
	}
	wg.isColumn = v
	wr.adjustNeeded = true
	return nil
}

// Calls update on all contained windows
func (wr *WindowRoot) Update(r *rpn.RPN, screenw, screenh int, updateInput bool) error {
	if wr.adjustNeeded {
		if err := wr.group.adjustChildren(screenw, screenh); err != nil {
			return err
		}
		wr.group.showBorder(screenw, screenh)
		wr.adjustNeeded = false
		// We want to give screens other than the input screen a chance to
		// redraw before getting locked into the input screen
		if err := wr.group.update(r, screenw, screenh); err != nil {
			return err
		}
	}
	if updateInput {
		// Update the input window first
		input := wr.FindWindow("i")
		if input == nil {
			return rpn.ErrNotFound
		}
		if err := input.Update(r, false); err != nil {
			return err
		}
	}
	return wr.group.update(r, screenw, screenh)
}

func (wr *WindowRoot) UpdateByName(r *rpn.RPN, name string, force bool) error {
	if name == "root" {
		return wr.Update(r, wr.group.w, wr.group.y, false)
	}
	wge := wr.group.findwindowGroupEntry(name)
	if wge == nil {
		return rpn.ErrNotFound
	}
	if wge.group != nil {
		return wge.group.update(r, wr.group.w, wr.group.h)
	}
	return wge.window.Update(r, force)
}

// Removes all children
func (wr *WindowRoot) RemoveAllChildren() {
	wr.group.removeAllChildren()
}

func (wr *WindowRoot) Dump(r *rpn.RPN) {
	wr.group.dump(r, "root", 0, 100)
	r.Print("\n.wtarget=")
	target, err := r.GetStringVariable(".wtarget")
	if err != nil {
		r.Print("root (unset)")
	} else {
		r.Print(target)
	}

	r.Print(" .wend=")
	fr, err := r.GetVariable(".wend")
	if err != nil {
		r.Print("true (unset)")
	} else if !fr.IsBool() {
		r.Print("true (not a bool)")
	} else {
		r.Print(fr.String(false))
	}

	r.Print(" .wweight=")
	fr, err = r.GetVariable(".wweight")
	if err != nil {
		r.Println("100 (unset)")
	} else if fr.IsNumber() {
		r.Println(fr.String(false))
	} else {
		r.Println("100 (bad type)")
	}
}

func (wr *WindowRoot) Snapshot(buff []byte, name string) ([]byte, error) {
	w := wr.FindWindow(name)
	var err error
	if w != nil {
		buff = SnapshotProps(buff, w, name)
	} else {
		wg, err := wr.FindwindowGroup(name)
		if err != nil {
			return nil, err
		}
		buff = wg.snapshot(buff, name)
	}
	return buff, err
}
