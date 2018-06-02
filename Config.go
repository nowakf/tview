package tview

type config struct {
	fontSize         float64
	fontPath         string
	adjustX, adjustY float64
}

func Config(fsize float64, fpath string, adjustX, adjustY float64) *config {
	return &config{
		fontSize: fsize,
		fontPath: fpath,
		adjustX:  adjustX,
		adjustY:  adjustY,
	}
}

func (c *config) FontSize() float64 {
	return c.fontSize
}

func (c *config) FontPath() string {
	return c.fontPath
}

func (c *config) AdjustXY() (float64, float64) {
	return c.adjustX, c.adjustY
}
