package hoverpopover

import (
	"github.com/diamondburned/gotk4/pkg/glib/v2"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
)

// PopoverController  provides a way to open and close a popover while also
// reusing the widget if it's been open recently.
type PopoverController struct {
	parent      *gtk.Widget
	popover     *gtk.Popover
	initPopover func(*gtk.Popover)
	hideTimeout glib.SourceHandle
}

// NewPopoverController creates a new PopoverController.
func NewPopoverController(parent gtk.Widgetter, initFn func(*gtk.Popover)) *PopoverController {
	return &PopoverController{
		parent:      gtk.BaseWidget(parent),
		initPopover: initFn,
	}
}

// Popup pops up the popover.
func (p *PopoverController) Popup() {
	if p.popover != nil {
		if p.hideTimeout != 0 {
			glib.SourceRemove(p.hideTimeout)
			p.hideTimeout = 0
		}

		p.popover.SetCSSClasses(nil)
		p.initPopover(p.popover)

		p.popover.Popup()
		return
	}

	p.popover = gtk.NewPopover()
	p.popover.SetCSSClasses(nil)
	p.popover.SetParent(p.parent)
	p.initPopover(p.popover)
	p.popover.Popup()
}

// Popdown pops down the popover.
func (p *PopoverController) Popdown() {
	if p.popover == nil {
		return
	}

	p.popover.Popdown()

	if p.hideTimeout != 0 {
		return
	}

	p.hideTimeout = glib.TimeoutSecondsAddPriority(3, glib.PriorityLow, func() {
		p.popover.Unparent()
		p.popover = nil
		p.hideTimeout = 0
	})
}

// IsPoppedUp returns whether the popover is popped up.
func (p *PopoverController) IsPoppedUp() bool {
	return p.popover != nil && p.popover.IsVisible()
}