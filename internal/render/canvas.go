package render

import(
	"strings"
)

const (
	dotsPerCellX = 2
	dotsPerCellY = 4
)

type Canvas struct {
	width, height int // in terminal character cells
	dots []bool // size (width * 2) x (height * 4)
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
	}
}

func (c *Canvas) SetDot(x, y int) {
	if (x > 0 && x < c.width*dotsPerCellX && y > 0 && y < c.height*dotsPerCellY) {
		*c.dotAt(x, y) = true
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
			output.WriteRune(rune(0x2800 + int(bits)))
		}
		output.WriteString("\n")
	}
	return output.String()
}

