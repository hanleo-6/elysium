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

type Canvas struct {
	width, height int // in terminal character cells
	dots          []bool // size (width * 2) x (height * 4)
	colors        []colorful.Color // current colour for each cell
}

// using unicode braille pixel map for detail
var pixelMap = [4][2]byte{
	{0x01, 0x08},
	{0x02, 0x10},
	{0x04, 0x20},
	{0x40, 0x80},
}

func New(width, height int) *Canvas {
	return &Canvas{
		width:  width,
		height: height,
		dots:   make([]bool, width*dotsPerCellX*height*dotsPerCellY),
		colors: make([]colorful.Color, width*height),
	}
}

func (c *Canvas) DotWidth() int { return dotsPerCellX * c.width }

func (c *Canvas) DotHeight() int { return dotsPerCellY * c.height }

func (c *Canvas) SetDot(x, y int) bool {
	if (x >= 0 && x < c.width*dotsPerCellX && y >= 0 && y < c.height*dotsPerCellY) {
		*c.dotAt(x, y) = true
		return true
	}
	return false
}

func (c *Canvas) SetDotWithColor(x, y int, col colorful.Color) {
	if c.SetDot(x, y) {
		// remember we can only colour full cells
		cellX := x / dotsPerCellX
		cellY := y / dotsPerCellY
		c.colors[cellY*c.width+cellX] = col
	}
}

func (c *Canvas) dotAt(x, y int) *bool {
	return &c.dots[y*c.width*dotsPerCellX + x]
}

func (c *Canvas) Clear() {
	clear(c.dots)
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
			if bits != 0 {
				col := c.colors[cellY*c.width+cellX]
				r, g, b := col.RGB255()
				fmt.Fprintf(&output, "\033[38;2;%d;%d;%dm", r, g, b)
				output.WriteRune(rune(0x2800 + int(bits)))
				output.WriteString("\033[38;2;0m")
			} else {
				output.WriteRune(' ')
			}
		}
		output.WriteString("\n")
	}
	return output.String()
}

