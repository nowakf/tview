package tview

import (
	"image/color"
	"image/gif"
	"os"
	"time"

	"github.com/nowakf/pixel"
	"github.com/nowakf/pixel/pixelgl"
	"github.com/nowakf/ubcell"
)

type VideoPlayer struct {
	*Box
	sequence *gif.GIF
	mask     color.Color
	control  chan controlSig
	delay    int

	finished func(ubcell.Screen)
}
type controlSig int

const (
	redraw controlSig = iota
	complete
)

func NewVideoPlayer() *VideoPlayer {
	return &VideoPlayer{
		Box:   NewBox(),
		delay: 60,
	}
}
func (v *VideoPlayer) play(screen ubcell.Screen, delay int) {

	if v.sequence == nil {
		panic("you didn't call load")
	}

	v.control = make(chan controlSig)

	pic := pixel.PictureDataFromImage(v.sequence.Image[0])

	sprite := pixel.NewSprite(pic, pic.Bounds())

	step := time.NewTicker(time.Millisecond * time.Duration(delay))

	x, y, w, h := v.GetRect()

	for i := 0; i < len(v.sequence.Image); i++ {

		pic = pixel.PictureDataFromImage(v.sequence.Image[i])

		sprite.Set(pic, pic.Bounds())

		<-step.C

		screen.Call(func(win *pixelgl.Window) {
			sprite.DrawColorMask(win, v.GetTransform(screen, x, y, w, h), v.mask)
		})

	}
	v.finished(screen)

}
func (v *VideoPlayer) GetTransform(screen ubcell.Screen, x, y, w, h int) pixel.Matrix {
	return screen.GetMatrix(x, y, w, h)
}
func (v *VideoPlayer) Load(path string) (*VideoPlayer, error) {
	var err error

	file, err := os.Open(path)

	if err != nil {
		return nil, err
	}

	defer file.Close()

	v.sequence, err = gif.DecodeAll(file)

	if err != nil {
		return nil, err
	}

	return v, nil
}

//because overlaying videos is sometimes desirable, videoplayer can sometimes
//be visible even if draw is not called.
func (v *VideoPlayer) Draw(screen ubcell.Screen) {

	if v.finished == nil {
		panic("finish func not set")
	}

	v.Box.Draw(screen)

	if v.control == nil {
		println("play!", v.delay)
		go v.play(screen, v.delay)
	}

}
func (v *VideoPlayer) SetDelay(delay int) {
	v.delay = delay
}

func (v *VideoPlayer) SetMask(m color.RGBA) {
	v.mask = m
}

//this is kind of kludgy, but since there's no way to control
//application focus from this level, it will have to do.
func (v *VideoPlayer) SetFinishedFunc(f func()) {
	v.finished = func(screen ubcell.Screen) {
		close(v.control)
		v.control = nil
		f()
		screen.Show()
	}
}
