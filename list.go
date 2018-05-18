package tview

import (
	"fmt"
	"image/color"

	"github.com/nowakf/pixel/pixelgl"
	"github.com/nowakf/ubcell"
)

// listItem represents one item in a List.
type listItem struct {
	MainText      string // The main text of the list item.
	SecondaryText string // A secondary text to be shown underneath the main text.
	Shortcut      rune   // The key to select the list item directly, 0 if there is no shortcut.
	Selected      func() // The optional function which is called when the item is selected.
}

// List displays rows of items, each of which can be selected.
//
// See https://github.com/rivo/tview/wiki/List for an example.
type List struct {
	*Box

	// The items of the list.
	items []*listItem

	// The index of the currently selected item.
	currentItem int

	// Whether or not to show the secondary item texts.
	showSecondaryText bool

	// The item main text color.
	mainTextColor color.RGBA

	// The item secondary text color.
	secondaryTextColor color.RGBA

	// The item shortcut text color.
	shortcutColor color.RGBA

	// The text color for selected items.
	selectedTextColor color.RGBA

	// The background color for selected items.
	selectedBackgroundColor color.RGBA

	// An optional function which is called when the user has navigated to a list
	// item.
	changed func(index int, mainText, secondaryText string, shortcut rune)

	// An optional function which is called when a list item was selected. This
	// function will be called even if the list item defines its own callback.
	selected func(index int, mainText, secondaryText string, shortcut rune)

	// An optional function which is called when the user presses the Escape key.
	done func()
}

// NewList returns a new form.
func NewList() *List {
	return &List{
		Box:                     NewBox(),
		showSecondaryText:       true,
		mainTextColor:           Styles.PrimaryTextColor,
		secondaryTextColor:      Styles.TertiaryTextColor,
		shortcutColor:           Styles.SecondaryTextColor,
		selectedTextColor:       Styles.PrimitiveBackgroundColor,
		selectedBackgroundColor: Styles.PrimaryTextColor,
	}
}

// SetCurrentItem sets the currently selected item by its index. This triggers
// a "changed" event.
func (l *List) SetCurrentItem(index int) *List {
	l.currentItem = index
	if l.currentItem < len(l.items) && l.changed != nil {
		item := l.items[l.currentItem]
		l.changed(l.currentItem, item.MainText, item.SecondaryText, item.Shortcut)
	}
	return l
}

// GetCurrentItem returns the index of the currently selected list item.
func (l *List) GetCurrentItem() int {
	return l.currentItem
}

// SetMainTextColor sets the color of the items' main text.
func (l *List) SetMainTextColor(color color.RGBA) *List {
	l.mainTextColor = color
	return l
}

// SetSecondaryTextColor sets the color of the items' secondary text.
func (l *List) SetSecondaryTextColor(color color.RGBA) *List {
	l.secondaryTextColor = color
	return l
}

// SetShortcutColor sets the color of the items' shortcut.
func (l *List) SetShortcutColor(color color.RGBA) *List {
	l.shortcutColor = color
	return l
}

// SetSelectedTextColor sets the text color of selected items.
func (l *List) SetSelectedTextColor(color color.RGBA) *List {
	l.selectedTextColor = color
	return l
}

// SetSelectedBackgroundColor sets the background color of selected items.
func (l *List) SetSelectedBackgroundColor(color color.RGBA) *List {
	l.selectedBackgroundColor = color
	return l
}

// ShowSecondaryText determines whether or not to show secondary item texts.
func (l *List) ShowSecondaryText(show bool) *List {
	l.showSecondaryText = show
	return l
}

// SetChangedFunc sets the function which is called when the user navigates to
// a list item. The function receives the item's index in the list of items
// (starting with 0), its main text, secondary text, and its shortcut rune.
//
// This function is also called when the first item is added or when
// SetCurrentItem() is called.
func (l *List) SetChangedFunc(handler func(int, string, string, rune)) *List {
	l.changed = handler
	return l
}

// SetSelectedFunc sets the function which is called when the user selects a
// list item by pressing Enter on the current selection. The function receives
// the item's index in the list of items (starting with 0), its main text,
// secondary text, and its shortcut rune.
func (l *List) SetSelectedFunc(handler func(int, string, string, rune)) *List {
	l.selected = handler
	return l
}

// SetDoneFunc sets a function which is called when the user presses the Escape
// key.
func (l *List) SetDoneFunc(handler func()) *List {
	l.done = handler
	return l
}

// AddItem adds a new item to the list. An item has a main text which will be
// highlighted when selected. It also has a secondary text which is shown
// underneath the main text (if it is set to visible) but which may remain
// empty.
//
// The shortcut is a key binding. If the specified rune is entered, the item
// is selected immediately. Set to 0 for no binding.
//
// The "selected" callback will be invoked when the user selects the item. You
// may provide nil if no such item is needed or if all events are handled
// through the selected callback set with SetSelectedFunc().
func (l *List) AddItem(mainText, secondaryText string, shortcut rune, selected func()) *List {
	l.items = append(l.items, &listItem{
		MainText:      mainText,
		SecondaryText: secondaryText,
		Shortcut:      shortcut,
		Selected:      selected,
	})
	if len(l.items) == 1 && l.changed != nil {
		item := l.items[0]
		l.changed(0, item.MainText, item.SecondaryText, item.Shortcut)
	}
	return l
}

// Clear removes all items from the list.
func (l *List) Clear() *List {
	l.items = nil
	l.currentItem = 0
	return l
}

// Draw draws this primitive onto the screen.
func (l *List) Draw(screen ubcell.Screen) {
	l.Box.Draw(screen)

	// Determine the dimensions.
	x, y, width, height := l.GetInnerRect()
	bottomLimit := y + height

	// Do we show any shortcuts?
	var showShortcuts bool
	for _, item := range l.items {
		if item.Shortcut != 0 {
			showShortcuts = true
			x += 4
			width -= 4
			break
		}
	}

	// We want to keep the current selection in view. What is our offset?
	var offset int
	if l.showSecondaryText {
		if l.currentItem >= height/2 {
			offset = l.currentItem + 1 - (height / 2)
		}
	} else {
		if l.currentItem >= height {
			offset = l.currentItem + 1 - height
		}
	}

	// Draw the list items.
	for index, item := range l.items {
		if index < offset {
			continue
		}

		if y >= bottomLimit {
			break
		}

		// Shortcuts.
		if showShortcuts && item.Shortcut != 0 {
			Print(screen, fmt.Sprintf("(%s)", string(item.Shortcut)), x-5, y, 4, AlignRight, l.shortcutColor)
		}

		// Main text.
		Print(screen, item.MainText, x, y, width, AlignLeft, l.mainTextColor)

		// Background color of selected text.
		if index == l.currentItem {
			textWidth := StringWidth(item.MainText)
			for bx := 0; bx < textWidth && bx < width; bx++ {
				m, style := screen.GetContent(x+bx, y)
				fg := style.Foreground
				if fg == l.mainTextColor {
					fg = l.selectedTextColor
				}
				style = ubcell.Style{l.selectedBackgroundColor, fg}
				screen.SetContent(x+bx, y, m, style)
			}
		}

		y++

		if y >= bottomLimit {
			break
		}

		// Secondary text.
		if l.showSecondaryText {
			Print(screen, item.SecondaryText, x, y, width, AlignLeft, l.secondaryTextColor)
			y++
		}
	}
}

// KeyHandler returns the handler for this primitive.
func (l *List) KeyHandler() func(event *pixelgl.KeyEv, setFocus func(p Primitive)) {
	return l.WrapKeyHandler(func(event *pixelgl.KeyEv, setFocus func(p Primitive)) {
		previousItem := l.currentItem

		switch event.Key {
		case pixelgl.KeyTab, pixelgl.KeyDown, pixelgl.KeyRight:
			l.currentItem++
		case pixelgl.KeyUp, pixelgl.KeyLeft:
			l.currentItem--
		case pixelgl.KeyHome:
			l.currentItem = 0
		case pixelgl.KeyEnd:
			l.currentItem = len(l.items) - 1
		case pixelgl.KeyPageDown:
			l.currentItem += 5
		case pixelgl.KeyPageUp:
			l.currentItem -= 5
		case pixelgl.KeyEnter:
			item := l.items[l.currentItem]
			if item.Selected != nil {
				item.Selected()
			}
			if l.selected != nil {
				l.selected(l.currentItem, item.MainText, item.SecondaryText, item.Shortcut)
			}
		case pixelgl.KeyEscape:
			if l.done != nil {
				l.done()
			}
		case pixelgl.KeyRune:
			ch := event.Ch
			if ch != ' ' {
				// It's not a space bar. Is it a shortcut?
				var found bool
				for index, item := range l.items {
					if item.Shortcut == ch {
						// We have a shortcut.
						found = true
						l.currentItem = index
						break
					}
				}
				if !found {
					break
				}
			}
			item := l.items[l.currentItem]
			if item.Selected != nil {
				item.Selected()
			}
			if l.selected != nil {
				l.selected(l.currentItem, item.MainText, item.SecondaryText, item.Shortcut)
			}
		}

		if l.currentItem < 0 {
			l.currentItem = len(l.items) - 1
		} else if l.currentItem >= len(l.items) {
			l.currentItem = 0
		}

		if l.currentItem != previousItem && l.currentItem < len(l.items) && l.changed != nil {
			item := l.items[l.currentItem]
			l.changed(l.currentItem, item.MainText, item.SecondaryText, item.Shortcut)
		}
	})
}
