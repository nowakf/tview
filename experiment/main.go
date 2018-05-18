package main

import (
	"github.com/nowakf/pixel/pixelgl"
	"github.com/nowakf/tview"
)

func run() {
	tview.NewApplication()
}
func main() {
	pixelgl.Run(run)
}
