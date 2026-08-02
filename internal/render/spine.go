package render

import (
	"github.com/hanleo-6/elysium/internal/physics"
	"github.com/lucasb-eyer/go-colorful"
)

func (c *Canvas) DrawSpine(skelly *physics.Skeleton, col colorful.Color) {
	for _, constraint := range skelly.Constraints {
		x0 := int(constraint.P1.Pos.X)
		y0 := int(constraint.P1.Pos.Y)
		x1 := int(constraint.P2.Pos.X)
		y1 := int(constraint.P2.Pos.Y)
		c.DrawLine(x0, y0, x1, y1, col)
	}

	// temporarily emphasising the head
	if len(skelly.Points) > 0 {
		head := skelly.Points[0]
		c.DrawCircle(int(head.Pos.X), int(head.Pos.Y), 1.0, col)
	}
}

