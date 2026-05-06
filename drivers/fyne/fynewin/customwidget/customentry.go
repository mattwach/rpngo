package customwidget

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

type CustomEntry struct {
	widget.Entry
	minEntryWidth float32
}

func (e *CustomEntry) MinSize() fyne.Size {
	ms := e.Entry.MinSize() // Get default minimum size
	if ms.Width < e.minEntryWidth {
		ms.Width = e.minEntryWidth // Set minimum width if it's less than the specified minimum
	}
	return ms
}

func NewCustomEntry(onSubmitted func(string), minEntryWidth float32) *CustomEntry {
	e := &CustomEntry{
		minEntryWidth: minEntryWidth,
	}
	e.ExtendBaseWidget(e)
	e.OnSubmitted = onSubmitted
	return e
}
