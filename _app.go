package tview

import (
	"github.com/golang/freetype/truetype"
	"github.com/nowakf/pixel"
	"github.com/nowakf/pixel/pixelgl"
	"github.com/nowakf/pixel/text"
	"github.com/nowakf/ubcell"
	"golang.org/x/image/font/gofont/gomono"
	"image/color"
	"unicode"
)

type app struct {
	ready chan bool

	win  *pixelgl.Window
	text *text.Text

	rootFullScreen bool
	focus          Primitive
	root           Primitive

	cells      *ubcell.CellBuffer
	xInc, yInc float64

	inputCapture func(event *pixelgl.KeyEv) *pixelgl.KeyEv

	// An optional callback function which is invoked just before the root
	// primitive is drawn.
	beforeDraw func(screen *ubcell.CellBuffer) bool

	// An optional callback function which is invoked after the root primitive
	// was drawn.
	afterDraw func(screen *ubcell.CellBuffer)
}

const (
	xAdjust = 1.1
	yAdjust = 1.1
)

func App(w *pixelgl.Window) *app {

	a := new(app)

	a.ready = make(chan bool)

	a.win = w

	//make the text

	ttf, _ := truetype.Parse(gomono.TTF)
	face := truetype.NewFace(ttf, &truetype.Options{Size: 18})
	atlas := text.NewAtlas(face, text.ASCII, text.RangeTable(unicode.S))
	a.text = text.New(w.Bounds().Center(), atlas)

	glyph := a.text.BoundsOf("A")

	a.xInc, a.yInc = glyph.W(), glyph.H()

	a.xInc -= xAdjust
	a.yInc -= yAdjust

	fh, fw := a.win.Bounds().H(), a.win.Bounds().W()

	height, width := int(fh/a.yInc), int(fw/a.xInc)

	a.cells = ubcell.NewCellBuffer(
		height,
		width,
		func(x int, y int, ch rune) {
			a.text.Add(ch, pixel.V(float64(x)*a.xInc, fh-((1+float64(y))*a.yInc)))
		},
		func(c color.RGBA) {
			a.text.Ink(c)
		},
		func() {
			a.text.Apply()
		},
	)
	return a

}

// SetInputCapture sets a function which captures all key events before they are
// forwarded to the key event handler of the primitive which currently has
// focus. This function can then choose to forward that key event (or a
// different one) by returning it or stop the key event processing by returning
// nil.
func (a *app) SetInputCapture(capture func(event *pixelgl.KeyEv) *pixelgl.KeyEv) *app {
	a.inputCapture = capture
	return a
}

// GetInputCapture returns the function installed with SetInputCapture() or nil
// if no such function has been installed.
func (a *app) GetInputCapture() func(event *pixelgl.KeyEv) *pixelgl.KeyEv {
	return a.inputCapture
}

//this is seperate because you can call it if the window resizes
func (a *app) SetRoot(root Primitive, fullscreen bool) *app {
	a.root = root
	a.rootFullScreen = fullscreen
	//some kind of setfocus?
	return a
}
func (a *app) Loop() {
	for {
		root := a.root
		cells := a.cells
		fullScreen := a.rootFullScreen

		if fullScreen && root != nil {
			height, width := a.cells.Size()
			root.SetRect(0, 0, width, height)
		}

		root.Draw(cells)

		a.ready <- cells.Draw()
		<-a.ready
	}
}

//so this should be called from the main thread...
func (a *app) Draw(t pixel.Target, m pixel.Matrix) {

	select {
	case bufferHasChanged := <-a.ready:
		if bufferHasChanged {
			a.text.Draw(t, m)
			a.text.Clear()
		}
		a.ready <- true
	default:
	}

}
