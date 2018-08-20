package tview

import (
	//	"github.com/nowakf/pixel/pixelgl"
	"github.com/nowakf/ubcell"
)

type CellMap struct {
	*Box
	delay int

	finished func(ubcell.Screen)
}

func (c *CellMap) Draw(screen ubcell.Screen) {
	c.Box.Draw(screen)
}
