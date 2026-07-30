package render

import (
	"math"
)

func (c *Canvas) DrawCircle(centerX, centerY int, radius float64) {
	for angle := 0.0; angle < 2*math.Pi; angle += 0.05 {
		x := int(float64(centerX) + radius*math.Cos(angle))
		y := int(float64(centerY) + radius*math.Sin(angle))
		c.SetDot(x, y)
	}
}

