package tview

import (
	"image/color"

	"github.com/nowakf/pixel/pixelgl"
	"github.com/nowakf/ubcell"
)

// Checkbox implements a simple box for boolean values which can be checked and
// unchecked.
//
// See https://github.com/rivo/tview/wiki/Checkbox for an example.
type Checkbox struct {
	*Box

	// Whether or not this box is checked.
	checked bool

	// The text to be displayed before the input area.
	label string

	// The label color.
	labelColor color.RGBA

	// The background color of the input area.
	fieldBackgroundColor color.RGBA

	// The text color of the input area.
	fieldTextColor color.RGBA

	// An optional function which is called when the user changes the checked
	// state of this checkbox.
	changed func(checked bool)

	// An optional function which is called when the user indicated that they
	// are done entering text. The key which was pressed is provided (tab,
	// shift-tab, or escape).
	done func(*pixelgl.KeyEv)
}

// NewCheckbox returns a new input field.
func NewCheckbox() *Checkbox {
	return &Checkbox{
		Box:                  NewBox(),
		labelColor:           Styles.SecondaryTextColor,
		fieldBackgroundColor: Styles.ContrastBackgroundColor,
		fieldTextColor:       Styles.PrimaryTextColor,
	}
}

// SetChecked sets the state of the checkbox.
func (c *Checkbox) SetChecked(checked bool) *Checkbox {
	c.checked = checked
	return c
}

// IsChecked returns whether or not the box is checked.
func (c *Checkbox) IsChecked() bool {
	return c.checked
}

// SetLabel sets the text to be displayed before the input area.
func (c *Checkbox) SetLabel(label string) *Checkbox {
	c.label = label
	return c
}

// GetLabel returns the text to be displayed before the input area.
func (c *Checkbox) GetLabel() string {
	return c.label
}

// SetLabelColor sets the color of the label.
func (c *Checkbox) SetLabelColor(color color.RGBA) *Checkbox {
	c.labelColor = color
	return c
}

// SetFieldBackgroundColor sets the background color of the input area.
func (c *Checkbox) SetFieldBackgroundColor(color color.RGBA) *Checkbox {
	c.fieldBackgroundColor = color
	return c
}

// SetFieldTextColor sets the text color of the input area.
func (c *Checkbox) SetFieldTextColor(color color.RGBA) *Checkbox {
	c.fieldTextColor = color
	return c
}

// SetFormAttributes sets attributes shared by all form items.
func (c *Checkbox) SetFormAttributes(label string, labelColor, bgColor, fieldTextColor, fieldBgColor color.RGBA) FormItem {
	c.label = label
	c.labelColor = labelColor
	c.backgroundColor = bgColor
	c.fieldTextColor = fieldTextColor
	c.fieldBackgroundColor = fieldBgColor
	return c
}

// GetFieldWidth returns this primitive's field width.
func (c *Checkbox) GetFieldWidth() int {
	return 1
}

// SetChangedFunc sets a handler which is called when the checked state of this
// checkbox was changed by the user. The handler function receives the new
// state.
func (c *Checkbox) SetChangedFunc(handler func(checked bool)) *Checkbox {
	c.changed = handler
	return c
}

// SetDoneFunc sets a handler which is called when the user is done entering
// text. The callback function is provided with the key that was pressed, which
// is one of the following:
//
//   - KeyEscape: Abort text input.
//   - KeyTab: Move to the next field.
//   - KeyBacktab: Move to the previous field.
func (c *Checkbox) SetDoneFunc(handler func(key *pixelgl.KeyEv)) *Checkbox {
	c.done = handler
	return c
}

// SetFinishedFunc calls SetDoneFunc().
func (c *Checkbox) SetFinishedFunc(handler func(key *pixelgl.KeyEv)) FormItem {
	return c.SetDoneFunc(handler)
}

// Draw draws this primitive onto the screen.
func (c *Checkbox) Draw(screen ubcell.Screen) {
	c.Box.Draw(screen)

	// Prepare
	x, y, width, height := c.GetInnerRect()
	rightLimit := x + width
	if height < 1 || rightLimit <= x {
		return
	}

	// Draw label.
	_, drawnWidth := Print(screen, c.label, x, y, rightLimit-x, AlignLeft, c.labelColor)
	x += drawnWidth

	// Draw checkbox.
	fieldStyle := ubcell.StyleDefault.Background(c.fieldBackgroundColor).Foreground(c.fieldTextColor)
	if c.focus.HasFocus() {
		fieldStyle = ubcell.StyleDefault.Background(c.fieldTextColor).Foreground(c.fieldBackgroundColor)
	}
	checkedRune := 'X'
	if !c.checked {
		checkedRune = ' '
	}
	screen.SetContent(x, y, checkedRune, fieldStyle)
}

// KeyHandler returns the handler for this primitive.
func (c *Checkbox) KeyHandler() func(event *pixelgl.KeyEv, setFocus func(p Primitive)) {
	return c.WrapKeyHandler(func(event *pixelgl.KeyEv, setFocus func(p Primitive)) {
		// Process key event.
		switch key := event.Key; key {
		case pixelgl.KeyRune, pixelgl.KeyEnter: // Check.
			if key == pixelgl.KeyRune && event.Ch != ' ' {
				break
			}
			c.checked = !c.checked
			if c.changed != nil {
				c.changed(c.checked)
			}
		case pixelgl.KeyTab, pixelgl.KeyEscape: // We're done.
			if c.done != nil {
				c.done(event)
			}
		}
	})
}
