package render

import (
	"github.com/lucasb-eyer/go-colorful"
)

// DrawCircle draws using the Midpoint Circle Algorithm
// Adapted from https://en.wikipedia.org/wiki/Midpoint_circle_algorithm
func (c *Canvas) DrawCircle(centerX, centerY int, r int, col colorful.Color) {
	x := r
	y := 0
	p := 1 - r
	deltas := [4][2]int { {1, 1}, {1, -1}, {-1, 1}, {-1, -1} }
	for x >= y {
		for _, d := range deltas {
			a := d[0]*x
			b := d[1]*y
			c.SetDotWithColor(centerX+a, centerY+b, col)
			c.SetDotWithColor(centerX+b, centerY+a, col)
		}
		if p > 0 {
			x--
			p += 1 - 2*x + 2*y
		} else if p <= 0 {
			y += 1
			p += 1 + 2*y
		}
	}
}

func (c* Canvas) DrawLine(x0, y0, x1, y1 int, col colorful.Color) {
	// Single quadrant:
	// dx := x1-x0
	// dy := y1-y0
	// err := 2*dy-dx // P_0, err for the next pixel
	// y := y0
	// for x := x0; x <= x1; x++ {
	// 	c.SetDotWithColor(x, y, col)
	// 	if (err < 0) {
	// 		err = err + 2*dy
	// 		// choose y_k
	// 	} else {
	// 		err = err + 2*dy - 2*dx
	// 		// choose y_k + 1
	// 		y += 1
	// 	}
	// }

	dx := abs(x1 - x0)
	dy := abs(y1 - y0)
	// set step direction for x and y
	sx := -1
	if x0 < x1 {
		sx = 1
	}
	sy := -1
	if y0 < y1 {
		sy = 1
	}
	err := dx - dy

	for {
		c.SetDotWithColor(x0, y0, col)
		if (x0 == x1 && y0 == y1) {
			break
		}
		e2 := 2 * err
		// determine which step,
		// then pay step cost for that direction
		if e2 > -dy {
			x0 += sx
			err -= dy
		}
		if e2 < dx { // NOTE: can use else if here for a more 'joined' line look
			y0 += sy
			err += dx
		}
	}
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

