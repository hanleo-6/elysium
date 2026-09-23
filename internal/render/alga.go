package render

import (
	"github.com/hanleo-6/elysium/internal/sim"
)

func (c *Canvas) DrawAlga(alga *sim.Alga) {
	if alga.Eaten {
		return
	}
	color := AlgaColor()
	c.DrawCircleFill(int(alga.Pos.X), int(alga.Pos.Y), int(alga.Radius), color)
}
