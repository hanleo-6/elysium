package render

import (
)

// DrawCircle draws using the Midpoint Circle Algorithm
// Adapted from https://en.wikipedia.org/wiki/Midpoint_circle_algorithm
func (c *Canvas) DrawCircle(centerX, centerY int, r int) {
	x := r
	y := 0
	p := 1 - r
	deltas := [4][2]int { {1, 1}, {1, -1}, {-1, 1}, {-1, -1} }
	for x >= y {
		for _, d := range deltas {
			a := d[0]*x+r
			b := d[1]*y+r
			c.SetDot(centerX+a, centerY+b)
			c.SetDot(centerX+b, centerY+a)
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

