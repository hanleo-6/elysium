package render

import(
	"fmt"
	"strings"

	"github.com/lucasb-eyer/go-colorful"
)

// Each terminal character cell is split into a 2 x 4 rectangle grid of dots
const (
	dotsPerCellX = 2
	dotsPerCellY = 4
)

// using unicode braille pixel map for detail
var pixelMap = [4][2]byte{
	{0x01, 0x08},
	{0x02, 0x10},
	{0x04, 0x20},
	{0x40, 0x80},
}

type Canvas struct {
	width, height int // in terminal character cells
	dots          []bool // size (width * 2) x (height * 4)
	foreground    []colorful.Color // foreground color per cell
	background    []colorful.Color // background color per cell
}

func NewCanvas(width, height int) *Canvas {
	return &Canvas{
		width:  width,
		height: height,
		dots:   make([]bool, width*dotsPerCellX*height*dotsPerCellY),
		foreground: make([]colorful.Color, width*height),
		background: make([]colorful.Color, width*height),
	}
}

func (c *Canvas) DotWidth() int { return dotsPerCellX * c.width }

func (c *Canvas) DotHeight() int { return dotsPerCellY * c.height }

func (c *Canvas) dotAt(x, y int) *bool {
	return &c.dots[y*c.width*dotsPerCellX + x]
}

func (c *Canvas) inBounds(x, y int) bool {
	return (x >= 0 && x < c.width*dotsPerCellX && y >= 0 && y < c.height*dotsPerCellY)
}

func (c *Canvas) SetDot(x, y int) bool {
	if c.inBounds(x, y) {
		*c.dotAt(x, y) = true
		return true
	}
	return false
}

type SetColourFunction func(x, y int, col colorful.Color)

func (c *Canvas) SetDotWithColor(x, y int, col colorful.Color) {
	if c.SetDot(x, y) {
		// remember we can only colour full cells
		cellX := x / dotsPerCellX
		cellY := y / dotsPerCellY
		c.foreground[cellY*c.width+cellX] = col
	}
}

func (c *Canvas) SetBackground(x, y int, col colorful.Color) {
	if c.inBounds(x, y) {
		cellX := x / dotsPerCellX
		cellY := y / dotsPerCellY
		c.background[cellY*c.width+cellX] = col
	}
}

func (c *Canvas) FillBackground() {
	for y := 0; y < c.DotHeight(); y++ {
		yRatio := float64(y) / float64(c.height)
		col := WaterColor(yRatio)
		for x := 0; x < c.DotWidth(); x++ {
			c.SetBackground(x, y, col)
		}
	}
}

func (c *Canvas) Clear() {
	clear(c.dots)
	c.FillBackground()
}

func (c *Canvas) Render() string {
	output := strings.Builder{}
	// output.Grow((c.width + 1) * c.height * 3) // necessary? (3B UTF)
	for cellY := 0; cellY < c.height; cellY++ {
		for cellX := 0; cellX < c.width; cellX++ {
			bits := byte(0)
			for row := 0; row < dotsPerCellY; row++ {
				for col := 0; col < dotsPerCellX; col++ {
					x := cellX*dotsPerCellX + col
					y := cellY*dotsPerCellY + row
					if *c.dotAt(x, y) {
						bits |= pixelMap[row][col]
					}
				}
			}
			var rn rune
			// foreground
			if bits != 0 {
				col := c.foreground[cellY*c.width+cellX]
				r, g, b := col.RGB255()
				fmt.Fprintf(&output, "\033[38;2;%d;%d;%dm", r, g, b) // Foreground RGB escape sequence
				rn = rune(0x2800 + int(bits))
			} else {
				rn = ' '
			}
			// background
			col := c.background[cellY*c.width+cellX]
			r, g, b := col.RGB255()
			if !(r == 0 && g == 0 && b == 0) {
				fmt.Fprintf(&output, "\033[48;2;%d;%d;%dm", r, g, b)
			}
			// write
			output.WriteRune(rn)
			output.WriteString("\033[38;2;0m") // Clear colour
		}
		output.WriteString("\n")
	}
	return output.String()
}

