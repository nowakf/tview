package tview

import (
	"image/color"

	"golang.org/x/image/colornames"
)

// Styles defines various colors used when primitives are initialized. These
// may be changed to accommodate a different look and feel.
//
// The default is for applications with a black background and basic colors:
// black, white, yellow, green, and blue.
var Styles = struct {
	PrimitiveBackgroundColor    color.RGBA // Main background color for primitives.
	ContrastBackgroundColor     color.RGBA // Background color for contrasting elements.
	MoreContrastBackgroundColor color.RGBA // Background color for even more contrasting elements.
	BorderColor                 color.RGBA // Box borders.
	TitleColor                  color.RGBA // Box titles.
	GraphicsColor               color.RGBA // Graphics.
	PrimaryTextColor            color.RGBA // Primary text.
	SecondaryTextColor          color.RGBA // Secondary text (e.g. labels).
	TertiaryTextColor           color.RGBA // Tertiary text (e.g. subtitles, notes).
	InverseTextColor            color.RGBA // Text on primary-colored backgrounds.
	ContrastSecondaryTextColor  color.RGBA // Secondary text on ContrastBackgroundColor-colored backgrounds.
}{
	PrimitiveBackgroundColor:    colornames.Dimgray,
	ContrastBackgroundColor:     colornames.Grey,
	MoreContrastBackgroundColor: colornames.Darkblue,
	BorderColor:                 colornames.Lightgrey,
	TitleColor:                  colornames.Lightgrey,
	GraphicsColor:               colornames.Lightgrey,
	PrimaryTextColor:            colornames.Lightgrey,
	SecondaryTextColor:          colornames.Lightgrey,
	TertiaryTextColor:           colornames.Lightgoldenrodyellow,
	InverseTextColor:            colornames.Yellow,
	ContrastSecondaryTextColor:  colornames.Pink,
}
