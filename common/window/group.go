package window

import (
	"mattwach/rpngo/common/rpn"
	"strconv"
)

type windowGroupEntry struct {
	name   string
	weight int
	// Only one of the following should be non-nil
	group  *windowGroup
	window WindowWithProps
}

func (wge *windowGroupEntry) resize(x, y, w, h int) {
	if wge.group != nil {
		wge.group.resize(x, y, w, h)
	} else if wge.window != nil {
		wge.window.ResizeWindow(x, y, w, h)
	}
}

func (wge *windowGroupEntry) showBorder(screenw, screenh int) error {
	if wge.group != nil {
		wge.group.showBorder(screenw, screenh)
	} else if wge.window != nil {
		if err := wge.window.ShowBorder(screenw, screenh); err != nil {
			return err
		}
	}
	return nil
}

type windowGroup struct {
	isColumn bool
	// Coordinates are in global screen coordinates
	x        int
	y        int
	w        int
	h        int
	children []*windowGroupEntry
}

// Generates commands that would be needed to recreate the group
func (wg *windowGroup) snapshot(buff []byte, name string) []byte {
	if wg.isColumn {
		buff = append(buff, '\'')
		buff = append(buff, []byte(name)...)
		buff = append(buff, []byte("' w.columns\n")...)
	}
	for _, c := range wg.children {
		buff = c.snapshot(buff, name)
	}
	return buff
}

func (wge *windowGroupEntry) snapshot(buff []byte, parent string) []byte {
	if wge.group != nil {
		return wge.snapshotGroup(buff, parent)
	}
	return wge.snapshotWindow(buff, parent)
}

func (wge *windowGroupEntry) snapshotGroup(buff []byte, parent string) []byte {
	buff = append(buff, '\'')
	buff = append(buff, []byte(wge.name)...)
	buff = append(buff, []byte("' w.new.group\n")...)
	buff = wge.appendWeight(buff)
	if parent != "root" {
		buff = wge.moveToEnd(buff, parent)
	}
	return wge.group.snapshot(buff, wge.name)
}

func (wge *windowGroupEntry) snapshotWindow(buff []byte, parent string) []byte {
	if wge.window.Type() != "input" {
		buff = append(buff, '\'')
		buff = append(buff, []byte(wge.name)...)
		buff = append(buff, []byte("' w.new.")...)
		buff = append(buff, []byte(wge.window.Type())...)
		buff = append(buff, '\n')
	}
	buff = wge.appendWeight(buff)
	if (parent != "root") || (wge.window.Type() == "input") {
		buff = wge.moveToEnd(buff, parent)
	}
	return SnapshotProps(buff, wge.window, wge.name)
}

func (wge *windowGroupEntry) appendWeight(buff []byte) []byte {
	if wge.weight == 100 {
		// it's the default value
		return buff
	}
	buff = append(buff, '\'')
	buff = append(buff, []byte(wge.name)...)
	buff = append(buff, []byte("' ")...)
	buff = append(buff, []byte(strconv.Itoa(wge.weight))...)
	buff = append(buff, []byte(" w.weight\n")...)
	return buff
}

func (wge *windowGroupEntry) moveToEnd(buff []byte, parent string) []byte {
	buff = append(buff, '\'')
	buff = append(buff, []byte(wge.name)...)
	buff = append(buff, []byte("' '")...)
	buff = append(buff, []byte(parent)...)
	buff = append(buff, []byte("' w.move.end\n")...)
	return buff
}

func SnapshotProps(buff []byte, w WindowWithProps, name string) []byte {
	for _, p := range w.ListProps() {
		f, _ := w.GetProp(p)
		buff = append(buff, '\'')
		buff = append(buff, []byte(name)...)
		buff = append(buff, []byte("' '")...)
		buff = append(buff, []byte(p)...)
		buff = append(buff, []byte("' ")...)
		buff = append(buff, []byte(f.String(true))...)
		buff = append(buff, []byte(" w.setp\n")...)
	}
	return buff
}

func (wg *windowGroup) removeChild(c *windowGroupEntry) {
	idx := -1
	for i, rc := range wg.children {
		if rc == c {
			idx = i
			break
		}
	}
	// idx may be -1 here, which will crash, but that is what we want
	wg.children = append(wg.children[:idx], wg.children[idx+1:]...)
}

func (wg *windowGroup) findwindowGroupEntry(name string) *windowGroupEntry {
	for _, c := range wg.children {
		if c.name == name {
			return c
		}
		if c.group != nil {
			wge := c.group.findwindowGroupEntry(name)
			if wge != nil {
				return wge
			}
		}
	}
	return nil
}

func (wg *windowGroup) findwindowGroupEntryAndParent(name string) (*windowGroup, *windowGroupEntry) {
	for _, c := range wg.children {
		if c.name == name {
			return wg, c
		}
		if c.group != nil {
			g, wge := c.group.findwindowGroupEntryAndParent(name)
			if wge != nil {
				return g, wge
			}
		}
	}
	return nil, nil
}

func (wg *windowGroup) resize(x, y, w, h int) {
	wg.x = x
	wg.y = y
	wg.w = w
	wg.h = h
}

func (wg *windowGroup) dump(r *rpn.RPN, name string, indent int, weight int) {
	pad(r, indent)
	r.Print(name)
	r.Print("(x=")
	r.Print(strconv.Itoa(wg.x))
	r.Print(", y=")
	r.Print(strconv.Itoa(wg.y))
	r.Print(", w=")
	r.Print(strconv.Itoa(wg.w))
	r.Print(", h=")
	r.Print(strconv.Itoa(wg.h))
	r.Print(", cols=")
	r.Print(strconv.FormatBool(wg.isColumn))
	r.Print(", weight=")
	r.Print(strconv.Itoa(weight))
	r.Print("):\n")
	indent++
	for _, c := range wg.children {
		if c.group != nil {
			c.group.dump(r, c.name, indent, c.weight)
		}
		if c.window != nil {
			x, y := c.window.WindowXY()
			w, h := c.window.WindowSize()
			pad(r, indent)
			r.Print(c.name)
			r.Print("(type=")
			r.Print(c.window.Type())
			r.Print(", x=")
			r.Print(strconv.Itoa(x))
			r.Print(", y=")
			r.Print(strconv.Itoa(y))
			r.Print(", w=")
			r.Print(strconv.Itoa(w))
			r.Print(", h=")
			r.Print(strconv.Itoa(h))
			r.Print(", weight=")
			r.Print(strconv.Itoa(c.weight))
			r.Print(")\n")
		}
	}
}

func pad(r *rpn.RPN, indent int) {
	for range indent {
		r.Print("  ")
	}
}

func (wg *windowGroup) removeAllChildren() {
	for i, c := range wg.children {
		if c.group != nil {
			c.group.removeAllChildren()
		}
		wg.children[i] = nil
	}
	wg.children = make([]*windowGroupEntry, 0)
}

func (wg *windowGroup) adjustChildren(screenw, screenh int) error {
	totalWeight := 0
	for _, c := range wg.children {
		totalWeight += c.weight
	}
	if wg.isColumn {
		wg.adjustChildrenColumn(screenw, screenh, totalWeight)
	} else {
		wg.adjustChildrenRow(screenw, screenh, totalWeight)
	}
	var err error
	return err
}

func (wg *windowGroup) adjustChildrenColumn(screenw, screenh, totalWeight int) {
	x1 := wg.x
	weightSum := 0
	for _, c := range wg.children {
		weightSum += c.weight
		x2 := wg.x + wg.w*weightSum/totalWeight
		c.resize(x1, wg.y, x2-x1, wg.h)
		if c.group != nil {
			c.group.adjustChildren(screenw, screenh)
		}
		x1 = x2 - 1
	}
}

func (wg *windowGroup) adjustChildrenRow(screenw, screenh, totalWeight int) {
	y1 := wg.y
	weightSum := 0
	for _, c := range wg.children {
		weightSum += c.weight
		y2 := wg.y + wg.h*weightSum/totalWeight
		c.resize(wg.x, y1, wg.w, y2-y1)
		if c.group != nil {
			c.group.adjustChildren(screenw, screenh)
		}
		y1 = y2 - 1
	}
}

func (wg *windowGroup) showBorder(screenw, screenh int) error {
	for _, c := range wg.children {
		if err := c.showBorder(screenw, screenh); err != nil {
			return err
		}
	}
	return nil
}

// Calls update on all contained windows
func (wg *windowGroup) update(rpn *rpn.RPN, screenw, screenh int) error {
	for _, c := range wg.children {
		if c.name == "i" {
			continue
		}
		if c.window != nil {
			if err := c.window.Update(rpn, false); err != nil {
				return err
			}
			continue
		}
		if err := c.group.update(rpn, screenw, screenh); err != nil {
			return err
		}
	}
	return nil
}
