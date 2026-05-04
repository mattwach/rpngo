package fyneplotwin

import (
	"fmt"
	"mattwach/rpngo/common/rpn"
	"mattwach/rpngo/common/window/plotwin"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
	"github.com/fogleman/gg"
)

// InteractiveImage is a custom fyne widget that renders a set of plots
// and supports mouse actions, such as dragging and scrolling to zoom.
type interactiveImage struct {
	widget.BaseWidget
	// fyne is picky about how you call into it, requiring Do(AndWait) when
	// calling from a different go routine.
	inMainContext bool
	// We use gg to handle the pixel work
	ggimg *gg.Context
	image *canvas.Image
	// PixelPlotWindow (shared with picocalc), handles most of the graphing
	// logic.
	parent *plotwin.PixelPlotWindow
	// rpnInstance is needed when a graph change is initiated by the UI,
	// such as panning or zooming.  We need to "check out" the RPN instance,
	// which may not be possible if it is busy.
	rpnInstance chan *rpn.RPN

	// mouse drag event
	mouseDown desktop.MouseButton
	// when the mouse left button is down, the mouse cursor should be positioned
	// at the given point
	anchorX float64
	anchorY float64
	// for middle button magnify, we also care about tracking the mouse
	// position in terms of the graph window so that a x/y magnify scale
	// can be stabily calculated.  There may be a way to do the math without
	// this, but it is easier to understand with it.
	posx float32
	posy float32
	// the plotw and ploth are used to recalculate the graph min and max
	// coordinates
	plotw float64
	ploth float64
}

// CreateRenderer is required to satisfy the fyne.Widget interface.
func (i *interactiveImage) CreateRenderer() fyne.WidgetRenderer {
	i.inMainContext = true
	defer func() {
		i.inMainContext = false
	}()
	return widget.NewSimpleRenderer(i.image)
}

// MouseDown captures the initial point of a mouse drag event, which is used to
// anchor the position of the graph during panning.
func (i *interactiveImage) MouseDown(ev *desktop.MouseEvent) {
	i.mouseDown = ev.Button
	i.posx = ev.Position.X
	i.posy = ev.Position.Y
	i.anchorX, i.anchorY = i.parent.PixelToCoord(int(ev.Position.X), int(ev.Position.Y))
	i.plotw = float64(i.parent.Common.MaxX - i.parent.Common.MinX)
	i.ploth = float64(i.parent.Common.MaxY - i.parent.Common.MinY)
}

// MouseIn is required to satisfy the desktop.Mouseable interface, but is not
// used in this implementation.
func (i *interactiveImage) MouseIn(ev *desktop.MouseEvent) {}

// MouseMoved handles mouse movement events. If the mouse is currently down,
// it is treated as part of a drag event and the graph is panned accordingly.
// If the mouse is up, it is treated as a hover event and the coordinates under
// the cursor are displayed.
func (i *interactiveImage) MouseMoved(ev *desktop.MouseEvent) {
	i.inMainContext = true
	defer func() {
		i.inMainContext = false
	}()
	switch i.mouseDown {
	case desktop.MouseButtonPrimary:
		i.plotDragged(ev)
	case desktop.MouseButtonTertiary:
		i.magnify(ev)
	default:
		i.drawPointerCoordinates(ev)
	}
}

// MouseOut is required to satisfy the desktop.Mouseable interface, but is not
// used in this implementation.
func (i *interactiveImage) MouseOut() {}

// MouseUp captures the end of a mouse drag event.
func (i *interactiveImage) MouseUp(ev *desktop.MouseEvent) {
	i.mouseDown = 0
}

// Resize satisfies the fyne.Widget interface and is called when the widget is
// resized.
func (i *interactiveImage) Resize(size fyne.Size) {
	i.inMainContext = true
	defer func() {
		i.inMainContext = false
	}()
	i.BaseWidget.Resize(size)
	if i.parent == nil {
		return
	}
	select {
	case r := <-i.rpnInstance:
		defer func() {
			i.rpnInstance <- r
		}()
		i.ggimg = gg.NewContext(int(size.Width), int(size.Height))
		i.image.Image = i.ggimg.Image()
		i.image.Resize(size)
		i.parent.ResizeWindow(0, 0, int(size.Width), int(size.Height))
		i.parent.Update(r, false)
	default:
		// no action if the rpn instance is not available
	}
}

const scrollPercent = 0.25

// Scrolled satisfies the desktop.Scroller interface and is called when the
// user scrolls while the mouse is over the widget. It zooms in or out of the
// graph depending on the scroll direction.
func (i *interactiveImage) Scrolled(ev *fyne.ScrollEvent) {
	i.inMainContext = true
	defer func() {
		i.inMainContext = false
	}()
	select {
	case r := <-i.rpnInstance:
		defer func() {
			i.rpnInstance <- r
		}()
		plotw := float64(i.parent.Common.MaxX - i.parent.Common.MinX)
		ploth := float64(i.parent.Common.MaxY - i.parent.Common.MinY)

		if ev.Scrolled.DY > 0 {
			// zoom in
			plotw *= (1 - scrollPercent)
			ploth *= (1 - scrollPercent)
		} else {
			// zoom out
			plotw *= (1 + scrollPercent)
			ploth *= (1 + scrollPercent)
		}

		// the goal is to increase the bounds of the plot window while keeping
		// the point under the cursor fixed
		xoffset := plotw * float64(ev.Position.X) / float64(i.ggimg.Width())
		yoffset := ploth * (float64(i.ggimg.Height()) - float64(ev.Position.Y)) / float64(i.ggimg.Height())
		anchorX, anchorY := i.parent.PixelToCoord(int(ev.Position.X), int(ev.Position.Y))
		if i.isAutoX() {
			i.parent.Common.MinV = anchorX - xoffset
			i.parent.Common.MaxV = i.parent.Common.MinV + plotw
		} else {
			i.parent.Common.MinX = anchorX - xoffset
			i.parent.Common.MaxX = i.parent.Common.MinX + plotw
		}
		if !i.parent.Common.AutoY {
			i.parent.Common.MinY = anchorY - yoffset
			i.parent.Common.MaxY = i.parent.Common.MinY + ploth
		}
		i.parent.Update(r, true)
	default:
		// no action if the rpn instance is not available
	}
}

// plotDragged handles mouse movement events that occur while the mouse
// is down, which are treated as part of a drag event. It pans the graph so that
// the point under the cursor remains fixed at the anchor point.
func (i *interactiveImage) plotDragged(ev *desktop.MouseEvent) {
	select {
	case r := <-i.rpnInstance:
		defer func() {
			i.rpnInstance <- r
		}()
		// The goal is to set minx, maxv, miny, maxy so that the point
		// under the cursor is still at anchorX, anchorY
		xoffset := i.plotw * float64(ev.Position.X) / float64(i.ggimg.Width())
		yoffset := i.ploth * (float64(i.ggimg.Height()) - float64(ev.Position.Y)) / float64(i.ggimg.Height())
		if i.isAutoX() {
			i.parent.Common.MinV = i.anchorX - xoffset
			i.parent.Common.MaxV = i.parent.Common.MinV + i.plotw
		} else {
			i.parent.Common.MinX = i.anchorX - xoffset
			i.parent.Common.MaxX = i.parent.Common.MinX + i.plotw
		}
		i.parent.Common.AutoY = false
		i.parent.Common.MinY = i.anchorY - yoffset
		i.parent.Common.MaxY = i.parent.Common.MinY + i.ploth
		i.parent.Update(r, true)
	default:
		// no action if the rpn instance is not available
	}
}

const magnifyFactor = 1.5

func (i *interactiveImage) magnify(ev *desktop.MouseEvent) {
	i.inMainContext = true
	defer func() {
		i.inMainContext = false
	}()
	select {
	case r := <-i.rpnInstance:
		defer func() {
			i.rpnInstance <- r
		}()
		plotw := float64(i.parent.Common.MaxX - i.parent.Common.MinX)
		ploth := float64(i.parent.Common.MaxY - i.parent.Common.MinY)

		xscale := (ev.Position.X - i.posx) / float32(i.ggimg.Width()) * magnifyFactor
		plotw *= (1 - float64(xscale))

		yscale := (ev.Position.Y - i.posy) / float32(i.ggimg.Height()) * magnifyFactor
		ploth *= (1 + float64(yscale))

		// the goal is to increase the bounds of the plot window while keeping
		// the point under the cursor fixed
		xoffset := plotw * float64(i.posx) / float64(i.ggimg.Width())
		yoffset := ploth * (float64(i.ggimg.Height()) - float64(i.posy)) / float64(i.ggimg.Height())
		if i.isAutoX() {
			i.parent.Common.MinV = i.anchorX - xoffset
			i.parent.Common.MaxV = i.parent.Common.MinV + plotw
		} else {
			i.parent.Common.MinX = i.anchorX - xoffset
			i.parent.Common.MaxX = i.parent.Common.MinX + plotw
		}
		i.parent.Common.AutoY = false
		i.parent.Common.MinY = i.anchorY - yoffset
		i.parent.Common.MaxY = i.parent.Common.MinY + ploth
		i.parent.Update(r, true)
		i.posx = ev.Position.X
		i.posy = ev.Position.Y
	default:
		// no action if the rpn instance is not available
	}
}

// drawPointerCoordinates handles mouse movement events that occur while the mouse
// is up, which are treated as hover events. It displays the coordinates under
// the cursor.
func (i *interactiveImage) drawPointerCoordinates(ev *desktop.MouseEvent) {
	x, y := i.parent.PixelToCoord(int(ev.Position.X), int(ev.Position.Y))
	s := fmt.Sprintf("(%.4f, %.4f)", x, y)
	w, h := i.ggimg.MeasureString(s)
	i.ggimg.SetRGB(0, 0, 0)
	i.ggimg.DrawRectangle(20, 20, w+20, h)
	i.ggimg.Fill()
	i.ggimg.SetRGB(1, 1, 1)
	i.ggimg.DrawString(s, 20, 20+h)
	i.image.Refresh()
}

// Returns true if the plot should be usnig auto x mode
func (i *interactiveImage) isAutoX() bool {
	if !i.parent.Common.AutoX {
		return false
	}
	if i.parent.Common.HasAParametricPlot() {
		i.parent.Common.AutoX = false
		return false
	}
	return true
}
