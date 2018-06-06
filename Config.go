package tview

import "github.com/nowakf/pixel/pixelgl"

type Config struct {
	//ubcell config
	FontSize         float64
	FontPath         string
	AdjustX, AdjustY float64
	DPI              float64
	//pixel config
	WindowConfig pixelgl.WindowConfig
}

func (c *Config) GetFontSize() float64 {
	return c.FontSize
}

func (c *Config) GetFontPath() string {
	return c.FontPath
}

func (c *Config) GetAdjustXY() (float64, float64) {
	return c.AdjustX, c.AdjustY
}

func (c *Config) GetDPI() float64 {
	return c.DPI
}

func (c *Config) GetWindowConfig() pixelgl.WindowConfig {
	return c.WindowConfig
}
