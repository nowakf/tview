package tview

import (
	"image/color"
	"math"
	"strings"
	"unicode/utf8"

	runewidth "github.com/mattn/go-runewidth"
	"github.com/nowakf/pixel/pixelgl"
	"github.com/nowakf/ubcell"
)

// InputField is a one-line box (three lines if there is a title) where the
// user can enter text.
//
// Use SetMaskCharacter() to hide input from onlookers (e.g. for password
// input).
//
// See https://github.com/rivo/tview/wiki/InputField for an example.
type InputField struct {
	*Box

	// The text that was entered.
	text string

	// The text to be displayed before the input area.
	label string

	// The text to be displayed in the input area when "text" is empty.
	placeholder string

	// The label color.
	labelColor color.RGBA

	// The background color of the input area.
	fieldBackgroundColor color.RGBA

	// The text color of the input area.
	fieldTextColor color.RGBA

	// The text color of the placeholder.
	placeholderTextColor color.RGBA

	// The screen width of the input area. A value of 0 means extend as much as
	// possible.
	fieldWidth int

	// A character to mask entered text (useful for password fields). A value of 0
	// disables masking.
	maskCharacter rune

	// An optional function which may reject the last character that was entered.
	accept func(text string, ch rune) bool

	// An optional function which is called when the input has changed.
	changed func(text string)

	// An optional function which is called when the user indicated that they
	// are done entering text. The key which was pressed is provided (tab,
	// shift-tab, enter, or escape).
	done func(*pixelgl.KeyEv)
}

// NewInputField returns a new input field.
func NewInputField() *InputField {
	return &InputField{
		Box:                  NewBox(),
		labelColor:           Styles.SecondaryTextColor,
		fieldBackgroundColor: Styles.ContrastBackgroundColor,
		fieldTextColor:       Styles.PrimaryTextColor,
		placeholderTextColor: Styles.ContrastSecondaryTextColor,
	}
}

// SetText sets the current text of the input field.
func (i *InputField) SetText(text string) *InputField {
	i.text = text
	if i.changed != nil {
		i.changed(text)
	}
	return i
}

// GetText returns the current text of the input field.
func (i *InputField) GetText() string {
	return i.text
}

// SetLabel sets the text to be displayed before the input area.
func (i *InputField) SetLabel(label string) *InputField {
	i.label = label
	return i
}

// GetLabel returns the text to be displayed before the input area.
func (i *InputField) GetLabel() string {
	return i.label
}

// SetPlaceholder sets the text to be displayed when the input text is empty.
func (i *InputField) SetPlaceholder(text string) *InputField {
	i.placeholder = text
	return i
}

// SetLabelColor sets the color of the label.
func (i *InputField) SetLabelColor(color color.RGBA) *InputField {
	i.labelColor = color
	return i
}

// SetFieldBackgroundColor sets the background color of the input area.
func (i *InputField) SetFieldBackgroundColor(color color.RGBA) *InputField {
	i.fieldBackgroundColor = color
	return i
}

// SetFieldTextColor sets the text color of the input area.
func (i *InputField) SetFieldTextColor(color color.RGBA) *InputField {
	i.fieldTextColor = color
	return i
}

// SetPlaceholderExtColor sets the text color of placeholder text.
func (i *InputField) SetPlaceholderExtColor(color color.RGBA) *InputField {
	i.placeholderTextColor = color
	return i
}

// SetFormAttributes sets attributes shared by all form items.
func (i *InputField) SetFormAttributes(label string, labelColor, bgColor, fieldTextColor, fieldBgColor color.RGBA) FormItem {
	i.label = label
	i.labelColor = labelColor
	i.backgroundColor = bgColor
	i.fieldTextColor = fieldTextColor
	i.fieldBackgroundColor = fieldBgColor
	return i
}

// SetFieldWidth sets the screen width of the input area. A value of 0 means
// extend as much as possible.
func (i *InputField) SetFieldWidth(width int) *InputField {
	i.fieldWidth = width
	return i
}

// GetFieldWidth returns this primitive's field width.
func (i *InputField) GetFieldWidth() int {
	return i.fieldWidth
}

// SetMaskCharacter sets a character that masks user input on a screen. A value
// of 0 disables masking.
func (i *InputField) SetMaskCharacter(mask rune) *InputField {
	i.maskCharacter = mask
	return i
}

// SetAcceptanceFunc sets a handler which may reject the last character that was
// entered (by returning false).
//
// This package defines a number of variables Prefixed with InputField which may
// be used for common input (e.g. numbers, maximum text length).
func (i *InputField) SetAcceptanceFunc(handler func(textToCheck string, lastChar rune) bool) *InputField {
	i.accept = handler
	return i
}

// SetChangedFunc sets a handler which is called whenever the text of the input
// field has changed. It receives the current text (after the change).
func (i *InputField) SetChangedFunc(handler func(text string)) *InputField {
	i.changed = handler
	return i
}

// SetDoneFunc sets a handler which is called when the user is done entering
// text. The callback function is provided with the key that was pressed, which
// is one of the following:
//
//   - KeyEnter: Done entering text.
//   - KeyEscape: Abort text input.
//   - KeyTab: Move to the next field.
//   - KeyBacktab: Move to the previous field.
func (i *InputField) SetDoneFunc(handler func(key *pixelgl.KeyEv)) *InputField {
	i.done = handler
	return i
}

// SetFinishedFunc calls SetDoneFunc().
func (i *InputField) SetFinishedFunc(handler func(key *pixelgl.KeyEv)) FormItem {
	return i.SetDoneFunc(handler)
}

// Draw draws this primitive onto the screen.
func (i *InputField) Draw(screen ubcell.Screen) {
	i.Box.Draw(screen)

	// Prepare
	x, y, width, height := i.GetInnerRect()
	rightLimit := x + width
	if height < 1 || rightLimit <= x {
		return
	}

	// Draw label.
	_, drawnWidth := Print(screen, i.label, x, y, rightLimit-x, AlignLeft, i.labelColor)
	x += drawnWidth

	// Draw input area.
	fieldWidth := i.fieldWidth
	if fieldWidth == 0 {
		fieldWidth = math.MaxInt32
	}
	if rightLimit-x < fieldWidth {
		fieldWidth = rightLimit - x
	}
	fieldStyle := ubcell.StyleDefault.Background(i.fieldBackgroundColor)
	for index := 0; index < fieldWidth; index++ {
		screen.SetContent(x+index, y, ' ', fieldStyle)
	}

	// Draw placeholder text.
	text := i.text
	if text == "" && i.placeholder != "" {
		Print(screen, i.placeholder, x, y, fieldWidth, AlignLeft, i.placeholderTextColor)
	}

	// Draw entered text.
	if i.maskCharacter > 0 {
		text = strings.Repeat(string(i.maskCharacter), utf8.RuneCountInString(i.text))
	}
	fieldWidth-- // We need one cell for the cursor.
	if fieldWidth < runewidth.StringWidth(text) {
		runes := []rune(text)
		for pos := len(runes) - 1; pos >= 0; pos-- {
			ch := runes[pos]
			w := runewidth.RuneWidth(ch)
			if fieldWidth-w < 0 {
				break
			}
			_, style := screen.GetContent(x+fieldWidth-w, y)
			style.Foreground(i.fieldTextColor)
			for w > 0 {
				fieldWidth--
				screen.SetContent(x+fieldWidth, y, ch, style)
				w--
			}
		}
	} else {
		pos := 0
		for _, ch := range text {
			w := runewidth.RuneWidth(ch)
			_, style := screen.GetContent(x+pos, y)
			style.Foreground(i.fieldTextColor)
			for w > 0 {
				screen.SetContent(x+pos, y, ch, style)
				pos++
				w--
			}
		}
	}

	// Set cursor.
	if i.focus.HasFocus() {
		i.setCursor(screen)
	}
}

// setCursor sets the cursor position.
func (i *InputField) setCursor(screen ubcell.Screen) {
	x := i.x
	y := i.y
	rightLimit := x + i.width
	if i.border {
		x++
		y++
		rightLimit -= 2
	}
	fieldWidth := runewidth.StringWidth(i.text)
	if i.fieldWidth > 0 && fieldWidth > i.fieldWidth-1 {
		fieldWidth = i.fieldWidth - 1
	}
	x += StringWidth(i.label) + fieldWidth
	if x >= rightLimit {
		x = rightLimit - 1
	}
	screen.ShowCursor(x, y)
}

// KeyHandler returns the handler for this primitive.
func (i *InputField) KeyHandler() func(event *pixelgl.KeyEv, setFocus func(p Primitive)) {
	return i.WrapKeyHandler(func(event *pixelgl.KeyEv, setFocus func(p Primitive)) {
		// Trigger changed events.
		currentText := i.text
		defer func() {
			if i.text != currentText && i.changed != nil {
				i.changed(i.text)
			}
		}()

		// Process key event.
		switch key := event.Key; key {
		case pixelgl.KeyRune: // Regular character.
			newText := i.text + string(event.Ch)
			if i.accept != nil {
				if !i.accept(newText, event.Ch) {
					break
				}
			}
			i.text = newText
		case pixelgl.KeyBackspace:
			if len(i.text) == 0 {
				break
			}
			runes := []rune(i.text)
			i.text = string(runes[:len(runes)-1])
		case pixelgl.KeyEnter, pixelgl.KeyTab, pixelgl.KeyEscape: // We're done.
			if i.done != nil {
				i.done(event)
			}
		}
	})
}
